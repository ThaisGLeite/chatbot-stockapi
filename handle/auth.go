package handle

import (
	"chatbot/redis"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	JWTKeyEnvVar = "JWT_KEY"
	LoginPath    = "../static/login.html"
	RegisterPath = "../static/register.html"
)

// Errors
var (
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	ErrTokenGenerationFailed     = errors.New("error generating token")
	ErrHashingPassword           = errors.New("error hashing password")
	ErrStoringData               = errors.New("error storing data in Redis")
)

// LoginHandler handles login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, LoginPath)
		return
	}

	username, password, err := parseForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !authenticate(username, password) {
		http.Error(w, ErrInvalidUsernameOrPassword.Error(), http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(username)
	if err != nil {
		http.Error(w, ErrTokenGenerationFailed.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}

// RegisterHandler handles registration requests
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, RegisterPath)
		return
	}

	username, password, err := parseForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = register(username, password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}

func parseForm(r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	return r.FormValue("username"), r.FormValue("password"), nil
}

func authenticate(username, password string) bool {
	hashedPassword, err := redis.GetHashedPassword(username)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv(JWTKeyEnvVar))
	return token.SignedString(jwtKey)
}

func register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(ErrHashingPassword, err)
		return ErrHashingPassword
	}

	err = redis.StoreUserData(username, string(hashedPassword))
	if err != nil {
		fmt.Println(ErrStoringData, err)
		return ErrStoringData
	}

	return nil
}
