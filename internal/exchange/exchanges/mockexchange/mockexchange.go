package mockexchange

import (
	"crypto-trading-bot/internal/exchange"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type MockExchange struct {
	//asyncMgr *exchange.AsyncManager
	dataQueues map[exchange.CommandID]*exchange.PriorityQueueManager[exchange.Candle]
	//orderQueues map[exchange.CommandID]*exchange.PriorityQueueManager[exchange.Order]
	results map[exchange.CommandID]interface{}
	mu      sync.RWMutex
	err     error

	// Настройки для тестов
	DelayMin time.Duration // Минимальная задержка имитации
	DelayMax time.Duration // Максимальная задержка имитации
	ErrRate  float64       // Вероятность ошибки (0.0 - 1.0)
}

func NewMockExchange() *MockExchange {
	return &MockExchange{
		//asyncMgr: exchange.NewAsyncManager(10 * time.Minute),
		dataQueues: make(map[exchange.CommandID]*exchange.PriorityQueueManager[exchange.Candle]),
		//orderQueues: make(map[exchange.CommandID]*exchange.PriorityQueueManager[exchange.Order]),
		results:  make(map[exchange.CommandID]interface{}),
		DelayMin: 50 * time.Millisecond,
		DelayMax: 500 * time.Millisecond,
		ErrRate:  0.0, // по умолчанию ошибок нет
	}
}

// Утилита для генерации случайной задержки
func (m *MockExchange) randomDelay() time.Duration {
	if m.DelayMin == m.DelayMax {
		return m.DelayMin
	}
	delta := m.DelayMax - m.DelayMin
	return m.DelayMin + time.Duration(rand.Int63n(int64(delta)))
}

// Утилита: симуляция ошибки
func (m *MockExchange) shouldError() bool {
	return rand.Float64() < m.ErrRate
}

// ————————————————————————————————————————————————————————————————
// Реализация интерфейса exchange.Exchange
// ————————————————————————————————————————————————————————————————

func (m *MockExchange) FetchCandlesAsync(symbol string, interval string, limit int) exchange.CommandID {
	cmdID := exchange.GetCmdID("mock_candles", symbol, interval)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				//m.asyncMgr.StoreResult(cmdID, fmt.Errorf("panic in mock: %v", r))
				m.err = fmt.Errorf("panic in mock: %v", r)
			}
		}()

		time.Sleep(m.randomDelay())

		if m.shouldError() {
			//m.asyncMgr.StoreResult(cmdID, fmt.Errorf("simulated error fetching candles for %s", symbol))
			m.err = fmt.Errorf("simulated error fetching candles for %s", symbol)
			return
		}

		var candles []*exchange.Record[exchange.Candle]
		now := time.Now()
		for i := 0; i < limit; i++ {
			ts := now.Add(time.Duration(-i) * time.Minute)
			candles = append(candles, &exchange.Record[exchange.Candle]{
				Timestamp: ts,
				Data: exchange.Candle{
					Symbol:    symbol,
					Interval:  interval,
					Timestamp: ts,
					Open:      50000.0 + float64(i)*10,
					High:      50100.0 + float64(i)*5,
					Low:       49900.0 + float64(i)*5,
					Close:     50050.0 + float64(i)*10,
					Volume:    100.0 + float64(i),
				},
			})
		}

		//m.asyncMgr.StoreResult(cmdID, candles)
		if m.dataQueues[cmdID] == nil {
			m.dataQueues[cmdID] = exchange.NewPriorityQueueManager[exchange.Candle]()
		}
		m.dataQueues[cmdID].PushBatch(candles...)
	}()

	return cmdID
}

func (m *MockExchange) PlaceOrderAsync(order exchange.Order) exchange.CommandID {
	cmdID := exchange.CommandID(fmt.Sprintf("mock_order_%s_%s", order.Symbol, time.Now().Format("20060102150405.999")))

	// go func() {
	// 	time.Sleep(m.randomDelay())

	// 	if m.shouldError() {
	// 		//m.asyncMgr.StoreResult(cmdID, fmt.Errorf("simulated order rejection"))
	// 		m.err = fmt.Errorf("simulated order rejection")
	// 		return
	// 	}

	// 	// Имитируем успешный ордер
	// 	filledOrder := order
	// 	filledOrder.ID = "mock_order_id_" + time.Now().Format("20060102150405")
	// 	filledOrder.Status = "filled"

	// 	//m.asyncMgr.StoreResult(cmdID, filledOrder)
	// 	if m.orderQueues[cmdID] == nil {
	// 		m.orderQueues[cmdID] = exchange.NewPriorityQueueManager[exchange.Order]()
	// 	}
	// 	m.orderQueues[cmdID].PushBatch(&exchange.Record[exchange.Order]{
	// 		Timestamp: time.Now(),
	// 		Data:      filledOrder,
	// 	})

	// }()

	return cmdID
}

func (m *MockExchange) FetchOpenPositionsAsync(symbol string) exchange.CommandID {
	cmdID := exchange.CommandID(fmt.Sprintf("mock_positions_%s_%d", symbol, time.Now().UnixNano()))

	// go func() {
	// 	time.Sleep(m.randomDelay())

	// 	if m.shouldError() {
	// 		//m.asyncMgr.StoreResult(cmdID, fmt.Errorf("simulated error fetching positions"))
	// 		m.err = fmt.Errorf("simulated error fetching positions")
	// 		return
	// 	}

	// 	positions := []exchange.Position{
	// 		{
	// 			Symbol:        symbol,
	// 			Side:          "long",
	// 			Size:          0.1,
	// 			Entry:         50000.0,
	// 			UnrealizedPnL: 250.0,
	// 		},
	// 	}

	// 	m.asyncMgr.StoreResult(cmdID, positions)
	// }()

	return cmdID
}

