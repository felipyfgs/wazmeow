package dto

import "time"

// SuccessResponse represents a generic success response
// @Description Resposta de sucesso padrão da API
type SuccessResponse struct {
	Success bool        `json:"success" example:"true" description:"Indica se a operação foi bem-sucedida"`
	Message string      `json:"message" example:"Operação realizada com sucesso" description:"Mensagem descritiva do resultado"`
	Data    interface{} `json:"data,omitempty" description:"Dados retornados pela operação (opcional)"`
}

// ErrorResponse represents a generic error response
// @Description Resposta de erro padrão da API
type ErrorResponse struct {
	Success bool        `json:"success" example:"false" description:"Sempre false para respostas de erro"`
	Error   string      `json:"error" example:"Erro interno do servidor" description:"Mensagem de erro"`
	Code    string      `json:"code,omitempty" example:"INTERNAL_ERROR" description:"Código do erro (opcional)"`
	Details string      `json:"details,omitempty" example:"Detalhes técnicos do erro" description:"Detalhes adicionais do erro (opcional)"`
	Context interface{} `json:"context,omitempty" description:"Contexto adicional do erro (opcional)"`
}

// ValidationErrorResponse represents a validation error response
type ValidationErrorResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Code    string                 `json:"code"`
	Fields  []ValidationFieldError `json:"fields"`
}

// ValidationFieldError represents a field validation error
type ValidationFieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// PaginationRequest represents pagination parameters
// @Description Parâmetros de paginação para listagens
type PaginationRequest struct {
	Limit  int `json:"limit" query:"limit" validate:"min=1,max=100" example:"10" description:"Número máximo de itens por página (1-100)"`
	Offset int `json:"offset" query:"offset" validate:"min=0" example:"0" description:"Número de itens a pular (para paginação)"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Pages  int `json:"pages"`
}

// HealthResponse represents the health check response
// @Description Resposta do health check da aplicação
type HealthResponse struct {
	Status    string                 `json:"status" example:"healthy" description:"Status geral da aplicação"`
	Timestamp time.Time              `json:"timestamp" example:"2024-01-01T12:00:00Z" description:"Timestamp da verificação"`
	Version   string                 `json:"version" example:"1.0.0" description:"Versão da aplicação"`
	Uptime    string                 `json:"uptime" example:"2h30m45s" description:"Tempo de atividade da aplicação"`
	Services  map[string]interface{} `json:"services" description:"Status dos serviços individuais"`
}

// ServiceHealth represents the health status of a service
type ServiceHealth struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// MetricsResponse represents application metrics
// @Description Métricas completas da aplicação
type MetricsResponse struct {
	Sessions  SessionMetrics  `json:"sessions" description:"Métricas das sessões"`
	WhatsApp  WhatsAppMetrics `json:"whatsapp" description:"Métricas do WhatsApp"`
	System    SystemMetrics   `json:"system" description:"Métricas do sistema"`
	Timestamp time.Time       `json:"timestamp" example:"2024-01-01T12:00:00Z" description:"Timestamp da coleta das métricas"`
}

// SessionMetrics represents session-related metrics
// @Description Métricas relacionadas às sessões WhatsApp
type SessionMetrics struct {
	Total        int `json:"total" example:"10" description:"Total de sessões"`
	Connected    int `json:"connected" example:"5" description:"Sessões conectadas"`
	Disconnected int `json:"disconnected" example:"3" description:"Sessões desconectadas"`
	Error        int `json:"error" example:"1" description:"Sessões com erro"`
	Active       int `json:"active" example:"4" description:"Sessões ativas"`
}

// WhatsAppMetrics represents WhatsApp-related metrics
// @Description Métricas relacionadas ao WhatsApp
type WhatsAppMetrics struct {
	TotalClients         int `json:"total_clients" example:"5" description:"Total de clientes WhatsApp"`
	ConnectedClients     int `json:"connected_clients" example:"3" description:"Clientes conectados"`
	AuthenticatedClients int `json:"authenticated_clients" example:"2" description:"Clientes autenticados"`
	ErrorClients         int `json:"error_clients" example:"1" description:"Clientes com erro"`
	MessagesSent         int `json:"messages_sent" example:"150" description:"Total de mensagens enviadas"`
	MessagesReceived     int `json:"messages_received" example:"75" description:"Total de mensagens recebidas"`
}

// SystemMetrics represents system-related metrics
// @Description Métricas relacionadas ao sistema
type SystemMetrics struct {
	Uptime              string `json:"uptime" example:"2h30m45s" description:"Tempo de atividade do sistema"`
	MemoryUsage         string `json:"memory_usage" example:"256MB" description:"Uso de memória"`
	CPUUsage            string `json:"cpu_usage" example:"15%" description:"Uso de CPU"`
	DatabaseStatus      string `json:"database_status" example:"healthy" description:"Status do banco de dados"`
	DatabaseConnections int    `json:"database_connections" example:"5" description:"Número de conexões ativas no banco"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(error, code, details string) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error:   error,
		Code:    code,
		Details: details,
	}
}

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(fields []ValidationFieldError) *ValidationErrorResponse {
	return &ValidationErrorResponse{
		Success: false,
		Error:   "Validation failed",
		Code:    "VALIDATION_ERROR",
		Fields:  fields,
	}
}

// CalculatePages calculates the number of pages for pagination
func (p *PaginationResponse) CalculatePages() {
	if p.Limit > 0 {
		p.Pages = (p.Total + p.Limit - 1) / p.Limit
	}
}

// NewPaginationResponse creates a new pagination response
func NewPaginationResponse(total, limit, offset int) *PaginationResponse {
	pagination := &PaginationResponse{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
	pagination.CalculatePages()
	return pagination
}
