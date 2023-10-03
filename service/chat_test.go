package service

import (
	"testing"

	"github.com/shoet/gpt-chat/interfaces"
	"github.com/shoet/gpt-chat/models"
	"github.com/shoet/gpt-chat/testutil"
)

func TestChatService_Chat(t *testing.T) {
	wantChatReq := &models.ChatMessage{
		Category: "golang",
		Content:  "hello go",
		Role:     "user",
	}
	wantChatRes := &models.ChatMessage{
		Id:       1,
		Category: "golang",
		Content:  "response message",
		Role:     "api",
	}
	mockGpt := &interfaces.ChatGPTMock{}
	mockGpt.ChatFunc = func(message *models.ChatMessage) (*models.ChatMessage, error) {
		if err := testutil.AssertObject(t, message, wantChatReq); err != nil {
			t.Errorf("failed to assert chat request: %v", err)
		}
		return wantChatRes, nil
	}

	mockStorage := &interfaces.StorageMock{}
	calledSaveChatHistory := 0
	mockStorage.SaveChatHistoryFunc = func(message *models.ChatMessage) error {
		calledSaveChatHistory += 1
		if message.Role == "user" {
			if err := testutil.AssertObject(t, message, wantChatReq); err != nil {
				t.Errorf("failed to assert chat request: %v", err)
			}
		} else {
			if err := testutil.AssertObject(t, message, wantChatRes); err != nil {
				t.Errorf("failed to assert chat request: %v", err)
			}
		}
		return nil
	}

	s, err := NewChatService(mockGpt, mockStorage)
	if err != nil {
		t.Fatalf("failed to create chat service: %v", err)
	}

	if err := s.Chat("golang", "hello go"); err != nil {
		t.Errorf("failed to chat: %v", err)
	}

	if calledSaveChatHistory != 2 {
		t.Errorf("failed to save chat history: %v", err)
	}

}
