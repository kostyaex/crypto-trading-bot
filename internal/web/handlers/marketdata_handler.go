package handlers

import (
	"crypto-trading-bot/internal/services"
	"crypto-trading-bot/internal/utils"
	"encoding/json"
	"net/http"
)

// MarketDataHandler обрабатывает запросы, связанные с рыночными данными.
type MarketDataHandler struct {
	marketDataService services.MarketDataService
	exchangeService   services.ExchangeService
	logger            *utils.Logger
}

// NewMarketDataHandler создает новый экземпляр MarketDataHandler.
func NewMarketDataHandler(marketDataService services.MarketDataService, exchangeService services.ExchangeService, logger *utils.Logger) *MarketDataHandler {
	return &MarketDataHandler{marketDataService: marketDataService, exchangeService: exchangeService, logger: logger}
}

// GetMarketData обрабатывает GET-запрос для получения рыночных данных.
func (h *MarketDataHandler) GetMarketData(w http.ResponseWriter, r *http.Request) {
	marketData := h.exchangeService.FetchData()

	// Сохранение рыночных данных в базу данных
	if err := h.marketDataService.SaveMarketData(marketData); err != nil {
		h.logger.Errorf("Failed to save market data: %v", err)
		http.Error(w, "Failed to save market data", http.StatusInternalServerError)
		return
	}

	// Отправка ответа в формате JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(marketData); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
