package server

import (
	"log"
	"net/http"
	"product-service/src/api"
	"product-service/src/config"
	"product-service/src/service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	cfg     *config.Config
	service service.ProductService
}

func NewServer(cfg *config.Config, productService service.ProductService) *Server {
	return &Server{
		cfg:     cfg,
		service: productService,
	}
}

func (s *Server) Run() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	apiHandler := api.NewHandler(s.service, s.cfg)

	// --- Configuração das Rotas ---
	// Rotas Públicas
	router.Post("/create", apiHandler.HandleCreate)
	router.Get("/{id}", apiHandler.HandleGet)
	router.Get("/list", apiHandler.HandleList)
	router.Put("/reduce-stock/{id}", apiHandler.HandleReduceStock)
	router.Put("/{id}", apiHandler.HandleUpdate)
	router.Delete("/{id}", apiHandler.HandleDelete)

	log.Printf("Servidor de Produtos iniciado em %s", s.cfg.ListenAddr)
	if err := http.ListenAndServe(s.cfg.ListenAddr, router); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
