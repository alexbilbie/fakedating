package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/mail"

	"fakedating/pkg/payload"
	"fakedating/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Get email + password from body
	body, readErr := io.ReadAll(io.LimitReader(r.Body, 500))
	defer r.Body.Close()
	if readErr != nil {
		util.WriteErrorResponse("Failed to read the request body", readErr, http.StatusInternalServerError, w)
		return
	}

	var loginPayload payload.LoginRequest
	if unmarshallErr := json.Unmarshal(body, &loginPayload); unmarshallErr != nil {
		util.WriteErrorResponse("Failed to decode the request body", readErr, http.StatusBadRequest, w)
		return
	}

	// Validate payload
	email, parseErr := mail.ParseAddress(loginPayload.Email)
	if parseErr != nil {
		util.WriteErrorResponse("Invalid email address", parseErr, http.StatusBadRequest, w)
		return
	}

	if len(loginPayload.Password) == 0 {
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
		[]byte(loginPayload.Password),
	); compareErr != nil {
		util.WriteErrorResponse("Invalid email or password", compareErr, http.StatusBadRequest, w)
		return
	}

	// Generate auth token
	token, createTokenErr := h.authRepository.CreateTokenForUser(user.ID)
	if createTokenErr != nil {
		util.WriteErrorResponse("Failed to create auth token", createTokenErr, http.StatusInternalServerError, w)
		return
	}

	// Return token + user
	log.Printf("User %q logged in successfully", user.ID)
	util.WriteJSONResponse(
		payload.LoginResponse{
			User:  user,
			Token: token,
		},
		http.StatusOK,
		w,
	)
}
