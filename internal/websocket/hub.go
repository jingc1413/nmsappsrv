package websocket

import (
	"sync"

	"nmsappsrv/pkg/logger"

	"github.com/gorilla/websocket"
)

// Hub manages all connected WebSocket clients and message broadcasting.
type Hub struct {
	// All connected clients
	clients map[*Client]bool
	// Broadcast channel for messages sent to all clients
	broadcast chan []byte
	// Register channel for new client connections
	register chan *Client
	// Unregister channel for client disconnections
	unregister chan *Client
	// Mutex for thread-safe access to clients map
	mu sync.RWMutex
}

// Client represents a single WebSocket connection.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	id   string
	// Topics this client is subscribed to
	topics map[string]bool
	// Mutex for thread-safe access to topics
	topicsMu sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main event loop. Should be called in a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logger.Infof("websocket: client %s connected, total clients: %d", client.id, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			logger.Infof("websocket: client %s disconnected, total clients: %d", client.id, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client's send buffer is full, disconnect it
					h.mu.RUnlock()
					h.mu.Lock()
					if _, ok := h.clients[client]; ok {
						delete(h.clients, client)
						close(client.send)
					}
					h.mu.Unlock()
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMessage sends a message to all connected clients.
func (h *Hub) BroadcastMessage(message []byte) {
	select {
	case h.broadcast <- message:
	default:
		logger.Errorf("websocket: broadcast channel full, dropping message")
	}
}

// BroadcastToTopic sends a message only to clients subscribed to the given topic.
func (h *Hub) BroadcastToTopic(topic string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		client.topicsMu.RLock()
		subscribed := client.topics[topic]
		client.topicsMu.RUnlock()

		if subscribed {
			select {
			case client.send <- message:
			default:
				logger.Errorf("websocket: client %s send buffer full, skipping", client.id)
			}
		}
	}
}

// SubscribeTopic subscribes a client to a topic.
func (c *Client) SubscribeTopic(topic string) {
	c.topicsMu.Lock()
	defer c.topicsMu.Unlock()
	c.topics[topic] = true
	logger.Infof("websocket: client %s subscribed to topic %s", c.id, topic)
}

// UnsubscribeTopic unsubscribes a client from a topic.
func (c *Client) UnsubscribeTopic(topic string) {
	c.topicsMu.Lock()
	defer c.topicsMu.Unlock()
	delete(c.topics, topic)
	logger.Infof("websocket: client %s unsubscribed from topic %s", c.id, topic)
}
