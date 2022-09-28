package middleware

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"fakedating/pkg/util"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	validToken = ksuid.New()
	validUser  = ksuid.New()
)

type mockAuthRepository struct{}

func (mockAuthRepository) GetUserIDByToken(token string) (ksuid.KSUID, error) {
	if token == validToken.String() {
		return validUser, nil
	}
	return ksuid.KSUID{}, errors.New("invalid token")
}

func TestAuthenticateRequest_MissingAuthToken(t *testing.T) {
	w := httptest.NewRecorder()
	middleware := AuthenticateRequest(mockAuthRepository{}, mockHandler{})

	req, _ := http.NewRequest("GET", "/", nil)
	middleware.ServeHTTP(w, req)

	assert.Equal(t, w.Result().StatusCode, http.StatusUnauthorized)
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"Msg":"Missing authorization token"}`, string(body))
	assert.Equal(t, w.Result().Header.Get("content-type"), "application/json")
}

func TestAuthenticateRequest_InvalidAuthToken(t *testing.T) {
	w := httptest.NewRecorder()
	middleware := AuthenticateRequest(mockAuthRepository{}, mockHandler{})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", ksuid.New().String())
	middleware.ServeHTTP(w, req)

	assert.Equal(t, w.Result().StatusCode, http.StatusUnauthorized)
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"Msg":"Invalid authorization token"}`, string(body))
	assert.Equal(t, w.Result().Header.Get("content-type"), "application/json")
}

func TestAuthenticateRequest_ValidAuthToken(t *testing.T) {
	w := httptest.NewRecorder()
	middleware := AuthenticateRequest(mockAuthRepository{}, mockHandler{})

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", validToken.String())
	middleware.ServeHTTP(w, req)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(`{"UserID":%q}`, validUser), string(body))
	assert.Equal(t, w.Result().Header.Get("content-type"), "application/json")
}

type mockHandler struct{}

func (h mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	util.WriteJSONResponse(
		struct {
			UserID ksuid.KSUID
		}{
			UserID: userID,
		},
		http.StatusOK,
		w,
	)
}
