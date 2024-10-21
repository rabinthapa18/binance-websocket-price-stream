package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

// for keeping track of WebSocket clients
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan float64)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // accept all requests
	},
}

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
func processMessage(conn *websocket.Conn) {
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

			averagePrice := calculateAverage(orderBook)

			// DEBUG: print the average price upto 2 decimal places
			log.Printf("Average Price: %.2f\n", averagePrice)

			// broadcasting average price to all clients
			broadcast <- averagePrice
		}
	}
}

// function to calculate average price
func calculateAverage(orderBook OrderBook) float64 {
	var totalPrice float64
	var totalCount int

	for _, bid := range orderBook.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		totalPrice += price
		totalCount++
	}

	for _, ask := range orderBook.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		totalPrice += price
		totalCount++
	}

	return totalPrice / float64(totalCount)
}

// function to handle new WebSocket clients
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading to WebSocket:", err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		_, _, err := ws.ReadMessage() // keeping the connection alive
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

// function to broadcast messages to clients
func handleMessages() {
	for {
		avgPrice := <-broadcast
		for client := range clients {
			err := client.WriteJSON(avgPrice)
			if err != nil {
				log.Printf("Error writing message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	conn := connectToBinance()

	go processMessage(conn)
	go handleMessages()

	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
