package interfaces

import "github.com/shoet/gpt-chat/models"

//go:generate go run github.com/matryer/moq -out storage_moq.go . Storage
type Storage interface {
	SaveChatHistory(message *models.ChatMessage) error
}
