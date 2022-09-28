package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"fakedating/pkg/model"
	"github.com/segmentio/ksuid"
)

// NewUser returns an initialised user repository
func NewUser(db *sql.DB) User {
	return User{db: db}
}

type User struct {
	db *sql.DB
}

// Create persists a user and returns the new record with its ID value set
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

// GetByID returns a user record found by the email value
func (repo User) GetByID(id ksuid.KSUID) (model.User, error) {
	var user model.User
	var gender string

	err := repo.db.
		QueryRow("SELECT id, email, full_name, password_hash, gender, age, location FROM users WHERE id = ?", id).
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

// GetByEmail returns a user record found by the email value
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

// ListMatches returns other users limited by search parameters and excluding existing matches
func (repo User) ListMatches(userID ksuid.KSUID, options ...model.SearchParameterOpt) ([]model.User, error) {
	queryParts := []string{
		"SELECT id, email, full_name, gender, age, location FROM users WHERE id <> ?",
		"AND id NOT IN (SELECT recipient_id FROM swipes WHERE swiper_id = ?)",
	}
	queryParams := []any{userID, userID}

	var searchParams model.SearchParameter
	for _, opt := range options {
		opt(&searchParams)
	}

	if searchParams.AgeLower != 0 {
		queryParts = append(queryParts, "AND age >= ?")
		queryParams = append(queryParams, searchParams.AgeLower)
	}

	if searchParams.AgeUpper != 0 {
		queryParts = append(queryParts, "AND age <= ?")
		queryParams = append(queryParams, searchParams.AgeUpper)
	}

	if searchParams.Latitude != 0 && searchParams.Longitude != 0 {
		queryParts = append(queryParts, "AND ST_Distance_Sphere(location, POINT(?,?)) <= ? * 1000")
		queryParams = append(queryParams, searchParams.Longitude)
		queryParams = append(queryParams, searchParams.Latitude)
		queryParams = append(queryParams, searchParams.Radius)
	}

	queryParts = append(queryParts, "LIMIT 25")

	if searchParams.Offset != 0 {
		queryParts = append(queryParts, "OFFSET ?")
		queryParams = append(queryParams, searchParams.Offset)
	}

	fmt.Println(strings.Join(queryParts, " "))
	fmt.Println(queryParams)

	rows, queryErr := repo.db.Query(strings.Join(queryParts, " "), queryParams...)
	if queryErr != nil {
		return nil, fmt.Errorf("failed to query for users: %w", queryErr)
	}

	var result []model.User
	for {
		if !rows.Next() {
			break
		}

		var user model.User
		var gender string
		if scanErr := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&gender,
			&user.Age,
			&user.Location,
		); scanErr != nil {
			return nil, fmt.Errorf("failed to read query results for users: %w", scanErr)
		}

		if gender == "female" {
			user.Gender = model.GenderFemale
		} else {
			user.Gender = model.GenderMale
		}
		result = append(result, user)
	}

	return result, nil
}

// SaveSwipe persists a swipe and returns a mutual or unrequited match
func (repo User) SaveSwipe(swiperID ksuid.KSUID, recipient ksuid.KSUID, matched bool) (model.ProfileMatch, error) {
	// Check if the recipient has matched with the swiper
	var mutualMatch bool
	mutualMatchErr := repo.db.
		QueryRow("SELECT matched FROM swipes WHERE swiper_id = ? AND recipient_id = ?", recipient, swiperID).
		Scan(&mutualMatch)
	if mutualMatchErr != nil && mutualMatchErr != sql.ErrNoRows {
		return model.ProfileMatchUnrequited, mutualMatchErr
	}

	// Save the new match (ignore a duplicate)
	_, saveErr := repo.db.Exec(
		"INSERT INTO swipes (swiper_id, recipient_id, matched) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE id = id",
		swiperID,
		recipient,
		matched,
	)
	if saveErr != nil {
		return model.ProfileMatchUnrequited, saveErr
	}

	if matched && mutualMatch {
		return model.ProfileMatchMutual, nil
	}
	return model.ProfileMatchUnrequited, nil
}
