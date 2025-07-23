package exchange

import (
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/models"
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

func NewBinance(apiKey, apiSecret string, logger *logger.Logger) *Binance {
	//client := binance.NewClient(apiKey, apiSecret)
	return &Binance{
		exchangeBase: exchangeBase{name: "Binance", logger: logger},
		baseURL:      "https://api.binance.com",
		//client:       client,
	}
}

// GetMarketData получает рыночные данные с Binance.
func (b *Binance) GetMarketData(symbol, interval string, startTime time.Time) (marketData []*models.MarketData, lastTime time.Time, err error) {

	//https://developers.binance.com/docs/binance-spot-api-docs/rest-api/market-data-endpoints#klinecandlestick-data

	irkutskLocation, err := time.LoadLocation("Asia/Irkutsk")
	if err != nil {
		panic(err)
	}

	// Преобразование startTime в миллисекунды с начала эпохи Unix
	// По умолчанию часовой пояс UTC
	startTimestamp := startTime.UTC().UnixNano() / int64(time.Millisecond)

	url := fmt.Sprintf("%s/api/v3/klines?symbol=%s&interval=%s&startTime=%d", b.baseURL, symbol, interval, startTimestamp)

	resp, err := http.Get(url)
	if err != nil {
		b.logger.Errorf("Failed to fetch market data: %v", err)
		return
	}
	defer resp.Body.Close()

	var klines [][]interface{}

	if err = json.NewDecoder(resp.Body).Decode(&klines); err != nil {
		b.logger.Errorf("Failed to decode market data: %v", err)
		return
	}

	for _, kline := range klines {
		//var _openTime int64 = kline[0].(float64)
		//openTime := int64(kline[0].(float64)) //strconv.ParseInt(kline[0].(string), 10, 64)
		closeTime := int64(kline[6].(float64))
		openPrice, _ := strconv.ParseFloat(kline[1].(string), 64)
		hightPrice, _ := strconv.ParseFloat(kline[2].(string), 64)
		lowPrice, _ := strconv.ParseFloat(kline[3].(string), 64)
		closePrice, _ := strconv.ParseFloat(kline[4].(string), 64)
		volume, _ := strconv.ParseFloat(kline[5].(string), 64)
		buyVolume, _ := strconv.ParseFloat(kline[9].(string), 64)

		// Преобразуем из UTC в локальное время
		lastTime = time.UnixMilli(closeTime).In(irkutskLocation)

		marketData = append(marketData, &models.MarketData{
			Symbol:     symbol,
			OpenPrice:  openPrice,
			HightPrice: hightPrice,
			LowPrice:   lowPrice,
			ClosePrice: closePrice,
			Volume:     volume,
			BuyVolume:  buyVolume,
			SellVolume: volume - buyVolume,
			TimeFrame:  interval,
			Timestamp:  lastTime,
		})
	}

	return
}
