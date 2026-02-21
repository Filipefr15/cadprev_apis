package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

func FetchDairCarteira(ctx context.Context, baseURL string, params map[string]string) ([]entity.Dair_carteira, error) {
	query := BuildQueryString(params)
	url := fmt.Sprintf("%s/DAIR_CARTEIRA?%s", baseURL, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var apiResp struct {
		Success bool                   `json:"success"`
		Data    []entity.Dair_carteira `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	return apiResp.Data, nil
}
