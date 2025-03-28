// internal/web/server.go

package web

import (
	"context"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/services"
	"crypto-trading-bot/internal/utils"
	"crypto-trading-bot/internal/web/handlers"

	//"crypto-trading-bot/internal/trading"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	port   string
	router *mux.Router
	repo   *repositories.Repository
	logger *utils.Logger
	//trader *trading.Trader
	server            *http.Server
	strategyService   services.StrategyService
	marketDataService services.MarketDataService
	exchangeService   services.ExchangeService
}

func NewServer(port string, repo *repositories.Repository, logger *utils.Logger, exchangeService services.ExchangeService, strategyService services.StrategyService, marketDataService services.MarketDataService) *Server {
	router := mux.NewRouter()

	s := &Server{
		port:   port,
		router: router,
		repo:   repo,
		//trader: trader,
		logger:            logger,
		strategyService:   strategyService,
		marketDataService: marketDataService,
		exchangeService:   exchangeService,
	}

	// Настройка маршрутов для веб-интерфейса
	fs := http.FileServer(http.Dir("./web/ui/static"))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	strategyHandler := handlers.NewStrategyHandler(strategyService, logger)
	marketDataHandler := handlers.NewMarketDataHandler(marketDataService, exchangeService, logger)

	// // Настройка маршрутов для API

	//s.router.Handle("/", templ.Handler(hello("Костя"))).Methods("GET")

	s.router.HandleFunc("/", strategyHandler.GetStrategiesListPage).Methods("GET")
	s.router.HandleFunc("/strategies/new", strategyHandler.GetNewStrategyPage).Methods("GET")
	s.router.HandleFunc("/strategies/{id}/edit", strategyHandler.GetEditStrategyPage).Methods("GET")
	s.router.HandleFunc("/strategies", strategyHandler.PostCreateStrategy).Methods("POST")
	s.router.HandleFunc("/strategies/{id}", strategyHandler.PostUpdateStrategy).Methods("POST")
	s.router.HandleFunc("/strategies/{id}/delete", strategyHandler.PostDeleteStrategy).Methods("POST")

	s.router.HandleFunc("/marketdata", marketDataHandler.GetMarketData).Methods("GET")

	return s
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
