package service

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shoet/gpt-chat/models"
	"github.com/shoet/gpt-chat/store"
)

type StorageRDB struct {
	db   *sqlx.DB
	repo *store.Repository
}

func NewStorageRDB(db *sqlx.DB, repo *store.Repository) (*StorageRDB, error) {
	return &StorageRDB{
		db:   db,
		repo: repo,
	}, nil
}

func (s *StorageRDB) AddChatMessage(message *models.ChatMessage) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	if _, err := s.repo.AddChatMessage(tx, message); err != nil {
		return fmt.Errorf("failed to save chat history: %w", err)
	}
	return tx.Commit()
}

func (s *StorageRDB) ListChatSummary(latest int) ([]string, error) {
	summaries, err := s.repo.ListChatSummary(s.db, latest)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat summaries: %w", err)
	}
	return summaries, nil
}
