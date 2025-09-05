package handlers

import (
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/utils"
	"crypto-trading-bot/internal/web/ui"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TraderHandler struct {
	// marketDataService marketdata.MarketDataService
	// exchangeService   exchange.ExchangeService
	conf   *config.Config
	logger *logger.Logger
}

func NewTraderHandler(conf *config.Config, logger *logger.Logger) *TraderHandler {
	return &TraderHandler{
		conf:   conf,
		logger: logger,
	}
}

func (h *TraderHandler) GetBacktestingPage(w http.ResponseWriter, r *http.Request) {

	if err := ui.BacktestingPage().Render(r.Context(), w); err != nil {
		h.logger.Errorf("Ошибка формирования страницы бектестинга: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

func (h *TraderHandler) PostRunBacktesting(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Start string `json:"start"`
		Stop  string `json:"stop"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	layout := "2006-01-02T15:04" // Формат даты и времени
	start, err := time.Parse(layout, requestData.Start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stop, err := time.Parse(layout, requestData.Stop)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Start:", start)
	fmt.Println("Stop:", stop)

	//results := h.marketDataService.RunBacktesting(start, stop)
	results := "{}"

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
}

//curl -X POST -H "Content-Type: application/json" -d '{"start": "2023-01-01T12:00:00", "stop": "2023-01-01T13:00:00"}' http://localhost:5000/api/runbacktesting

// Получить список URL сохраненных серий в каталоге data. Возвращает в виде JSON.
func (h *TraderHandler) GetSeriesDumpsList(w http.ResponseWriter, r *http.Request) {

	names, err := utils.FileList(h.conf.Data.Dir+"/series", "/series/")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(names); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Получить список URL сохраненных серий в каталоге data. Возвращает в виде JSON.
func (h *TraderHandler) GetBacktestsDumpsList(w http.ResponseWriter, r *http.Request) {

	names, err := utils.FileList(h.conf.Data.Dir+"/backtests", "/backtests/")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Предупреждение кэширования:
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if err := json.NewEncoder(w).Encode(names); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
