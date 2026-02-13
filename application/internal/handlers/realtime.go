package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "sync"
    "time"
)

// clientConnection is a simple wrapper for a channel used to send bytes to a client
type clientConnection struct {
    ch chan []byte
}

var (
    // map VIN -> list of client connections
    vinClients = make(map[string][]clientConnection)
    clientsMu  sync.Mutex
)

// Event is the required event envelope
type Event struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

// BroadcastToVIN sends an event to all connected clients subscribed to vin.
// Uses a mutex for thread-safety per requirements.
func BroadcastToVIN(vin string, eventType string, payload interface{}) {
    evt := Event{Type: eventType, Payload: payload}
    data, err := json.Marshal(evt)
    if err != nil {
        return
    }

    clientsMu.Lock()
    conns := make([]clientConnection, len(vinClients[vin]))
    copy(conns, vinClients[vin])
    clientsMu.Unlock()

    for _, c := range conns {
        // non-blocking send with goroutine to avoid blocking broadcaster
        go func(ch chan []byte) {
            select {
            case ch <- data:
            case <-time.After(1 * time.Second):
            }
        }(c.ch)
    }
}

// addClient registers a new client channel for a VIN
func addClient(vin string, c clientConnection) {
    clientsMu.Lock()
    defer clientsMu.Unlock()
    vinClients[vin] = append(vinClients[vin], c)
}

// removeClient removes a client channel for a VIN
func removeClient(vin string, ch chan []byte) {
    clientsMu.Lock()
    defer clientsMu.Unlock()
    conns := vinClients[vin]
    for i := 0; i < len(conns); i++ {
        if conns[i].ch == ch {
            // remove
            vinClients[vin] = append(conns[:i], conns[i+1:]...)
            break
        }
    }
    if len(vinClients[vin]) == 0 {
        delete(vinClients, vin)
    }
}

// SSEHandler handles Server-Sent Events subscriptions for a single VIN.
// Route: GET /ws/cars/{vin}
func SSEHandler(w http.ResponseWriter, r *http.Request) {
    vin := strings.TrimPrefix(r.URL.Path, "/ws/cars/")
    if vin == "" {
        http.Error(w, "missing vin", http.StatusBadRequest)
        return
    }

    // Set SSE headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming unsupported", http.StatusInternalServerError)
        return
    }

    ch := make(chan []byte, 10)
    conn := clientConnection{ch: ch}
    addClient(vin, conn)
    defer removeClient(vin, ch)

    // Close channel when done
    notify := r.Context().Done()

    // Send a welcome message
    welcome := Event{Type: "VEHICLE_STATE_UPDATED", Payload: map[string]string{"status": "connected"}}
    if b, err := json.Marshal(welcome); err == nil {
        fmt.Fprintf(w, "data: %s\n\n", b)
        flusher.Flush()
    }

    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-notify:
            return
        case msg := <-ch:
            fmt.Fprintf(w, "data: %s\n\n", msg)
            flusher.Flush()
        case <-ticker.C:
            // keep-alive comment
            fmt.Fprintf(w, ": ping\n\n")
            flusher.Flush()
        }
    }
}
