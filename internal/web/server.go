// internal/web/server.go

package web

import (
	"context"
	"crypto-trading-bot/internal/data"

	//"crypto-trading-bot/internal/trading"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	port   string
	router *mux.Router
	repo   *data.PostgresRepository
	//trader *trading.Trader
	server *http.Server
}

func NewServer(port string, repo *data.PostgresRepository) *Server {
	router := mux.NewRouter()
	s := &Server{
		port:   port,
		router: router,
		repo:   repo,
		//trader: trader,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Настройка маршрутов для API
	s.router.HandleFunc("/api/strategies", s.createStrategyHandler).Methods("POST")
	s.router.HandleFunc("/api/strategies/{id}", s.updateStrategyHandler).Methods("PUT")
	s.router.HandleFunc("/api/strategies/{id}", s.deleteStrategyHandler).Methods("DELETE")
	s.router.HandleFunc("/api/strategies/{id}", s.getStrategyHandler).Methods("GET")

	// Настройка маршрутов для веб-интерфейса
	fs := http.FileServer(http.Dir("./web/ui/static"))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	s.router.HandleFunc("/", s.indexHandler).Methods("GET")
	// Добавьте другие маршруты по необходимости
}

func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	go func() {
		<-ctx.Done()
		s.server.Shutdown(ctx)
	}()

	log.Printf("Starting web server on port %s", s.port)
	return s.server.ListenAndServe()
}

// Обработчики маршрутов
func (s *Server) createStrategyHandler(w http.ResponseWriter, r *http.Request) {
	// Логика создания новой стратегии
}

func (s *Server) updateStrategyHandler(w http.ResponseWriter, r *http.Request) {
	// Логика обновления существующей стратегии
}

func (s *Server) deleteStrategyHandler(w http.ResponseWriter, r *http.Request) {
	// Логика удаления стратегии
}

func (s *Server) getStrategyHandler(w http.ResponseWriter, r *http.Request) {
	// Логика получения информации о стратегии
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Логика обработки главной страницы веб-интерфейса
}
