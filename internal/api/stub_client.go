package api

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

// StubAPIClient is a stub implementation of ExternalAPIClient
// In production, this would make real HTTP calls to external APIs
type StubAPIClient struct {
	baseURL string
}

// NewStubAPIClient creates a new stub API client
func NewStubAPIClient(baseURL string) *StubAPIClient {
	return &StubAPIClient{
		baseURL: baseURL,
	}
}

// FetchData simulates fetching data from an external API
func (c *StubAPIClient) FetchData(ctx context.Context) ([]*entity.Data, error) {
	// Simulate API call delay
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	// Return stub data simulating API response
	data := make([]*entity.Data, 0, 5)
	for i := 1; i <= 5; i++ {
		data = append(data, entity.NewData(
			fmt.Sprintf("ext-api-%d", i),
			fmt.Sprintf("Data Item %d", i),
			fmt.Sprintf("Description for item %d from external API", i),
			c.baseURL,
			rand.Float64()*1000,
		))
	}

	return data, nil
}

// FetchDataByID simulates fetching a specific data item by ID from external API
func (c *StubAPIClient) FetchDataByID(ctx context.Context, id string) (*entity.Data, error) {
	// Simulate API call delay
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(50 * time.Millisecond):
	}

	// Return stub data for the specific ID
	return entity.NewData(
		id,
		fmt.Sprintf("Data Item %s", id),
		fmt.Sprintf("Description for item %s from external API", id),
		c.baseURL,
		rand.Float64()*1000,
	), nil
}
