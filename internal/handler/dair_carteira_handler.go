package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Filipefr15/cadprev_apis/internal/usecase"
)

type DairCarteiraHandler struct {
	usecase *usecase.DairCarteiraUseCase
}

func NewDairCarteiraHandler(usecase *usecase.DairCarteiraUseCase) *DairCarteiraHandler {
	return &DairCarteiraHandler{
		usecase: usecase,
	}
}

// GetDairCarteiraHandler handles GET /dair_carteira
// @Summary Retrieve Dair Carteira data
// @Description Fetch data from the database with optional filters
// @Tags dair_carteira
// @Accept json
// @Produce json
// @Param dt_ano query string true "Year"
// @Param nr_cnpj query string false "CNPJ"
// @Param no_ente query string false "Entity Name"
// @Param mes query string false "Month"
// @Param segmento query string false "Segment"
// @Param sg_uf query string false "State"
// @Success 200 {array} entity.Dair_carteira
// @Failure 400 {string} string "Missing required parameters"
// @Failure 500 {string} string "Database error"
// @Router /dair_carteira [get]
func (h *DairCarteiraHandler) GetDairCarteiraHandler(w http.ResponseWriter, r *http.Request) {
	dtAno := r.URL.Query().Get("dt_ano")
	nrCnpj := r.URL.Query().Get("nr_cnpj")
	noEnte := r.URL.Query().Get("no_ente")
	mes := r.URL.Query().Get("mes")
	segmento := r.URL.Query().Get("segmento")
	sgUf := r.URL.Query().Get("sg_uf")

	if dtAno == "" || (nrCnpj == "" && noEnte == "") {
		http.Error(w, "Parâmetros obrigatórios ausentes: dt_ano e nr_cnpj ou no_ente", http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"dt_ano":           dtAno,
		"nr_cnpj_entidade": nrCnpj,
		"no_ente":          noEnte,
		"dt_mes_bimestre":  mes,
		"no_segmento":      segmento,
		"sg_uf":            sgUf,
	}

	results, err := h.usecase.GetDairCarteira(context.Background(), params)
	if err != nil {
		http.Error(w, "Erro ao consultar o banco de dados", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Erro ao processar os dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// BuscarDairCarteiraHandler handles POST /buscar_dair_carteira
// @Summary Trigger data fetching from external API
// @Description Initiates a process to fetch and store data from an external API
// @Tags dair_carteira
// @Accept json
// @Produce plain
// @Param sg_uf query string false "State"
// @Param no_ente query string false "Entity Name"
// @Success 200 {string} string "Fetch initiated successfully"
// @Failure 500 {string} string "API fetch error"
// @Router /buscar_dair_carteira [post]
func (h *DairCarteiraHandler) BuscarDairCarteiraHandler(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{
		"sg_uf":   r.URL.Query().Get("sg_uf"),
		"no_ente": r.URL.Query().Get("no_ente"),
	}

	err := h.usecase.BuscarDairCarteira(context.Background(), params)
	if err != nil {
		http.Error(w, "Erro ao buscar dados da API", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Busca iniciada com sucesso!"))
}
