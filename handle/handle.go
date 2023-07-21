package handle

import (
	"chatbot/redis"
	"context"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Handle() {
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)
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

		// User is authenticated, return a placeholder token
		// Note: In a real application, you would generate and return a JWT or similar token here
		w.Write([]byte("token"))
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
