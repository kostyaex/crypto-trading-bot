package strategy

import (
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/models"
)

type StrategyService interface {
	GetActiveStrategies() ([]*models.Strategy, error)
	// GetStrategyByID(id int) (*models.Strategy, error)
	// SaveStrategy(strat *models.Strategy) error
	// UpdateStrategy(strat *models.Strategy) error
	// DeleteStrategy(id int) error
}

type strategyService struct {
	repo *repositories.Repository
}

func NewStrategyService(repo *repositories.Repository) StrategyService {
	return &strategyService{repo: repo}
}

func (s *strategyService) GetActiveStrategies() ([]*models.Strategy, error) {
	return s.repo.Strategy.GetActiveStrategies()
}

// func (s *strategyService) GetStrategyByID(id int) (*models.Strategy, error) {
// 	return s.repo.Strategy.GetStrategyByID(id)
// }

// func (s *strategyService) SaveStrategy(strat *models.Strategy) error {
// 	return s.repo.Strategy.SaveStrategy(strat)
// }

// func (s *strategyService) UpdateStrategy(strat *models.Strategy) error {
// 	return s.repo.Strategy.UpdateStrategy(strat)
// }

// func (s *strategyService) DeleteStrategy(id int) error {
// 	return s.repo.Strategy.DeleteStrategy(id)
// }
