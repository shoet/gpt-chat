package models

type ChatMessage struct {
	Id       int    `db:"id"`
	Category string `db:"category"`
	Content  string `db:"content"`
	Role     string `db:"role"`
}
