// api.go
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

func sendEventToWebSocketServer(event string, message string) {
	url := "http://localhost:8080/send-event"
	payload := []byte(fmt.Sprintf(`{"event": "%s", "message": "%s"}`, event, message))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println("HTTP request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("HTTP request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("HTTP request failed:", resp.Status)
		return
	}
}

func main() {
	// Example: Send an event to WebSocket server
	sendEventToWebSocketServer("event-channel", "Notification message for event1")
}
