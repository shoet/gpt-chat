package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

func (c *ChatGPTService) Chat(input *models.ChatMessage, option *models.ChatMessageOption) (*models.ChatMessage, error) {
	req, err := c.buildRequestWithStream(input, option)
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

func (c *ChatGPTService) buildRequestWithStream(input *models.ChatMessage, option *models.ChatMessageOption) (*http.Request, error) {
	systemTemplate, err := LoadChatSystemTemplate(option.Summaries)
	if err != nil {
		return nil, fmt.Errorf("failed to load system template: %w", err)
	}
	messages := []models.ChatGPTRequestMessage{
		{Role: "user", Content: input.Message},
		{Role: "system", Content: systemTemplate},
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

func LoadChatSystemTemplate(chatSummaries []string) (string, error) {
	systemTemplate := struct {
		ChatHistory string
	}{
		ChatHistory: strings.Join(chatSummaries, "\n"),
	}
	templateTxt, err := LoadChatTemplate("chatgpt/system.txt")
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}
	t, err := template.New("system").Parse(templateTxt)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}
	var b bytes.Buffer
	t.Execute(&b, systemTemplate)
	return b.String(), nil
}

func LoadChatTemplate(templateName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	templateDir := filepath.Join(cwd, "templates")
	b, err := os.ReadFile(filepath.Join(templateDir, templateName))
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}
	return string(b), nil
}
