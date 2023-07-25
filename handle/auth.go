package handle

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	JWTKeyEnvVar = "JWT_KEY"
	LoginPath    = "./static/login.html"
	RegisterPath = "./static/register.html"
)

// Define errors
var (
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	ErrTokenGenerationFailed     = errors.New("error generating token")
	ErrHashingPassword           = errors.New("error hashing password")
	ErrStoringData               = errors.New("error storing data in Redis")
)

// Define the interfaces for dependencies
type PasswordHasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

type TokenGenerator interface {
	NewWithClaims(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token
	SignedString(key []byte) (string, error)
	GenerateJWT(username string) (string, error) // Added this
}

type UserDataStore interface {
	GetHashedPassword(username string) (string, error)
	StoreUserData(username, password string) error
}

// Define the handlers
type Handlers struct {
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
	userDataStore  UserDataStore
}

func NewHandlers(passwordHasher PasswordHasher, tokenGenerator TokenGenerator, userDataStore UserDataStore) *Handlers {
	return &Handlers{
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
		userDataStore:  userDataStore,
	}
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, LoginPath)
		return
	}

	username, password, err := parseForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.authenticate(username, password) {
		http.Error(w, ErrInvalidUsernameOrPassword.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.tokenGenerator.GenerateJWT(username) // Modified this
	if err != nil {
		http.Error(w, ErrTokenGenerationFailed.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, RegisterPath)
		return
	}

	username, password, err := parseForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.register(username, password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a 201 status code on successful user creation
	// Set the status code in the 'Location' header.
	w.Header().Set("Location", fmt.Sprintf("/login?status=%d", http.StatusCreated))

	// Redirect the client to the login page.
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func parseForm(r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	return r.FormValue("username"), r.FormValue("password"), nil
}

func (h *Handlers) authenticate(username, password string) bool {
	hashedPassword, err := h.userDataStore.GetHashedPassword(username)
	if err != nil {
		return false
	}

	err = h.passwordHasher.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (h *Handlers) register(username, password string) error {
	hashedPassword, err := h.passwordHasher.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrHashingPassword
	}

	err = h.userDataStore.StoreUserData(username, string(hashedPassword))
	if err != nil {
		return ErrStoringData
	}

	return nil
}
