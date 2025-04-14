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
	// repo   *repositories.Repository
	logger *utils.Logger
	//trader *trading.Trader
	server *http.Server
	// strategyService   services.StrategyService
	// marketDataService services.MarketDataService
	// exchangeService   services.ExchangeService
	resourcesHandler  *handlers.ResourcesHandler
	strategyHandler   *handlers.StrategyHandler
	marketDataHandler *handlers.MarketDataHandler
}

func NewServer(port string, repo *repositories.Repository, logger *utils.Logger, exchangeService services.ExchangeService, strategyService services.StrategyService, marketDataService services.MarketDataService) *Server {
	router := mux.NewRouter()

	resourcesHandler := handlers.NewResourcesHandler(logger)
	strategyHandler := handlers.NewStrategyHandler(strategyService, logger)
	marketDataHandler := handlers.NewMarketDataHandler(marketDataService, exchangeService, logger)

	s := &Server{
		port:   port,
		router: router,
		// repo:   repo,
		//trader: trader,
		logger: logger,
		// strategyService:   strategyService,
		// marketDataService: marketDataService,
		// exchangeService:   exchangeService,
		resourcesHandler:  resourcesHandler,
		strategyHandler:   strategyHandler,
		marketDataHandler: marketDataHandler,
	}

	s.routes()

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
