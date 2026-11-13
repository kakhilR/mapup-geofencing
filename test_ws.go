package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(
		"ws://localhost:8080/ws/alerts",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Println("Connected to WebSocket")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Received:", string(message))
	}
}