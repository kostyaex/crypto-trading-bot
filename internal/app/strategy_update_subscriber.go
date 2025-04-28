package app

import (
	"crypto-trading-bot/internal/services"
	"crypto-trading-bot/internal/utils"
)

// StrategyUpdateSubscriber представляет подписчика, который обновляет состояния стратегий
type StrategyUpdateSubscriber struct {
	logger              *utils.Logger
	strategyService     services.StrategyService
	behaviorTreeService services.BehaviorTreeService
}

// NewStrategyUpdateSubscriber создает нового подписчика для обновления состояний стратегий
func NewStrategyUpdateSubscriber(logger *utils.Logger, strategyService services.StrategyService, behaviorTreeService services.BehaviorTreeService) *StrategyUpdateSubscriber {
	return &StrategyUpdateSubscriber{
		logger:              logger,
		strategyService:     strategyService,
		behaviorTreeService: behaviorTreeService,
	}
}

// Handle обрабатывает событие
func (sus *StrategyUpdateSubscriber) Handle(event Event) {
	// if event.Type() != "AnalysisCompleted" {
	// 	return
	// }

	// analysisEvent, ok := event.Payload().(AnalysisCompletedEvent)
	// if !ok {
	// 	sus.logger.Errorf("Invalid payload type for AnalysisCompleted event")
	// 	return
	// }

	// sus.logger.Infof("Updating strategy states for symbol: %s, timestamp: %v", analysisEvent.Symbol, analysisEvent.Timestamp)

	// // Получение всех активных стратегий
	// strategies, err := sus.strategyService.GetActiveStrategies()
	// if err != nil {
	// 	sus.logger.Errorf("Failed to get active strategies: %v", err)
	// 	return
	// }

	// for _, strat := range strategies {
	// 	sus.logger.Infof("Updating strategy: %s", strat.Name)

	// 	// Получение текущего состояния поведенческого дерева для стратегии
	// 	btState, err := sus.behaviorTreeService.GetBehaviorTreeState(strat.ID)
	// 	if err != nil {
	// 		sus.logger.Errorf("Failed to get behavior tree state for strategy %s: %v", strat.Name, err)
	// 		continue
	// 	}

	// 	// Создание поведенческого дерева из состояния
	// 	bt, err := strategy.NewBehaviorTree(btState)
	// 	if err != nil {
	// 		sus.logger.Errorf("Failed to create behavior tree for strategy %s: %v", strat.Name, err)
	// 		continue
	// 	}

	// 	// Обновление состояния поведенческого дерева на основе индикаторов
	// 	updatedState, err := bt.UpdateState(analysisEvent.Symbol, analysisEvent.Indicators, analysisEvent.Timestamp)
	// 	if err != nil {
	// 		sus.logger.Errorf("Failed to update behavior tree state for strategy %s: %v", strat.Name, err)
	// 		continue
	// 	}

	// 	// Сохранение обновленного состояния поведенческого дерева в базу данных
	// 	if err := sus.behaviorTreeService.SaveBehaviorTreeState(strat.ID, updatedState); err != nil {
	// 		sus.logger.Errorf("Failed to save behavior tree state for strategy %s: %v", strat.Name, err)
	// 		continue
	// 	}

	// 	sus.logger.Infof("Behavior tree state updated for strategy %s", strat.Name)
	// }
}
