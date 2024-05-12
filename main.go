package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, continuing with system env")
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func DataProcessor(input []byte) ([]byte, error) {
	return input, nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		processedMessage, err := DataProcessor(message)
		if err != nil {
			log.Printf("error processing data: %s", err)
			continue
		}

		if err := ws.WriteMessage(websocket.TextMessage, processedMessage); err != nil {
			log.Printf("error sending message: %s", err)
			continue
		}
	}
}

func main() {
	loadEnv()
	http.HandleFunc("/", handleConnections)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}