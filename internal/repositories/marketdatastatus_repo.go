// crypto-trading-bot/repositories/MarketDataStatus_repo.go

package repositories

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
	"database/sql"
	"fmt"
)

// MarketDataStatusRepository определяет интерфейс для работы с MarketDataStatus.
type MarketDataStatusRepository interface {
	GetMarketDataStatus(id int) (*models.MarketDataStatus, error)
	SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error
	GetMarketDataStatusList() ([]*models.MarketDataStatus, error)
}

// marketdatastatusRepository реализует MarketDataStatusRepository.
type marketdatastatusRepository struct {
	db     *DB
	logger *utils.Logger
}

// NewMarketDataStatusRepository создает новый экземпляр MarketDataStatusRepository.
func NewMarketDataStatusRepository(db *DB, logger *utils.Logger) MarketDataStatusRepository {
	return &marketdatastatusRepository{db: db, logger: logger}
}

// GetMarketDataStatus находит marketdatastatus по ID.
func (r *marketdatastatusRepository) GetMarketDataStatus(id int) (*models.MarketDataStatus, error) {
	var marketdatastatus models.MarketDataStatus
	query := "SELECT id, exchange, symbol, time_frame, active, actual_time, status FROM market_data_statuss WHERE id = $1"

	err := r.db.Get(&marketdatastatus, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Errorf("marketdatastatus with ID %d not found", id)
			return nil, fmt.Errorf("marketdatastatus not found")
		}
		r.logger.Errorf("Failed to get marketdatastatus: %v", err)
		return nil, err
	}
	return &marketdatastatus, nil
}

// SaveMarketDataStatus сохраняет marketdatastatus в базу данных.
func (r *marketdatastatusRepository) SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error {
	query := "INSERT INTO market_data_statuss (exchange, symbol, time_frame, active, actual_time, status) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (exchange,symbol,time_frame) DO UPDATE SET active = $4, actual_time = $5, status = $6;"
	_, err := r.db.Exec(query,
		marketdatastatus.Exchange,
		marketdatastatus.Symbol,
		marketdatastatus.TimeFrame,
		marketdatastatus.Active,
		marketdatastatus.ActualTime,
		marketdatastatus.Status,
	)
	if err != nil {
		r.logger.Errorf("Failed to save marketdatastatus: %v", err)
		return err
	}
	return nil
}

// GetMarketDataStatusList выбирает список marketdatastatus из базы данных.
func (r *marketdatastatusRepository) GetMarketDataStatusList() ([]*models.MarketDataStatus, error) {
	var marketdatastatus []*models.MarketDataStatus
	query := "SELECT id, exchange, symbol, time_frame, active, actual_time, status FROM market_data_statuss"

	err := r.db.Select(&marketdatastatus, query)
	if err != nil {
		r.logger.Errorf("Failed to get data from market_data_statuss: %v", err)
		return nil, err
	}

	return marketdatastatus, nil
}
