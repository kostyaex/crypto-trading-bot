package repositories

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
	"sort"
)

type MarketDataRepository interface {
	SaveMarketData(data *models.MarketData) error
	GetMarketData(symbol string, limit int) ([]*models.MarketData, error)
}

type marketDataRepository struct {
	db     *DB
	logger *utils.Logger
}

func NewMarketDataRepository(db *DB, logger *utils.Logger) MarketDataRepository {
	return &marketDataRepository{db: db, logger: logger}
}

func (r *marketDataRepository) SaveMarketData(data *models.MarketData) error {
	query := `
        INSERT INTO market_data (symbol, price, timestamp)
        VALUES ($1, $2, $3)
        ON CONFLICT (symbol, timestamp) DO UPDATE
        SET price = EXCLUDED.price;
    `

	_, err := r.db.Exec(query, data.Symbol, data.Price, data.Timestamp)
	if err != nil {
		r.logger.Errorf("Failed to save market data: %v", err)
		return err
	}

	r.logger.Infof("Market data saved: %v", data)
	return nil
}

func (r *marketDataRepository) GetMarketData(symbol string, limit int) ([]*models.MarketData, error) {
	query := `
        SELECT symbol, price, timestamp
        FROM market_data
        WHERE symbol = $1
        ORDER BY timestamp DESC
        LIMIT $2;
    `

	var marketData []*models.MarketData
	err := r.db.Select(&marketData, query, symbol, limit)
	if err != nil {
		r.logger.Errorf("Failed to get market data for symbol %s: %v", symbol, err)
		return nil, err
	}

	// Сортируем данные по возрастанию времени
	sort.Slice(marketData, func(i, j int) bool {
		return marketData[i].Timestamp.Before(marketData[j].Timestamp)
	})

	r.logger.Infof("Market data retrieved for symbol %s: %v", symbol, marketData)
	return marketData, nil
}
