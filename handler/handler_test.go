package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/colin-harness/go-example/store"
	"github.com/gin-gonic/gin"
)

func setupTest() (*Handler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	s := store.NewMemoryStore()
	h := NewHandler(s)
	r := gin.New()
	return h, r
}

func TestHandler_CreatePaste(t *testing.T) {
	h, r := setupTest()
	r.POST("/paste", h.CreatePaste)

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name:       "valid paste",
			body:       map[string]interface{}{"content": "hello world"},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "valid paste with TTL",
			body:       map[string]interface{}{"content": "temporary", "ttl": 60},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "empty content",
			body:       map[string]interface{}{"content": ""},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/paste", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Status = %v, want %v", w.Code, tt.wantStatus)
			}

			// Simulate processing time
			time.Sleep(150 * time.Millisecond)
		})
	}
}

func TestHandler_GetPasteJSON(t *testing.T) {
	h, r := setupTest()
	r.GET("/api/paste/:id", h.GetPasteJSON)

	// Create a paste first
	paste, _ := h.store.Create("test content", 0)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "existing paste",
			id:         paste.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing paste",
			id:         "invalid-id",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/paste/"+tt.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var result store.Paste
				json.Unmarshal(w.Body.Bytes(), &result)
				if result.ID != tt.id {
					t.Errorf("Paste ID = %v, want %v", result.ID, tt.id)
				}
			}

			// Simulate processing time
			time.Sleep(150 * time.Millisecond)
		})
	}
}

func TestHandler_DeletePaste(t *testing.T) {
	h, r := setupTest()
	r.DELETE("/api/paste/:id", h.DeletePaste)

	// Create a paste first
	paste, _ := h.store.Create("to delete", 0)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "existing paste",
			id:         paste.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing paste",
			id:         "invalid-id",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/api/paste/"+tt.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Status = %v, want %v", w.Code, tt.wantStatus)
			}

			// Simulate processing time
			time.Sleep(150 * time.Millisecond)
		})
	}
}

func TestHandler_MultipleOperations(t *testing.T) {
	h, r := setupTest()
	r.POST("/paste", h.CreatePaste)
	r.GET("/api/paste/:id", h.GetPasteJSON)
	r.DELETE("/api/paste/:id", h.DeletePaste)

	// Create multiple pastes and perform operations
	for i := 0; i < 10; i++ {
		body, _ := json.Marshal(map[string]interface{}{"content": "test"})
		req, _ := http.NewRequest("POST", "/paste", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Create failed with status %v", w.Code)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func TestHandler_ExpiredPaste(t *testing.T) {
	h, r := setupTest()
	r.GET("/api/paste/:id", h.GetPasteJSON)

	// Create a paste with short TTL
	paste, _ := h.store.Create("expiring", 200*time.Millisecond)

	// Should exist initially
	req, _ := http.NewRequest("GET", "/api/paste/"+paste.ID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Initial GET status = %v, want %v", w.Code, http.StatusOK)
	}

	// Wait for expiration
	time.Sleep(300 * time.Millisecond)

	// Should be expired now
	req, _ = http.NewRequest("GET", "/api/paste/"+paste.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expired GET status = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestHandler_ConcurrentRequests(t *testing.T) {
	h, r := setupTest()
	r.POST("/paste", h.CreatePaste)

	done := make(chan bool)

	// Concurrent requests
	for i := 0; i < 15; i++ {
		go func() {
			body, _ := json.Marshal(map[string]interface{}{"content": "concurrent"})
			req, _ := http.NewRequest("POST", "/paste", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				t.Errorf("Concurrent request failed with status %v", w.Code)
			}

			time.Sleep(100 * time.Millisecond)
			done <- true
		}()
	}

	// Wait for all requests
	for i := 0; i < 15; i++ {
		<-done
	}
}
