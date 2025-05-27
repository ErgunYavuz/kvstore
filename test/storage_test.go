package test

import (
	"fmt"
	"kvstore/storage"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	storage := storage.NewMemoryStorage()

	err := storage.Put("test-key", "test-value")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	value, err := storage.Get("test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", value)
	}

	if !storage.Has("test-key") {
		t.Error("Expected key to exist")
	}

	if storage.Size() != 1 {
		t.Errorf("Expected size 1, got %d", storage.Size())
	}

	deleted, err := storage.Delete("test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !deleted {
		t.Error("Expected key to be deleted")
	}

	if storage.Size() != 0 {
		t.Errorf("Expected size 0 after delete, got %d", storage.Size())
	}
}

func TestStorageConcurrency(t *testing.T) {
	storage := storage.NewMemoryStorage()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			key := fmt.Sprintf("key%d", id)
			value := fmt.Sprintf("value%d", id)

			err := storage.Put(key, value)
			if err != nil {
				t.Errorf("Put failed for %s: %v", key, err)
			}

			retrievedValue, err := storage.Get(key)
			if err != nil {
				t.Errorf("Get failed for %s: %v", key, err)
			}
			if retrievedValue != value {
				t.Errorf("Expected %s, got %s", value, retrievedValue)
			}

			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if storage.Size() != 10 {
		t.Fatalf("Expected 10 keys, got %d", storage.Size())
	}
}
