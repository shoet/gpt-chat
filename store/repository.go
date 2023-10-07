package store

import (
	"database/sql"
	"fmt"

	"github.com/shoet/gpt-chat/clocker"
	"github.com/shoet/gpt-chat/models"
)

type Repository struct {
	Clocker clocker.Clocker
}

func NewRepository(c clocker.Clocker) (*Repository, error) {
	return &Repository{
		Clocker: c,
	}, nil
}

func (r *Repository) AddChatMessage(db Execer, message *models.ChatMessage) (models.ChatMessageId, error) {
	sql := `
	INSERT INTO chat_message 
		(category, message, role, created, modified)
	VALUES 
		(?, ?, ?, ?, ?);
	`
	now := r.Clocker.Now()
	res, err := db.Exec(sql, message.Category, message.Message, message.Role, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to insert chat message: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return models.ChatMessageId(id), nil
}

type Execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}
