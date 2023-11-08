package store

import (
	"fmt"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/shoet/gpt-chat/clocker"
	"github.com/shoet/gpt-chat/models"
	"github.com/shoet/gpt-chat/testutil"
)

func prepareDB(t *testing.T) (*sqlx.DB, func() error) {
	t.Helper()
	c := mysql.Config{
		DBName:               "gpt",
		User:                 "gpt",
		Passwd:               "gpt",
		Addr:                 "localhost:33306",
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db, closer, err := NewRDB(c.FormatDSN())
	if err != nil {
		t.Fatalf("failed create new db connection: %v", err)
	}
	return db, closer
}

func Test_Repository_AddChatMessage(t *testing.T) {
	db, closer := prepareDB(t)
	tx, err := db.Beginx()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	clocker := &clocker.FixedClocker{}
	sut := &Repository{
		Clocker: clocker,
	}

	want := &models.ChatMessage{
		Category: "golang",
		Role:     "user",
		Message:  "test message",
	}
	id, err := sut.AddChatMessage(tx, want)
	if err != nil {
		t.Fatalf("failed to add chat message: %v", err)
	}
	want.Id = id
	want.Created = clocker.Now()
	want.Modified = clocker.Now()

	got := []*models.ChatMessage{}
	sq := `
	SELECT
		id, category, message, role, created, modified
	FROM
		chat_message
	WHERE id = ?;
	`
	if err := tx.Select(&got, sq, id); err != nil {
		t.Fatalf("failed to select chat message: %v", err)
	}

	if err := testutil.AssertObject(t, want, got[0]); err != nil {
		t.Fatalf("failed to assert object: %v", err)
	}

	t.Cleanup(func() {
		tx.Rollback()
		_ = closer()
	})
}

func Test_Repository_ListChatHistory(t *testing.T) {
	db, closer := prepareDB(t)
	tx, err := db.Beginx()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	clocker := &clocker.FixedClocker{}
	sut := &Repository{
		Clocker: clocker,
	}

	messages, err := sut.ListChatHistory(tx, "golang", 10)
	if err != nil {
		t.Fatalf("failed to list chat message: %v", err)
	}

	if len(messages) == 0 {
		t.Fatalf("failed len chat messages: %v", err)
	}

	sorted, err := messages.SortByCreated(false)
	if err != nil {
		t.Fatalf("failed to sort chat messages: %v", err)
	}

	for _, m := range *sorted {
		fmt.Println(m.Created, m.Role)
	}

	t.Cleanup(func() {
		tx.Rollback()
		_ = closer()
	})
}
