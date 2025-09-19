package exchange

import (
	"fmt"
	"time"
)

// Exchange — интерфейс, который должны реализовать все биржи
type Exchange interface {
	// Асинхронные команды — возвращают ID для получения результата позже
	FetchCandlesAsync(symbol string, interval string, limit int) CommandID
	PlaceOrderAsync(order Order) CommandID
	FetchOpenPositionsAsync(symbol string) CommandID
	ClosePositionAsync(symbol string, side string) CommandID
	FetchBalanceAsync(asset string) CommandID

	// Получение данных по ID — неблокирующее
	PopCandle(cmdID CommandID) (Candle, bool, error) // (результат, есть_ли_результат)

	// Опционально: подписка на стримы через WebSocket
	SubscribeCandles(symbol string, interval string) CommandID
	UnsubscribeCandles(symbol string, interval string)
}

func GetCmdID(prefix, symbol, interval string) CommandID {
	return CommandID(fmt.Sprintf("%s_%s_%s_%d", prefix, symbol, interval, time.Now().UnixNano()))
}
