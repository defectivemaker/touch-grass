package db

import (
    "context"
    "fmt"
	"time"
)

// User represents a user in the users table
type UserID struct {
    UUID string
}
func PeriodicDBUpdate() {
    // Start the background job to update user statistics every 3 hours
    go updateUserStatistics()

    select {}
}

func updateUserStatistics() {
    for {
        // Calculate and update user statistics
        err := calculateAndUpdateUserStatistics()
        if err != nil {
            fmt.Printf("Error updating user statistics: %v\n", err)
        } else {
            fmt.Println("User statistics updated successfully")
        }

        // Sleep for 3 hours before the next update
        time.Sleep(1 * time.Hour)
    }
}

func getAllUsers() ([]UserID, error) {
    query := "SELECT uuid FROM users"
    rows, err := Pool.Query(context.Background(), query)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %v", err)
    }
    defer rows.Close()

    var users []UserID
    for rows.Next() {
        var user UserID
        if err := rows.Scan(&user.UUID); err != nil {
            return nil, fmt.Errorf("failed to scan row: %v", err)
        }
        users = append(users, user)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to iterate over rows: %v", err)
    }

    return users, nil
}

func calculateTotalCounts(userUUID string) (int, int, int, error) {
    query := `
        SELECT 
            COUNT(DISTINCT payphoneID) AS total_payphones,
            COUNT(*) AS total_entries,
            COUNT(DISTINCT mapLocation) AS total_maps
        FROM entries
        WHERE deviceUUID = $1
    `
    var totalPayphones, totalEntries, totalMaps int
    err := Pool.QueryRow(context.Background(), query, userUUID).Scan(&totalPayphones, &totalEntries, &totalMaps)
    if err != nil {
        return 0, 0, 0, fmt.Errorf("failed to execute query: %v", err)
    }

    return totalPayphones, totalEntries, totalMaps, nil
}

func calculateRanks(userUUID string, totalPayphones, totalEntries, totalMaps int) (int, int, int, error) {
    query := `
        SELECT
            (SELECT COUNT(*) + 1 FROM user_statistics WHERE total_payphones > $1) AS payphone_rank,
            (SELECT COUNT(*) + 1 FROM user_statistics WHERE total_entries > $2) AS entry_rank,
            (SELECT COUNT(*) + 1 FROM user_statistics WHERE total_maps > $3) AS map_rank
    `
    var payphoneRank, entryRank, mapRank int
    err := Pool.QueryRow(context.Background(), query, totalPayphones, totalEntries, totalMaps).Scan(&payphoneRank, &entryRank, &mapRank)
    if err != nil {
        return 0, 0, 0, fmt.Errorf("failed to execute query: %v", err)
    }

    return payphoneRank, entryRank, mapRank, nil
}

func updateUserStatisticsInDB(userUUID string, totalPayphones, totalEntries, totalMaps, payphoneRank, entryRank, mapRank int) error {
    query := `
        INSERT INTO user_statistics (user_id, total_payphones, total_entries, total_maps, payphone_rank, entry_rank, map_rank, last_updated)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
        ON CONFLICT (user_id) DO UPDATE SET
            total_payphones = $2,
            total_entries = $3,
            total_maps = $4,
            payphone_rank = $5,
            entry_rank = $6,
            map_rank = $7,
            last_updated = NOW()
    `
    _, err := Pool.Exec(context.Background(), query, userUUID, totalPayphones, totalEntries, totalMaps, payphoneRank, entryRank, mapRank)
    if err != nil {
        return fmt.Errorf("failed to update user statistics: %v", err)
    }

    return nil
}

func calculateAndUpdateUserStatistics() error {
    // Get all users from the users table
    users, err := getAllUsers()
    if err != nil {
        return fmt.Errorf("failed to get users: %v", err)
    }

    // Iterate over each user
    for _, user := range users {
        // Calculate total counts for the user
        totalPayphones, totalEntries, totalMaps, err := calculateTotalCounts(user.UUID)
        if err != nil {
            return fmt.Errorf("failed to calculate total counts for user %s: %v", user.UUID, err)
        }

        // Calculate ranks for the user
        payphoneRank, entryRank, mapRank, err := calculateRanks(user.UUID, totalPayphones, totalEntries, totalMaps)
        if err != nil {
            return fmt.Errorf("failed to calculate ranks for user %s: %v", user.UUID, err)
        }

        // Update the user's statistics in the user_statistics table
        err = updateUserStatisticsInDB(user.UUID, totalPayphones, totalEntries, totalMaps, payphoneRank, entryRank, mapRank)
        if err != nil {
            return fmt.Errorf("failed to update user statistics for user %s: %v", user.UUID, err)
        }
    }

    return nil
}
