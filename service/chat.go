package service

import (
	"fmt"

	"github.com/shoet/gpt-chat/interfaces"
	"github.com/shoet/gpt-chat/models"
)

type ChatService struct {
	ChatGPT interfaces.ChatGPT
	Storage interfaces.Storage
}

func NewChatService(
	chatGpt interfaces.ChatGPT,
	storage interfaces.Storage,
) (*ChatService, error) {
	chat := &ChatService{}
	return chat, nil
}

func (c *ChatService) Chat(category string, message string) error {
	// TODO: start chat intaractive
	userMsg := &models.ChatMessage{
		Category: category,
		Content:  message,
		Role:     "user",
	}
	apiMsg, err := c.ChatGPT.Chat(userMsg)
	if err != nil {
		return fmt.Errorf("failed to call ChatGPT API: %w", err)
	}

	// add storage
	if err := c.Storage.SaveChatHistory(userMsg); err != nil {
		return fmt.Errorf("failed to save chat history: %w", err)
	}
	if err := c.Storage.SaveChatHistory(apiMsg); err != nil {
		return fmt.Errorf("failed to save chat history: %w", err)
	}
	return nil

}
