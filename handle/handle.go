package handle

import (
	"chatbot/redis"
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Chatroom represents a chatroom with a unique ID and a list of users in the chatroom
var chatroomCounter int64

func Handle() {
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)
	http.HandleFunc("/createChatroom", CreateChatroomHandler)
}

// LoginHandler handles login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Get Redis client
		redisClient := redis.GetRedisClient()

		// Retrieve hashed password from Redis for the submitted username
		hashedPassword, err := redisClient.Get(context.Background(), username).Result()
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Check if submitted password matches stored hashed password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Create the JWT token
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &jwt.StandardClaims{
			Subject:   username,
			ExpiresAt: expirationTime.Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		var jwtKey = []byte(os.Getenv("JWT_KEY"))
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		// Write the token to the response
		w.Write([]byte(tokenString))
		return
	}

	http.ServeFile(w, r, "../static/login.html")
}

// RegisterHandler handles registration requests
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Hash and salt password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Get Redis client
		redisClient := redis.GetRedisClient()

		// Store form data in Redis
		err = redisClient.Set(context.Background(), username, string(hashedPassword), 0).Err()
		if err != nil {
			http.Error(w, "Error storing data in Redis", http.StatusInternalServerError)
			return
		}

		// Redirect the client to the login.html page
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "../static/register.html")
}

// CreateChatroomHandler handles chatroom creation requests
func CreateChatroomHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		chatroomName := r.FormValue("chatroomName")

		// Increment chatroomCounter
		chatroomCounter++

		// Store chatroomName with chatroomCounter as ID in Redis
		redisClient := redis.GetRedisClient()
		err := redisClient.Set(context.Background(), strconv.FormatInt(chatroomCounter, 10), chatroomName, 0).Err()
		if err != nil {
			http.Error(w, "Error creating chatroom", http.StatusInternalServerError)
			return
		}

		// Write the chatroom ID to the response
		w.Write([]byte(strconv.FormatInt(chatroomCounter, 10)))
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
