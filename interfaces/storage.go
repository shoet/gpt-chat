package interfaces

import "github.com/shoet/gpt-chat/models"

type Storage interface {
	SaveChatHistory(message *models.ChatMessage) error
}
