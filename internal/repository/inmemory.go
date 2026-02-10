package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

// InMemoryRepository is a stub implementation of DataRepository using in-memory storage
// In production, this would use a real database (PostgreSQL, MongoDB, etc.)
type InMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*entity.Data
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]*entity.Data),
	}
}

// Save persists a data entity (creates or updates)
func (r *InMemoryRepository) Save(ctx context.Context, data *entity.Data) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	if data.ID == "" {
		return fmt.Errorf("data ID cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Clone the data to avoid external modifications
	stored := *data
	r.data[data.ID] = &stored

	return nil
}

// FindByID retrieves a data entity by ID
func (r *InMemoryRepository) FindByID(ctx context.Context, id string) (*entity.Data, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data, exists := r.data[id]
	if !exists {
		return nil, fmt.Errorf("data with ID %s not found", id)
	}

	// Return a copy to avoid external modifications
	result := *data
	return &result, nil
}

// FindAll retrieves all data entities
func (r *InMemoryRepository) FindAll(ctx context.Context) ([]*entity.Data, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	result := make([]*entity.Data, 0, len(r.data))
	for _, data := range r.data {
		// Return copies to avoid external modifications
		copy := *data
		result = append(result, &copy)
	}

	return result, nil
}

// Delete removes a data entity by ID
func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, exists := r.data[id]; !exists {
		return fmt.Errorf("data with ID %s not found", id)
	}

	delete(r.data, id)
	return nil
}
