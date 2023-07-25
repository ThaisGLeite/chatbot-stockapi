package handle_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"chatbot/handle"

	"github.com/stretchr/testify/assert"
)

func TestStaticFilesHandler(t *testing.T) {
	// Create the static directory for testing
	os.Mkdir("./static", os.ModePerm)

	h := handle.StaticFilesHandler()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)

	// Delete the static directory after testing
	os.RemoveAll("./static")
}
