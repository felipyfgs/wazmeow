package infra_repository_sqlite_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wazmeow/internal/domain/session"
)

func TestSessionIntegration(t *testing.T) {
	t.Run("session creation and basic operations", func(t *testing.T) {
		// Test basic session operations
		sess := session.NewSession("test-session")
		
		assert.NotNil(t, sess, "Session should be created")
		assert.NotEmpty(t, sess.Name(), "Session should have a name")
		assert.False(t, sess.ID().IsEmpty(), "Session should have a valid ID")
		
		// Test session methods
		assert.True(t, sess.CanConnect(), "New session should be able to connect")
		assert.False(t, sess.IsConnected(), "New session should not be connected")
		assert.False(t, sess.IsActive(), "New session should not be active")
	})
}

func TestSessionStatusTransitions(t *testing.T) {
	t.Run("status transitions work correctly", func(t *testing.T) {
		sess := session.NewSession("test-session")
		
		// Initial state
		assert.True(t, sess.CanConnect(), "New session should be able to connect")
		assert.False(t, sess.IsConnected(), "New session should not be connected")
		
		// Set connecting
		sess.SetConnecting()
		assert.True(t, sess.CanConnect(), "Connecting session should be able to connect")
		assert.False(t, sess.IsConnected(), "Connecting session should not be connected")
		assert.True(t, sess.IsConnecting(), "Session should be in connecting state")
		
		// Connect (simulate successful connection)
		err := sess.Connect("test@s.whatsapp.net")
		assert.NoError(t, err)
		assert.False(t, sess.CanConnect(), "Connected session should not be able to connect")
		assert.True(t, sess.IsConnected(), "Connected session should be connected")
		assert.True(t, sess.IsActive(), "Connected session should be active")
		assert.False(t, sess.IsConnecting(), "Connected session should not be connecting")
		
		// Disconnect
		sess.Disconnect()
		assert.True(t, sess.CanConnect(), "Disconnected session should be able to connect")
		assert.False(t, sess.IsConnected(), "Disconnected session should not be connected")
		assert.False(t, sess.IsActive(), "Disconnected session should not be active")
	})
}