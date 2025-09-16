package exchange

// Exchange — интерфейс, который должны реализовать все биржи
type Exchange interface {
	// Асинхронные команды — возвращают ID для получения результата позже
	FetchCandlesAsync(symbol string, interval string, limit int) CommandID
	PlaceOrderAsync(order Order) CommandID
	FetchOpenPositionsAsync(symbol string) CommandID
	ClosePositionAsync(symbol string, side string) CommandID
	FetchBalanceAsync(asset string) CommandID

	// Получение результата по ID — неблокирующее
	GetResult(cmdID CommandID) (interface{}, bool) // (результат, есть_ли_результат)

	// Опционально: подписка на стримы через WebSocket
	SubscribeCandles(symbol string, interval string, handler func(Candle))
	UnsubscribeCandles(symbol string, interval string)
}
