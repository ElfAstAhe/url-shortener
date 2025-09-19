package model

import (
	"errors"

	"github.com/google/uuid"
)

type ShortURIUser struct {
	ID         string `db:"id" json:"id"`
	ShortURIID string `db:"short_uri_id" json:"short_uri_id"`
	UserID     string `db:"user_id" json:"user_id"`
}

func NewShortURIUser(shortURIID string, userID string) (*ShortURIUser, error) {
	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &ShortURIUser{
		ID:         newID.String(),
		ShortURIID: shortURIID,
		UserID:     userID,
	}, nil
}

func ValidateShortURIUser(entity *ShortURIUser) error {
	if entity == nil {
		return errors.New("entity is nil")
	}
	if entity.ShortURIID == "" {
		return errors.New("short_uri_id is empty")
	}
	if entity.UserID == "" {
		return errors.New("user_id is empty")
	}

	return nil
}
