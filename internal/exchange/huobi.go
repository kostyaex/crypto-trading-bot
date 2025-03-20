package exchange

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Huobi struct {
	exchangeBase
	apiKey    string
	apiSecret string
}

func NewHuobi(apiKey, apiSecret string, logger *utils.Logger) *Huobi {
	return &Huobi{
		exchangeBase: exchangeBase{name: "Huobi", logger: logger},
		apiKey:       apiKey,
		apiSecret:    apiSecret,
	}
}

func (h *Huobi) GetMarketData() ([]*models.MarketData, error) {
	url := "https://api.huobi.pro/market/tickers"

	resp, err := http.Get(url)
	if err != nil {
		h.logError(err, "Failed to fetch market data")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logError(err, "Failed to read response body")
		return nil, err
	}

	var result struct {
		Status string `json:"status"`
		Data   []struct {
			Symbol string `json:"symbol"`
			Close  string `json:"close"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		h.logError(err, "Failed to unmarshal response")
		return nil, err
	}

	if result.Status != "ok" {
		h.logError(fmt.Errorf("invalid status: %s", result.Status), "Failed to fetch market data")
		return nil, fmt.Errorf("invalid status: %s", result.Status)
	}

	var marketData []*models.MarketData
	for _, item := range result.Data {
		price, err := strconv.ParseFloat(item.Close, 64)
		if err != nil {
			h.logError(err, "Failed to parse price for symbol %s", item.Symbol)
			continue
		}

		marketData = append(marketData, &models.MarketData{
			Symbol:    item.Symbol,
			Price:     price,
			Timestamp: time.Now(),
		})
	}

	h.logInfo("Fetched market data: %v", marketData)
	return marketData, nil
}
