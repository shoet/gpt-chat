package service

import (
	"fmt"
	"net/http"
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
