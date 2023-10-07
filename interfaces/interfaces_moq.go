// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package interfaces

import (
	"github.com/shoet/gpt-chat/models"
	"sync"
)

// Ensure, that ChatGPTMock does implement ChatGPT.
// If this is not the case, regenerate this file with moq.
var _ ChatGPT = &ChatGPTMock{}

// ChatGPTMock is a mock implementation of ChatGPT.
//
//	func TestSomethingThatUsesChatGPT(t *testing.T) {
//
//		// make and configure a mocked ChatGPT
//		mockedChatGPT := &ChatGPTMock{
//			ChatFunc: func(message *models.ChatMessage, option *models.ChatMessageOption) (*models.ChatMessage, error) {
//				panic("mock out the Chat method")
//			},
//		}
//
//		// use mockedChatGPT in code that requires ChatGPT
//		// and then make assertions.
//
//	}
type ChatGPTMock struct {
	// ChatFunc mocks the Chat method.
	ChatFunc func(message *models.ChatMessage, option *models.ChatMessageOption) (*models.ChatMessage, error)

	// calls tracks calls to the methods.
	calls struct {
		// Chat holds details about calls to the Chat method.
		Chat []struct {
			// Message is the message argument value.
			Message *models.ChatMessage
			// Option is the option argument value.
			Option *models.ChatMessageOption
		}
	}
	lockChat sync.RWMutex
}

// Chat calls ChatFunc.
func (mock *ChatGPTMock) Chat(message *models.ChatMessage, option *models.ChatMessageOption) (*models.ChatMessage, error) {
	if mock.ChatFunc == nil {
		panic("ChatGPTMock.ChatFunc: method is nil but ChatGPT.Chat was just called")
	}
	callInfo := struct {
		Message *models.ChatMessage
		Option  *models.ChatMessageOption
	}{
		Message: message,
		Option:  option,
	}
	mock.lockChat.Lock()
	mock.calls.Chat = append(mock.calls.Chat, callInfo)
	mock.lockChat.Unlock()
	return mock.ChatFunc(message, option)
}

// ChatCalls gets all the calls that were made to Chat.
// Check the length with:
//
//	len(mockedChatGPT.ChatCalls())
func (mock *ChatGPTMock) ChatCalls() []struct {
	Message *models.ChatMessage
	Option  *models.ChatMessageOption
} {
	var calls []struct {
		Message *models.ChatMessage
		Option  *models.ChatMessageOption
	}
	mock.lockChat.RLock()
	calls = mock.calls.Chat
	mock.lockChat.RUnlock()
	return calls
}

// Ensure, that StorageMock does implement Storage.
// If this is not the case, regenerate this file with moq.
var _ Storage = &StorageMock{}

// StorageMock is a mock implementation of Storage.
//
//	func TestSomethingThatUsesStorage(t *testing.T) {
//
//		// make and configure a mocked Storage
//		mockedStorage := &StorageMock{
//			AddChatMessageFunc: func(message *models.ChatMessage) error {
//				panic("mock out the AddChatMessage method")
//			},
//			ListChatSummaryFunc: func(latest int) ([]string, error) {
//				panic("mock out the ListChatSummary method")
//			},
//		}
//
//		// use mockedStorage in code that requires Storage
//		// and then make assertions.
//
//	}
type StorageMock struct {
	// AddChatMessageFunc mocks the AddChatMessage method.
	AddChatMessageFunc func(message *models.ChatMessage) error

	// ListChatSummaryFunc mocks the ListChatSummary method.
	ListChatSummaryFunc func(latest int) ([]string, error)

	// calls tracks calls to the methods.
	calls struct {
		// AddChatMessage holds details about calls to the AddChatMessage method.
		AddChatMessage []struct {
			// Message is the message argument value.
			Message *models.ChatMessage
		}
		// ListChatSummary holds details about calls to the ListChatSummary method.
		ListChatSummary []struct {
			// Latest is the latest argument value.
			Latest int
		}
	}
	lockAddChatMessage  sync.RWMutex
	lockListChatSummary sync.RWMutex
}

// AddChatMessage calls AddChatMessageFunc.
func (mock *StorageMock) AddChatMessage(message *models.ChatMessage) error {
	if mock.AddChatMessageFunc == nil {
		panic("StorageMock.AddChatMessageFunc: method is nil but Storage.AddChatMessage was just called")
	}
	callInfo := struct {
		Message *models.ChatMessage
	}{
		Message: message,
	}
	mock.lockAddChatMessage.Lock()
	mock.calls.AddChatMessage = append(mock.calls.AddChatMessage, callInfo)
	mock.lockAddChatMessage.Unlock()
	return mock.AddChatMessageFunc(message)
}

// AddChatMessageCalls gets all the calls that were made to AddChatMessage.
// Check the length with:
//
//	len(mockedStorage.AddChatMessageCalls())
func (mock *StorageMock) AddChatMessageCalls() []struct {
	Message *models.ChatMessage
} {
	var calls []struct {
		Message *models.ChatMessage
	}
	mock.lockAddChatMessage.RLock()
	calls = mock.calls.AddChatMessage
	mock.lockAddChatMessage.RUnlock()
	return calls
}

// ListChatSummary calls ListChatSummaryFunc.
func (mock *StorageMock) ListChatSummary(latest int) ([]string, error) {
	if mock.ListChatSummaryFunc == nil {
		panic("StorageMock.ListChatSummaryFunc: method is nil but Storage.ListChatSummary was just called")
	}
	callInfo := struct {
		Latest int
	}{
		Latest: latest,
	}
	mock.lockListChatSummary.Lock()
	mock.calls.ListChatSummary = append(mock.calls.ListChatSummary, callInfo)
	mock.lockListChatSummary.Unlock()
	return mock.ListChatSummaryFunc(latest)
}

// ListChatSummaryCalls gets all the calls that were made to ListChatSummary.
// Check the length with:
//
//	len(mockedStorage.ListChatSummaryCalls())
func (mock *StorageMock) ListChatSummaryCalls() []struct {
	Latest int
} {
	var calls []struct {
		Latest int
	}
	mock.lockListChatSummary.RLock()
	calls = mock.calls.ListChatSummary
	mock.lockListChatSummary.RUnlock()
	return calls
}
