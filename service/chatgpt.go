package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shoet/gpt-chat/interfaces"
	"github.com/shoet/gpt-chat/models"
	"github.com/shoet/gpt-chat/util"
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
	req, err := c.buildChatRequestWithStream(input, option)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	resp, err := c.executeChatRequestWithStream(req)
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

func (c *ChatGPTService) Summary(request *models.ChatMessage, answer *models.ChatMessage) (string, error) {
	req, err := c.buildSummaryRequest(request, answer)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}
	resp, err := c.client.Do(req)
	defer resp.Body.Close()
	chatGptResponse := &models.ChatGPTResponse{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	if err := json.Unmarshal(b, chatGptResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return chatGptResponse.Choices[0].Message.Content, nil
}

func (c *ChatGPTService) buildChatRequestWithStream(
	input *models.ChatMessage, option *models.ChatMessageOption,
) (*http.Request, error) {
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

func (c *ChatGPTService) executeChatRequestWithStream(req *http.Request) ([]byte, error) {
	header := "data:"
	var buffer []byte
	chunkedCallback := func(b []byte) error {
		w := os.Stdout
		if strings.HasPrefix(string(b), header) {
			b = []byte(strings.TrimPrefix(string(b), header))
		}
		b = []byte(strings.TrimSpace(string(b)))
		if string(b) == "[DONE]" {
			w.Write([]byte("\n"))
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

func (c *ChatGPTService) buildSummaryRequest(
	request *models.ChatMessage, answer *models.ChatMessage,
) (*http.Request, error) {
	summaryTemplate, err := LoadSummaryTemplate()
	if err != nil {
		return nil, fmt.Errorf("failed to load summary template: %w", err)
	}
	input := struct {
		User   string `json:"user"`
		System string `json:"system"`
	}{
		User:   request.Message,
		System: answer.Message,
	}
	jsonB, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	messages := []models.ChatGPTRequestMessage{
		{Role: "user", Content: string(jsonB)},
		{Role: "system", Content: summaryTemplate},
	}
	requestBody := models.ChatGPTRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
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

func LoadChatSystemTemplate(chatSummaries []*models.ChatSummary) (string, error) {
	histories := []string{}
	for i, _ := range chatSummaries {
		histories = append(
			histories,
			fmt.Sprintf("%d. %s", i+1, chatSummaries[len(chatSummaries)-i-1].Summary),
		)
	}
	systemTemplate := struct {
		ChatHistory string
	}{
		ChatHistory: strings.Join(histories, "\n"),
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

func LoadSummaryTemplate() (string, error) {
	templateTxt, err := LoadChatTemplate("chatgpt/summary.txt")
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}
	return templateTxt, nil
}

func LoadChatTemplate(templateName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	rootDir, err := util.GetProjectRoot(cwd)
	if err != nil {
		return "", fmt.Errorf("failed to get project root: %w", err)
	}

	templateDir := filepath.Join(rootDir, "templates")
	b, err := os.ReadFile(filepath.Join(templateDir, templateName))
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}
	return string(b), nil
}
