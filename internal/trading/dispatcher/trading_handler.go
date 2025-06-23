package dispatcher

import (
	"fmt"
	"log"
	"time"
)

type TradingHandler struct{}

func (e *TradingHandler) Handle(signal TradeSignal) {
	switch signal.Type {
	case "buy":
		fmt.Printf("[BUY] %s | %.2f x %.2f\n", signal.Symbol, signal.Price, signal.Volume)
		// здесь можно вызвать API биржи
	case "sell":
		fmt.Printf("[SELL] %s | %.2f x %.2f\n", signal.Symbol, signal.Price, signal.Volume)
		// здесь можно вызвать API биржи
	}
}

type Position struct {
	Symbol     string
	Quantity   float64 // количество купленных активов
	EntryPrice float64 // цена входа
	OpenTime   time.Time
}

type TradingContext struct {
	CurrentPosition *Position
}

type PositionAwareExecutor struct {
	ctx *TradingContext
}

func (e *PositionAwareExecutor) Handle(signal TradeSignal) {
	switch signal.Type {
	case SignalBuy:
		if e.ctx.CurrentPosition == nil || e.ctx.CurrentPosition.Quantity <= 0 {
			e.executeBuy(signal)
		} else {
			log.Printf("[IGNORE BUY] Позиция уже открыта: %.2f @ %.2f", e.ctx.CurrentPosition.Quantity, e.ctx.CurrentPosition.EntryPrice)
		}
	case SignalSell:
		if e.ctx.CurrentPosition != nil && e.ctx.CurrentPosition.Quantity > 0 {
			sellQty := min(e.ctx.CurrentPosition.Quantity, signal.Volume)
			e.executeSell(signal, sellQty)
		} else {
			log.Printf("[IGNORE SELL] Нет открытой позиции для закрытия")
		}
	}
}

func (e *PositionAwareExecutor) executeBuy(signal TradeSignal) {
	log.Printf("[BUY] %s @ %.2f x %.2f", signal.Symbol, signal.Price, signal.Volume)

	e.ctx.CurrentPosition = &Position{
		Symbol:     signal.Symbol,
		Quantity:   signal.Volume,
		EntryPrice: signal.Price,
		OpenTime:   signal.Timestamp,
	}
}

func (e *PositionAwareExecutor) executeSell(signal TradeSignal, quantity float64) {
	log.Printf("[SELL] %s @ %.2f x %.2f", signal.Symbol, signal.Price, quantity)

	e.ctx.CurrentPosition.Quantity -= quantity
	if e.ctx.CurrentPosition.Quantity <= 0 {
		log.Printf("[POSITION CLOSED] Прибыль: %.2f", calculatePnL(e.ctx.CurrentPosition, signal.Price))
		e.ctx.CurrentPosition = nil
	}
}

func calculatePnL(pos *Position, exitPrice float64) float64 {
	return (exitPrice - pos.EntryPrice) * pos.Quantity
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
