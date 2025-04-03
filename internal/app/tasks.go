package app

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/services"
	"crypto-trading-bot/internal/utils"
)

type DataFetchingTask struct {
	repo            *repositories.Repository
	exchangeService services.ExchangeService
	logger          *utils.Logger
	eventPublisher  *EventPublisher
	prevIntervalMD  []*models.MarketData // здесь сохраняем данные за предудыщий интервал
}

func NewDataFetchingTask(repo *repositories.Repository, exchangeService services.ExchangeService, logger *utils.Logger, eventPublisher *EventPublisher) *DataFetchingTask {
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
