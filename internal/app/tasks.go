package app

import (
	"crypto-trading-bot/internal/data"
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/utils"
)

type DataFetchingTask struct {
	repo           *data.PostgresRepository
	exchanges      []exchange.Exchange
	logger         *utils.Logger
	eventPublisher *EventPublisher
}

func NewDataFetchingTask(repo *data.PostgresRepository, exchanges []exchange.Exchange, logger *utils.Logger, eventPublisher *EventPublisher) *DataFetchingTask {
	return &DataFetchingTask{
		repo:           repo,
		exchanges:      exchanges,
		logger:         logger,
		eventPublisher: eventPublisher,
	}
}

func (t *DataFetchingTask) Run() {
	t.logger.Infof("Starting data fetching task")
	var allMarketData []*data.MarketData
	for _, ex := range t.exchanges {
		t.logger.Infof("Fetching data from exchange: %s", ex.GetName())
		marketData, err := ex.GetMarketData()
		if err != nil {
			t.logger.Errorf("Failed to fetch data from exchange %s: %v", ex.GetName(), err)
			continue
		}
		for _, data := range marketData {
			if err := t.repo.SaveMarketData(data); err != nil {
				t.logger.Errorf("Failed to save market data for exchange %s: %v", ex.GetName(), err)
			} else {
				t.logger.Infof("Market data saved for exchange %s: %v", ex.GetName(), data)
			}
		}
		allMarketData = append(allMarketData, marketData...)
	}
	t.logger.Infof("Data fetching task completed")

	// Публикация события о загрузке рыночных данных
	event := MarketDataLoadedEvent{
		MarketData: allMarketData,
	}
	t.eventPublisher.Publish(event)
}
