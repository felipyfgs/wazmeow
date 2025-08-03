package domain_session_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wazmeow/internal/domain/session"
)

func TestNewSession(t *testing.T) {
	t.Run("should create session with valid name", func(t *testing.T) {
		name := "test-session"
		sess := session.NewSession(name)

		assert.NotNil(t, sess)
		assert.NotEmpty(t, sess.ID())
		assert.Equal(t, name, sess.Name())
		assert.Equal(t, session.StatusDisconnected, sess.Status())
		assert.Empty(t, sess.WaJID())
		assert.Empty(t, sess.QRCode())
		assert.False(t, sess.IsActive())
		assert.False(t, sess.CreatedAt().IsZero())
		assert.False(t, sess.UpdatedAt().IsZero())
	})
}

func TestCanConnect(t *testing.T) {
	testCases := []struct {
		name     string
		status   session.Status
		expected bool
	}{
		{"disconnected", session.StatusDisconnected, true},
		{"connecting", session.StatusConnecting, true},
		{"connected", session.StatusConnected, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sess := session.RestoreSession(
				session.NewSessionID(),
				"test-session",
				tc.status,
				"",
				"",
				false,
				time.Now(),
				time.Now(),
			)

			result := sess.CanConnect()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSessionValidation(t *testing.T) {
	t.Run("should validate correct session", func(t *testing.T) {
		sess := session.NewSession("valid-session")
		err := sess.Validate()
		assert.NoError(t, err)
	})

	t.Run("should reject empty session name", func(t *testing.T) {
		sess := session.RestoreSession(
			session.NewSessionID(),
			"",
			session.StatusDisconnected,
			"",
			"",
			false,
			time.Now(),
			time.Now(),
		)
		err := sess.Validate()
		assert.Error(t, err)
	})
}

// Test basic methods to ensure they work
func TestSessionMethods(t *testing.T) {
	sess := session.NewSession("test-session")

	t.Run("ID method should work", func(t *testing.T) {
		id := sess.ID()
		assert.False(t, id.IsEmpty())
	})

	t.Run("Name method should work", func(t *testing.T) {
		name := sess.Name()
		assert.Equal(t, "test-session", name)
	})

	t.Run("Status method should work", func(t *testing.T) {
		status := sess.Status()
		assert.Equal(t, session.StatusDisconnected, status)
	})

	t.Run("CreatedAt and UpdatedAt should work", func(t *testing.T) {
		createdAt := sess.CreatedAt()
		updatedAt := sess.UpdatedAt()
		
		assert.False(t, createdAt.IsZero())
		assert.False(t, updatedAt.IsZero())
		assert.True(t, updatedAt.After(createdAt) || updatedAt.Equal(createdAt))
	})
}