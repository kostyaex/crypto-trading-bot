// internal/web/server.go

package web

import (
	"context"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/service/exchange"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/trading/strategy"
	"crypto-trading-bot/internal/web/handlers"

	//"crypto-trading-bot/internal/trading"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	port   string
	router *mux.Router
	// repo   *repositories.Repository
	logger *logger.Logger
	//trader *trading.Trader
	server *http.Server
	// strategyService   services.StrategyService
	// marketDataService services.MarketDataService
	// exchangeService   services.ExchangeService
	resourcesHandler  *handlers.ResourcesHandler
	strategyHandler   *handlers.StrategyHandler
	marketDataHandler *handlers.MarketDataHandler
}

func NewServer(port string, repo *repositories.Repository, logger *logger.Logger, exchangeService exchange.ExchangeService, strategyService strategy.StrategyService, marketDataService marketdata.MarketDataService) *Server {
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

// startMetricsServer запускает HTTP-сервер для метрик Prometheus и ждёт ctx.Done()
func StartMetricsServer(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Metrics server started on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start metrics server: %v", err)
		}
	}()

	go func() {
		// Ждём сигнал о завершении
		<-ctx.Done()
		server.Shutdown(ctx)
	}()

	return nil
}
