package main

import (
	"context"
	"log"
	"os"

	"github.com/Filipefr15/cadprev_apis/internal/api"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env: ", err)
	}

	baseURL := getEnv("API_BASE_URL", "https://apicadprev.trabalho.gov.br")
	dbPath := getEnv("SQLITE_DB_PATH", "./cadprev.db")

	db, err := api.OpenSQLite(dbPath)
	if err != nil {
		log.Fatal("Erro ao abrir SQLite: ", err)
	}
	defer db.Close()

	if err := api.CreateDairCarteiraTable(db); err != nil {
		log.Fatal("Erro ao criar tabela: ", err)
	}

	// Parâmetros dinâmicos (exemplo)
	params := map[string]string{
		"sg_uf":   "MT",
		"no_ente": "Sinop",
		// Adicione outros conforme necessário
	}

	ctx := context.Background()
	items, err := api.FetchDairCarteira(ctx, baseURL, params)
	if err != nil {
		log.Fatal("Erro ao buscar API: ", err)
	}

	for _, item := range items {
		if err := api.InsertDairCarteira(db, item); err != nil {
			log.Println("Erro ao inserir: ", err)
		}
	}

	log.Println("Finalizado!")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
