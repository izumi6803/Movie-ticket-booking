package websocket

import (
	"cinema-backend/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Cho phép tất cả origins trong development
		// Trong production, nên giới hạn origins
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client đại diện cho một WebSocket connection
type Client struct {
	ID         string
	Socket     *websocket.Conn
	Send       chan []byte
	ShowtimeID string // Showtime mà client đang xem
	Room       string // Room ID = showtimeID
}

// Hub quản lý tất cả clients và rooms
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	rooms      map[string]map[*Client]bool // roomID -> clients
	mutex      sync.RWMutex
}

// Message struct cho WebSocket messages
type Message struct {
	Type       string      `json:"type"`
	ShowtimeID string      `json:"showtimeId"`
	Data       interface{} `json:"data"`
	Timestamp  int64       `json:"timestamp"`
}

// SeatUpdateMessage gửi khi có thay đổi về ghế
type SeatUpdateMessage struct {
	SeatID    string `json:"seatId"`
	SeatLabel string `json:"seatLabel"`
	Status    string `json:"status"` // available, locked, booked
	UserID    string `json:"userId,omitempty"`
}

// NewHub tạo Hub mới
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

// Run bắt đầu Hub loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			// Thêm client vào room
			if client.Room != "" {
				if h.rooms[client.Room] == nil {
					h.rooms[client.Room] = make(map[*Client]bool)
				}
				h.rooms[client.Room][client] = true
			}
			h.mutex.Unlock()
			log.Printf("Client %s joined room %s", client.ID, client.Room)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				// Xóa client khỏi room
				if client.Room != "" {
					if room, ok := h.rooms[client.Room]; ok {
						delete(room, client)
						if len(room) == 0 {
							delete(h.rooms, client.Room)
						}
					}
				}
			}
			h.mutex.Unlock()
			log.Printf("Client %s left room %s", client.ID, client.Room)

		case message := <-h.broadcast:
			h.broadcastToRoom(message)
		}
	}
}

// broadcastToRoom gửi message đến tất cả clients trong room
func (h *Hub) broadcastToRoom(message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	room, ok := h.rooms[message.ShowtimeID]
	if !ok {
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range room {
		select {
		case client.Send <- data:
		default:
			// Channel đầy, xóa client
			close(client.Send)
			delete(room, client)
			delete(h.clients, client)
		}
	}
}

// BroadcastSeatUpdate gửi cập nhật ghế đến tất cả clients trong showtime
func (h *Hub) BroadcastSeatUpdate(showtimeID string, seatUpdate SeatUpdateMessage) {
	message := Message{
		Type:       "SEAT_UPDATE",
		ShowtimeID: showtimeID,
		Data:       seatUpdate,
		Timestamp:  0, // Sẽ được set khi marshal
	}
	h.broadcast <- message
}

// BroadcastSeatLock gửi thông báo ghế bị lock
func (h *Hub) BroadcastSeatLock(showtimeID string, seatIDs []string, seatLabels []string, userID string) {
	message := Message{
		Type:       "SEATS_LOCKED",
		ShowtimeID: showtimeID,
		Data: map[string]interface{}{
			"seatIds":    seatIDs,
			"seatLabels": seatLabels,
			"userId":     userID,
		},
	}
	h.broadcast <- message
}

// BroadcastSeatUnlock gửi thông báo ghế được unlock
func (h *Hub) BroadcastSeatUnlock(showtimeID string, seatIDs []string, seatLabels []string) {
	message := Message{
		Type:       "SEATS_UNLOCKED",
		ShowtimeID: showtimeID,
		Data: map[string]interface{}{
			"seatIds":    seatIDs,
			"seatLabels": seatLabels,
		},
	}
	h.broadcast <- message
}

// BroadcastBookingCreated gửi thông báo booking mới
func (h *Hub) BroadcastBookingCreated(showtimeID string, bookingID string, seatIDs []string) {
	message := Message{
		Type:       "BOOKING_CREATED",
		ShowtimeID: showtimeID,
		Data: map[string]interface{}{
			"bookingId": bookingID,
			"seatIds":   seatIDs,
		},
	}
	h.broadcast <- message
}

// GetRoomClientCount trả về số clients trong room
func (h *Hub) GetRoomClientCount(roomID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if room, ok := h.rooms[roomID]; ok {
		return len(room)
	}
	return 0
}

// WebSocketHandler xử lý WebSocket connections
type WebSocketHandler struct {
	hub         *Hub
	seatService *services.SeatService
}

func NewWebSocketHandler(hub *Hub, seatService *services.SeatService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		seatService: seatService,
	}
}

// HandleWebSocket xử lý upgrade HTTP -> WebSocket
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	showtimeID := c.Query("showtimeId")
	if showtimeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "showtimeId is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:         fmt.Sprintf("client_%d", len(h.hub.clients)),
		Socket:     conn,
		Send:       make(chan []byte, 256),
		ShowtimeID: showtimeID,
		Room:       showtimeID,
	}

	h.hub.register <- client

	// Start goroutines cho read và write
	go client.writePump()
	go client.readPump(h.hub, h.seatService)
}

// writePump gửi messages từ channel đến WebSocket
func (c *Client) writePump() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// readPump đọc messages từ WebSocket
func (c *Client) readPump(hub *Hub, seatService *services.SeatService) {
	defer func() {
		hub.unregister <- c
		c.Socket.Close()
	}()

	c.Socket.SetReadLimit(512)
	// c.Socket.SetReadDeadline(time.Now().Add(60 * time.Second))
	// c.Socket.SetPongHandler(func(string) error {
	// 	c.Socket.SetReadDeadline(time.Now().Add(60 * time.Second))
	// 	return nil
	// })

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Xử lý message từ client
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Xử lý các loại message
		switch msg.Type {
		case "PING":
			// Gửi PONG response
			pongMsg := Message{Type: "PONG", ShowtimeID: c.ShowtimeID}
			data, _ := json.Marshal(pongMsg)
			c.Send <- data

		case "GET_SEATS_STATUS":
			// Client yêu cầu status của tất cả ghế
			if seatService != nil && c.ShowtimeID != "" {
				// Lấy screenID từ showtimeID (cần implement)
				// Hiện tại gửi message xác nhận
				confirmMsg := Message{
					Type:       "SEATS_STATUS",
					ShowtimeID: c.ShowtimeID,
					Data:       map[string]string{"status": "subscribed"},
				}
				data, _ := json.Marshal(confirmMsg)
				c.Send <- data
			}
		}
	}
}
