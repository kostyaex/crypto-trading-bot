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

// GetMarketData получает рыночные данные с Huobi по указанному символу, интервалу и начальной дате.
func (h *Huobi) GetMarketData(symbol, interval string, startTime time.Time) (marketData []*models.MarketData, lastTime time.Time, err error) {
	// Преобразование startTime в миллисекунды с начала эпохи Unix
	startTimestamp := startTime.UnixNano() / int64(time.Millisecond)

	url := fmt.Sprintf("%s/market/history/kline?symbol=%s&period=%s&from=%d", h.baseURL, symbol, interval, startTimestamp)

	resp, err := http.Get(url)
	if err != nil {
		h.logger.Errorf("Failed to fetch market data: %v", err)
		return
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

	if err = json.NewDecoder(resp.Body).Decode(&huobiResponse); err != nil {
		h.logger.Errorf("Failed to decode market data: %v", err)
		return
	}

	for _, kline := range huobiResponse.Data {
		timestamp := time.Unix(int64(kline.ID), 0)
		lastTime = timestamp

		marketData = append(marketData, &models.MarketData{
			Symbol:     symbol,
			OpenPrice:  kline.Open,
			ClosePrice: kline.Close,
			Volume:     kline.Amount,
			TimeFrame:  interval,
			Timestamp:  timestamp,
		})
	}

	return
}
