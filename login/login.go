package login

import (
	"chatbot/redis"
	"context"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Get Redis client
		redisClient := redis.GetRedisClient()

		// Retrieve password from Redis for the submitted username
		storedPassword, err := redisClient.Get(context.Background(), username).Result()
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Check if submitted password matches stored password
		if password != storedPassword {
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
