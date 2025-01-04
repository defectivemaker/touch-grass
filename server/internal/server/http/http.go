// This file will have all the functions that handle the HTTP requests that are coming from the react server
package http

import (
	"net/http"
	"fmt"
    "os"
    "log"
    "encoding/json"
    "time"
    "strconv"
    "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/httplog/v2"
    "github.com/golang-jwt/jwt/v5"
    "context"
    "strings"

    "server-indicum/internal/server/db"
    "server-indicum/internal/server/ws"
)

type User struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
}

var logger = httplog.Options{
    Concise:        false,
    RequestHeaders: true,
    HideRequestHeaders: []string{
        "accept",
        "accept-encoding",
        "accept-language",
        "accept-ranges",
        "connection",
        "cookie",
        "user-agent",
    },
    QuietDownRoutes: []string{
        "/",
        "/ping",
    },
    QuietDownPeriod: 10 * time.Second,
}




func CORSHandler(allowedOrigin string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Println("request", r)
            w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
func JWTVerifier(jwtSecret []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tokenString := ""

            // Extract the token from the Authorization header
            authHeader := r.Header.Get("Authorization")
            if authHeader != "" {
                // Split the header on the space to separate "Bearer" and the token
                split := strings.SplitN(authHeader, " ", 2)
                if len(split) == 2 && split[0] == "Bearer" {
                    tokenString = strings.TrimSpace(split[1])
                }
            }

            // If not found in the Authorization header, try the query string
            if tokenString == "" {
                tokenString = r.URL.Query().Get("jwt")
                if tokenString != "" {
                    fmt.Println("Token retrieved from query string")
                }
            }

            
            // Validate token if it is present
            if tokenString != "" {
                token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                    // Make sure token's algorithm is what you expect
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
                    }
                    return jwtSecret, nil
                })

                if err != nil {
                    if (err.Error() == "token has invalid claims: token is expired") {
                        fmt.Println("Token expired")
                        http.Error(w, "Token expired", http.StatusUnauthorized)
                        return
                    }
                    fmt.Println("Failed to parse token:", err)
                    http.Error(w, "Invalid token", http.StatusUnauthorized)
                    return
                }

                if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
                    // Pass the claims to the next handler, for example
                    ctx := context.WithValue(r.Context(), "claims", claims)
                    next.ServeHTTP(w, r.WithContext(ctx))
                    return
                }
            }

            fmt.Println("Authorization failed, no valid token provided")
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
        })
    }
}


func HandleHTTPServer() {

    logLocation := os.Getenv("LOG_LOCATION")

    fileLog, err := os.OpenFile(logLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    defer fileLog.Close()


    logger := httplog.NewLogger("http-logger", httplog.Options{
        Concise:        true,
        RequestHeaders: true,
        HideRequestHeaders: []string{
            "accept",
            "accept-encoding",
            "accept-language",
            "accept-ranges",
            "connection",
            "cookie",
        },
        QuietDownRoutes: []string{
            "/",
            "/ping",
        },
        QuietDownPeriod: 10 * time.Second,
        Writer: fileLog,
    })

    // allowedOrigin := os.Getenv("FRONTEND_SERVER_ADDRESS")
    allowedOrigin := "*"

    jwtSecret := []byte(os.Getenv("SUPABASE_JWT_SECRET"))

    r := chi.NewRouter()
    // requestlogger automatically has requestID and recoverer middleware
    r.Use(httplog.RequestLogger(logger, []string{"/ping"}))
    r.Use(CORSHandler(allowedOrigin))
    r.Use(middleware.Heartbeat("/ping"))


    // this is insecure $$
    // need to fix by either having an API key in ansible script (this is still vulnerable to people on device hacking)
    // or something else. I'm tireed rn can't think
    r.Post("/map-token-pub-key", mapTokenPubKey)


    r.Group(func(r chi.Router) {
        r.Use(JWTVerifier(jwtSecret))

        r.Get("/protected", protectedEndpoint)
        // a uuid will be passed in, and its corresponding token will be retreived
        // if no link exists, a token will be generated
        r.Get("/get-profile", getProfile)
        r.Get("/random-point", returnRandomPoint)
        r.Get("/get-entries", getEntriesUUID)
        r.Get("/leaderboad", getLeaderboard)
        r.Get("/get-recent-entry", getRecentEntry)
        r.Get("/statistics", getStatistics)
        // a request will be sent with token and pubkey
        // this will then be saved in the db
        r.Post("/add-mapuuid", addMapUUIDEntry)
        r.Post("/add-location" , addLocationEntry)
        r.Post("/add-user" , addUser)
        r.Post("/nearby-hotspots", getNearbyHotspots)

        r.Get("/ws", ws.WebSocketHandler)
    })

    httpAddress := os.Getenv("HTTPLISTENADDRESS")

    certFile := os.Getenv("TLS_CERT_FILE")
    keyFile := os.Getenv("TLS_PRIV_KEY")

    _, err = os.ReadFile(certFile)
    if err != nil {
        log.Fatalf("Failed to read certificate file: %v", err)
    }

    _, err = os.ReadFile(keyFile)
    if err != nil {
        log.Fatalf("Failed to read key file: %v", err)
    }

    fmt.Println("HTTPS server listening on", httpAddress)
    if err := http.ListenAndServeTLS(httpAddress, certFile, keyFile, r); err != nil {
        log.Fatalf("Failed to start HTTPS server: %v", err)
    }

}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Protected endpoint"))
}

