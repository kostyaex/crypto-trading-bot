// crypto-trading-bot/services/MarketDataStatus_service.go

package services

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/utils"
)

// MarketDataStatusService определяет интерфейс для бизнес-логики MarketDataStatus.
type MarketDataStatusService interface {
	GetMarketDataStatus(id int) (*models.MarketDataStatus, error)
	SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error
	GetMarketDataStatusList() ([]*models.MarketDataStatus, error)
}

// marketdatastatusService реализует MarketDataStatusService.
type marketdatastatusService struct {
	marketdatastatusRepo repositories.MarketDataStatusRepository
	logger               *utils.Logger
}

// NewMarketDataStatusService создает новый экземпляр MarketDataStatusService.
func NewMarketDataStatusService(marketdatastatusRepo repositories.MarketDataStatusRepository, logger *utils.Logger) MarketDataStatusService {
	return &marketdatastatusService{
		marketdatastatusRepo: marketdatastatusRepo,
		logger:               logger,
	}
}

// GetMarketDataStatus получает информацию о marketdatastatus по ID.
func (s *marketdatastatusService) GetMarketDataStatus(id int) (*models.MarketDataStatus, error) {
	marketdatastatus, err := s.marketdatastatusRepo.GetMarketDataStatus(id)
	if err != nil {
		s.logger.Errorf("Failed to get marketdatastatus: %v", err)
		return nil, err
	}
	return marketdatastatus, nil
}

// SaveMarketDataStatus сохраняет информацию о marketdatastatus.
func (s *marketdatastatusService) SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error {
	if err := s.marketdatastatusRepo.SaveMarketDataStatus(marketdatastatus); err != nil {
		s.logger.Errorf("Failed to save marketdatastatus: %v", err)
		return err
	}
	return nil
}

// GetMarketDataStatus получает список marketdatastatus.
func (s *marketdatastatusService) GetMarketDataStatusList() ([]*models.MarketDataStatus, error) {
	marketdatastatus, err := s.marketdatastatusRepo.GetMarketDataStatusList()
	if err != nil {
		s.logger.Errorf("Failed to get marketdatastatus list: %v", err)
		return nil, err
	}
	return marketdatastatus, nil
}
