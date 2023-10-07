package service

import (
	"encoding/json"
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
	summaries, err := c.storage.ListChatSummary(10)
	if err != nil {
		return fmt.Errorf("failed to load chat history: %w", err)
	}
	userMsg := &models.ChatMessage{
		Category: category,
		Message:  message,
		Role:     "user",
	}
	option := &models.ChatMessageOption{
		Summaries: summaries,
	}
	apiMsg, err := c.chatGpt.Chat(userMsg, option)
	if err != nil {
		return fmt.Errorf("failed to call ChatGPT API: %w", err)
	}
	if err := c.storage.AddChatMessage(userMsg); err != nil {
		return fmt.Errorf("failed to save chat history: %w", err)
	}

	summary, err := ParseSummary(apiMsg)
	if err != nil {
		return fmt.Errorf("failed to parse summary: %w", err)
	}
	apiMsg.Summary = summary
	if err := c.storage.AddChatMessage(apiMsg); err != nil {
		return fmt.Errorf("failed to save chat history: %w", err)
	}
	return nil
}

func ParseSummary(m *models.ChatMessage) (string, error) {
	s := struct {
		Summary string `json:"summary"`
	}{}
	if err := json.Unmarshal([]byte(m.Message), &s); err != nil {
		return "", fmt.Errorf("failed to parse summary: %w", err)
	}
	return s.Summary, nil
}
