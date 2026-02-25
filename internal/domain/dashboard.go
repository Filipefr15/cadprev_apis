package domain

// PatrimonioMensalResponse representa o patrimônio total de um mês
type PatrimonioMensalResponse struct {
	DtAno           string  `json:"dt_ano"`
	DtMesBimestre   string  `json:"dt_mes_bimestre"`
	NoEnte          string  `json:"no_ente"`
	PatrimonioTotal float64 `json:"patrimonio_total"`
}

// ComposicaoSegmentoResponse representa a composição de investimentos por segmento
type ComposicaoSegmentoResponse struct {
	DtAno         string  `json:"dt_ano"`
	DtMesBimestre string  `json:"dt_mes_bimestre"`
	NoSegmento    string  `json:"no_segmento"`
	ValorSegmento float64 `json:"valor_segmento"`
	Percentual    float64 `json:"percentual"`
}

// EvolucaoAnualResponse representa a evolução do patrimônio anual
type EvolucaoAnualResponse struct {
	DtAno           string  `json:"dt_ano"`
	NoEnte          string  `json:"no_ente"`
	PatrimonioTotal float64 `json:"patrimonio_total"`
}

// VariacaoMensalResponse representa a variação mensal de patrimônio (para rentabilidade)
type VariacaoMensalResponse struct {
	DtAno              string  `json:"dt_ano"`
	DtMesBimestre      string  `json:"dt_mes_bimestre"`
	NoEnte             string  `json:"no_ente"`
	PatrimonioAtual    float64 `json:"patrimonio_atual"`
	PatrimonioAnterior float64 `json:"patrimonio_anterior"`
	VariacaoAbsoluta   float64 `json:"variacao_absoluta"`
	VariacaoPercentual float64 `json:"variacao_percentual"`
}

// ResumoDashboardResponse agrupa todos os dados do dashboard
type ResumoDashboardResponse struct {
	PatrimonioMensal   []PatrimonioMensalResponse   `json:"patrimonio_mensal"`
	ComposicaoSegmento []ComposicaoSegmentoResponse `json:"composicao_segmento"`
	EvolucaoAnual      []EvolucaoAnualResponse      `json:"evolucao_anual"`
	VariacaoMensal     []VariacaoMensalResponse     `json:"variacao_mensal"`
}
