package redis

import (
	"chatbot/model"
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	// Setting up test environment
	os.Setenv("REDIS_URL", "localhost:6379")
	os.Setenv("REDIS_USERNAME", "")
	os.Setenv("REDIS_PASSWORD", "")

	err := InitializeRedisClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize redis client: %v", err)
	}

	// Creating a test chatroom before running the tests
	err = StoreChatroomData("testChatroom")
	if err != nil {
		log.Fatalf("Failed to create test chatroom: %v", err)
	}
}
func TestInitializeRedisClient(t *testing.T) {
	err := InitializeRedisClient(ctx)
	assert.Nil(t, err)
}

func TestCheckChatroomExist(t *testing.T) {
	// Here, we are assuming the chatroom name "testChatroom" already exists.
	exists, err := CheckChatroomExist("testChatroom")
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestStoreMessageInChatroom(t *testing.T) {
	err := StoreMessageInChatroom("testChatroom", "testUser", "test message")
	assert.Nil(t, err)
}

func TestStoreStockDataInChatroom(t *testing.T) {
	stockData := model.StockData{
		StockCode: "TST",
		Price:     150.50,
	}
	err := StoreStockDataInChatroom("testChatroom", stockData)
	assert.Nil(t, err)
}

func TestRetrieveChatroomMessages(t *testing.T) {
	_, err := RetrieveChatroomMessages("testChatroom")
	assert.Nil(t, err)
}

func TestStoreUserData(t *testing.T) {
	err := StoreUserData("testUser", "hashedPassword")
	assert.Nil(t, err)
}

func TestStoreChatroomData(t *testing.T) {
	err := StoreChatroomData("testChatroom")
	assert.Nil(t, err)
}

func TestStoreChatroomMessage(t *testing.T) {
	err := StoreChatroomMessage("testChatroom", "testUser", "test message")
	assert.Nil(t, err)
}

func TestGetHashedPassword(t *testing.T) {
	// We assume that "testUser" has a hashed password "hashedPassword"
	password, err := GetHashedPassword("testUser")
	assert.Nil(t, err)
	assert.Equal(t, "hashedPassword", password)
}

func TestGetAllChatrooms(t *testing.T) {
	_, err := GetAllChatrooms()
	assert.Nil(t, err)
}
