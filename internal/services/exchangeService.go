package services

import (
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/utils"
)

type ExchangeService interface {
	FetchData() []*models.MarketData
}

type exchangeService struct {
	repo      *repositories.Repository
	exchanges []exchange.Exchange
	logger    *utils.Logger
}

func NewEchangeService(repo *repositories.Repository, logger *utils.Logger, exchanges []exchange.Exchange) ExchangeService {
	return &exchangeService{repo: repo, logger: logger, exchanges: exchanges}
}

func (s *exchangeService) FetchData() []*models.MarketData {
	s.logger.Infof("Starting data fetching")
	var allMarketData []*models.MarketData
	for _, ex := range s.exchanges {
		s.logger.Infof("Fetching data from exchange: %s", ex.GetName())
		marketData, err := ex.GetMarketData()
		if err != nil {
			s.logger.Errorf("Failed to fetch data from exchange %s: %v", ex.GetName(), err)
			continue
		}
		if err := s.repo.MarketData.SaveMarketData(marketData); err != nil {
			s.logger.Errorf("Failed to save market data for exchange %s: %v", ex.GetName(), err)
		} else {
			s.logger.Infof("Market data saved for exchange %s", ex.GetName())
		}
		allMarketData = append(allMarketData, marketData...)
	}
	s.logger.Infof("Data fetching task completed")

	return allMarketData
}
