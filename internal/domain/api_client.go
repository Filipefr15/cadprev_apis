package domain

import (
	"context"

	"github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

// ExternalAPIClient defines the interface for external API communication
// This interface belongs to the domain layer but is implemented in infrastructure
type ExternalAPIClient interface {
	FetchData(ctx context.Context) ([]*entity.Data, error)
	FetchDataByID(ctx context.Context, id string) (*entity.Data, error)
}
