package data

import (
	"crypto-trading-bot/internal/strategy"
	"crypto-trading-bot/internal/utils"
	"database/sql"
	"encoding/json"
	"sort"
	"time"
)

type PostgresRepository struct {
	db     *DB
	logger *utils.Logger
}

func NewPostgresRepository(db *DB, logger *utils.Logger) *PostgresRepository {
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

func (r *PostgresRepository) SaveIndicator(symbol, indicatorName string, value float64, timestamp time.Time) error {
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

func (r *PostgresRepository) GetMarketData(symbol string, limit int) ([]*MarketData, error) {
	query := `
        SELECT symbol, price, timestamp
        FROM market_data
        WHERE symbol = $1
        ORDER BY timestamp DESC
        LIMIT $2;
    `

	var marketData []*MarketData
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

func (r *PostgresRepository) GetActiveStrategies() ([]*strategy.Strategy, error) {
	query := `
        SELECT id, name, description, config
        FROM strategies
        WHERE active = true;
    `

	var strategies []*strategy.Strategy
	err := r.db.Select(&strategies, query)
	if err != nil {
		r.logger.Errorf("Failed to get active strategies: %v", err)
		return nil, err
	}

	r.logger.Infof("Active strategies retrieved: %v", strategies)
	return strategies, nil
}

func (r *PostgresRepository) GetBehaviorTreeState(strategyID int) (map[string]interface{}, error) {
	query := `
        SELECT state
        FROM behavior_trees
        WHERE strategy_id = $1
        ORDER BY last_executed DESC
        LIMIT 1;
    `

	var stateJSON json.RawMessage
	err := r.db.QueryRow(query, strategyID).Scan(&stateJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warnf("No behavior tree state found for strategy ID %d", strategyID)
			return nil, nil
		}
		r.logger.Errorf("Failed to get behavior tree state for strategy ID %d: %v", strategyID, err)
		return nil, err
	}

	var state map[string]interface{}
	if err := json.Unmarshal(stateJSON, &state); err != nil {
		r.logger.Errorf("Failed to unmarshal behavior tree state for strategy ID %d: %v", strategyID, err)
		return nil, err
	}

	r.logger.Infof("Behavior tree state retrieved for strategy ID %d: %v", strategyID, state)
	return state, nil
}

func (r *PostgresRepository) SaveBehaviorTreeState(strategyID int, state map[string]interface{}) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		r.logger.Errorf("Failed to marshal behavior tree state for strategy ID %d: %v", strategyID, err)
		return err
	}

	query := `
        INSERT INTO behavior_trees (strategy_id, state, last_executed)
        VALUES ($1, $2, NOW());
    `

	_, err = r.db.Exec(query, strategyID, stateJSON)
	if err != nil {
		r.logger.Errorf("Failed to save behavior tree state for strategy ID %d: %v", strategyID, err)
		return err
	}

	r.logger.Infof("Behavior tree state saved for strategy ID %d: %v", strategyID, state)
	return nil
}
