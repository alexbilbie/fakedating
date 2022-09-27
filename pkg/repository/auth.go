package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/segmentio/ksuid"
)

func NewAuth(db *sql.DB) Auth {
	return Auth{db: db}
}

type Auth struct {
	db *sql.DB
}

func (repo Auth) CreateTokenForUser(userID ksuid.KSUID) (string, error) {
	token := ksuid.New()
	_, err := repo.db.Exec("INSERT INTO auth_tokens (id, user_id) VALUES (?, ?)", token.String(), userID.String())
	if err != nil {
		return "", fmt.Errorf("failed to create token for user: %w", err)
	}
	return token.String(), nil
}

func (repo Auth) GetUserIDByToken(token string) (ksuid.KSUID, error) {
	var _userID string
	if err := repo.db.QueryRow("SELECT user_id FROM auth_tokens WHERE id = ?", token).Scan(&_userID); err != nil {
		if err == sql.ErrNoRows {
			return ksuid.KSUID{}, errors.New("unknown auth token")
		}
		return ksuid.KSUID{}, fmt.Errorf("failed to lookup auth token: %w", err)
	}

	userID, parseErr := ksuid.Parse(_userID)
	if parseErr != nil {
		return ksuid.KSUID{}, fmt.Errorf("failed to parse user ID from auth token: %w", parseErr)
	}

	return userID, nil
}
