package models

import (
	"sort"
	"time"
)

type ChatMessageId int

type ChatMessage struct {
	Id       ChatMessageId `db:"id"`
	Category string        `db:"category"`
	Message  string        `db:"message"`
	Role     string        `db:"role"`
	Created  time.Time     `db:"created"`
	Modified time.Time     `db:"modified"`
}

type ChatMessages []*ChatMessage

func (c *ChatMessages) SortByCreated(desc bool) (*ChatMessages, error) {
	if desc {
		sort.SliceStable(*c, func(i, j int) bool {
			return (*c)[i].Created.After((*c)[j].Created)
		})
	} else {
		sort.SliceStable(*c, func(i, j int) bool {
			return (*c)[i].Created.Before((*c)[j].Created)
		})
	}
	return c, nil
}

func (c *ChatMessage) GPTRequestMessage() ChatGPTRequestMessage {
	return ChatGPTRequestMessage{
		Role:    c.Role,
		Content: c.Message,
	}
}

type ChatSummaryId int

type ChatSummary struct {
	Id       ChatSummaryId `db:"id"`
	Category string        `db:"category"`
	Summary  string        `db:"summary"`
	Created  time.Time     `db:"created"`
	Modified time.Time     `db:"modified"`
}

type ChatMessageOption struct {
	Summaries     []*ChatSummary
	LatestHistory ChatMessages
}
