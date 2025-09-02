package app

import (
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/repositories"
	"crypto-trading-bot/internal/service/exchange"
	"crypto-trading-bot/pkg/types"
)

type DataFetchingTask struct {
	repo            *repositories.Repository
	exchangeService exchange.ExchangeService
	logger          *logger.Logger
	eventPublisher  *EventPublisher
	prevIntervalMD  []*types.MarketData // здесь сохраняем данные за предудыщий интервал
}

func NewDataFetchingTask(repo *repositories.Repository, exchangeService exchange.ExchangeService, logger *logger.Logger, eventPublisher *EventPublisher) *DataFetchingTask {
	return &DataFetchingTask{
		repo:            repo,
		exchangeService: exchangeService,
		logger:          logger,
		eventPublisher:  eventPublisher,
	}
}

func (t *DataFetchingTask) Run() {
	// allMarketData := t.exchangeService.LoadData()
	// // Публикация события о загрузке рыночных данных
	// event := MarketDataLoadedEvent{
	// 	MarketData: allMarketData,
	// }
	// t.eventPublisher.Publish(event)
}
