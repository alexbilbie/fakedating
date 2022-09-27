package handler

import (
	"net/http"
)

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Get email + password from body

	// Validate email address

	// Validate password

	// Fetch user from database by email

	// Compare password hashes

	// Generate + persist auth token

	// Return token + user
}
