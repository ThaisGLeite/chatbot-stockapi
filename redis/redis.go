package redis

import (
	"chatbot/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

var (
	ErrInvalidChatroomName = errors.New("invalid chatroom name")
	ErrInvalidUsername     = errors.New("invalid username")
	ErrInvalidMessage      = errors.New("invalid message")

	chatroomNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	usernamePattern     = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	messagePattern      = regexp.MustCompile(`^[^<>]+$`) // Rejecting any HTML tags
)

type RedisUserDataStore struct{}

func (r *RedisUserDataStore) GetHashedPassword(username string) (string, error) {
	return GetHashedPassword(username)
}

func (r *RedisUserDataStore) StoreUserData(username, password string) error {
	return StoreUserData(username, password)
}

// InitializeRedisClient creates a new Redis client and tests the connection
func InitializeRedisClient(ctx context.Context) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("Redis client: Ping?", pong)

	return nil
}

// Check if chatroomName exists in Redis Set
func CheckChatroomExist(chatroomName string) (bool, error) {
	exists, err := redisClient.SIsMember(context.Background(), "chatrooms", chatroomName).Result()
	return exists, err
}

// StoreMessageInChatroom stores a message in the specified chatroom
func StoreMessageInChatroom(chatroomName string, username string, message string) error {
	// Validate inputs
	if !chatroomNamePattern.MatchString(chatroomName) {
		return ErrInvalidChatroomName
	}
	if !usernamePattern.MatchString(username) {
		return ErrInvalidUsername
	}
	if !messagePattern.MatchString(message) {
		return ErrInvalidMessage
	}

	// Prepare message with sender and timestamp
	timestamp := time.Now().Unix()
	messageWithSender := fmt.Sprintf(`{ "username": "%s", "message": "%s", "timestamp": %d }`, username, message, timestamp)

	// Append the message to the list of messages in this chatroom
	err := redisClient.RPush(context.Background(), chatroomName, messageWithSender).Err()
	if err != nil {
		return err
	}

	// Limit the list to the last 50 messages
	err = redisClient.LTrim(context.Background(), chatroomName, -50, -1).Err()
	return err
}

// Store stock data in chatroom
func StoreStockDataInChatroom(chatroomName string, stockData model.StockData) error {
	// Prepare stock data message
	stockDataMessage := fmt.Sprintf("Stock code: %s, price: %f", stockData.StockCode, stockData.Price)

	// Store stock data as a message in the chatroom
	return StoreMessageInChatroom(chatroomName, "Bot", stockDataMessage)
}

// Retrieve all messages from a chatroom
func RetrieveChatroomMessages(chatroomName string) ([]model.ChatMessage, error) {
	// Get all messages from this chatroom

	messagesWithSenders, err := redisClient.LRange(context.Background(), chatroomName, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Parse the messages into ChatMessage structs and sort them by timestamp
	var messages []model.ChatMessage
	for _, messageWithSender := range messagesWithSenders {
		var msg model.ChatMessage
		err := json.Unmarshal([]byte(messageWithSender), &msg)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp < messages[j].Timestamp
	})
	return messages, nil
}

// Updated: StoreUserData stores user's hashed password as a field in a hash with the key 'users'
func StoreUserData(username string, hashedPassword string) error {
	// Store form data in Redis hash with the key 'users'
	err := redisClient.HSet(context.Background(), "users", username, hashedPassword).Err()
	return err
}

// This function stores chatroom data in Redis
func StoreChatroomData(chatroomName string) error {
	// Store chatroom data in a Redis Set
	err := redisClient.SAdd(context.Background(), "chatrooms", chatroomName).Err()
	return err
}

// StoreChatroomMessage stores a user's message in a chatroom
func StoreChatroomMessage(chatroomName string, username string, message string) error {
	// Prepare the message to store username along with the message
	messageToStore := username + ": " + message

	// Append the message to the list of messages in this chatroom
	err := redisClient.RPush(context.Background(), chatroomName, messageToStore).Err()
	return err
}

// Updated: This function retrieves hashed password from Redis for the submitted username
func GetHashedPassword(username string) (string, error) {
	hashedPassword, err := redisClient.HGet(context.Background(), "users", username).Result()
	return hashedPassword, err
}

// Fetch all chatrooms from the Redis Set
func GetAllChatrooms() ([]string, error) {
	chatrooms, err := redisClient.SMembers(context.Background(), "chatrooms").Result()
	return chatrooms, err
}

// Close connection to Redis
func Close() {
	redisClient.Close()
}
