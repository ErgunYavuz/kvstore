package storage

import (
	"fmt"
	"sync"
)

type Storage interface {
	put(key, value string) error
	get(key string) (string, error)
	delete(key string) (bool, error)
	has(key string) bool
	keys() []string
	size() int
}

type MemoryStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: make(map[string]string)}
}

func (s *MemoryStorage) Put(key, value string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	return nil
}

func (s *MemoryStorage) Get(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	v := s.data[key]

	return v, nil
}

func (s *MemoryStorage) Delete(key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key cannot be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)

	return true, nil
}

func (s *MemoryStorage) Has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[key]

	return ok
}

func (s *MemoryStorage) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *MemoryStorage) Print() {
	for k, v := range s.data {
		fmt.Print(k + " : " + v + ", ")
	}
	fmt.Println()
}
