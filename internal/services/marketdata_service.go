package services

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/utils"
)

type MarketDataService interface {
	SaveMarketData(data []*models.MarketData) error
	GetMarketData(symbol string, limit int) ([]*models.MarketData, error)
}

type marketDataService struct {
	repo   *repositories.Repository
	logger *utils.Logger
}

func NewMarketDataService(repo *repositories.Repository, logger *utils.Logger) MarketDataService {
	return &marketDataService{repo: repo, logger: logger}
}

// SaveMarketData сохраняет рыночные данные.
func (s *marketDataService) SaveMarketData(data []*models.MarketData) error {
	if err := s.repo.MarketData.SaveMarketData(data); err != nil {
		s.logger.Errorf("Failed to save market data: %v", err)
		return err
	}
	return nil
}

func (s *marketDataService) GetMarketData(symbol string, limit int) ([]*models.MarketData, error) {
	return s.repo.MarketData.GetMarketData(symbol, limit)
}
