package models

type ChatHistory struct {
	Id      int    `db:"id"`
	Content string `db:"content"`
}
