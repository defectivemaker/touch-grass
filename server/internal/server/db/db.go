package db

import (
    "context"
    "database/sql"
    "fmt"
    "os"
    "strconv"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    // "github.com/joho/godotenv"
    "github.com/gorilla/websocket"
    "time"
    
    "server-indicum/internal/common"
    "server-indicum/internal/server/ws"
)

// Use capital letter so it can be used in db-update.go
var Pool *pgxpool.Pool

func InitDB() error {

    config, err := pgx.ParseConfig("")
    if err != nil {
        return fmt.Errorf("Failed to parse config: %s", err)
    }
    config.User = os.Getenv("DB_USER")
    config.Password = os.Getenv("DB_PASS")
    config.Host = os.Getenv("DB_HOST")

    port, err := strconv.ParseUint(os.Getenv("DB_PORT"), 10, 16)
    if err != nil {
        return fmt.Errorf("Failed to parse port: %s", err)
    }
    config.Port = uint16(port)
    config.Database = os.Getenv("PGDATABASE")
    config.TLSConfig = nil // Set to nil to disable SSL, or set up proper TLS config

    pgxpoolConfig, err := pgxpool.ParseConfig("")
    if err != nil {
        return fmt.Errorf("Failed to parse Pool config: %s", err)
    }
    pgxpoolConfig.ConnConfig = config

    ctx := context.Background()
    Pool, err = pgxpool.NewWithConfig(ctx, pgxpoolConfig)
    if err != nil {
        return fmt.Errorf("Failed to connect to DB: %s", err)
    }

    if err := Pool.Ping(ctx); err != nil {
        return fmt.Errorf("Failed to ping DB: %s", err)
    } else {
        fmt.Println("Successfully pinged the DB")
    }

    // db = stdlib.OpenDB(*config)

    return nil
}

// I have created a notification channel in the database that is triggered
// when something is added to the entries table. This waits for a notification
// and then sends a message to the websocket connection with the UUID that was added in the DB
func ListenForDBInserts() {
    conn, err := Pool.Acquire(context.Background())
    if err != nil {
        fmt.Printf("Failed to acquire connection: %v", err)
        return
    }
    defer conn.Release()
    _, err = conn.Exec(context.Background(), "LISTEN new_entry")
    if err != nil {
        fmt.Printf("Failed to execute LISTEN command: %v", err)
        return
    }

    for {
        notification, err := conn.Conn().WaitForNotification(context.Background())
        if err != nil {
            fmt.Printf("Error waiting for notification: %v", err)
            time.Sleep(5 * time.Second)
            continue
        }
        fmt.Printf("New entry added: Channel %s Payload %s\n", notification.Channel, notification.Payload)
        
        // check if there is a websocket with the key of the uuid
        // if so this means that the user is online and they just added a new entry
        ws.WSConnMutex.Lock()
        if ws, ok := ws.WSConnections[notification.Payload]; ok {
            if err := ws.WriteMessage(websocket.TextMessage, []byte(notification.Payload)); err != nil {
                fmt.Println("Error sending message:", err)
                // Handle error, possibly close and delete connection
            }
        } else {
            fmt.Println("No connection found for UUID:", notification.Payload)
        }
        ws.WSConnMutex.Unlock()

    }
}


func waitForNotification() (string, error) {
    var notification string
    // Wait for a notification
    err := Pool.QueryRow(context.Background(), "SELECT pg_notify").Scan(&notification)
    fmt.Println("notification:", notification)
    return notification, err
}

func DBInitializeUserStatistics(uuid string) error {
    _, err := Pool.Exec(context.Background(), `
        INSERT INTO user_statistics (
            user_id, total_payphones, total_entries, total_maps,
            payphone_rank, entry_rank, map_rank, last_updated
        ) VALUES ($1, 0, 0, 0, 0, 0, 0, NOW())
    `, uuid)
    return err
}

func DBGetStatistics(uuid string) (common.Statistics, error) {
    var stats common.Statistics

    var LastUpdated float64
    err := Pool.QueryRow(context.Background(), `
        SELECT 
            total_payphones,
            total_entries,
            total_maps,
            payphone_rank,
            entry_rank,
            map_rank,
            EXTRACT(EPOCH from last_updated)
        FROM user_statistics
        WHERE user_id = $1
    `, uuid).Scan(
        &stats.TotalPayphones,
        &stats.TotalEntries,
        &stats.TotalMaps,
        &stats.PayphoneRank,
        &stats.EntryRank,
        &stats.MapRank,
        &LastUpdated,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return common.Statistics{}, fmt.Errorf("User statistics not found for UUID: %s", uuid)
        }
        return common.Statistics{}, fmt.Errorf("Failed to get statistics: %v", err)
    }
    stats.LastUpdated = int64(LastUpdated)
    // get total number of users
    err = Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
    if err != nil {
        return common.Statistics{}, fmt.Errorf("Failed to get total number of users: %v", err)
    }
    return stats, nil
}

