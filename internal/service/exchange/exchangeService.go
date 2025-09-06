package exchange

import (
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/types"
	"time"
)

type ExchangeService interface {
	LoadData(exchange Exchange, symbol string, timeFrame string, startTime time.Time) (marketData []*types.MarketData, lastTime time.Time, err error)
}

type exchangeService struct {
	repo      *repositories.Repository
	logger    *logger.Logger
	exchanges []Exchange
}

func NewEchangeService(repo *repositories.Repository, logger *logger.Logger, exchanges []Exchange) ExchangeService {
	return &exchangeService{
		repo:      repo,
		logger:    logger,
		exchanges: exchanges,
	}
}

// Загружает данные с указанной бирже по указанной паре и интервалу
func (s *exchangeService) LoadData(exchange Exchange, symbol string, timeFrame string, startTime time.Time) (marketData []*types.MarketData, lastTime time.Time, err error) {

	s.logger.Infof("Fetching data from exchange: %s %s %v", exchange.GetName(), symbol, startTime)

	marketData, lastTime, err = exchange.GetMarketData(symbol, timeFrame, startTime)
	if err != nil {
		s.logger.Errorf("Failed to fetch data from exchange %s: %v", exchange.GetName(), err)
		return
	}

	if lastTime.Before(startTime) {
		lastTime = startTime
	}

	s.logger.Infof("Data fetching task completed")

	return
}
