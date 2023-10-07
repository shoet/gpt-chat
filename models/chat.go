package models

import "time"

type ChatMessageId int

type ChatMessage struct {
	Id       ChatMessageId `db:"id"`
	Category string        `db:"category"`
	Message  string        `db:"message"`
	Role     string        `db:"role"`
	Summary  string        `db:"summary"`
	Created  time.Time     `db:"created"`
	Modified time.Time     `db:"modified"`
}

type ChatMessageOption struct {
	Summaries []string
}
