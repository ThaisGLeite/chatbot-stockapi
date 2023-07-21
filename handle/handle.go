package handle

import (
	"chatbot/redis"
	"net/http"
	"os"
	"strconv"
	"strings"
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

		// Retrieve hashed password from Redis for the submitted username
		hashedPassword, err := redis.GetHashedPassword(username)
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

		// Store form data in Redis
		err = redis.StoreUserData(username, string(hashedPassword))
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
		err := redis.StoreChatroomData(strconv.FormatInt(chatroomCounter, 10), chatroomName)
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

// SendMessageHandler handles sending messages
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		chatroomID := r.FormValue("chatroomID")
		username := r.FormValue("username")
		message := r.FormValue("message")

		// Store message in chatroom
		err := redis.StoreChatroomMessage(chatroomID, username, message)
		if err != nil {
			http.Error(w, "Error sending message", http.StatusInternalServerError)
			return
		}

		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

// RetrieveMessagesHandler handles retrieving all messages from a chatroom
func RetrieveMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		r.ParseForm()

		chatroomID := r.FormValue("chatroomID")

		// Retrieve all messages from the chatroom
		messages, err := redis.RetrieveChatroomMessages(chatroomID)
		if err != nil {
			http.Error(w, "Error retrieving messages", http.StatusInternalServerError)
			return
		}

		// Convert the slice of messages to a single string with each message on a new line
		messagesString := strings.Join(messages, "\n")

		// Write the messages to the response
		w.Write([]byte(messagesString))
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
