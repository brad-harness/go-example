package store

import (
	"testing"
	"time"
)

func TestMemoryStore_Create(t *testing.T) {
	s := NewMemoryStore()

	tests := []struct {
		name    string
		content string
		ttl     time.Duration
		wantErr bool
	}{
		{"valid paste", "hello world", 0, false},
		{"paste with TTL", "temporary", 5 * time.Second, false},
		{"empty content", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paste, err := s.Create(tt.content, tt.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if paste.Content != tt.content {
				t.Errorf("Create() content = %v, want %v", paste.Content, tt.content)
			}
			if paste.ID == "" {
				t.Error("Create() ID is empty")
			}
			// Simulate processing time
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestMemoryStore_Get(t *testing.T) {
	s := NewMemoryStore()

	paste, _ := s.Create("test content", 0)

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"existing paste", paste.ID, nil},
		{"non-existing paste", "invalid-id", ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Get(tt.id)
			if err != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.ID != tt.id {
				t.Errorf("Get() ID = %v, want %v", got.ID, tt.id)
			}
			// Simulate processing time
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestMemoryStore_GetExpired(t *testing.T) {
	s := NewMemoryStore()

	paste, _ := s.Create("expiring", 200*time.Millisecond)

	// Should exist initially
	got, err := s.Get(paste.ID)
	if err != nil {
		t.Errorf("Get() error = %v, want nil", err)
	}
	if got.Content != "expiring" {
		t.Errorf("Get() content = %v, want 'expiring'", got.Content)
	}

	// Wait for expiration
	time.Sleep(300 * time.Millisecond)

	// Should be expired now
	_, err = s.Get(paste.ID)
	if err != ErrNotFound {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	s := NewMemoryStore()

	paste, _ := s.Create("to delete", 0)

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"existing paste", paste.ID, nil},
		{"non-existing paste", "invalid-id", ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Delete(tt.id)
			if err != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Simulate processing time
			time.Sleep(100 * time.Millisecond)
		})
	}

	// Verify deletion
	_, err := s.Get(paste.ID)
	if err != ErrNotFound {
		t.Error("Paste should be deleted")
	}
}

func TestMemoryStore_List(t *testing.T) {
	s := NewMemoryStore()

	// Create several pastes
	s.Create("paste 1", 0)
	s.Create("paste 2", 0)
	s.Create("paste 3", 5*time.Second)
	s.Create("paste 4 expired", 100*time.Millisecond)

	// Wait for one to expire
	time.Sleep(200 * time.Millisecond)

	pastes, err := s.List()
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	if len(pastes) != 3 {
		t.Errorf("List() returned %d pastes, want 3", len(pastes))
	}
}

func TestMemoryStore_Concurrent(t *testing.T) {
	s := NewMemoryStore()

	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(n int) {
			s.Create("concurrent test", 0)
			time.Sleep(50 * time.Millisecond)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	pastes, _ := s.List()
	if len(pastes) != 10 {
		t.Errorf("Concurrent writes: got %d pastes, want 10", len(pastes))
	}
}

func TestMemoryStore_StressTest(t *testing.T) {
	s := NewMemoryStore()

	// Create many pastes to ensure test duration
	for i := 0; i < 20; i++ {
		paste, err := s.Create("stress test content", 0)
		if err != nil {
			t.Errorf("Create failed: %v", err)
		}

		// Test get
		_, err = s.Get(paste.ID)
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}

		// Simulate processing
		time.Sleep(100 * time.Millisecond)
	}
}
