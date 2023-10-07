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
		Summary:  "test summary",
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
		id, category, message, role, summary, created, modified
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

func Test_Repository_ListChatSummary(t *testing.T) {
	db, closer := prepareDB(t)
	t.Cleanup(func() { closer() })

	tx, err := db.Beginx()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	clocker := &clocker.FixedClocker{}
	sut := &Repository{
		Clocker: clocker,
	}

	sql := `
	INSERT INTO chat_message 
		(category, message, role, summary, created, modified)
	VALUES 
		(?, ?, ?, ?, ?, ?);
	`
	m := &models.ChatMessage{
		Category: "golang",
		Role:     "user",
		Message:  "test message",
		Created:  clocker.Now(),
		Modified: clocker.Now(),
	}

	for i := 0; i < 11; i++ {
		m.Summary = fmt.Sprintf("summary %d", i+1)
		_, err := tx.Exec(sql, m.Category, m.Message, m.Role, m.Summary, m.Created, m.Modified)
		if err != nil {
			t.Fatalf("failed to insert chat message: %v", err)
		}
	}

	summaries, err := sut.ListChatSummary(tx, 10)
	if err != nil {
		t.Errorf("failed to list chat summary: %v", err)
	}

	fmt.Println(summaries)

}
