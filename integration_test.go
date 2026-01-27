package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/colin-harness/go-example/handler"
	"github.com/colin-harness/go-example/store"
)

func TestIntegration_FullWorkflow(t *testing.T) {
	s := store.NewMemoryStore()
	h := handler.NewHandler(s)
	r := setupRouter(h)

	// Test 1: Create a paste
	createBody, _ := json.Marshal(map[string]interface{}{
		"content": "integration test content",
		"ttl":     0,
	})

	req := httptest.NewRequest("POST", "/paste", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Create paste failed: status = %v", w.Code)
	}

	var createdPaste store.Paste
	json.Unmarshal(w.Body.Bytes(), &createdPaste)

	if createdPaste.ID == "" {
		t.Fatal("Created paste has no ID")
	}

	time.Sleep(200 * time.Millisecond)

	// Test 2: Retrieve the paste
	req = httptest.NewRequest("GET", "/api/paste/"+createdPaste.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Get paste failed: status = %v", w.Code)
	}

	var retrievedPaste store.Paste
	json.Unmarshal(w.Body.Bytes(), &retrievedPaste)

	if retrievedPaste.Content != "integration test content" {
		t.Errorf("Content = %v, want 'integration test content'", retrievedPaste.Content)
	}

	time.Sleep(200 * time.Millisecond)

	// Test 3: Delete the paste
	req = httptest.NewRequest("DELETE", "/api/paste/"+createdPaste.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Delete paste failed: status = %v", w.Code)
	}

	time.Sleep(200 * time.Millisecond)

	// Test 4: Verify deletion
	req = httptest.NewRequest("GET", "/api/paste/"+createdPaste.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Deleted paste still accessible: status = %v", w.Code)
	}

	time.Sleep(200 * time.Millisecond)
}

func TestIntegration_PasteExpiration(t *testing.T) {
	s := store.NewMemoryStore()
	h := handler.NewHandler(s)
	r := setupRouter(h)

	// Create paste with 1 second TTL
	createBody, _ := json.Marshal(map[string]interface{}{
		"content": "expiring content",
		"ttl":     1,
	})

	req := httptest.NewRequest("POST", "/paste", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Create paste failed: status = %v", w.Code)
	}

	var paste store.Paste
	json.Unmarshal(w.Body.Bytes(), &paste)

	time.Sleep(300 * time.Millisecond)

	// Should exist initially
	req = httptest.NewRequest("GET", "/api/paste/"+paste.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Paste should exist: status = %v", w.Code)
	}

	// Wait for expiration
	time.Sleep(1 * time.Second)

	// Should be expired
	req = httptest.NewRequest("GET", "/api/paste/"+paste.ID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Paste should be expired: status = %v", w.Code)
	}
}

func TestIntegration_MultiplePastes(t *testing.T) {
	s := store.NewMemoryStore()
	h := handler.NewHandler(s)
	r := setupRouter(h)

	pasteIDs := make([]string, 0)

	// Create multiple pastes
	for i := 0; i < 10; i++ {
		createBody, _ := json.Marshal(map[string]interface{}{
			"content": "paste content",
			"ttl":     0,
		})

		req := httptest.NewRequest("POST", "/paste", bytes.NewBuffer(createBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Create paste %d failed: status = %v", i, w.Code)
		}

		var paste store.Paste
		json.Unmarshal(w.Body.Bytes(), &paste)
		pasteIDs = append(pasteIDs, paste.ID)

		time.Sleep(100 * time.Millisecond)
	}

	// Verify all pastes exist
	for i, id := range pasteIDs {
		req := httptest.NewRequest("GET", "/api/paste/"+id, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Paste %d not found: status = %v", i, w.Code)
		}

		time.Sleep(50 * time.Millisecond)
	}

	// Delete all pastes
	for _, id := range pasteIDs {
		req := httptest.NewRequest("DELETE", "/api/paste/"+id, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Delete failed for paste %s: status = %v", id, w.Code)
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func TestIntegration_InvalidRequests(t *testing.T) {
	s := store.NewMemoryStore()
	h := handler.NewHandler(s)
	r := setupRouter(h)

	tests := []struct {
		name       string
		method     string
		path       string
		body       interface{}
		wantStatus int
	}{
		{
			name:       "empty content",
			method:     "POST",
			path:       "/paste",
			body:       map[string]interface{}{"content": ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid JSON",
			method:     "POST",
			path:       "/paste",
			body:       "not json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "non-existent paste",
			method:     "GET",
			path:       "/api/paste/does-not-exist",
			body:       nil,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "delete non-existent",
			method:     "DELETE",
			path:       "/api/paste/does-not-exist",
			body:       nil,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				if str, ok := tt.body.(string); ok {
					req = httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(str))
				} else {
					body, _ := json.Marshal(tt.body)
					req = httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(body))
				}
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Status = %v, want %v", w.Code, tt.wantStatus)
			}

			time.Sleep(100 * time.Millisecond)
		})
	}
}
