package service

import (
	"github.com/jmoiron/sqlx"
	"github.com/shoet/gpt-chat/models"
)

type ChatHistoryAdder interface {
	AddChatHistory(db *sqlx.DB, chatHistory *models.ChatHistory) error
}
