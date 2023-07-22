package redis

import (
	"chatbot/model"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitializeRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(pong, err)
	}
}

// Store message in chatroom
func StoreMessageInChatroom(chatroomID string, username string, message string) error {
	// Prepare message with sender and timestamp
	timestamp := time.Now().Unix()
	messageWithSender := fmt.Sprintf(`{ "username": "%s", "message": "%s", "timestamp": %d }`, username, message, timestamp)

	// Append the message to the list of messages in this chatroom
	err := redisClient.RPush(context.Background(), chatroomID, messageWithSender).Err()
	if err != nil {
		return err
	}

	// Limit the list to the last 50 messages
	err = redisClient.LTrim(context.Background(), chatroomID, -50, -1).Err()
	return err
}

// Retrieve all messages from a chatroom
func RetrieveChatroomMessages(chatroomID string) ([]model.ChatMessage, error) {
	// Get all messages from this chatroom
	messagesWithSenders, err := redisClient.LRange(context.Background(), chatroomID, 0, -1).Result()
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

func StoreUserData(username string, hashedPassword string) error {
	// Store form data in Redis
	err := redisClient.Set(context.Background(), username, hashedPassword, 0).Err()
	return err
}

// This function stores chatroom data in Redis
func StoreChatroomData(chatroomID string, chatroomName string) error {
	err := redisClient.Set(context.Background(), chatroomID, chatroomName, 0).Err()
	return err
}

// This function retrieves hashed password from Redis for the submitted username
func GetHashedPassword(username string) (string, error) {
	hashedPassword, err := redisClient.Get(context.Background(), username).Result()
	return hashedPassword, err
}

// StoreChatroomMessage stores a user's message in a chatroom
func StoreChatroomMessage(chatroomID string, username string, message string) error {
	// Prepare the message to store username along with the message
	messageToStore := username + ": " + message

	// Append the message to the list of messages in this chatroom
	err := redisClient.RPush(context.Background(), chatroomID, messageToStore).Err()
	return err
}

// Close connection to Redis
func Close() {
	redisClient.Close()
}
