package api

import (
	"database/sql"
	"fmt"

	"github.com/Filipefr15/cadprev_apis/internal/domain"
)

// GetPatrimonioMensal retorna o patrimônio total por mês de um ente
func GetPatrimonioMensal(db *sql.DB, nrCnpj string, dtAno string) ([]domain.PatrimonioMensalResponse, error) {
	query := `
	SELECT 
		dt_ano,
		dt_mes_bimestre,
		no_ente,
		COALESCE(SUM(CAST(vl_total_atual AS REAL)), 0) as patrimonio_total
	FROM dair_carteira
	WHERE nr_cnpj_entidade = ? AND dt_ano = ?
	GROUP BY dt_ano, dt_mes_bimestre, no_ente
	ORDER BY dt_ano, dt_mes_bimestre
	`

	rows, err := db.Query(query, nrCnpj, dtAno)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query patrimonio mensal: %w", err)
	}
	defer rows.Close()

	var results []domain.PatrimonioMensalResponse
	for rows.Next() {
		var record domain.PatrimonioMensalResponse
		if err := rows.Scan(&record.DtAno, &record.DtMesBimestre, &record.NoEnte, &record.PatrimonioTotal); err != nil {
			return nil, fmt.Errorf("erro ao fazer scan do resultado: %w", err)
		}
		results = append(results, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar resultados: %w", err)
	}

	return results, nil
}

// GetComposicaoSegmento retorna a composição de investimentos por segmento
func GetComposicaoSegmento(db *sql.DB, nrCnpj string, dtAno string, dtMesBimestre string) ([]domain.ComposicaoSegmentoResponse, error) {
	query := `
	SELECT 
		dt_ano,
		dt_mes_bimestre,
		no_segmento,
		COALESCE(SUM(CAST(vl_total_atual AS REAL)), 0) as valor_segmento
	FROM dair_carteira
	WHERE nr_cnpj_entidade = ? AND dt_ano = ? AND dt_mes_bimestre = ?
	GROUP BY dt_ano, dt_mes_bimestre, no_segmento
	ORDER BY valor_segmento DESC
	`

	rows, err := db.Query(query, nrCnpj, dtAno, dtMesBimestre)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query composição segmento: %w", err)
	}
	defer rows.Close()

	var results []domain.ComposicaoSegmentoResponse
	var totalValor float64

	// Primeira passagem para calcular o total
	tempResults := make([]domain.ComposicaoSegmentoResponse, 0)
	for rows.Next() {
		var record domain.ComposicaoSegmentoResponse
		if err := rows.Scan(&record.DtAno, &record.DtMesBimestre, &record.NoSegmento, &record.ValorSegmento); err != nil {
			return nil, fmt.Errorf("erro ao fazer scan do resultado: %w", err)
		}
		tempResults = append(tempResults, record)
		totalValor += record.ValorSegmento
	}

	// Segunda passagem para calcular percentuais
	for _, record := range tempResults {
		if totalValor > 0 {
			record.Percentual = (record.ValorSegmento / totalValor) * 100
		}
		results = append(results, record)
	}

	return results, nil
}

// GetEvolucaoAnual retorna a evolução do patrimônio anual
func GetEvolucaoAnual(db *sql.DB, nrCnpj string) ([]domain.EvolucaoAnualResponse, error) {
	query := `
	SELECT 
		dt_ano,
		no_ente,
		COALESCE(SUM(CAST(vl_total_atual AS REAL)), 0) as patrimonio_total
	FROM dair_carteira
	WHERE nr_cnpj_entidade = ?
	GROUP BY dt_ano, no_ente
	ORDER BY dt_ano
	`

	rows, err := db.Query(query, nrCnpj)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query evolução anual: %w", err)
	}
	defer rows.Close()

	var results []domain.EvolucaoAnualResponse
	for rows.Next() {
		var record domain.EvolucaoAnualResponse
		if err := rows.Scan(&record.DtAno, &record.NoEnte, &record.PatrimonioTotal); err != nil {
			return nil, fmt.Errorf("erro ao fazer scan do resultado: %w", err)
		}
		results = append(results, record)
	}

	return results, nil
}

// GetVariacaoMensal retorna a variação percentual mensal de patrimônio (para rentabilidade)
func GetVariacaoMensal(db *sql.DB, nrCnpj string, dtAno string) ([]domain.VariacaoMensalResponse, error) {
	// Obter patrimônios mensais
	patrimonios, err := GetPatrimonioMensal(db, nrCnpj, dtAno)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter patrimônios mensais: %w", err)
	}

	var results []domain.VariacaoMensalResponse

	// Calcular variações
	for i := 0; i < len(patrimonios); i++ {
		record := domain.VariacaoMensalResponse{
			DtAno:           patrimonios[i].DtAno,
			DtMesBimestre:   patrimonios[i].DtMesBimestre,
			NoEnte:          patrimonios[i].NoEnte,
			PatrimonioAtual: patrimonios[i].PatrimonioTotal,
		}

		// Se houver período anterior, calcular variação
		if i > 0 && patrimonios[i].NoEnte == patrimonios[i-1].NoEnte {
			record.PatrimonioAnterior = patrimonios[i-1].PatrimonioTotal
			record.VariacaoAbsoluta = record.PatrimonioAtual - record.PatrimonioAnterior

			if record.PatrimonioAnterior != 0 {
				record.VariacaoPercentual = (record.VariacaoAbsoluta / record.PatrimonioAnterior) * 100
			}
		}

		results = append(results, record)
	}

	return results, nil
}

// GetResumoDashboard retorna todos os dados do dashboard agregados
func GetResumoDashboard(db *sql.DB, nrCnpj string, dtAno string, dtMesBimestre string) (*domain.ResumoDashboardResponse, error) {
	resumo := &domain.ResumoDashboardResponse{
		PatrimonioMensal:   make([]domain.PatrimonioMensalResponse, 0),
		ComposicaoSegmento: make([]domain.ComposicaoSegmentoResponse, 0),
		EvolucaoAnual:      make([]domain.EvolucaoAnualResponse, 0),
		VariacaoMensal:     make([]domain.VariacaoMensalResponse, 0),
	}

	// Patrimônio mensal
	patrimonio, err := GetPatrimonioMensal(db, nrCnpj, dtAno)
	if err == nil {
		resumo.PatrimonioMensal = patrimonio
	}

	// Composição por segmento (usa o mês mais recente se não especificado)
	if dtMesBimestre == "" {
		// Pegar o mês mais recente
		var mesMaisRecente string
		queryMes := `SELECT MAX(dt_mes_bimestre) FROM dair_carteira WHERE nr_cnpj_entidade = ? AND dt_ano = ?`
		err := db.QueryRow(queryMes, nrCnpj, dtAno).Scan(&mesMaisRecente)
		if err == nil && mesMaisRecente != "" {
			dtMesBimestre = mesMaisRecente
		}
	}

	if dtMesBimestre != "" {
		composicao, err := GetComposicaoSegmento(db, nrCnpj, dtAno, dtMesBimestre)
		if err == nil {
			resumo.ComposicaoSegmento = composicao
		}
	}

	// Evolução anual
	evolucao, err := GetEvolucaoAnual(db, nrCnpj)
	if err == nil {
		resumo.EvolucaoAnual = evolucao
	}

	// Variação mensal
	variacao, err := GetVariacaoMensal(db, nrCnpj, dtAno)
	if err == nil {
		resumo.VariacaoMensal = variacao
	}

	return resumo, nil
}
