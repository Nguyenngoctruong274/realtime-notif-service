package websocket_service

import (
	"github.com/gorilla/websocket"
	"time"
)

// Client is a type that describe the clients' username and their connection
type Clients map[string]map[*websocket.Conn]bool

// Notification is a data for feature Notification-Managements
type Notification struct {
	ID        uint      `json:"id"`
	ObjectID  string    `json:"objectID"`
	ProfileID string    `json:"profileID"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Priority  string    `json:"priority"`
	IsMarked  bool      `json:"isMarked"`
	FromUser  string    `json:"fromUser"`
	CreatedAt time.Time `json:"createdAt"`
}
