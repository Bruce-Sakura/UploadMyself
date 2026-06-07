package mapper

import (
	"context"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/messages/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MessageMapper is the data-access layer for the messages table.
type MessageMapper struct {
	pool *pgxpool.Pool
}

func NewMessageMapper(pool *pgxpool.Pool) *MessageMapper {
	return &MessageMapper{pool: pool}
}

func (m *MessageMapper) Insert(ctx context.Context, convID, role, content string) error {
	_, err := m.pool.Exec(ctx,
		`INSERT INTO messages (conversation_id, role, content) VALUES ($1, $2, $3)`,
		convID, role, content)
	return err
}

// LoadHistory returns the latest `limit` messages for a conversation in
// chronological (ascending) order.
func (m *MessageMapper) LoadHistory(ctx context.Context, convID string, limit int) ([]entity.Message, error) {
	rows, err := m.pool.Query(ctx,
		`SELECT id, conversation_id, role, content, created_at
		 FROM messages WHERE conversation_id = $1
		 ORDER BY created_at DESC LIMIT $2`, convID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	msgs, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (entity.Message, error) {
		var msg entity.Message
		err := r.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.CreatedAt)
		return msg, err
	})
	if err != nil {
		return nil, err
	}
	// Reverse into chronological order.
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}
