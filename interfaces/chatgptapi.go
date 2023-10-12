package interfaces

import "github.com/shoet/gpt-chat/models"

type ChatGPT interface {
	Chat(message *models.ChatMessage, option *models.ChatMessageOption) (*models.ChatMessage, error)
	Summary(request *models.ChatMessage, answer *models.ChatMessage) (string, error)
}
