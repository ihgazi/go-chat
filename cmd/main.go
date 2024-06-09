package main

import (
    "log"
    "github.com/ihgazi/go-chat/db"
)

func main() {
    _, err := db.NewDatabase()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
}
