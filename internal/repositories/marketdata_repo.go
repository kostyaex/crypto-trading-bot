package repositories

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
	"database/sql"
	"fmt"
	"sort"
)

type MarketDataRepository interface {
	SaveMarketData(data []*models.MarketData) error
	GetMarketData(symbol string, limit int) ([]*models.MarketData, error)
	GetMarketDataStatus(id int) (*models.MarketDataStatus, error)
	SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error
	GetMarketDataStatusList() ([]*models.MarketDataStatus, error)
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

	stmt, err := tx.Prepare("INSERT INTO market_data (exchange, symbol, open_price, close_price, volume, buy_volume, sell_volume, time_frame, timestamp, vwap, buy_sell_volume_ratio, buy_sell_price_ratio, net_volume, net_volume_ratio, net_buy_sell_volume, buy_cluster, sell_cluster) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)")
	if err != nil {
		r.logger.Errorf("Failed to prepare statement: %v", err)
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, d := range data {
		_, err := stmt.Exec(d.Exchange, d.Symbol, d.OpenPrice, d.ClosePrice, d.Volume, d.BuyVolume, d.SellVolume, d.TimeFrame, d.Timestamp, d.VWAP, d.BuySellVolumeRatio, d.BuySellPriceRatio, d.NetVolume, d.NetVolumeRatio, d.NetBuySellVolume, d.BuyCluster, d.SellCluster)
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
        SELECT symbol, open_price, close_price, volume, buy_volume, sell_volume, time_frame, timestamp, vwap, buy_sell_volume_ratio, buy_sell_price_ratio, net_volume, net_volume_ratio, net_buy_sell_volume, buy_cluster, sell_cluster
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

// GetMarketDataStatus находит marketdatastatus по ID.
func (r *marketDataRepository) GetMarketDataStatus(id int) (*models.MarketDataStatus, error) {
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
func (r *marketDataRepository) SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error {
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
func (r *marketDataRepository) GetMarketDataStatusList() ([]*models.MarketDataStatus, error) {
	var marketdatastatus []*models.MarketDataStatus
	query := "SELECT id, exchange, symbol, time_frame, active, actual_time, status FROM market_data_statuss"

	err := r.db.Select(&marketdatastatus, query)
	if err != nil {
		r.logger.Errorf("Failed to get data from market_data_statuss: %v", err)
		return nil, err
	}

	return marketdatastatus, nil
}
