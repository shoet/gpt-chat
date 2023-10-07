package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/shoet/gpt-chat/interfaces"
	"github.com/shoet/gpt-chat/models"
)

type ChatGPTService struct {
	apiKey string
	client interfaces.Client
}

func NewChatGPTService(apiKey string, client interfaces.Client) *ChatGPTService {
	return &ChatGPTService{
		apiKey: apiKey,
		client: client,
	}
}

func (c *ChatGPTService) Chat(input *models.ChatMessage) (*models.ChatMessage, error) {
	req, err := c.buildRequestWithStream(input.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	resp, err := c.executeRequestWithStream(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	respMessage := models.ChatMessage{
		Category: input.Category,
		Message:  string(resp),
		Role:     "system",
	}
	return &respMessage, nil
}

func (c *ChatGPTService) buildRequestWithStream(input string) (*http.Request, error) {
	messages := []models.ChatGPTRequestMessage{
		{Role: "user", Content: input},
		{Role: "system", Content: "ユーザーからの要求分に最も適した回答を提供して下さい。"},
	}
	requestBody := models.ChatGPTRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
		Stream:   true,
	}
	b, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer([]byte(b)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	return req, nil
}

func (c *ChatGPTService) executeRequestWithStream(req *http.Request) ([]byte, error) {
	header := "data:"
	var buffer []byte
	chunkedCallback := func(b []byte) error {
		w := os.Stdout
		if strings.HasPrefix(string(b), header) {
			b = []byte(strings.TrimPrefix(string(b), header))
		}
		b = []byte(strings.TrimSpace(string(b)))
		if string(b) == "[DONE]" {
			return nil
		}

		var resp models.ChatGPTResponse
		if err := json.Unmarshal(b, &resp); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		wb := []byte(resp.Choices[0].Delta.Content)
		if _, err := w.Write(wb); err != nil {
			return err
		}
		buffer = append(buffer, wb...)
		return nil
	}

	sseClient := &SSEClient{client: c.client}
	if err := sseClient.Request(req, "\n\n", chunkedCallback); err != nil {
		return nil, fmt.Errorf("failed to read SSE: %w", err)
	}

	return buffer, nil
}