func DBGetRecentEntry(uuid string) (common.Entry, error) {
    var entry common.Entry

    row := Pool.QueryRow(context.Background(), `
        SELECT id, deviceUUID, payphoneID, payphoneMAC, payphoneTime, EXTRACT(EPOCH FROM recordedTime) 
        FROM entries
        WHERE recordedTime > NOW() - INTERVAL '10 minutes'
        AND recordedTime < NOW()
        AND deviceUUID = $1
        ORDER BY id DESC
        LIMIT 1`, uuid)

    var recordedTime float64
    if err := row.Scan(&entry.ID, &entry.DeviceUUID, &entry.PayphoneID, &entry.PayphoneMAC, &entry.PayphoneTime, &recordedTime); err != nil {
        return common.Entry{}, fmt.Errorf("No entries in the last 10 minutes")
    }
    entry.RecordedTime = int64(recordedTime)

    return entry, nil
}

func DBGetNearbyHotspots(lat, lon float64) ([]common.DataPoint, error) {
    var dataPoints []common.DataPoint
    fmt.Println("lat and lon", lat, lon)
    query := `SELECT ST_X(location::geometry) as Latitude, ST_Y(location::geometry) as Longitude, uuid, street_address 
              FROM telstra_hotspots 
              WHERE ST_DWithin(location::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 1000);`

    rows, err := Pool.Query(context.Background(), query, lon, lat)
    if err != nil {
        return nil, fmt.Errorf("Failed to retrieve nearby hotspots: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var dp common.DataPoint
        if err := rows.Scan(&dp.Point.Lat, &dp.Point.Long, &dp.UUID, &dp.Address); err != nil {
            return nil, fmt.Errorf("Failed to scan nearby hotspot: %v", err)
        }
        dataPoints = append(dataPoints, dp)
    }

    return dataPoints, nil
}

func AddEntryToDB(entry common.Payload, deviceUUID string) (int64, error) {
    var id int64
    err := Pool.QueryRow(context.Background(), `INSERT INTO entries (deviceUUID, payphoneID, payphoneMAC, payphoneTime, recordedTime) 
                        VALUES ($1, $2, $3, $4, TO_TIMESTAMP($5)) RETURNING id`, deviceUUID, entry.PayphoneID, entry.PayphoneMAC, entry.PayphoneTime, entry.Time).Scan(&id)
    if err != nil {
        return 0, fmt.Errorf("Failed to insert entry: %v\n", err)
    }
    return id, nil
}

func DBGetLeaderboard() ([]common.LeaderboardVal, error) {
    var leaderboard []common.LeaderboardVal

    sqlQuery := `SELECT deviceUUID, COUNT(DISTINCT payphoneID) AS uniquePayphoneCount
    FROM entries
    GROUP BY deviceUUID
    ORDER BY uniquePayphoneCount DESC;`

    rows, err := Pool.Query(context.Background(), sqlQuery)
    if err != nil {
        return nil, fmt.Errorf("Failed to retrieve leaderboard: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var val common.LeaderboardVal
        if err := rows.Scan(&val.UUID, &val.TotalPoints); err != nil {
            return nil, fmt.Errorf("Failed to scan entry: %v", err)
        }
        leaderboard = append(leaderboard, val)
    }

    fmt.Println("printing leaderboard")
    fmt.Println(leaderboard)

    return leaderboard, nil
}

func DBGetEntriesWithUUID(deviceUUID string) ([]common.Entry, error) {
    var entries []common.Entry
    rows, err := Pool.Query(context.Background(), `
        SELECT id, deviceUUID, payphoneID, payphoneMAC, payphoneTime, 
               EXTRACT(EPOCH FROM recordedTime), mapUUID, 
               mapLatitude, mapLongitude, ST_AsText(mapLocation) as mapLocationText
        FROM entries 
        WHERE deviceUUID = $1`, deviceUUID)
    if err != nil {
        return nil, fmt.Errorf("Failed to retrieve entries: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var e common.Entry
        var recordedTime float64
        var sqlMapUUID, sqlMapLatitude, sqlMapLongitude, sqlMapLocationText sql.NullString
        if err := rows.Scan(
            &e.ID, 
            &e.DeviceUUID, 
            &e.PayphoneID, 
            &e.PayphoneMAC, 
            &e.PayphoneTime, 
            &recordedTime, 
            &sqlMapUUID, 
            &sqlMapLatitude,
            &sqlMapLongitude,
            &sqlMapLocationText); err != nil {
            return nil, fmt.Errorf("Failed to scan entry: %v", err)
        }
        e.RecordedTime = int64(recordedTime)
        if sqlMapUUID.Valid {
            e.MapUUID = sqlMapUUID.String
        }
        if sqlMapLatitude.Valid {
            e.MapLatitude = sqlMapLatitude.String
        }
        if sqlMapLongitude.Valid {
            e.MapLongitude = sqlMapLongitude.String
        }
        if sqlMapLocationText.Valid {
            e.MapLocation = sqlMapLocationText.String
        }
        entries = append(entries, e)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("Failed during rows iteration: %v", err)
    }
    return entries, nil
}

func DBFindDeviceRSAPub(deviceUUID string) ([]byte, error) {
    var deviceRSAPub []byte
    err := Pool.QueryRow(context.Background(), `SELECT pub_key FROM users WHERE uuid = $1`, deviceUUID).Scan(&deviceRSAPub)
    if err != nil {
        return nil, fmt.Errorf("Can't retrieve pub key %v\n", err)
    }

    return deviceRSAPub, nil
}

func DBGetRandomPoint() (common.DataPoint, error) {
    var dataPoint common.DataPoint

    query := `SELECT ST_X(location::geometry) as Latitude, ST_Y(location::geometry) as Longitude, uuid, street_address 
              FROM telstra_hotspots 
              ORDER BY RANDOM() 
              LIMIT 1;`

    row := Pool.QueryRow(context.Background(), query)

    var lat, lon float64
    err := row.Scan(&lat, &lon, &dataPoint.UUID, &dataPoint.Address)
    if err != nil {
        return common.DataPoint{}, err
    }

    dataPoint.Point = common.Coord{Lat: lat, Long: lon}

    return dataPoint, nil
}

func DBAddUser(uuid, email string) error {
    _, err := Pool.Exec(context.Background(), "INSERT INTO users (uuid, email) VALUES ($1, $2)", uuid, email)
    if err != nil {
        return fmt.Errorf("Failed to insert user: %v", err)
    }
    fmt.Println("This was successful")
    return nil
}

func DBGetUserProfile(uuid string) (map[string]string, error) {
    var (
        email    sql.NullString
        username sql.NullString
        token    sql.NullString
    )
    
    err := Pool.QueryRow(context.Background(), "SELECT email, username, token FROM users WHERE uuid = $1", uuid).Scan(&email, &username, &token)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("no user found with uuid %s", uuid)
    } else if err != nil {
        return nil, fmt.Errorf("failed to query profile for user %s: %v", uuid, err)
    }
    
    profile := make(map[string]string)
    profile["email"] = email.String
    profile["username"] = username.String
    profile["uuid"] = uuid
    
    if !token.Valid {
        newToken, err := common.GenerateRandomString(20)
        if err != nil {
            return nil, fmt.Errorf("failed to generate token for user %s: %v", uuid, err)
        }
        
        _, err = Pool.Exec(context.Background(), "UPDATE users SET token = $1 WHERE uuid = $2", newToken, uuid)
        if err != nil {
            return nil, fmt.Errorf("failed to update token for user %s: %v", uuid, err)
        }
        
        profile["token"] = newToken
        fmt.Printf("Generated and updated token for user %s\n", uuid)
    } else {
        profile["token"] = token.String
    }
    
    return profile, nil
}

func DBSavePubKey(token string, pubKey string) (string, error) {
    fmt.Println("pubkey", pubKey)
    fmt.Println("token", token)
    // pem decode string into BLOB
    pubKeyByte := []byte(pubKey)
    // save pubkey to table users where uuid = uuid
    query := "UPDATE users SET pub_key = $1 WHERE token = $2"

    // Execute the query with the provided public key and token
    _, err := Pool.Exec(context.Background(), query, pubKeyByte, token)
    if err != nil {
        fmt.Println(err)
        return "", fmt.Errorf("failed to save public key: %v", err)
    }

    // get uuid
    var uuid string
    query = "SELECT uuid FROM users WHERE token = $1"
    err = Pool.QueryRow(context.Background(), query, token).Scan(&uuid)

    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("No UUID found with the given token")
        } else {
            return "", fmt.Errorf("%v", err)
        }
    }

    return uuid, nil
}

