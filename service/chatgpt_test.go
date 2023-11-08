package service

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/shoet/gpt-chat/models"
)

func Test_ChatGPTService_Summary(t *testing.T) {
	apiKey, ok := os.LookupEnv("CHATGPT_API_SECRET")
	if !ok {
		t.Fatal("CHATGPT_API_SECRET is not set.")
	}

	client := http.Client{}

	sut := NewChatGPTService(apiKey, &client)
	category := "挨拶"
	request := models.ChatMessage{
		Category: category,
		Message:  "あなたの名前はなんですか？",
	}
	answer := models.ChatMessage{
		Category: category,
		Message:  "私はGPTです。",
	}
	summary, err := sut.Summary(&request, &answer)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(summary)
}

func Test_buildChatRequestWithStream(t *testing.T) {

	client := http.Client{}
	c := NewChatGPTService("test", &client)
	m := &models.ChatMessage{
		Category: "golang",
		Message:  "test",
		Role:     "user",
	}
	o := &models.ChatMessageOption{
		LatestHistory: []*models.ChatMessage{
			&models.ChatMessage{
				Category: "golang",
				Message:  "request",
				Role:     "user",
			},
			&models.ChatMessage{
				Category: "golang",
				Message:  "response",
				Role:     "system",
			},
		},
	}
	req, err := c.buildChatRequestWithStream(m, o)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		t.Fatalf("failed to dump request: %v", err)
	}
	fmt.Println(string(dump))

}
