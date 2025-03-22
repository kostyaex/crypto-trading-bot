package repositories

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
	"sort"
)

type MarketDataRepository interface {
	SaveMarketData(data []*models.MarketData) error
	GetMarketData(symbol string, limit int) ([]*models.MarketData, error)
}

type marketDataRepository struct {
	db     *DB
	logger *utils.Logger
}

func NewMarketDataRepository(db *DB, logger *utils.Logger) MarketDataRepository {
	return &marketDataRepository{db: db, logger: logger}
}

// SaveMarketData сохраняет рыночные данные в базу данных.
func (r *marketDataRepository) SaveMarketData(data []*models.MarketData) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Errorf("Failed to begin transaction: %v", err)
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO market_data (exchange, symbol, open_price, close_price, volume, time_frame, timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		r.logger.Errorf("Failed to prepare statement: %v", err)
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, d := range data {
		_, err := stmt.Exec(d.Exchange, d.Symbol, d.OpenPrice, d.ClosePrice, d.Volume, d.TimeFrame, d.Timestamp)
		if err != nil {
			r.logger.Errorf("Failed to insert market data: %v", err)
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("Failed to commit transaction: %v", err)
		tx.Rollback()
		return err
	}

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
