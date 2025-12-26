package user

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrUserNotFound = errors.New("User not found")
	ErrUserExists   = errors.New("User already exists")
)

type MemoryRepository struct {
	mu      sync.RWMutex
	byID    map[string]*User
	byEmail map[string]*User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		byID:    make(map[string]*User),
		byEmail: make(map[string]*User),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[user.Email]; exists {
		return ErrUserExists
	}

	r.byID[user.ID] = user
	r.byEmail[user.Email] = user
	return nil
}

func (r *MemoryRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byEmail[email]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *MemoryRepository) FindByID(ctx context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byID[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}
