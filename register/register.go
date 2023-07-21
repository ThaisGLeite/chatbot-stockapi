package register

import (
	"chatbot/redis"
	"context"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
