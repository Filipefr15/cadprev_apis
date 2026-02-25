package usecase

import (
	"context"
	"database/sql"

	"github.com/Filipefr15/cadprev_apis/internal/api"
	"github.com/Filipefr15/cadprev_apis/internal/domain"
)

type DashboardUseCase struct {
	db *sql.DB
}

func NewDashboardUseCase(db *sql.DB) *DashboardUseCase {
	return &DashboardUseCase{
		db: db,
	}
}

// GetPatrimonioMensal retorna o patrimônio total mensal
func (uc *DashboardUseCase) GetPatrimonioMensal(ctx context.Context, nrCnpj string, dtAno string) ([]domain.PatrimonioMensalResponse, error) {
	return api.GetPatrimonioMensal(uc.db, nrCnpj, dtAno)
}

// GetComposicaoSegmento retorna a composição por segmento
func (uc *DashboardUseCase) GetComposicaoSegmento(ctx context.Context, nrCnpj string, dtAno string, dtMesBimestre string) ([]domain.ComposicaoSegmentoResponse, error) {
	return api.GetComposicaoSegmento(uc.db, nrCnpj, dtAno, dtMesBimestre)
}

// GetEvolucaoAnual retorna a evolução anual
func (uc *DashboardUseCase) GetEvolucaoAnual(ctx context.Context, nrCnpj string) ([]domain.EvolucaoAnualResponse, error) {
	return api.GetEvolucaoAnual(uc.db, nrCnpj)
}

// GetVariacaoMensal retorna a variação mensal (para rentabilidade)
func (uc *DashboardUseCase) GetVariacaoMensal(ctx context.Context, nrCnpj string, dtAno string) ([]domain.VariacaoMensalResponse, error) {
	return api.GetVariacaoMensal(uc.db, nrCnpj, dtAno)
}

// GetResumoDashboard retorna um resumo completo do dashboard
func (uc *DashboardUseCase) GetResumoDashboard(ctx context.Context, nrCnpj string, dtAno string, dtMesBimestre string) (*domain.ResumoDashboardResponse, error) {
	return api.GetResumoDashboard(uc.db, nrCnpj, dtAno, dtMesBimestre)
}
