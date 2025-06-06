package repositories

import (
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/models"
)

type StrategyRepository interface {
	GetActiveStrategies() ([]*models.Strategy, error)
	// GetStrategyByID(id int) (*models.Strategy, error)
	// SaveStrategy(strat *models.Strategy) error
	// UpdateStrategy(strat *models.Strategy) error
	// DeleteStrategy(id int) error
}

type strategyRepository struct {
	db     *DB
	logger *logger.Logger
}

func NewStrategyRepository(db *DB, logger *logger.Logger) StrategyRepository {
	return &strategyRepository{db: db, logger: logger}
}

func (r *strategyRepository) GetActiveStrategies() ([]*models.Strategy, error) {
	query := `
		 SELECT id, name, description, config
		 FROM strategies
		 WHERE active = true;
	 `

	var strategies []*models.Strategy
	err := r.db.Select(&strategies, query)
	if err != nil {
		r.logger.Errorf("Failed to get active strategies: %v", err)
		return nil, err
	}

	r.logger.Infof("Active strategies retrieved: %v", strategies)
	return strategies, nil
}

// func (r *strategyRepository) GetStrategyByID(id int) (*models.Strategy, error) {
// 	query := `
//         SELECT id, name, description, config, active
//         FROM strategies
//         WHERE id = $1;
//     `

// 	var strat models.Strategy
// 	err := r.db.Get(&strat, query, id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			r.logger.Warnf("Strategy with ID %d not found", id)
// 			return nil, err
// 		}
// 		r.logger.Errorf("Failed to get strategy by ID %d: %v", id, err)
// 		return nil, err
// 	}

// 	r.logger.Infof("Strategy retrieved: %v", strat)
// 	return &strat, nil
// }

// func (r *strategyRepository) SaveStrategy(strat *models.Strategy) error {
// 	query := `
//         INSERT INTO strategies (name, description, config, active)
//         VALUES ($1, $2, $3, $4)
//         RETURNING id;
//     `

// 	err := r.db.QueryRow(query, strat.Name, strat.Description, strat.Config, strat.Active).Scan(&strat.ID)
// 	if err != nil {
// 		r.logger.Errorf("Failed to save strategy: %v", err)
// 		return err
// 	}

// 	r.logger.Infof("Strategy saved: %v", strat)
// 	return nil
// }

// func (r *strategyRepository) UpdateStrategy(strat *models.Strategy) error {
// 	query := `
//         UPDATE strategies
//         SET name = $2, description = $3, config = $4, active = $5
//         WHERE id = $1;
//     `

// 	_, err := r.db.Exec(query, strat.ID, strat.Name, strat.Description, strat.Config, strat.Active)
// 	if err != nil {
// 		r.logger.Errorf("Failed to update strategy: %v", err)
// 		return err
// 	}

// 	r.logger.Infof("Strategy updated: %v", strat)
// 	return nil
// }

// func (r *strategyRepository) DeleteStrategy(id int) error {
// 	query := `
//         DELETE FROM strategies
//         WHERE id = $1;
//     `

// 	_, err := r.db.Exec(query, id)
// 	if err != nil {
// 		r.logger.Errorf("Failed to delete strategy: %v", err)
// 		return err
// 	}

// 	r.logger.Infof("Strategy deleted: ID %d", id)
// 	return nil
// }
