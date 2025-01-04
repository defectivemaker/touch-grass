package main

import (
    "os"
    "log"
    "os/signal"
    "server-indicum/internal/server/http"
    "server-indicum/internal/server/device"
    "server-indicum/internal/server/db"
)

func main(){
    
    err := db.InitDB()
    if err != nil { log.Fatalf("Failed to init DB: %v\n", err) } else {
        log.Println("DB initialized")
    }

    // need to set INDICUM_ENV=prod when running in prod
    // env := os.Getenv("INDICUM_ENV")

    // if env == "" {
    //     env = "dev"
    // }
    // // this is useful for when I need different vars for dev and prod
    // godotenv.Load(".env." + env + ".local")

    // envErr := godotenv.Load()
	// if envErr != nil {
	// 	log.Fatalf("Failed to load env %v\n" ,envErr)
	// }

    go http.HandleHTTPServer()
    go device.InitDeviceServer()
    go db.ListenForDBInserts()
    go db.PeriodicDBUpdate()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c
    log.Println("Graceful shutdown")
}

