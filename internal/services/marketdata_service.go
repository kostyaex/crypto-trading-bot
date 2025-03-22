package services

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
)

type MarketDataService interface {
	SaveMarketData(data *models.MarketData) error
	GetMarketData(symbol string, limit int) ([]*models.MarketData, error)
}

type marketDataService struct {
	repo *repositories.Repository
}

func NewMarketDataService(repo *repositories.Repository) MarketDataService {
	return &marketDataService{repo: repo}
}

func (s *marketDataService) SaveMarketData(data *models.MarketData) error {
	return s.repo.MarketData.SaveMarketData(data)
}

func (s *marketDataService) GetMarketData(symbol string, limit int) ([]*models.MarketData, error) {
	return s.repo.MarketData.GetMarketData(symbol, limit)
}
