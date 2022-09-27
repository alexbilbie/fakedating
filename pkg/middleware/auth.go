package middleware

import (
	"context"
	"net/http"
	"strings"

	"fakedating/pkg/util"
	"github.com/segmentio/ksuid"
)

type key int

const (
	contextUserID key = iota
)

// GetUserIDFromContext returns the ID of the authenticated user
func GetUserIDFromContext(ctx context.Context) ksuid.KSUID {
	u, ok := ctx.Value(contextUserID).(ksuid.KSUID)
	if ok {
		return u
	}
	return ksuid.KSUID{}
}

// A list of API endpoints which don't require authentication
var openRoutes = []string{
	"/login",
	"/user/create",
}

type AuthRepository interface {
	GetUserIDByToken(token string) (ksuid.KSUID, error)
}

// AuthenticateRequest will validate an access token if one is present in the headers for non-whitelisted endpoints
func AuthenticateRequest(repo AuthRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Determine if the current endpoint is open and doesn't require auth
			if util.InSlice(r.URL.Path, openRoutes) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract the auth token from the Authorization header
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				util.WriteErrorResponse("Missing authorization token", nil, http.StatusUnauthorized, w)
				return
			}

			// Lookup user in database
			userID, lookupErr := repo.GetUserIDByToken(authHeader)
			if lookupErr != nil {
				// TODO
			}

			// Update the request context with the user ID for use later on
			newCtx := context.WithValue(r.Context(), contextUserID, userID)

			// Serve the actual API request with the authorized context
			next.ServeHTTP(w, r.WithContext(newCtx))
		},
	)
}
