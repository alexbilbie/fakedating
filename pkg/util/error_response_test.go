package util

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	WriteErrorResponse("Something went wrong", errors.New("oops"), http.StatusBadRequest, w)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"Msg":"Something went wrong","Error":"oops"}`, string(body))
	assert.Equal(t, w.Result().Header.Get("content-type"), "application/json")
}
