package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Filipefr15/cadprev_apis/internal/usecase"
)

// DataHandler handles HTTP requests for data ingestion and retrieval
type DataHandler struct {
	ingestUseCase *usecase.IngestUseCase
}

// NewDataHandler creates a new HTTP handler with dependency injection
func NewDataHandler(ingestUseCase *usecase.IngestUseCase) *DataHandler {
	return &DataHandler{
		ingestUseCase: ingestUseCase,
	}
}

// IngestAll handles POST /api/ingest - ingests all data from external API
func (h *DataHandler) IngestAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	err := h.ingestUseCase.IngestAll(ctx)
	if err != nil {
		log.Printf("Error ingesting data: %v", err)
		http.Error(w, "Failed to ingest data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Data ingestion completed successfully",
	}); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// IngestByID handles POST /api/ingest/{id} - ingests specific data by ID
func (h *DataHandler) IngestByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/ingest/")
	id := strings.TrimSpace(path)

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	data, err := h.ingestUseCase.IngestByID(ctx, id)
	if err != nil {
		log.Printf("Error ingesting data by ID: %v", err)
		http.Error(w, "Failed to ingest data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// GetAll handles GET /api/data - retrieves all data
func (h *DataHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	dataList, err := h.ingestUseCase.GetAll(ctx)
	if err != nil {
		log.Printf("Error retrieving data: %v", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dataList); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// GetByID handles GET /api/data/{id} - retrieves specific data by ID
func (h *DataHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/data/")
	id := strings.TrimSpace(path)

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	data, err := h.ingestUseCase.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error retrieving data by ID: %v", err)
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// Health handles GET /health - health check endpoint
func (h *DataHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	}); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// RegisterRoutes sets up all HTTP routes
func (h *DataHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/api/ingest", h.IngestAll)
	mux.HandleFunc("/api/ingest/", h.IngestByID)
	mux.HandleFunc("/api/data", h.GetAll)
	mux.HandleFunc("/api/data/", h.GetByID)
}
