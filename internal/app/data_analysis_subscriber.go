package app

import (
	"crypto-trading-bot/internal/analysis"
	"crypto-trading-bot/internal/data"
	"crypto-trading-bot/internal/utils"
)

// DataAnalysisSubscriber представляет подписчика, который анализирует данные
type DataAnalysisSubscriber struct {
	repo           *data.PostgresRepository
	logger         *utils.Logger
	eventPublisher *EventPublisher
}

// NewDataAnalysisSubscriber создает нового подписчика для анализа данных
func NewDataAnalysisSubscriber(repo *data.PostgresRepository, logger *utils.Logger, eventPublisher *EventPublisher) *DataAnalysisSubscriber {
	return &DataAnalysisSubscriber{
		repo:           repo,
		logger:         logger,
		eventPublisher: eventPublisher,
	}
}

// Handle обрабатывает событие
func (das *DataAnalysisSubscriber) Handle(event Event) {
	if event.Type() != "MarketDataLoaded" {
		return
	}

	marketData, ok := event.Payload().([]*data.MarketData)
	if !ok {
		das.logger.Errorf("Invalid payload type for MarketDataLoaded event")
		return
	}

	das.logger.Infof("Analyzing market data: %v", marketData)

	for _, md := range marketData {
		// Пример анализа данных
		rsi, err := analysis.CalculateRSI(das.repo, md.Symbol, 14)
		if err != nil {
			das.logger.Errorf("Failed to calculate RSI for symbol %s: %v", md.Symbol, err)
			continue
		}

		// macd, signal, err := analysis.CalculateMACD(das.repo, md.Symbol, 12, 26, 9)
		// if err != nil {
		// 	das.logger.Errorf("Failed to calculate MACD for symbol %s: %v", md.Symbol, err)
		// 	continue
		// }

		das.logger.Infof("RSI for %s: %.2f", md.Symbol, rsi)
		//das.logger.Infof("MACD for %s: %.2f, Signal: %.2f", md.Symbol, macd, signal)

		// Сохранение результатов анализа в базу данных
		if err := das.repo.SaveIndicator(md.Symbol, "RSI", rsi, md.Timestamp); err != nil {
			das.logger.Errorf("Failed to save RSI for symbol %s: %v", md.Symbol, err)
		}

		// if err := das.repo.SaveIndicator(md.Symbol, "MACD", macd, md.Timestamp); err != nil {
		// 	das.logger.Errorf("Failed to save MACD for symbol %s: %v", md.Symbol, err)
		// }

		// if err := das.repo.SaveIndicator(md.Symbol, "MACDSignal", signal, md.Timestamp); err != nil {
		// 	das.logger.Errorf("Failed to save MACD Signal for symbol %s: %v", md.Symbol, err)
		// }

		// Публикация события о завершении анализа данных
		analysisEvent := AnalysisCompletedEvent{
			Symbol: md.Symbol,
			Indicators: map[string]float64{
				"RSI": rsi,
				//"MACD":       macd,
				//"MACDSignal": signal,
			},
			Timestamp: md.Timestamp,
		}
		das.eventPublisher.Publish(analysisEvent)
	}
}
