package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

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

func (c *ChatService) ChatInteractive(category string) error {
	for {
		message := Input(fmt.Sprintf("[%s] >> ", category))
		if err := c.Chat(category, message); err != nil {
			return fmt.Errorf("failed to chat: %w", err)
		}
	}
}

func (c *ChatService) Chat(category string, message string) error {
	summaries, err := c.storage.ListChatSummary(category, 10)
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

	if err := c.storage.AddChatMessage(apiMsg); err != nil {
		return fmt.Errorf("failed to save chat history: %w", err)
	}

	go func() {
		summary, err := c.chatGpt.Summary(userMsg, apiMsg)
		if err != nil {
			fmt.Printf("failed to call ChatGPT API: %v\n", err)
			return
		}
		s := &models.ChatSummary{
			Category: category,
			Summary:  summary,
		}
		if err := c.storage.AddSummary(s); err != nil {
			fmt.Printf("failed to save summary: %v\n", err)
			return
		}
	}()

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

func Input(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(nnSplitFunc)
	scanner.Scan()
	text := scanner.Text()
	return text
}

func nnSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	nn := "\n\n"
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	idx := bytes.Index(data, []byte(nn))
	if idx >= 0 {
		return idx + len(nn), data[0:idx], nil
	}
	return 0, nil, nil
}
