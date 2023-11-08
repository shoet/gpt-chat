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

func (r *Repository) AddSummary(db Execer, summary *models.ChatSummary) (models.ChatSummaryId, error) {
	sql := `
	INSERT INTO chat_summary 
		(category, summary, created, modified)
	VALUES 
		(?, ?, ?, ?);
	`
	now := r.Clocker.Now()
	res, err := db.Exec(sql, summary.Category, summary.Summary, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to insert chat summary: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return models.ChatSummaryId(id), nil
}

func (r *Repository) ListChatSummary(db Queryer, category string, latest int) ([]*models.ChatSummary, error) {
	sql := `
	SELECT
		id, category, summary, created, modified
	FROM
		chat_summary
	WHERE
		summary IS NOT NULL AND summary != ''
		AND category = ?
	ORDER BY 
		id DESC
	LIMIT ?
	;
	`
	var rows []*models.ChatSummary
	err := db.Select(&rows, sql, category, latest)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat summaries: %w", err)
	}
	return rows, nil
}

func (r *Repository) ListChatHistory(
	db Queryer, category string, latest int,
) (models.ChatMessages, error) {
	sql := `
	SELECT
		id, category, message, role, created, modified
	FROM
		chat_message
	WHERE
		category = ?
	ORDER BY
		created DESC
	LIMIT ?
	;
	`
	var rows models.ChatMessages
	err := db.Select(&rows, sql, category, latest)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat message: %w", err)
	}
	return rows, nil
}

type Execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

type Queryer interface {
	Select(dest interface{}, query string, args ...interface{}) error
}
