package interfaces

import "github.com/shoet/gpt-chat/models"

type Storage interface {
	AddChatMessage(message *models.ChatMessage) error
	ListChatSummary(category string, latest int) ([]*models.ChatSummary, error)
	ListChatHistory(category string, latest int) (models.ChatMessages, error)
	AddSummary(summary *models.ChatSummary) error
}
