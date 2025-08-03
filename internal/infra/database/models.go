package database

import (
	"time"

	"wazmeow/internal/domain/session"

	"github.com/uptrace/bun"
)

// WazMeowSessionModel represents the database model for sessions
type WazMeowSessionModel struct {
	bun.BaseModel `bun:"table:wazmeow_sessions"`

	ID        string    `bun:"id,pk,type:varchar(36)" json:"id"`
	Name      string    `bun:"name,unique,notnull,type:varchar(50)" json:"name"`
	Status    string    `bun:"status,notnull,type:varchar(20),default:'disconnected'" json:"status"`
	WaJID     string    `bun:"wa_jid,type:varchar(100)" json:"wa_jid,omitempty"`
	QRCode    string    `bun:"qr_code,type:text" json:"qr_code,omitempty"`
	IsActive  bool      `bun:"is_active,notnull,default:false" json:"is_active"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

// ToWazMeowSessionModel converts a domain session to database model
func ToWazMeowSessionModel(sess *session.Session) *WazMeowSessionModel {
	return &WazMeowSessionModel{
		ID:        sess.ID().String(),
		Name:      sess.Name(),
		Status:    sess.Status().String(),
		WaJID:     sess.WaJID(),
		QRCode:    sess.QRCode(),
		IsActive:  sess.IsActive(),
		CreatedAt: sess.CreatedAt(),
		UpdatedAt: sess.UpdatedAt(),
	}
}

// FromWazMeowSessionModel converts a database model to domain session
func FromWazMeowSessionModel(model *WazMeowSessionModel) (*session.Session, error) {
	status, err := session.StatusFromString(model.Status)
	if err != nil {
		return nil, err
	}

	sessionID, err := session.SessionIDFromString(model.ID)
	if err != nil {
		return nil, err
	}

	return session.RestoreSession(
		sessionID,
		model.Name,
		status,
		model.WaJID,
		model.QRCode,
		model.IsActive,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}