func DBAddMapUUIDEntry(entryID int, mapUUID string) error {
    fmt.Println("testing", entryID, mapUUID)
    _, err := Pool.Exec(context.Background(), `UPDATE entries SET mapUUID = $1 WHERE id = $2`, mapUUID, entryID)

    if err != nil {return fmt.Errorf("Failed to update map %s :%v", mapUUID, err) }
    fmt.Println("Succesfully added to db")
    return nil
}

func DBUpdateLocation(entryID int, latitude, longitude string) error {
    if latitude == "0" && longitude == "0" {
        _, err := Pool.Exec(context.Background(), `UPDATE entries SET mapLatitude = NULL, mapLongitude = NULL, mapLocation = NULL WHERE id = $1`, entryID)
        if err != nil {
            return fmt.Errorf("error clearing location for entry: %v", err)
        }
        fmt.Println("Location cleared for entry ID:", entryID)
    } else {
        point := fmt.Sprintf("POINT(%s %s)", longitude, latitude) 
        _, err := Pool.Exec(context.Background(), `UPDATE entries SET mapLatitude = $1, mapLongitude = $2, mapLocation = ST_PointFromText($3) WHERE id = $4`, latitude, longitude, point, entryID)
        if err != nil {
            return fmt.Errorf("error updating entry: %v", err)
        }
        fmt.Println("Location updated for entry ID:", entryID)
    }
    return nil
}