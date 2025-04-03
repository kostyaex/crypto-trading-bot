package services

import (
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/utils"
	"time"
)

type ExchangeService interface {
	LoadData(exchange exchange.Exchange, symbol string, timeFrame string, startTime time.Time) (marketData []*models.MarketData, lastTime time.Time, err error)
}

type exchangeService struct {
	repo      *repositories.Repository
	logger    *utils.Logger
	exchanges []exchange.Exchange
}

func NewEchangeService(repo *repositories.Repository, logger *utils.Logger, exchanges []exchange.Exchange) ExchangeService {
	return &exchangeService{
		repo:      repo,
		logger:    logger,
		exchanges: exchanges,
	}
}

// Загружает данные с указанной бирже по указанной паре и интервалу
func (s *exchangeService) LoadData(exchange exchange.Exchange, symbol string, timeFrame string, startTime time.Time) (marketData []*models.MarketData, lastTime time.Time, err error) {

	s.logger.Infof("Fetching data from exchange: %s %s %v", exchange.GetName(), symbol, startTime)

	marketData, lastTime, err = exchange.GetMarketData(symbol, timeFrame, startTime)
	if err != nil {
		s.logger.Errorf("Failed to fetch data from exchange %s: %v", exchange.GetName(), err)
		return
	}
	s.logger.Infof("Data fetching task completed")

	return
}
