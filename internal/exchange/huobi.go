package exchange

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Huobi struct {
	exchangeBase
	baseURL   string
	apiKey    string
	apiSecret string
}

func NewHuobi(apiKey, apiSecret string, logger *utils.Logger) *Huobi {
	return &Huobi{
		exchangeBase: exchangeBase{name: "Huobi", logger: logger},
		baseURL:      "https://api.huobi.pro",
		apiKey:       apiKey,
		apiSecret:    apiSecret,
	}
}

// GetMarketData получает рыночные данные с Huobi.
func (h *Huobi) GetMarketData() ([]*models.MarketData, error) {
	url := fmt.Sprintf("%s/market/history/kline?symbol=btcusdt&period=1day", h.baseURL)

	resp, err := http.Get(url)
	if err != nil {
		h.logger.Errorf("Failed to fetch market data: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var huobiResponse struct {
		Data []struct {
			ID     int     `json:"id"`
			Open   float64 `json:"open"`
			Close  float64 `json:"close"`
			Amount float64 `json:"amount"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&huobiResponse); err != nil {
		h.logger.Errorf("Failed to decode market data: %v", err)
		return nil, err
	}

	var marketData []*models.MarketData
	for _, kline := range huobiResponse.Data {
		timestamp := time.Unix(int64(kline.ID), 0)

		marketData = append(marketData, &models.MarketData{
			Exchange:   "huobi",
			Symbol:     "BTCUSDT",
			OpenPrice:  kline.Open,
			ClosePrice: kline.Close,
			Volume:     kline.Amount,
			TimeFrame:  "1d", // Пример таймфрейма
			Timestamp:  timestamp,
		})
	}

	return marketData, nil
}
