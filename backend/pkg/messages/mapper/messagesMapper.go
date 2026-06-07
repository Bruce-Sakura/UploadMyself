package mapper

import (
	"context"
	"database/sql"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/entity"
)

// MessageMapper is the data-access layer for the messages table.
type MessageMapper struct {
	db *sql.DB
}

func NewMessageMapper(db *sql.DB) *MessageMapper {
	return &MessageMapper{db: db}
}

func (m *MessageMapper) Insert(ctx context.Context, convID, role, content string) error {
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO messages (conversation_id, role, content) VALUES (?, ?, ?)`,
		convID, role, content)
	return err
}

// LoadHistory returns the latest `limit` messages for a conversation in
// chronological (ascending) order.
func (m *MessageMapper) LoadHistory(ctx context.Context, convID string, limit int) ([]entity.Message, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, conversation_id, role, content, created_at
		 FROM messages WHERE conversation_id = ?
		 ORDER BY created_at DESC LIMIT ?`, convID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Reverse into chronological order.
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}
