package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Filipefr15/cadprev_apis/internal/usecase"
)

type DashboardHandler struct {
	usecase *usecase.DashboardUseCase
}

func NewDashboardHandler(usecase *usecase.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{
		usecase: usecase,
	}
}

// GetPatrimonioMensalHandler handles GET /dashboard/patrimonio-mensal
// @Summary Get monthly patrimony
// @Description Retrieve total patrimony values grouped by month for a specific entity
// @Tags dashboard
// @Accept json
// @Produce json
// @Param nr_cnpj query string true "CNPJ"
// @Param dt_ano query string true "Year"
// @Success 200 {array} domain.PatrimonioMensalResponse
// @Failure 400 {string} string "Missing required parameters"
// @Failure 500 {string} string "Database error"
// @Router /dashboard/patrimonio-mensal [get]
func (h *DashboardHandler) GetPatrimonioMensalHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recebida requisição para /dashboard/patrimonio-mensal")

	nrCnpj := r.URL.Query().Get("nr_cnpj")
	dtAno := r.URL.Query().Get("dt_ano")

	if nrCnpj == "" || dtAno == "" {
		log.Println("Parâmetros obrigatórios ausentes: nr_cnpj e dt_ano")
		http.Error(w, "Parâmetros obrigatórios ausentes: nr_cnpj e dt_ano", http.StatusBadRequest)
		return
	}

	log.Printf("Consultando patrimônio mensal para CNPJ=%s, Ano=%s", nrCnpj, dtAno)
	results, err := h.usecase.GetPatrimonioMensal(context.Background(), nrCnpj, dtAno)
	if err != nil {
		log.Printf("Erro ao consultar patrimônio mensal: %v", err)
		http.Error(w, "Erro ao consultar dados do banco de dados", http.StatusInternalServerError)
		return
	}

	log.Printf("Consulta bem-sucedida. Número de resultados: %d", len(results))
	response, err := json.Marshal(results)
	if err != nil {
		log.Printf("Erro ao processar os dados: %v", err)
		http.Error(w, "Erro ao processar os dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Println("Resposta enviada com sucesso.")
}

// GetComposicaoSegmentoHandler handles GET /dashboard/composicao-segmento
// @Summary Get composition by segment
// @Description Retrieve asset composition grouped by segment for a specific month
// @Tags dashboard
// @Accept json
// @Produce json
// @Param nr_cnpj query string true "CNPJ"
// @Param dt_ano query string true "Year"
// @Param dt_mes_bimestre query string false "Month/Bimonthly"
// @Success 200 {array} domain.ComposicaoSegmentoResponse
// @Failure 400 {string} string "Missing required parameters"
// @Failure 500 {string} string "Database error"
// @Router /dashboard/composicao-segmento [get]
func (h *DashboardHandler) GetComposicaoSegmentoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recebida requisição para /dashboard/composicao-segmento")

	nrCnpj := r.URL.Query().Get("nr_cnpj")
	dtAno := r.URL.Query().Get("dt_ano")
	dtMesBimestre := r.URL.Query().Get("dt_mes_bimestre")

	if nrCnpj == "" || dtAno == "" {
		log.Println("Parâmetros obrigatórios ausentes: nr_cnpj e dt_ano")
		http.Error(w, "Parâmetros obrigatórios ausentes: nr_cnpj e dt_ano", http.StatusBadRequest)
		return
	}

	log.Printf("Consultando composição por segmento para CNPJ=%s, Ano=%s, Mês=%s", nrCnpj, dtAno, dtMesBimestre)
	results, err := h.usecase.GetComposicaoSegmento(context.Background(), nrCnpj, dtAno, dtMesBimestre)
	if err != nil {
		log.Printf("Erro ao consultar composição por segmento: %v", err)
		http.Error(w, "Erro ao consultar dados do banco de dados", http.StatusInternalServerError)
		return
	}

	log.Printf("Consulta bem-sucedida. Número de segmentos: %d", len(results))
	response, err := json.Marshal(results)
	if err != nil {
		log.Printf("Erro ao processar os dados: %v", err)
		http.Error(w, "Erro ao processar os dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Println("Resposta enviada com sucesso.")
}

// GetEvolucaoAnualHandler handles GET /dashboard/evolucao-anual
// @Summary Get annual evolution
// @Description Retrieve annual patrimony evolution for a specific entity
// @Tags dashboard
// @Accept json
// @Produce json
// @Param nr_cnpj query string true "CNPJ"
// @Success 200 {array} domain.EvolucaoAnualResponse
// @Failure 400 {string} string "Missing required parameters"
// @Failure 500 {string} string "Database error"
// @Router /dashboard/evolucao-anual [get]
func (h *DashboardHandler) GetEvolucaoAnualHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recebida requisição para /dashboard/evolucao-anual")

	nrCnpj := r.URL.Query().Get("nr_cnpj")

	if nrCnpj == "" {
		log.Println("Parâmetro obrigatório ausente: nr_cnpj")
		http.Error(w, "Parâmetro obrigatório ausente: nr_cnpj", http.StatusBadRequest)
		return
	}

	log.Printf("Consultando evolução anual para CNPJ=%s", nrCnpj)
	results, err := h.usecase.GetEvolucaoAnual(context.Background(), nrCnpj)
	if err != nil {
		log.Printf("Erro ao consultar evolução anual: %v", err)
		http.Error(w, "Erro ao consultar dados do banco de dados", http.StatusInternalServerError)
		return
	}

	log.Printf("Consulta bem-sucedida. Número de anos: %d", len(results))
	response, err := json.Marshal(results)
	if err != nil {
		log.Printf("Erro ao processar os dados: %v", err)
		http.Error(w, "Erro ao processar os dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Println("Resposta enviada com sucesso.")
}

// GetVariacaoMensalHandler handles GET /dashboard/variacao-mensal
// @Summary Get monthly variation (profitability)
// @Description Retrieve monthly variation percentages for profitability analysis
// @Tags dashboard
// @Accept json
// @Produce json
// @Param nr_cnpj query string true "CNPJ"
// @Param dt_ano query string true "Year"
// @Success 200 {array} domain.VariacaoMensalResponse
// @Failure 400 {string} string "Missing required parameters"
// @Failure 500 {string} string "Database error"
// @Router /dashboard/variacao-mensal [get]
func (h *DashboardHandler) GetVariacaoMensalHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recebida requisição para /dashboard/variacao-mensal")

	nrCnpj := r.URL.Query().Get("nr_cnpj")
	dtAno := r.URL.Query().Get("dt_ano")

	if nrCnpj == "" || dtAno == "" {
		log.Println("Parâmetros obrigatórios ausentes: nr_cnpj e dt_ano")
		http.Error(w, "Parâmetros obrigatórios ausentes: nr_cnpj e dt_ano", http.StatusBadRequest)
		return
	}

	log.Printf("Consultando variação mensal para CNPJ=%s, Ano=%s", nrCnpj, dtAno)
	results, err := h.usecase.GetVariacaoMensal(context.Background(), nrCnpj, dtAno)
	if err != nil {
		log.Printf("Erro ao consultar variação mensal: %v", err)
		http.Error(w, "Erro ao consultar dados do banco de dados", http.StatusInternalServerError)
		return
	}

	log.Printf("Consulta bem-sucedida. Número de períodos: %d", len(results))
	response, err := json.Marshal(results)
	if err != nil {
		log.Printf("Erro ao processar os dados: %v", err)
		http.Error(w, "Erro ao processar os dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Println("Resposta enviada com sucesso.")
}

// GetResumoDashboardHandler handles GET /dashboard/resumo
// @Summary Get complete dashboard summary
// @Description Retrieve all dashboard data aggregated (patrimony, composition, evolution, variation)
// @Tags dashboard
// @Accept json
// @Produce json
// @Param nr_cnpj query string true "CNPJ"
// @Param dt_ano query string true "Year"
// @Param dt_mes_bimestre query string false "Month/Bimonthly"
// @Success 200 {object} domain.ResumoDashboardResponse
// @Failure 400 {string} string "Missing required parameters"
// @Failure 500 {string} string "Database error"
// @Router /dashboard/resumo [get]
func (h *DashboardHandler) GetResumoDashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recebida requisição para /dashboard/resumo")

	nrCnpj := r.URL.Query().Get("nr_cnpj")
	dtAno := r.URL.Query().Get("dt_ano")
	dtMesBimestre := r.URL.Query().Get("dt_mes_bimestre")

	if nrCnpj == "" || dtAno == "" {
		log.Println("Parâmetros obrigatórios ausentes: nr_cnpj e dt_ano")
		http.Error(w, "Parâmetros obrigatórios ausentes: nr_cnpj e dt_ano", http.StatusBadRequest)
		return
	}

	log.Printf("Consultando resumo dashboard para CNPJ=%s, Ano=%s, Mês=%s", nrCnpj, dtAno, dtMesBimestre)
	result, err := h.usecase.GetResumoDashboard(context.Background(), nrCnpj, dtAno, dtMesBimestre)
	if err != nil {
		log.Printf("Erro ao consultar resumo dashboard: %v", err)
		http.Error(w, "Erro ao consultar dados do banco de dados", http.StatusInternalServerError)
		return
	}

	log.Println("Consulta bem-sucedida.")
	response, err := json.Marshal(result)
	if err != nil {
		log.Printf("Erro ao processar os dados: %v", err)
		http.Error(w, "Erro ao processar os dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	log.Println("Resposta enviada com sucesso.")
}
