package exchange

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Binance struct {
	exchangeBase
	baseURL string
	//client *binance.Client
}

func NewBinance(apiKey, apiSecret string, logger *utils.Logger) *Binance {
	//client := binance.NewClient(apiKey, apiSecret)
	return &Binance{
		exchangeBase: exchangeBase{name: "Binance", logger: logger},
		baseURL:      "https://api.binance.com",
		//client:       client,
	}
}

// GetMarketData получает рыночные данные с Binance.
func (b *Binance) GetMarketData() ([]*models.MarketData, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/24hr", b.baseURL)

	resp, err := http.Get(url)
	if err != nil {
		b.logger.Errorf("Failed to fetch market data: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var tickers []struct {
		Symbol     string `json:"symbol"`
		OpenPrice  string `json:"openPrice"`
		ClosePrice string `json:"lastPrice"`
		Volume     string `json:"volume"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tickers); err != nil {
		b.logger.Errorf("Failed to decode market data: %v", err)
		return nil, err
	}

	var marketData []*models.MarketData
	for _, ticker := range tickers {
		openPrice, _ := strconv.ParseFloat(ticker.OpenPrice, 64)
		closePrice, _ := strconv.ParseFloat(ticker.ClosePrice, 64)
		volume, _ := strconv.ParseFloat(ticker.Volume, 64)

		marketData = append(marketData, &models.MarketData{
			Exchange:   "binance",
			Symbol:     ticker.Symbol,
			OpenPrice:  openPrice,
			ClosePrice: closePrice,
			Volume:     volume,
			TimeFrame:  "1d", // Пример таймфрейма
			Timestamp:  time.Now(),
		})
	}

	return marketData, nil
}
