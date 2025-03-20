package app

import (
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/strategy"
	"crypto-trading-bot/internal/utils"
)

// StrategyUpdateSubscriber представляет подписчика, который обновляет состояния стратегий
type StrategyUpdateSubscriber struct {
	repo   *repositories.Repository
	logger *utils.Logger
}

// NewStrategyUpdateSubscriber создает нового подписчика для обновления состояний стратегий
func NewStrategyUpdateSubscriber(repo *repositories.Repository, logger *utils.Logger) *StrategyUpdateSubscriber {
	return &StrategyUpdateSubscriber{
		repo:   repo,
		logger: logger,
	}
}

// Handle обрабатывает событие
func (sus *StrategyUpdateSubscriber) Handle(event Event) {
	if event.Type() != "AnalysisCompleted" {
		return
	}

	analysisEvent, ok := event.Payload().(AnalysisCompletedEvent)
	if !ok {
		sus.logger.Errorf("Invalid payload type for AnalysisCompleted event")
		return
	}

	sus.logger.Infof("Updating strategy states for symbol: %s, timestamp: %v", analysisEvent.Symbol, analysisEvent.Timestamp)

	// Получение всех активных стратегий
	strategies, err := sus.repo.Strategy.GetActiveStrategies()
	if err != nil {
		sus.logger.Errorf("Failed to get active strategies: %v", err)
		return
	}

	for _, strat := range strategies {
		sus.logger.Infof("Updating strategy: %s", strat.Name)

		// Получение текущего состояния поведенческого дерева для стратегии
		btState, err := sus.repo.GetBehaviorTreeState(strat.ID)
		if err != nil {
			sus.logger.Errorf("Failed to get behavior tree state for strategy %s: %v", strat.Name, err)
			continue
		}

		// Создание поведенческого дерева из состояния
		bt, err := strategy.NewBehaviorTree(btState)
		if err != nil {
			sus.logger.Errorf("Failed to create behavior tree for strategy %s: %v", strat.Name, err)
			continue
		}

		// Обновление состояния поведенческого дерева на основе индикаторов
		updatedState, err := bt.UpdateState(analysisEvent.Symbol, analysisEvent.Indicators, analysisEvent.Timestamp)
		if err != nil {
			sus.logger.Errorf("Failed to update behavior tree state for strategy %s: %v", strat.Name, err)
			continue
		}

		// Сохранение обновленного состояния поведенческого дерева в базу данных
		if err := sus.repo.SaveBehaviorTreeState(strat.ID, updatedState); err != nil {
			sus.logger.Errorf("Failed to save behavior tree state for strategy %s: %v", strat.Name, err)
			continue
		}

		sus.logger.Infof("Behavior tree state updated for strategy %s", strat.Name)
	}
}
