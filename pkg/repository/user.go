package repository

import (
	"database/sql"
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
		user.LocationLatitude,
		user.LocationLongitude,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to persist new user: %w", err)
	}

	return user, nil
}

func (repo User) GetByEmail(email string) (model.User, error) {
	// TODO implement me
	panic("implement me")
}

func (repo User) ListMatches(userID ksuid.KSUID, options ...model.SearchParameterOpt) ([]model.User, error) {
	// TODO implement me
	panic("implement me")
}

func (repo User) Swipe(swiperID ksuid.KSUID, recipient ksuid.KSUID, match bool) (model.ProfileMatch, error) {
	// TODO implement me
	panic("implement me")
}