func (m *MockExchange) ClosePositionAsync(symbol string, side string) exchange.CommandID {
	cmdID := exchange.CommandID(fmt.Sprintf("mock_close_%s_%s_%d", symbol, side, time.Now().UnixNano()))

	// go func() {
	// 	time.Sleep(m.randomDelay())

	// 	if m.shouldError() {
	// 		m.asyncMgr.StoreResult(cmdID, fmt.Errorf("simulated error closing position"))
	// 		return
	// 	}

	// 	result := map[string]interface{}{
	// 		"symbol": symbol,
	// 		"side":   side,
	// 		"status": "closed",
	// 		"time":   time.Now(),
	// 	}

	// 	m.asyncMgr.StoreResult(cmdID, result)
	// }()

	return cmdID
}

func (m *MockExchange) FetchBalanceAsync(asset string) exchange.CommandID {
	cmdID := exchange.CommandID(fmt.Sprintf("mock_balance_%s_%d", asset, time.Now().UnixNano()))

	// go func() {
	// 	time.Sleep(m.randomDelay())

	// 	if m.shouldError() {
	// 		m.asyncMgr.StoreResult(cmdID, fmt.Errorf("simulated balance fetch error"))
	// 		return
	// 	}

	// 	balance := exchange.Balance{
	// 		Asset:  asset,
	// 		Free:   1000.0,
	// 		Locked: 50.0,
	// 	}

	// 	m.asyncMgr.StoreResult(cmdID, balance)
	// }()

	return cmdID
}

//	func (m *MockExchange) GetResult(cmdID exchange.CommandID) (interface{}, bool) {
//		return m.asyncMgr.GetResult(cmdID)
//	}
func (m *MockExchange) PopCandle(cmdID exchange.CommandID) (exchange.Candle, bool, error) {

	if m.err != nil {
		return exchange.Candle{}, false, m.err
	}

	q, ok := m.dataQueues[cmdID]
	if !ok {
		return exchange.Candle{}, false, fmt.Errorf("Не нашёл очередь для %s", cmdID)
	}

	record, ok := q.PopOne()

	if !ok {
		return exchange.Candle{}, ok, nil
	}

	return record.Data, ok, nil
}

// ————————————————————————————————————————————————————————————————
// WebSocket имитация (опционально, если нужно тестировать стримы)
// ————————————————————————————————————————————————————————————————

type candleHandler struct {
	symbol   string
	interval string
	handler  func(exchange.Candle, error)
}

var mockCandleStreamHandlers []candleHandler
var mockCandleStreamMu sync.RWMutex

func (m *MockExchange) SubscribeCandles(symbol string, interval string) exchange.CommandID {
	mockCandleStreamMu.Lock()
	defer mockCandleStreamMu.Unlock()

	cmdID := exchange.GetCmdID("mock_candles", symbol, interval)

	handler := func(candle exchange.Candle, err error) {
		if err != nil {
			m.err = err
		} else {
			if m.dataQueues[cmdID] == nil {
				m.dataQueues[cmdID] = exchange.NewPriorityQueueManager[exchange.Candle]()
			}
			m.dataQueues[cmdID].PushBatch(&exchange.Record[exchange.Candle]{
				Timestamp: candle.Timestamp,
				Data:      candle,
			})
		}
	}

	mockCandleStreamHandlers = append(mockCandleStreamHandlers, candleHandler{
		symbol:   symbol,
		interval: interval,
		handler:  handler,
	})

	// Запускаем генерацию свечей в фоне, если ещё не запущена
	if len(mockCandleStreamHandlers) == 1 {
		go m.mockCandleStream(symbol)
	}

	return cmdID
}

func (m *MockExchange) UnsubscribeCandles(symbol string, interval string) {
	mockCandleStreamMu.Lock()
	defer mockCandleStreamMu.Unlock()

	var newHandlers []candleHandler
	for _, h := range mockCandleStreamHandlers {
		if h.symbol != symbol || h.interval != interval {
			newHandlers = append(newHandlers, h)
		}
	}
	mockCandleStreamHandlers = newHandlers
}

func (m *MockExchange) mockCandleStream(symbol string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mockCandleStreamMu.RLock()
		handlers := make([]candleHandler, len(mockCandleStreamHandlers))
		copy(handlers, mockCandleStreamHandlers)
		mockCandleStreamMu.RUnlock()

		if len(handlers) == 0 {
			return // больше подписчиков нет
		}

		candle := exchange.Candle{
			Symbol:    symbol, // можно рандомизировать
			Timestamp: time.Now(),
			Open:      50000 + rand.Float64()*100,
			High:      50100 + rand.Float64()*50,
			Low:       49900 + rand.Float64()*50,
			Close:     50050 + rand.Float64()*80,
			Volume:    100 + rand.Float64()*200,
		}

		for _, h := range handlers {
			h.handler(candle, nil)
		}
	}
}

var _ exchange.Exchange = (*MockExchange)(nil)
