package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/Filipefr15/cadprev_apis/internal/domain"
	"github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

// IngestUseCase orchestrates the data ingestion process
// It fetches data from external API and persists it to the repository
type IngestUseCase struct {
	apiClient  domain.ExternalAPIClient
	repository domain.DataRepository
}

// NewIngestUseCase creates a new IngestUseCase with dependency injection
func NewIngestUseCase(apiClient domain.ExternalAPIClient, repository domain.DataRepository) *IngestUseCase {
	return &IngestUseCase{
		apiClient:  apiClient,
		repository: repository,
	}
}

// IngestAll fetches all data from external API and persists to repository
func (uc *IngestUseCase) IngestAll(ctx context.Context) error {
	// Fetch data from external API
	dataList, err := uc.apiClient.FetchData(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch data from external API: %w", err)
	}

	log.Printf("Fetched %d records from external API", len(dataList))

	// Persist each data item to repository
	for _, data := range dataList {
		if err := uc.repository.Save(ctx, data); err != nil {
			log.Printf("Failed to save data with ID %s: %v", data.ID, err)
			// Continue with other records even if one fails
			continue
		}
		log.Printf("Successfully saved data with ID: %s", data.ID)
	}

	return nil
}

// IngestByID fetches a specific data item by ID and persists it
func (uc *IngestUseCase) IngestByID(ctx context.Context, id string) (*entity.Data, error) {
	// Fetch data from external API
	data, err := uc.apiClient.FetchDataByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data by ID from external API: %w", err)
	}

	// Persist to repository
	if err := uc.repository.Save(ctx, data); err != nil {
		return nil, fmt.Errorf("failed to save data: %w", err)
	}

	log.Printf("Successfully ingested data with ID: %s", data.ID)
	return data, nil
}

// GetAll retrieves all data from repository
func (uc *IngestUseCase) GetAll(ctx context.Context) ([]*entity.Data, error) {
	return uc.repository.FindAll(ctx)
}

// GetByID retrieves a specific data item by ID from repository
func (uc *IngestUseCase) GetByID(ctx context.Context, id string) (*entity.Data, error) {
	return uc.repository.FindByID(ctx, id)
}