func getStatistics(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(jwt.MapClaims)
    if !ok {
        http.Error(w, "Could not get claims from context", http.StatusInternalServerError)
        return
    }
    uuid := claims["sub"].(string)

    stats, err := db.DBGetStatistics(uuid)
    if err != nil {
        http.Error(w, fmt.Sprintf("%v", err), http.StatusNotFound)
    }

    jsonResponse, err := json.Marshal(stats)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling JSON: %v", err), http.StatusInternalServerError)
    }
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func getRecentEntry(w http.ResponseWriter, r *http.Request) {

    claims, ok := r.Context().Value("claims").(jwt.MapClaims)
    if !ok {
        http.Error(w, "Could not get claims from context", http.StatusInternalServerError)
        return
    }
	uuid := claims["sub"].(string)

    entry, err := db.DBGetRecentEntry(uuid)
    if err != nil { 
        if err == fmt.Errorf("No entries in the last 10 minutes") {
            w.WriteHeader(http.StatusNoContent)
        }
    }

    jsonResponse, err := json.Marshal(entry)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling JSON: %v", err), http.StatusInternalServerError)
    }
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

func getNearbyHotspots(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")

    type requestBody struct {
        Latitude  string `json:"latitude"`
        Longitude string `json:"longitude"`
    }

    var body requestBody
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    lat, err := strconv.ParseFloat(body.Latitude, 64)
    if err != nil {
        http.Error(w, "Invalid latitude value", http.StatusBadRequest)
        return
    }

    lon, err := strconv.ParseFloat(body.Longitude, 64)
    if err != nil {
        http.Error(w, "Invalid longitude value", http.StatusBadRequest)
        return
    }
    fmt.Println("getting nearby hotspots")
    hotspots, err := db.DBGetNearbyHotspots(lat, lon)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    jsonResponse, err := json.Marshal(hotspots)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write(jsonResponse)
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
    
    w.Header().Set("Content-Type", "application/json") 

    leaderboard, err := db.DBGetLeaderboard()
    if err != nil {
        fmt.Println("can't get leaderboard %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println("Successfully got leaderboard %+v\n", leaderboard)
    jsonResponse, err := json.Marshal(leaderboard)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write(jsonResponse)

}

func addUser(w http.ResponseWriter, r *http.Request) {
    type requestBody struct {
        UUID  string `json:"id"`
        Email string `json:"email"`
    }

    var body requestBody
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Check if user already exists
    _, err = db.DBGetStatistics(body.UUID)
    if err == nil {
        // User already exists
        http.Error(w, "User already exists", http.StatusConflict)
        return
    } 
    // Add user
    err = db.DBAddUser(body.UUID, body.Email)
    if err != nil {
        http.Error(w, "Failed to add user", http.StatusInternalServerError)
        return
    }

    // Initialize user statistics
    err = db.DBInitializeUserStatistics(body.UUID)
    if err != nil {
        http.Error(w, "Failed to initialize user statistics", http.StatusInternalServerError)
        return
    }

    // Send success response
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User added successfully"})
}


func getProfile(w http.ResponseWriter, r *http.Request) {
    type ProfileResponse struct {
        Token    string `json:"token"`
        Email    string `json:"email"`
        Username string `json:"username"`
        UUID     string `json:"uuid"`
    }
    
    // Set CORS headers first
    w.Header().Set("Content-Type", "application/json") // Set the Content-Type as application/json
    
    claims, ok := r.Context().Value("claims").(jwt.MapClaims)
    if !ok {
        http.Error(w, "Could not get claims from context", http.StatusInternalServerError)
        return
    }
    
    uuid := claims["sub"].(string)
    profile, err := db.DBGetUserProfile(uuid)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        response := map[string]string{"error": "Error retrieving profile"}
        json.NewEncoder(w).Encode(response) // Send JSON error message
        log.Printf("Failed to retrieve profile: %v", err)
        return
    }
    
    response := ProfileResponse{
        Token:    profile["token"],
        Email:    profile["email"],
        Username: profile["username"],
        UUID:     profile["uuid"],
    }
    
    json.NewEncoder(w).Encode(response) // Encode the response as JSON and write it
}


func mapTokenPubKey(w http.ResponseWriter, r *http.Request) {

    
    w.Header().Set("Content-Type", "application/json") // Set the Content-Type as application/json


    if r.Method != "POST" {
        fmt.Println("Only accept POST requests", r)
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }

    // Define a structure to match the expected request body
    type requestBody struct {
        Token    string `json:"token"`
        PubKey  string `json:"pub_key"`
    }

    // Decode the JSON body of the request
    var body requestBody
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    // Call the db.SavePubKey method with the token and public key
    uuid, err := db.DBSavePubKey(body.Token, body.PubKey)
    if err != nil {
        // Log the error for internal debugging.
        log.Printf("Failed to save public key: %v", err)

        // Send a more generic error message to the client, or include the error if you decide it's safe.
        errMsg := fmt.Sprintf("Failed to save public key: %v", err) // Consider the security implications
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }

    type responseBody struct {
        UUID string `json:"uuid"`
    }

    w.WriteHeader(http.StatusOK)
    response := responseBody{UUID: uuid}
    json.NewEncoder(w).Encode(response)
    // Respond to the request indicating success
    // w.Write([]byte("Public key successfully mapped to token"))
}

func addMapUUIDEntry(w http.ResponseWriter, r *http.Request) {
    fmt.Println("adding map to entry");

    if r.Method != "POST" {
        fmt.Println("Only accept POST requests", r)
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }
    type inputData struct {
        EntryID int
        MapUUID string
    }

    var new inputData
    err := json.NewDecoder(r.Body).Decode(&new)
    fmt.Println("getting there")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    err = db.DBAddMapUUIDEntry(new.EntryID, new.MapUUID)
    if err != nil { 
        fmt.Println("Can't add to DB %v\n", err)
        http.Error(w, fmt.Sprintf("Failed to add to DB: %v", err), http.StatusBadRequest)
    }


    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message": "Entry updated successfully"}`))
}


func addLocationEntry(w http.ResponseWriter, r *http.Request) {
    fmt.Println("adding location to entry");

    if r.Method != "POST" {
        fmt.Println("Only accept POST requests", r)
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }
    type inputData struct {
        EntryID int
        Latitude string
        Longitude string
    }

    var new inputData
    err := json.NewDecoder(r.Body).Decode(&new)
    fmt.Println("getting there")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    err = db.DBUpdateLocation(new.EntryID, new.Latitude, new.Longitude)
    if err != nil { 
        fmt.Println("Can't add to DB %v\n", err)
        http.Error(w, fmt.Sprintf("Failed to add to DB: %v", err), http.StatusBadRequest)
    }


    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message": "Entry updated successfully"}`))
}

func getEntriesUUID(w http.ResponseWriter, r *http.Request) {

    claims, ok := r.Context().Value("claims").(jwt.MapClaims)
    if !ok {
        http.Error(w, "Could not get claims from context", http.StatusInternalServerError)
        return
    }
	uuid := claims["sub"].(string)

    if uuid == "" {
        http.Error(w, "uuid is required", http.StatusBadRequest)
        return
    }

    entries, err := db.DBGetEntriesWithUUID(uuid)
    if err != nil {
        fmt.Println("can't get entries %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println("Successfully got entries %+v\n", entries)
    jsonResponse, err := json.Marshal(entries)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println("entriees", jsonResponse)

    
    w.Write(jsonResponse)
}

func returnRandomPoint(w http.ResponseWriter, r *http.Request) {
    dataPoint, err := db.DBGetRandomPoint()
    if err != nil { 
        http.Error(w, fmt.Sprintf("Can't get point: %v", err), http.StatusInternalServerError)
    }

    jsonResponse, err := json.Marshal(dataPoint)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling JSON: %v", err), http.StatusInternalServerError)
    }
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)

}