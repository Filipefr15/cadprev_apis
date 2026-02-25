package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Filipefr15/cadprev_apis/docs"
	"github.com/Filipefr15/cadprev_apis/internal/api"
	"github.com/Filipefr15/cadprev_apis/internal/handler"
	"github.com/Filipefr15/cadprev_apis/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Carregar .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env: ", err)
	}

	dbPath := getEnv("SQLITE_DB_PATH", "./cadprev.db")

	db, err := api.OpenSQLite(dbPath)
	if err != nil {
		log.Fatal("Erro ao abrir SQLite: ", err)
	}
	defer db.Close()

	if err := api.CreateDairCarteiraTable(db); err != nil {
		log.Fatal("Erro ao criar tabela: ", err)
	}

	// Inicializar Swagger
	docs.SwaggerInfo.Title = "API Documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	startServer(db)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// HTTP server and router setup
func setupRouter(db *sql.DB) *chi.Mux {
	router := chi.NewRouter()

	// Middleware para habilitar CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Substitua por domínios específicos, se necessário
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Cache pré-voo por 5 minutos
	}))

	// Rotas para Dair Carteira
	dairHandler := handler.NewDairCarteiraHandler(usecase.NewDairCarteiraUseCase(db))
	router.Get("/dair_carteira", dairHandler.GetDairCarteiraHandler)
	router.Post("/buscar_dair_carteira", dairHandler.BuscarDairCarteiraHandler)

	// Rotas para Dashboard
	dashboardHandler := handler.NewDashboardHandler(usecase.NewDashboardUseCase(db))
	router.Get("/dashboard/patrimonio-mensal", dashboardHandler.GetPatrimonioMensalHandler)
	router.Get("/dashboard/composicao-segmento", dashboardHandler.GetComposicaoSegmentoHandler)
	router.Get("/dashboard/evolucao-anual", dashboardHandler.GetEvolucaoAnualHandler)
	router.Get("/dashboard/variacao-mensal", dashboardHandler.GetVariacaoMensalHandler)
	router.Get("/dashboard/resumo", dashboardHandler.GetResumoDashboardHandler)

	// Swagger
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	return router
}

// Função principal para iniciar o servidor HTTP
func startServer(db *sql.DB) {
	router := setupRouter(db)

	port := getEnv("PORT", "8080")
	log.Println("Iniciando servidor na porta:", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Erro ao iniciar servidor: ", err)
	}
}
