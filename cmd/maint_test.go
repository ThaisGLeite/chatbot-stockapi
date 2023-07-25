package main

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Mock implementation of BcryptHasher
type MockBcryptHasher struct {
	err error
}

func (h *MockBcryptHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	if h.err != nil {
		return nil, h.err
	}
	return password, nil
}

func (h *MockBcryptHasher) CompareHashAndPassword(hashedPassword, password []byte) error {
	return h.err
}

// Mock implementation of JWTGenerator
type MockJWTGenerator struct {
	err error
}

func (t *MockJWTGenerator) NewWithClaims(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token {
	return jwt.NewWithClaims(method, claims)
}

func (t *MockJWTGenerator) SignedString(key []byte) (string, error) {
	if t.err != nil {
		return "", t.err
	}
	return "mockSignedString", nil
}

func (t *MockJWTGenerator) GenerateJWT(username string) (string, error) {
	if t.err != nil {
		return "", t.err
	}
	return "mockJWT", nil
}

func TestBcryptHasher(t *testing.T) {
	hasher := &BcryptHasher{}
	password := []byte("testpassword")
	hashedPassword, err := hasher.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = hasher.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestJWTGenerator(t *testing.T) {
	generator := &JWTGenerator{}
	token, err := generator.GenerateJWT("testuser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(token) == 0 {
		t.Fatalf("Expected a token, got an empty string")
	}
}

func TestMockBcryptHasherWithError(t *testing.T) {
	hasher := &MockBcryptHasher{err: errors.New("mock error")}
	_, err := hasher.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}

func TestMockJWTGeneratorWithError(t *testing.T) {
	generator := &MockJWTGenerator{err: errors.New("mock error")}
	_, err := generator.GenerateJWT("testuser")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}
