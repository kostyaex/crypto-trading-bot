package web

import (
	"crypto-trading-bot/internal/web/ui"
	"net/http"

	"github.com/a-h/templ"
)

func (s *Server) routes() {
	// Настройка маршрутов для веб-интерфейса
	fs := http.FileServer(http.Dir("./internal/web/ui/assets"))
	s.router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	fs2 := http.FileServer(http.Dir("./data"))
	s.router.PathPrefix("/data/").Handler(http.StripPrefix("/data/", fs2))

	//
	// Определяем маршрут для проксирования запросов к PostgREST
	//

	// маршрут для получения метаданных
	//router.HandleFunc("/resources", GetHandler).Methods("GET")

	s.router.HandleFunc("/api/resources/{resource}", GetHandler).Methods("GET")
	s.router.HandleFunc("/api/resources/{resource}", PostHandler).Methods("POST")
	s.router.HandleFunc("/api/resources/{resource}", DeleteHandler).Methods("DELETE")
	s.router.HandleFunc("/api/resources/{resource}", PatchHandler).Methods("PATCH")

	///

	// // Настройка маршрутов для API

	s.router.Handle("/", templ.Handler(ui.IndexComponent())).Methods("GET")

	s.router.HandleFunc("/resources/{resource}", s.resourcesHandler.GetResourcesListPage).Methods("GET")

	// s.router.HandleFunc("/strategies", s.strategyHandler.GetStrategiesListPage).Methods("GET")
	// s.router.HandleFunc("/strategies", s.strategyHandler.GetStrategiesListPage).Methods("GET")
	// s.router.HandleFunc("/strategies/new", s.strategyHandler.GetNewStrategyPage).Methods("GET")
	// s.router.HandleFunc("/strategies/{id}/edit", s.strategyHandler.GetEditStrategyPage).Methods("GET")
	// s.router.HandleFunc("/strategies", s.strategyHandler.PostCreateStrategy).Methods("POST")
	// s.router.HandleFunc("/strategies/{id}", s.strategyHandler.PostUpdateStrategy).Methods("POST")
	// s.router.HandleFunc("/strategies/{id}/delete", s.strategyHandler.PostDeleteStrategy).Methods("POST")

	s.router.HandleFunc("/backtesting", s.traderHandler.GetBacktestingPage).Methods("GET")
	s.router.HandleFunc("/api/runbacktesting", s.traderHandler.PostRunBacktesting).Methods("POST")
	s.router.HandleFunc("/api/seriesdumpslist", s.traderHandler.GetSeriesDumpsList).Methods("GET")

	//s.router.HandleFunc("/marketdata", s.marketDataHandler.GetMarketData).Methods("GET")
}
