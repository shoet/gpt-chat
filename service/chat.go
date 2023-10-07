package service

import (
	"fmt"

	"github.com/shoet/gpt-chat/interfaces"
	"github.com/shoet/gpt-chat/models"
)

type ChatService struct {
	chatGpt interfaces.ChatGPT
	storage interfaces.Storage
}

func NewChatService(
	chatGpt interfaces.ChatGPT,
	storage interfaces.Storage,
) (*ChatService, error) {
	chat := &ChatService{
		chatGpt: chatGpt,
		storage: storage,
	}
	return chat, nil
}

func (c *ChatService) Chat(category string, message string) error {
	// TODO: start chat intaractive
	userMsg := &models.ChatMessage{
		Category: category,
		Message:  message,
		Role:     "user",
	}
	apiMsg, err := c.chatGpt.Chat(userMsg)
	if err != nil {
		return fmt.Errorf("failed to call ChatGPT API: %w", err)
	}

	// add storage
	if err := c.storage.SaveChatHistory(userMsg); err != nil {
		return fmt.Errorf("failed to save chat history: %w", err)
	}
	if err := c.storage.SaveChatHistory(apiMsg); err != nil {
		return fmt.Errorf("failed to save chat history: %w", err)
	}
	return nil

}
