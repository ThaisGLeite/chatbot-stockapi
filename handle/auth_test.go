package handle_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"chatbot/handle"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockPasswordHasher struct{}
type mockTokenGenerator struct{}
type mockUserDataStore struct{}

func (m mockPasswordHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return []byte("hashedPassword"), nil
}

func (m mockPasswordHasher) CompareHashAndPassword(hashedPassword, password []byte) error {
	if string(hashedPassword) == "hashedPassword" && string(password) == "password" {
		return nil
	}
	return bcrypt.ErrMismatchedHashAndPassword
}

func (m mockTokenGenerator) NewWithClaims(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func (m mockTokenGenerator) SignedString(key []byte) (string, error) {
	return "testToken", nil
}

func (m mockTokenGenerator) GenerateJWT(username string) (string, error) { // Added this
	return "testToken", nil
}

func (m mockUserDataStore) GetHashedPassword(username string) (string, error) {
	if username == "user" {
		return "hashedPassword", nil
	}
	return "", errors.New("user not found")
}

func (m mockUserDataStore) StoreUserData(username, password string) error {
	if username == "newUser" {
		return nil
	}
	return errors.New("user already exists")
}

func TestHandlers(t *testing.T) {
	handlers := handle.NewHandlers(mockPasswordHasher{}, mockTokenGenerator{}, mockUserDataStore{})

	// Server for LoginHandler
	loginServer := httptest.NewServer(http.HandlerFunc(handlers.LoginHandler))
	defer loginServer.Close()

	// Test successful login
	resp, err := http.Post(loginServer.URL, "application/x-www-form-urlencoded", strings.NewReader("username=user&password=password"))
	assert.NoError(t, err)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "testToken", string(body))

	// Server for RegisterHandler
	registerServer := httptest.NewServer(http.HandlerFunc(handlers.RegisterHandler))
	defer registerServer.Close()

	// Test registration
	resp, err = http.Post(registerServer.URL, "application/x-www-form-urlencoded", strings.NewReader("username=newUser&password=newPassword"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
