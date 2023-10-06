package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (c *ChatGPTService) Chat(input string) ([]byte, error) {
	req, err := c.buildRequestWithStream(input)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	resp, err := c.executeRequestWithStream(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	return resp, nil
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

type SSEClient struct {
	client interfaces.Client
}

func (c *SSEClient) Request(req *http.Request, chunkSep string, chunkHandler func(b []byte) error) error {
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	chEve := make(chan string)
	chErr := make(chan error)
	scanner := SplitScanner(res.Body, chunkSep)
	go func() {
		for scanner.Scan() {
			b := scanner.Bytes()
			chEve <- string(b)
		}
		if err := scanner.Err(); err != nil {
			chErr <- err
			return
		}
		chErr <- io.EOF
	}()

	for {
		select {
		case err := <-chErr:
			if err == io.EOF {
				return nil
			}
			return err
		case event := <-chEve:
			if err := chunkHandler([]byte(event)); err != nil {
				return err
			}
		}
	}
}

func SplitScanner(r io.Reader, sep string) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	initBufferSize := 1024
	maxBufferSize := 4096
	scanner.Buffer(make([]byte, initBufferSize), maxBufferSize)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		beforeSep := bytes.Index(data, []byte(sep)) // 最初sepの直前
		if beforeSep >= 0 {
			// 最初のsepの位置, dataのsepの直前までのスライス, nil
			return beforeSep + len(sep), data[0:beforeSep], nil
		}
		if atEOF {
			// 残りのすべて
			return len(data), data, nil
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	return scanner
}
