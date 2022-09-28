package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/segmentio/ksuid"
)

// NewAuth returns an initialised auth repository
func NewAuth(db *sql.DB) Auth {
	return Auth{db: db}
}

type Auth struct {
	db *sql.DB
}

// CreateTokenForUser generates an authentication token and persists it to the database
func (repo Auth) CreateTokenForUser(userID ksuid.KSUID) (string, error) {
	token := ksuid.New()
	_, err := repo.db.Exec("INSERT INTO auth_tokens (id, user_id) VALUES (?, ?)", token.String(), userID.String())
	if err != nil {
		return "", fmt.Errorf("failed to create token for user: %w", err)
	}
	return token.String(), nil
}

// GetUserIDByToken returns the user ID represented by an authentication token
func (repo Auth) GetUserIDByToken(token string) (ksuid.KSUID, error) {
	var userID ksuid.KSUID
	if err := repo.db.QueryRow("SELECT user_id FROM auth_tokens WHERE id = ?", token).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return ksuid.KSUID{}, errors.New("unknown auth token")
		}
		return ksuid.KSUID{}, fmt.Errorf("failed to lookup auth token: %w", err)
	}

	return userID, nil
}
