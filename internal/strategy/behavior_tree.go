package strategy

import (
	"fmt"
	"time"
)

// BehaviorTree представляет поведенческое дерево
type BehaviorTree struct {
	RootNode Node
	State    map[string]interface{}
}

// Node представляет узел поведенческого дерева
type Node interface {
	Execute(symbol string, indicators map[string]float64, timestamp time.Time) (bool, error)
}

// NewBehaviorTree создает новое поведенческое дерево из состояния
func NewBehaviorTree(state map[string]interface{}) (*BehaviorTree, error) {
	// Здесь можно добавить логику создания поведенческого дерева из состояния
	// Например, парсинг JSON и создание узлов
	rootNode := &ExampleNode{}
	return &BehaviorTree{
		RootNode: rootNode,
		State:    state,
	}, nil
}

// UpdateState обновляет состояние поведенческого дерева на основе индикаторов
func (bt *BehaviorTree) UpdateState(symbol string, indicators map[string]float64, timestamp time.Time) (map[string]interface{}, error) {
	success, err := bt.RootNode.Execute(symbol, indicators, timestamp)
	if err != nil {
		return nil, err
	}

	// Обновляем состояние поведенческого дерева
	bt.State["last_execution"] = timestamp
	bt.State["success"] = success

	return bt.State, nil
}

// ExampleNode представляет пример узла поведенческого дерева
type ExampleNode struct{}

// Execute выполняет логику узла
func (en *ExampleNode) Execute(symbol string, indicators map[string]float64, timestamp time.Time) (bool, error) {
	rsi, ok := indicators["RSI"]
	if !ok {
		return false, fmt.Errorf("RSI indicator not found")
	}

	macd, ok := indicators["MACD"]
	if !ok {
		return false, fmt.Errorf("MACD indicator not found")
	}

	macdSignal, ok := indicators["MACDSignal"]
	if !ok {
		return false, fmt.Errorf("MACD Signal indicator not found")
	}

	// Пример логики узла
	if rsi < 30 && macd > macdSignal {
		return true, nil // Бай
	}

	if rsi > 70 && macd < macdSignal {
		return true, nil // Селл
	}

	return false, nil
}
