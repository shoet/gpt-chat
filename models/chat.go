package models

import "time"

type ChatMessage struct {
	Id       int       `db:"id"`
	Category string    `db:"category"`
	Content  string    `db:"content"`
	Role     string    `db:"role"`
	Created  time.Time `db:"created"`
	Modified time.Time `db:"modified"`
}
