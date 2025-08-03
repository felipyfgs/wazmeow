package session

import (
	"time"
)

// Session represents a WhatsApp session entity
type Session struct {
	id        SessionID
	name      string
	status    Status
	waJID     string
	qrCode    string
	isActive  bool
	createdAt time.Time
	updatedAt time.Time
}

// NewSession creates a new session with the given name
func NewSession(name string) *Session {
	if name == "" {
		panic("session name cannot be empty")
	}

	return &Session{
		id:        NewSessionID(),
		name:      name,
		status:    StatusDisconnected,
		waJID:     "",
		qrCode:    "",
		isActive:  false,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

// RestoreSession restores a session from persistence
func RestoreSession(id SessionID, name string, status Status, waJID string, qrCode string, isActive bool, createdAt, updatedAt time.Time) *Session {
	return &Session{
		id:        id,
		name:      name,
		status:    status,
		waJID:     waJID,
		qrCode:    qrCode,
		isActive:  isActive,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// Connect marks the session as connected with the given WhatsApp JID
func (s *Session) Connect(waJID string) error {
	if s.status == StatusConnected {
		return ErrSessionAlreadyConnected
	}

	if waJID == "" {
		return ErrInvalidWhatsAppJID
	}

	s.waJID = waJID
	s.status = StatusConnected
	s.isActive = true
	s.updatedAt = time.Now()

	return nil
}

// Disconnect marks the session as disconnected
func (s *Session) Disconnect() {
	s.status = StatusDisconnected
	s.isActive = false
	s.updatedAt = time.Now()
}

// SetConnecting marks the session as connecting
func (s *Session) SetConnecting() {
	s.status = StatusConnecting
	s.updatedAt = time.Now()
}

// SetQRCode updates the session QR code
func (s *Session) SetQRCode(qrCode string) {
	s.qrCode = qrCode
	s.updatedAt = time.Now()
}

// ClearQRCode clears the session QR code
func (s *Session) ClearQRCode() {
	s.qrCode = ""
	s.updatedAt = time.Now()
}

// UpdateName updates the session name
func (s *Session) UpdateName(name string) error {
	if name == "" {
		return ErrInvalidSessionName
	}

	s.name = name
	s.updatedAt = time.Now()
	return nil
}

// CanConnect returns true if the session can be connected
func (s *Session) CanConnect() bool {
	// Allow connection for disconnected and connecting states
	// Only prevent connection if already connected
	return s.status != StatusConnected
}

// IsConnected returns true if the session is connected
func (s *Session) IsConnected() bool {
	return s.status == StatusConnected && s.isActive
}

// IsConnecting returns true if the session is in connecting state
func (s *Session) IsConnecting() bool {
	return s.status == StatusConnecting
}

// Getters
func (s *Session) ID() SessionID {
	return s.id
}

func (s *Session) Name() string {
	return s.name
}

func (s *Session) Status() Status {
	return s.status
}

func (s *Session) WaJID() string {
	return s.waJID
}

func (s *Session) QRCode() string {
	return s.qrCode
}

func (s *Session) IsActive() bool {
	return s.isActive
}

func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Session) UpdatedAt() time.Time {
	return s.updatedAt
}

// Validate validates the session entity
func (s *Session) Validate() error {
	if s.name == "" {
		return ErrInvalidSessionName
	}

	if len(s.name) < 3 || len(s.name) > 50 {
		return ErrInvalidSessionName
	}

	return nil
}
