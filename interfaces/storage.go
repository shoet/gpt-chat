package interfaces

import "github.com/shoet/gpt-chat/models"

type Storage interface {
	AddChatMessage(message *models.ChatMessage) error
}
