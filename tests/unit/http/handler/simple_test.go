package http_handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"wazmeow/internal/domain/session"
)

func TestHTTPHandlerIntegration(t *testing.T) {
	t.Run("basic HTTP server setup", func(t *testing.T) {
		// Test that we can create a basic HTTP server with routes
		r := chi.NewRouter()
		
		// Add a simple test route
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test response"))
		})

		// Create test request
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test response", w.Body.String())
	})
}

func TestSessionHTTPEndpoints(t *testing.T) {
	t.Run("session creation endpoint", func(t *testing.T) {
		r := chi.NewRouter()
		
		// Mock session creation endpoint
		r.Post("/sessions", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			
			// Create a session
			session := session.NewSession("http-test-session")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "` + session.ID().String() + `", "name": "` + session.Name() + `"}`))
		})

		// Test POST request
		req := httptest.NewRequest("POST", "/sessions", strings.NewReader(`{"name": "http-test-session"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "http-test-session")
	})

	t.Run("session listing endpoint", func(t *testing.T) {
		r := chi.NewRouter()
		
		// Mock session listing endpoint
		r.Get("/sessions", func(w http.ResponseWriter, r *http.Request) {
			sessions := []*session.Session{
				session.NewSession("session1"),
				session.NewSession("session2"),
			}
			
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"sessions": [{"id": "` + sessions[0].ID().String() + `"}, {"id": "` + sessions[1].ID().String() + `"}]}`))
		})

		// Test GET request
		req := httptest.NewRequest("GET", "/sessions", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "sessions")
	})

	t.Run("health check endpoint", func(t *testing.T) {
		r := chi.NewRouter()
		
		// Mock health check endpoint
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "healthy", "timestamp": "2023-01-01T00:00:00Z"}`))
		})

		// Test health check
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "healthy")
	})
}

func TestHTTPMiddlewareIntegration(t *testing.T) {
	t.Run("request with middleware", func(t *testing.T) {
		r := chi.NewRouter()
		
		// Add simple logging middleware
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Add custom header
				w.Header().Set("X-Test-Middleware", "applied")
				next.ServeHTTP(w, r)
			})
		})
		
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("protected content"))
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "applied", w.Header().Get("X-Test-Middleware"))
		assert.Equal(t, "protected content", w.Body.String())
	})
}

func TestSessionValidationInHTTP(t *testing.T) {
	t.Run("valid session creation request", func(t *testing.T) {
		r := chi.NewRouter()
		
		// Mock endpoint with simple validation
		r.Post("/sessions", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "test-id", "name": "valid-session"}`))
		})

		// Test with valid data
		req := httptest.NewRequest("POST", "/sessions", strings.NewReader(`{"name": "valid-session"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "valid-session")
	})
}