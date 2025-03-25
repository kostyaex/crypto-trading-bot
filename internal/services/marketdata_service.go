package services

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/utils"
)

type MarketDataService interface {
	SaveMarketData(data []*models.MarketData) error
	GetMarketData(symbol string, limit int) ([]*models.MarketData, error)
	GetMarketDataStatus(id int) (*models.MarketDataStatus, error)
	SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error
	GetMarketDataStatusList() ([]*models.MarketDataStatus, error)
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

// GetMarketDataStatus получает информацию о marketdatastatus по ID.
func (s *marketDataService) GetMarketDataStatus(id int) (*models.MarketDataStatus, error) {
	marketdatastatus, err := s.repo.MarketData.GetMarketDataStatus(id)
	if err != nil {
		s.logger.Errorf("Failed to get marketdatastatus: %v", err)
		return nil, err
	}
	return marketdatastatus, nil
}

// SaveMarketDataStatus сохраняет информацию о marketdatastatus.
func (s *marketDataService) SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error {
	if err := s.repo.MarketData.SaveMarketDataStatus(marketdatastatus); err != nil {
		s.logger.Errorf("Failed to save marketdatastatus: %v", err)
		return err
	}
	return nil
}

// GetMarketDataStatus получает список marketdatastatus.
func (s *marketDataService) GetMarketDataStatusList() ([]*models.MarketDataStatus, error) {
	marketdatastatus, err := s.repo.MarketData.GetMarketDataStatusList()
	if err != nil {
		s.logger.Errorf("Failed to get marketdatastatus list: %v", err)
		return nil, err
	}
	return marketdatastatus, nil
}
