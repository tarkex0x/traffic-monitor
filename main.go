package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/websocket"
    "github.com/joho/godotenv"
)

func loadEnvironmentVariables() {
    err := godotenv.Load()
    if err != nil {
        logOutput("Error loading .env file, continuing with environment variables from the system")
    }
}

var websocketUpgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func processReceivedData(input []byte) ([]byte, error) {
    return input, nil
}

func websocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
    wsConn, err := websocketUpgrader.Upgrade(w, r, nil)
    if err != nil {
        logOutput("WebSocket Upgrade error: %v", err)
        return // changed from log.Fatal to prevent server stop
    }
    defer wsConn.Close()

    for {
        _, message, err := wsConn.ReadMessage()
        if err != nil {
            logOutput("WebSocket read error: %v", err)
            break
        }

        processedMessage, err := processReceivedData(message)
        if err != nil {
            logOutput("Error processing data: %s", err)
            continue
        }

        if err := wsConn.WriteMessage(websocket.TextMessage, processedMessage); err != nil {
            logOutput("WebSocket write error: %s", err)
            continue
        }
    }
}

func logOutput(format string, v ...interface{}) {
    // Custom log formatting can be added here
    log.Printf(format, v...)
}

func main() {
    loadEnvironmentVariables()
    http.HandleFunc("/", websocketConnectionHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    logOutput("Starting WebSocket server on :%s", port)
    err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
    if err != nil {
        logOutput("ListenAndServe error: %v", err)
    }
}