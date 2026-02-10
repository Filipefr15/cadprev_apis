package domain

import (
	"context"

	"github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

// DataRepository defines the interface for data persistence
// This interface belongs to the domain layer but is implemented in infrastructure
type DataRepository interface {
	Save(ctx context.Context, data *entity.Data) error
	FindByID(ctx context.Context, id string) (*entity.Data, error)
	FindAll(ctx context.Context) ([]*entity.Data, error)
	Delete(ctx context.Context, id string) error
}
