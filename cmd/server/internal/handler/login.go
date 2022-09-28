package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/mail"

	"fakedating/pkg/model"
	"fakedating/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

type loginPayload struct {
	Email    string
	Password string
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Get email + password from body
	body, readErr := io.ReadAll(io.LimitReader(r.Body, 500))
	if readErr != nil {
		util.WriteErrorResponse("Failed to read the request body", readErr, http.StatusInternalServerError, w)
		return
	}

	var payload loginPayload
	if unmarshallErr := json.Unmarshal(body, &payload); unmarshallErr != nil {
		util.WriteErrorResponse("Failed to decode the request body", readErr, http.StatusBadRequest, w)
		return
	}

	// Validate payload
	email, parseErr := mail.ParseAddress(payload.Email)
	if parseErr != nil {
		util.WriteErrorResponse("Invalid email address", parseErr, http.StatusBadRequest, w)
		return
	}

	if len(payload.Password) == 0 {
		util.WriteErrorResponse("Empty password", nil, http.StatusBadRequest, w)
		return
	}

	// Fetch user from database by email
	user, getUserErr := h.userRepository.GetByEmail(email.Address)
	if getUserErr != nil {
		util.WriteErrorResponse("Invalid email or password", getUserErr, http.StatusBadRequest, w)
		return
	}

	// Compare password hashes
	if compareErr := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(payload.Password),
	); compareErr != nil {
		util.WriteErrorResponse("Invalid email or password", compareErr, http.StatusBadRequest, w)
		return
	}

	// Generate + persist auth token
	token, createTokenErr := h.authRepository.CreateTokenForUser(user.ID)
	if createTokenErr != nil {
		util.WriteErrorResponse("Failed to create auth token", createTokenErr, http.StatusInternalServerError, w)
		return
	}

	// Return token + user
	encodedResp, marshalErr := json.Marshal(
		struct {
			User  model.User
			Token string
		}{
			User:  user,
			Token: token,
		},
	)
	if marshalErr != nil {
		util.WriteErrorResponse("Failed to marshal response to JSON", marshalErr, http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(encodedResp)
}
