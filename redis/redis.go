package redis

import (
	"context"
	"fmt"
	"os"

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
	// Prepare message with sender
	messageWithSender := fmt.Sprintf(`{ "username": "%s", "message": "%s" }`, username, message)

	// Append the message to the list of messages in this chatroom
	err := redisClient.RPush(context.Background(), chatroomID, messageWithSender).Err()
	return err
}

// Retrieve all messages from a chatroom
func RetrieveChatroomMessages(chatroomID string) ([]string, error) {
	// Get all messages from this chatroom
	messagesWithSenders, err := redisClient.LRange(context.Background(), chatroomID, 0, -1).Result()
	return messagesWithSenders, err
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
