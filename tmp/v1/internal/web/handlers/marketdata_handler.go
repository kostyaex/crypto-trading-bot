package handlers

import (
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/service/exchange"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/web/ui"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MarketDataHandler обрабатывает запросы, связанные с рыночными данными.
type MarketDataHandler struct {
	marketDataService marketdata.MarketDataService
	exchangeService   exchange.ExchangeService
	logger            *logger.Logger
}

// NewMarketDataHandler создает новый экземпляр MarketDataHandler.
func NewMarketDataHandler(marketDataService marketdata.MarketDataService, exchangeService exchange.ExchangeService, logger *logger.Logger) *MarketDataHandler {
	return &MarketDataHandler{
		marketDataService: marketDataService,
		exchangeService:   exchangeService,
		logger:            logger,
	}
}

// GetMarketData обрабатывает GET-запрос для получения рыночных данных.
func (h *MarketDataHandler) GetMarketData(w http.ResponseWriter, r *http.Request) {
	// marketData := h.exchangeService.LoadData()

	// // Отправка ответа в формате JSON
	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(marketData); err != nil {
	// 	h.logger.Errorf("Failed to encode response: %v", err)
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// }
}

func (h *MarketDataHandler) GetBacktestingPage(w http.ResponseWriter, r *http.Request) {

	if err := ui.BacktestingPage().Render(r.Context(), w); err != nil {
		h.logger.Errorf("Ошибка формирования страницы бектестинга: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

func (h *MarketDataHandler) PostRunBacktesting(w http.ResponseWriter, r *http.Request) {
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
