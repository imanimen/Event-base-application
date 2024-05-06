package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {return true},
	}

	// Mutex for thread-safe access to clients map
	mu      sync.Mutex
	clients = make(map[*websocket.Conn]struct{})

	// Slice to store events
	events []EEvent
)

type EEvent struct {
	ID        int       `json:"id"`
	Event     string    `json:"event"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Add the connection to the clients map
	mu.Lock()
	clients[conn] = struct{}{}
	mu.Unlock()

	// Send existing events to the client
	for _, event := range events {
		if err := conn.WriteJSON(event); err != nil {
			log.Println("Write error:", err)
			break
		}
	}

	// Handle incoming messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		// Additional handling of incoming messages if needed
	}

	// Remove the connection from the clients map when it's closed
	mu.Lock()
	delete(clients, conn)
	mu.Unlock()
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	type Event struct {
		Event   string `json:"event"`
		Message string `json:"message"`
	}

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Process the event (e.g., send it to WebSocket clients)
	log.Printf("Received event: %s - %s\n", event.Event, event.Message)

	// Prepare the JSON response
	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",                                                                                                                                                                                                
		Message: "Event received successfully",
	}

	// Encode the response into JSON and write it to the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Assign a unique ID to the event
	eventID := len(events) + 1
	e := EEvent{
		Event:   event.Event,
		Message: event.Message,
	}
	e.ID = eventID

	// Set the timestamp to the current time
	e.Timestamp = time.Now()

	// Store the event
	events = append(events, e)

	// Broadcast the event to WebSocket clients
	broadcastEvent(e)
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	// Return the list of events as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func broadcastEvent(event EEvent) {
	// Broadcast the event to all WebSocket clients
	for client := range clients {
		if err := client.WriteJSON(event); err != nil {
			log.Println("Write error:", err)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/send-event", handleEvent)
	http.HandleFunc("/list-events", listEvents)

	log.Println("WebSocket server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// )

// var (
// 	upgrader = websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin: func(r *http.Request) bool { return true },
// 	}

// 	// Mutex for thread-safe access to clients map
// 	mu      sync.Mutex
// 	clients = make(map[*websocket.Conn]struct{})

// 	// Slice to store events
// 	events []EEvent
// )

// type EEvent struct {
// 	ID        int       `json:"id"`
// 	Event     string    `json:"event"`
// 	Message   string    `json:"message"`
// 	Timestamp time.Time `json:"timestamp"`
// }

// func handleWebSocket(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("Upgrade error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// Add the connection to the clients map
// 	mu.Lock()
// 	clients[conn] = struct{}{}
// 	mu.Unlock()

// 	// Send existing events to the client
// 	for _, event := range events {
// 		if err := conn.WriteJSON(event); err != nil {
// 			log.Println("Write error:", err)
// 			break
// 		}
// 	}

// 	// Handle incoming messages
// 	for {
// 		_, _, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Read error:", err)
// 			break
// 		}
// 		// Additional handling of incoming messages if needed
// 	}

// 	// Remove the connection from the clients map when it's closed
// 	mu.Lock()
// 	delete(clients, conn)
// 	mu.Unlock()
// }

// func handleEvent(c *gin.Context) {
// 	type Event struct {
// 		Event   string `json:"event"`
// 		Message string `json:"message"`
// 	}

// 	var event Event
// 	if err := c.ShouldBindJSON(&event); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  "error",
// 			"message": "Invalid JSON",
// 		})
// 		return
// 	}

// 	// Process the event (e.g., send it to WebSocket clients)
// 	log.Printf("Received event: %s - %s\n", event.Event, event.Message)

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "success",
// 		"message": "Event received successfully",
// 	})

// 	// Assign a unique ID to the event
// 	eventID := len(events) + 1
// 	e := EEvent{
// 		Event:   event.Event,
// 		Message: event.Message,
// 	}
// 	e.ID = eventID

// 	// Set the timestamp to the current time
// 	e.Timestamp = time.Now()

// 	// Store the event
// 	events = append(events, e)

// 	// Broadcast the event to WebSocket clients
// 	broadcastEvent(e)
// }

// func listEvents(c *gin.Context) {
// 	// Return the list of events as JSON
// 	c.JSON(http.StatusOK, events)
// }

// func broadcastEvent(event EEvent) {
// 	// Broadcast the event to all WebSocket clients
// 	for client := range clients {
// 		if err := client.WriteJSON(event); err != nil {
// 			log.Println("Write error:", err)
// 		}
// 	}
// }

// func main() {
// 	r := gin.Default()

// 	r.GET("/ws", handleWebSocket)
// 	r.POST("/send-event", handleEvent)
// 	r.GET("/list-events", listEvents)

// 	log.Println("WebSocket server listening on port 8080...")
// 	r.Run(":8080")
// }