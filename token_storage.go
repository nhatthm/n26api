package n26api

import (
	"context"
	"errors"

	"github.com/nhatthm/n26api/pkg/auth"
)

// ErrTokenKeyEmpty indicates that the key of token is empty and we can not persist that.
var ErrTokenKeyEmpty = errors.New("token key is empty")

var _ auth.TokenStorage = (*InMemoryTokenStorage)(nil)

// InMemoryTokenStorage persists auth.OAuthToken into its memory.
type InMemoryTokenStorage struct {
	storage map[string]auth.OAuthToken
}

// Get gets OAuthToken from memory.
func (s *InMemoryTokenStorage) Get(_ context.Context, key string) (auth.OAuthToken, error) {
	return s.storage[key], nil
}

// Set sets OAuthToken to memory.
func (s *InMemoryTokenStorage) Set(_ context.Context, key string, token auth.OAuthToken) error {
	if key == "" {
		return ErrTokenKeyEmpty
	}

	s.storage[key] = token

	return nil
}

// NewInMemoryTokenStorage initiates a new InMemoryTokenStorage.
func NewInMemoryTokenStorage() *InMemoryTokenStorage {
	return &InMemoryTokenStorage{
		storage: make(map[string]auth.OAuthToken),
	}
}
