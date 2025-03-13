package data

import (
	"crypto-trading-bot/internal/utils"
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db     *sql.DB
	logger *utils.Logger
}

func NewPostgresRepository(db *sql.DB, logger *utils.Logger) *PostgresRepository {
	return &PostgresRepository{db: db, logger: logger}
}

func (r *PostgresRepository) SaveMarketData(data *MarketData) error {
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
