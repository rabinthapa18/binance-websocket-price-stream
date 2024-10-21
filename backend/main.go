package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// OrderBook struct to parse the message from Binance WebSocket
type OrderBook struct {
	LastUpdateID int64      `json:"u"`
	Bids         [][]string `json:"b"` // [["price", "quantity"]]
	Asks         [][]string `json:"a"` // [["price", "quantity"]]
}

// constants
const (
	binanceWebSocketEndpoint = "wss://stream.binance.com:9443/ws/btcusdt@depth"
	streamName               = "btcusdt@depth"
	responseStreamName       = "depthUpdate"
)

// function to connect to Binance WebSocket
func connectToBinance() *websocket.Conn {
	log.Printf("Connecting to Binance WebSocket at %s", binanceWebSocketEndpoint)

	conn, _, err := websocket.DefaultDialer.Dial(binanceWebSocketEndpoint, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	} else {
		log.Println("Connected to Binance WebSocket successfully.")
	}

	return conn
}

// function to read messages from Binance WebSocket
func readMessages(conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var data map[string]interface{}
		json.Unmarshal(message, &data)

		// checking if the message is coming from the stream we subscribed to
		if stream, ok := data["e"].(string); ok && stream == responseStreamName {

			// parsing the message into OrderBook struct
			var orderBook OrderBook
			json.Unmarshal(message, &orderBook)

			// DEBUG: printing the order book
			log.Println("Bids quantities:")
			for _, bid := range orderBook.Bids {
				price := bid[0]
				quantity := bid[1]
				log.Printf("Price: %s, Quantity: %s\n", price, quantity)
			}

			// DEBUG: printing the order book
			log.Println("Asks quantities:")
			for _, ask := range orderBook.Asks {
				price := ask[0]
				quantity := ask[1]
				log.Printf("Price: %s, Quantity: %s\n", price, quantity)
			}

		}
	}
}

func main() {
	conn := connectToBinance()
	readMessages(conn)
}
