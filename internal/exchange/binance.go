package exchange

import (
	"context"
	"crypto-trading-bot/internal/data"
	"crypto-trading-bot/internal/utils"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

type Binance struct {
	exchangeBase
	client *binance.Client
}

func NewBinance(apiKey, apiSecret string, logger *utils.Logger) *Binance {
	client := binance.NewClient(apiKey, apiSecret)
	return &Binance{
		exchangeBase: exchangeBase{name: "Binance", logger: logger},
		client:       client,
	}
}

func (b *Binance) GetMarketData() ([]*data.MarketData, error) {
	//symbols, err := b.client.NewListSymbolsService().Do(context.Background())
	symbols, err := b.client.NewListSymbolTickerService().Do(context.Background())
	if err != nil {
		b.logError(err, "Failed to list symbols")
		return nil, err
	}

	var marketData []*data.MarketData
	for _, symbol := range symbols {
		ticker, err := b.client.NewListPricesService().Symbol(symbol.Symbol).Do(context.Background())
		if err != nil {
			b.logError(err, "Failed to get ticker for symbol %s", symbol.Symbol)
			continue
		}

		for _, sp := range ticker {
			price, err := strconv.ParseFloat(sp.Price, 64)
			if err != nil {
				b.logError(err, "Failed to convert price %s", sp.Price)
				continue
			}

			marketData = append(marketData, &data.MarketData{
				Symbol:    sp.Symbol,
				Price:     price,
				Timestamp: time.Now(),
			})
		}

	}

	b.logInfo("Fetched market data: %v", marketData)
	return marketData, nil
}
