package entity

import "time"

// Message maps the messages table — agent conversation history.
type Message struct {
	ID             int64
	ConversationID string
	Role           string // system | user | assistant | tool
	Content        string
	CreatedAt      time.Time
}
