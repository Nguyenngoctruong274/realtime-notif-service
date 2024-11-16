package request

import "time"

type Notification struct {
	ID             uint      `json:"id"`
	NotificationID uint      `json:"notificationID"`
	ObjectID       string    `json:"objectID"`
	ProfileID      string    `json:"profileID"`
	Content        string    `json:"content"`
	Type           string    `json:"type"`
	Priority       string    `json:"priority"`
	IsMarked       bool      `json:"isMarked"`
	FromUser       string    `json:"fromUser"`
	IsLink         bool      `json:"isLink"`
	CreatedAt      time.Time `json:"createdAt"`
	ImageUrl       string    `json:"imageUrl"`
	Description    string    `json:"description"`
}

type MessageRequest struct {
	Username string         `json:"username"`
	Action   string         `json:"action"`
	Data     []Notification `json:"data"`
	Message  *string        `json:"message"`
	Unread   *int           `json:"unread,omitempty"`
}

type MarkAsReadRequest struct {
	Username string       `json:"username"`
	Action   string       `json:"action"`
	Data     Notification `json:"data"`
}

type JWTRequest struct {
	Username string `json:"username"`
}
