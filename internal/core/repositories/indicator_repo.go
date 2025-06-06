package repositories

import (
	"crypto-trading-bot/internal/core/logger"
	"time"
)

type IndicatorRepository interface {
	SaveIndicator(symbol, indicatorName string, value float64, timestamp time.Time) error
}

type indicatorRepository struct {
	db     *DB
	logger *logger.Logger
}

func NewIndicatorRepository(db *DB, logger *logger.Logger) IndicatorRepository {
	return &indicatorRepository{db: db, logger: logger}
}

func (r *indicatorRepository) SaveIndicator(symbol, indicatorName string, value float64, timestamp time.Time) error {
	query := `
        INSERT INTO indicators (symbol, indicator_name, value, timestamp)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (symbol, indicator_name, timestamp) DO UPDATE
        SET value = EXCLUDED.value;
    `

	_, err := r.db.Exec(query, symbol, indicatorName, value, timestamp)
	if err != nil {
		r.logger.Errorf("Failed to save indicator: %v", err)
		return err
	}

	r.logger.Infof("Indicator saved: symbol=%s, indicator_name=%s, value=%.2f, timestamp=%v", symbol, indicatorName, value, timestamp)
	return nil
}
