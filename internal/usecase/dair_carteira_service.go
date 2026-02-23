package usecase

import (
	"context"
	"database/sql"

	"github.com/Filipefr15/cadprev_apis/internal/api"
	"github.com/Filipefr15/cadprev_apis/internal/domain/entity"
)

type DairCarteiraUseCase struct {
	db *sql.DB
}

func NewDairCarteiraUseCase(db *sql.DB) *DairCarteiraUseCase {
	return &DairCarteiraUseCase{
		db: db,
	}
}

func (uc *DairCarteiraUseCase) GetDairCarteira(ctx context.Context, params map[string]interface{}) ([]entity.Dair_carteira, error) {
	return api.QueryDairCarteira(uc.db, params)
}

func (uc *DairCarteiraUseCase) BuscarDairCarteira(ctx context.Context, params map[string]string) error {
	baseURL := "https://apicadprev.trabalho.gov.br"
	items, err := api.FetchDairCarteira(ctx, baseURL, params)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := api.InsertDairCarteira(uc.db, item); err != nil {
			return err
		}
	}

	return nil
}
