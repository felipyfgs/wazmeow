package usecases_session_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wazmeow/internal/domain/session"
)

func TestUseCaseIntegration(t *testing.T) {
	t.Run("usecase creation and basic operations", func(t *testing.T) {
		// Test that we can create and manipulate sessions for usecases
		session := session.NewSession("test-usecase-session")
		
		assert.NotNil(t, session, "Session should be created successfully")
		assert.Equal(t, "test-usecase-session", session.Name(), "Session name should match")
		assert.False(t, session.ID().IsEmpty(), "Session should have a valid ID")
	})

	t.Run("session validation for usecases", func(t *testing.T) {
		validSession := session.NewSession("valid-session")
		err := validSession.Validate()
		assert.NoError(t, err, "Valid session should not have validation errors")

		// Test session operations that usecases would use
		assert.True(t, validSession.CanConnect(), "Valid session should be connectable")
		assert.False(t, validSession.IsConnected(), "New session should not be connected")
	})

	t.Run("session state management for usecases", func(t *testing.T) {
		session := session.NewSession("state-management-test")
		
		// Test transitions that usecases would trigger
		session.SetConnecting()
		assert.True(t, session.IsConnecting(), "Session should be in connecting state")
		
		err := session.Connect("test@s.whatsapp.net")
		assert.NoError(t, err, "Session should connect successfully")
		assert.True(t, session.IsConnected(), "Session should be connected")
		assert.True(t, session.IsActive(), "Session should be active")
		
		session.Disconnect()
		assert.False(t, session.IsConnected(), "Session should be disconnected")
		assert.False(t, session.IsActive(), "Session should be inactive")
	})

	t.Run("session qr code management", func(t *testing.T) {
		session := session.NewSession("qr-test")
		
		// Test QR code operations that usecases would manage
		session.SetQRCode("test-qr-code")
		assert.Equal(t, "test-qr-code", session.QRCode(), "QR code should be set correctly")
		
		session.ClearQRCode()
		assert.Empty(t, session.QRCode(), "QR code should be cleared")
	})
}

func TestSessionRepositoryIntegration(t *testing.T) {
	t.Run("repository interface compatibility", func(t *testing.T) {
		session := session.NewSession("repo-test")
		
		// Test that session has all the methods that repository would need
		id := session.ID()
		name := session.Name()
		status := session.Status()
		
		assert.False(t, id.IsEmpty(), "Session ID should be valid for repository")
		assert.NotEmpty(t, name, "Session name should be set for repository")
		assert.NotNil(t, status, "Session status should be available for repository")
		
		// Test session validation (important for repository operations)
		err := session.Validate()
		assert.NoError(t, err, "Session should pass validation for repository storage")
	})
}