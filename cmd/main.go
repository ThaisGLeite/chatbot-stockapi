package main

import (
	"chatbot/handle"
	"chatbot/natsclient"
	"chatbot/redis"
	"chatbot/ws"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func (h *BcryptHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (h *BcryptHasher) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

type JWTGenerator struct{}

func (t *JWTGenerator) NewWithClaims(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token {
	return jwt.NewWithClaims(method, claims)
}

func (t *JWTGenerator) SignedString(key []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	return token.SignedString(key)
}

func (t *JWTGenerator) GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func main() {
	// Connect to the NATS server
	natsConn := &natsclient.NatsConn{}
	if err := natsConn.Connect(); err != nil {
		fmt.Println("Failed to connect to NATS server: ", err)
	}

	// Create redis client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := redis.InitializeRedisClient(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize Handlers with the real dependencies
	hasher := &BcryptHasher{}
	tokenGenerator := &JWTGenerator{}
	userDataStore := redis.RedisUserDataStore{}

	handlers := handle.NewHandlers(hasher, tokenGenerator, &userDataStore)

	// Listen to stock data and update redis cache with new data
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	// Run in a goroutine
	go handle.ListenStockData(ctx, ws.Broadcaster)

	http.Handle("/", handle.StaticFilesHandler())
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/createChatroom", handle.CreateChatroomHandler)
	http.HandleFunc("/sendMessage", handle.SendMessageHandler)
	http.HandleFunc("/retrieveMessages", handle.RetrieveMessagesHandler)
	http.HandleFunc("/getAllChatrooms", handle.GetAllChatroomsHandler)
	http.HandleFunc("/checkChatroomExist", handle.CheckChatroomExistHandler)
	http.HandleFunc("/stockUpdates", ws.Connect)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		fmt.Println("Server stopped")
		defer natsConn.Close()
		defer redis.Close()
		log.Fatalf("Failed to start server: %s", err.Error())
		<-ctx.Done()
	}
}
