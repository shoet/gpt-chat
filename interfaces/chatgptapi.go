package interfaces

import "github.com/shoet/gpt-chat/models"

//go:generate go run github.com/matryer/moq -out chatgptapi_moq.go . ChatGPT
type ChatGPT interface {
	Chat(message *models.ChatMessage) (*models.ChatMessage, error)
}
