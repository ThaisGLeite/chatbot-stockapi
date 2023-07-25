package natsclient_test

import (
	"chatbot/natsclient"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	// Set NATS_URL for the testing environment
	os.Setenv("NATS_URL", "nats://localhost:4222")

	// Create a new NatsConn instance
	natsConn := &natsclient.NatsConn{}

	// Call Connect function
	err := natsConn.Connect()
	assert.Nil(t, err)

	// The Client should be non-nil if connected successfully
	assert.NotNil(t, natsclient.Client)
	assert.True(t, natsConn.Conn.IsConnected()) // conn is unexported, you might want to add an IsConnected method in your NatsConn struct

	// Test Close function
	natsConn.Close()
	assert.False(t, natsConn.Conn.IsConnected()) // Same comment as above
}
