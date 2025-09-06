package main

import (
	"crypto-trading-bot/internal/components"
	"crypto-trading-bot/internal/types"
	"encoding/json"
	"fmt"
)

func main() {
	registry := initRegistry()

	histSourceJson := json.RawMessage(`{
		"symbol" : "BTCUSDT",
		"interval" : "1s",
		"start_time" : "2025-07-23T00:00:00Z",
		"end_time" : "2025-07-23T01:00:00Z"
	}`)

	comp, err := registry.Build("source", histSourceJson)
	if err != nil {
		fmt.Printf("Ошибка формирования компоненты %s", err)
		return
	}

	source, ok := comp.(*components.HistoricalSourceSettings)
	if ok {
		fmt.Println(source.Symbol)
	}
}

func initRegistry() *components.ComponentRegistry {
	reg := components.NewComponentRegistry()

	reg.Register("source", func() types.Component {
		return &components.HistoricalSourceSettings{}
	})

	// Добавляй сюда новые компоненты — система сама их подхватит
	return reg
}
