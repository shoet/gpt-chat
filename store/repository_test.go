package store

import (
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/shoet/gpt-chat/clocker"
	"github.com/shoet/gpt-chat/models"
	"github.com/shoet/gpt-chat/testutil"
)

func Test_Repository_AddChatMessage(t *testing.T) {
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
	tx, err := db.Beginx()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	clocker := &clocker.FixedClocker{}
	r := &Repository{
		Clocker: clocker,
	}

	want := &models.ChatMessage{
		Category: "golang",
		Role:     "user",
		Message:  "test message",
	}
	id, err := r.AddChatMessage(tx, want)
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
