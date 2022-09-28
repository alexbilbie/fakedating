package util

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSONResponse(t *testing.T) {
	payload := struct {
		Foo string
	}{
		Foo: "Bar",
	}

	w := httptest.NewRecorder()
	WriteJSONResponse(payload, http.StatusTeapot, w)

	assert.Equal(t, w.Result().StatusCode, http.StatusTeapot)
	body, err := io.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"Foo":"Bar"}`, string(body))
	assert.Equal(t, w.Result().Header.Get("content-type"), "application/json")
}
