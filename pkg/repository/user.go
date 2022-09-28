package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"fakedating/pkg/model"
	"github.com/segmentio/ksuid"
)

func NewUser(db *sql.DB) User {
	return User{db: db}
}

type User struct {
	db *sql.DB
}

func (repo User) Create(user model.User) (model.User, error) {
	const query = `INSERT INTO 
    	users (id, email, full_name, password_hash, gender, age, location) 
		VALUES (?, ?, ?, ?, ?, ?, POINT(?, ?))`

	user.ID = ksuid.New()
	_, err := repo.db.Exec(
		query,
		user.ID.String(),
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Gender,
		user.Age,
		user.Location.Longitude, // MariaDB stores location as longitude first in POINT type
		user.Location.Latitude,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to persist new user: %w", err)
	}

	return user, nil
}

func (repo User) GetByEmail(email string) (model.User, error) {
	var user model.User
	var gender string

	err := repo.db.
		QueryRow("SELECT id, email, full_name, password_hash, gender, age, location FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &gender, &user.Age, &user.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, errors.New("unknown user")
		}
		return model.User{}, fmt.Errorf("failed to lookup user: %w", err)
	}

	if gender == "female" {
		user.Gender = model.GenderFemale
	} else {
		user.Gender = model.GenderMale
	}

	return user, nil
}

func (repo User) ListMatches(userID ksuid.KSUID, options ...model.SearchParameterOpt) ([]model.User, error) {
	// TODO implement me
	panic("implement me")
}

func (repo User) Swipe(swiperID ksuid.KSUID, recipient ksuid.KSUID, match bool) (model.ProfileMatch, error) {
	// TODO implement me
	panic("implement me")
}
