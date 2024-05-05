# Event-base-application

# Run server
`go run server/server.go`

# Send event via client
`go run client/api.go`

# Send via API
`http://localhost:8080/send-event`

# Paylaod
``
{
	"event": "event-channel-1", "message": "sample event"
}
``

# Server the html via VScode live server to show the realtime events