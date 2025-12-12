package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"api/internal/sse"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SSEController struct {
	clients map[uuid.UUID]chan interface{}
	mutex   sync.RWMutex
}

// Ensure SSEController implements sse.Service interface
var _ sse.Service = (*SSEController)(nil)

func NewSSEController() *SSEController {
	return &SSEController{
		clients: make(map[uuid.UUID]chan interface{}),
	}
}

func (c *SSEController) RegisterClient(ctx *fiber.Ctx) error {
	clientID := uuid.New()

	// Create a new channel for this client
	ch := make(chan interface{}, 100)

	// Register the client
	c.mutex.Lock()
	c.clients[clientID] = ch
	c.mutex.Unlock()

	// Set headers for SSE
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("Access-Control-Allow-Origin", "*")

	// Send initial connection message
	c.sendEvent(ch, "connected", map[string]interface{}{
		"client_id": clientID,
		"timestamp": time.Now(),
	})

	// Keep connection alive and send events
	go func() {
		for {
			select {
			case data := <-ch:
				c.sendEventToContext(ctx, "update", data)
			case <-ctx.Context().Done():
				// Client disconnected
				c.mutex.Lock()
				delete(c.clients, clientID)
				c.mutex.Unlock()
				close(ch)
				return
			}
		}
	}()

	return ctx.SendStatus(http.StatusOK)
}

func (c *SSEController) BroadcastEvent(eventType sse.EventType, data interface{}) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, ch := range c.clients {
		select {
		case ch <- map[string]interface{}{
			"type": eventType,
			"data": data,
		}:
		default:
			// Channel is full, skip this client
			log.Printf("SSE channel full, skipping client")
		}
	}
}

func (c *SSEController) sendEvent(ch chan interface{}, eventType string, data interface{}) {
	select {
	case ch <- map[string]interface{}{
		"type": eventType,
		"data": data,
	}:
	default:
		log.Printf("SSE channel full, dropping event")
	}
}

func (c *SSEController) sendEventToContext(ctx *fiber.Ctx, eventType string, data interface{}) {
	event, ok := data.(map[string]interface{})
	if !ok {
		event = map[string]interface{}{
			"type": eventType,
			"data": data,
		}
	}

	// Format as SSE event
	ctx.WriteString("event: ")
	ctx.WriteString(eventType)
	ctx.WriteString("\ndata: ")
	jsonData, _ := json.Marshal(event)
	ctx.Write(jsonData)
	ctx.WriteString("\n\n")
}
