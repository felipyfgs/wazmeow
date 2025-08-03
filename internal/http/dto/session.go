package dto

import (
	"time"

	"wazmeow/internal/domain/session"
)

// CreateSessionRequest represents the HTTP request to create a session
// @Description Dados para criação de uma nova sessão WhatsApp
type CreateSessionRequest struct {
	Name string `json:"name" validate:"required,session_name" example:"minha-sessao" description:"Nome único da sessão (apenas letras, números, hífens e underscores)"`
}

// SessionResponse represents the HTTP response for a session
// @Description Dados de uma sessão WhatsApp
type SessionResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" description:"ID único da sessão (UUID)"`
	Name      string    `json:"name" example:"minha-sessao" description:"Nome da sessão"`
	Status    string    `json:"status" example:"connected" enums:"disconnected,connecting,connected" description:"Status atual da sessão"`
	WaJID     string    `json:"wa_jid,omitempty" example:"5511999999999@s.whatsapp.net" description:"JID do WhatsApp (quando conectado)"`
	IsActive  bool      `json:"is_active" example:"true" description:"Indica se a sessão está ativa"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z" description:"Data de criação da sessão"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T12:30:00Z" description:"Data da última atualização"`
}

// SessionListResponse represents the HTTP response for listing sessions
// @Description Lista de sessões WhatsApp
type SessionListResponse struct {
	Sessions []*SessionResponse `json:"sessions" description:"Lista de sessões"`
	Total    int                `json:"total" example:"5" description:"Total de sessões encontradas"`
}

// ConnectSessionRequest represents the HTTP request to connect a session
type ConnectSessionRequest struct {
	// No additional fields needed - session ID comes from URL
}

// ConnectSessionResponse represents the HTTP response for connecting a session
// @Description Resposta da operação de conexão de sessão
type ConnectSessionResponse struct {
	Session   *SessionResponse `json:"session" description:"Dados atualizados da sessão"`
	QRCode    string           `json:"qr_code,omitempty" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..." description:"QR Code em base64 (quando necessário)"`
	NeedsAuth bool             `json:"needs_auth" example:"true" description:"Indica se é necessário escanear QR Code"`
	Message   string           `json:"message" example:"QR Code gerado. Escaneie com seu WhatsApp." description:"Mensagem informativa"`
}

// DisconnectSessionRequest represents the HTTP request to disconnect a session
type DisconnectSessionRequest struct {
	// No additional fields needed - session ID comes from URL
}

// DisconnectSessionResponse represents the HTTP response for disconnecting a session
type DisconnectSessionResponse struct {
	Session *SessionResponse `json:"session"`
	Message string           `json:"message"`
}

// DeleteSessionRequest represents the HTTP request to delete a session
type DeleteSessionRequest struct {
	// No fields needed - deletion always forces
}

// DeleteSessionResponse represents the HTTP response for deleting a session
type DeleteSessionResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// QRCodeResponse represents the HTTP response for QR code generation
// @Description Resposta com QR Code para autenticação
type QRCodeResponse struct {
	SessionID string `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000" description:"ID da sessão"`
	QRCode    string `json:"qr_code" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..." description:"QR Code em base64"`
	Message   string `json:"message" example:"QR Code gerado com sucesso" description:"Mensagem informativa"`
}

// PairPhoneRequest represents the HTTP request to pair with a phone number
// @Description Dados para emparelhamento com número de telefone
type PairPhoneRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required" example:"5511999999999" description:"Número de telefone para emparelhar"`
}

// PairPhoneResponse represents the HTTP response for phone pairing
// @Description Resposta do emparelhamento com telefone
type PairPhoneResponse struct {
	SessionID   string `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000" description:"ID da sessão"`
	PhoneNumber string `json:"phone_number" example:"5511999999999" description:"Número emparelhado"`
	Success     bool   `json:"success" example:"true" description:"Indica se o emparelhamento foi bem-sucedido"`
	Message     string `json:"message" example:"Telefone emparelhado com sucesso" description:"Mensagem informativa"`
}

// ProxySetRequest represents the HTTP request to set proxy configuration
// @Description Configuração de proxy para a sessão
type ProxySetRequest struct {
	ProxyURL string `json:"proxy_url" validate:"required" example:"http://proxy.example.com:8080" description:"URL do proxy"`
	Username string `json:"username,omitempty" example:"user" description:"Usuário do proxy (opcional)"`
	Password string `json:"password,omitempty" example:"pass" description:"Senha do proxy (opcional)"`
}

// ProxySetResponse represents the HTTP response for proxy configuration
// @Description Resposta da configuração de proxy
type ProxySetResponse struct {
	SessionID string `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000" description:"ID da sessão"`
	ProxyURL  string `json:"proxy_url" example:"http://proxy.example.com:8080" description:"URL do proxy configurado"`
	Success   bool   `json:"success" example:"true" description:"Indica se a configuração foi bem-sucedida"`
	Message   string `json:"message" example:"Proxy configurado com sucesso" description:"Mensagem informativa"`
}

// ToSessionResponse converts a domain session to HTTP response
func ToSessionResponse(sess *session.Session) *SessionResponse {
	return &SessionResponse{
		ID:        sess.ID().String(),
		Name:      sess.Name(),
		Status:    sess.Status().String(),
		WaJID:     sess.WaJID(),
		IsActive:  sess.IsActive(),
		CreatedAt: sess.CreatedAt(),
		UpdatedAt: sess.UpdatedAt(),
	}
}

// ToSessionListResponse converts a list of domain sessions to HTTP response
func ToSessionListResponse(sessions []*session.Session, total int) *SessionListResponse {
	sessionResponses := make([]*SessionResponse, len(sessions))
	for i, sess := range sessions {
		sessionResponses[i] = ToSessionResponse(sess)
	}

	return &SessionListResponse{
		Sessions: sessionResponses,
		Total:    total,
	}
}
