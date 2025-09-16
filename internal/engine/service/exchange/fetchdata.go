package exchange

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Config представляет настройку для одной пары/интервала
type FetchDataCommand struct {
	Symbol   string
	Interval string
}

// Result — структура, которую возвращает goroutine
type FetchDataResult struct {
	Command   FetchDataCommand
	Data      []byte // или ваш тип данных: []OHLCV, map[string]interface{}, и т.д.
	Error     error
	Timestamp time.Time
}

// управляет сбором данных от всех goroutines
type FetchDataExecutor struct {
	results  map[FetchDataCommand]*FetchDataResult
	mu       sync.RWMutex
	commands chan FetchDataCommand // входной канал для новых команд
	done     chan struct{}         // сигнал завершения
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewFetchDataExecutor(ctx context.Context) *FetchDataExecutor {
	ctx, cancel := context.WithCancel(ctx)
	te := &FetchDataExecutor{
		results:  make(map[FetchDataCommand]*FetchDataResult),
		commands: make(chan FetchDataCommand, 100), // буферизованный, чтобы не блокировать отправителя
		done:     make(chan struct{}),
		ctx:      ctx,
		cancel:   cancel,
	}

	go te.worker() // запускаем внутренний воркер
	return te
}

// SubmitCommand — безопасный способ отправить команду из любого места
func (te *FetchDataExecutor) SubmitCommand(cmd FetchDataCommand) {
	select {
	case te.commands <- cmd:
		// Команда принята
	default:
		fmt.Printf("⚠️ Trade queue full, command dropped: %v\n", cmd)
	}
}

// worker — внутренняя горутина, которая обрабатывает команды
func (te *FetchDataExecutor) worker() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-te.ctx.Done():
			fmt.Println("🛑 TradeExecutor shutting down...")
			close(te.done)
			return

		case cmd := <-te.commands:
			// Запускаем горутину для каждой команды — НЕ БЛОКИРУЕМ worker!
			go te.executeCommand(cmd)

			// case <-ticker.C:
			// 	// Периодическая проверка состояния ордеров (опционально)
			// 	// Например: опрос API на статус pending-ордеров
			// 	te.checkPendingOrders()
		}
	}
}

// executeTrade — выполняет одну торговую операцию (в своей goroutine)
func (te *FetchDataExecutor) executeCommand(cmd FetchDataCommand) {
	// Имитируем вызов биржевого API
	// В реальности: client.PlaceOrder(...)
	time.Sleep(200 * time.Millisecond) // задержка сети

	var result FetchDataResult
	result.Command = cmd
	result.Timestamp = time.Now()

	// // Симуляция успеха/ошибки
	// if cmd.Amount < 0.001 {
	// 	result.Status = "rejected"
	// 	result.Error = fmt.Errorf("min order size is 0.001")
	// } else if cmd.Pair == "SOL/USDT" && cmd.Side == "sell" {
	// 	// Искусственная ошибка для теста
	// 	result.Status = "rejected"
	// 	result.Error = fmt.Errorf("sell blocked for SOL/USDT due to risk policy")
	// } else {
	// 	// Успешное исполнение
	// 	result.OrderID = fmt.Sprintf("ord_%d", time.Now().UnixNano())
	// 	result.ExecutedQty = cmd.Amount
	// 	result.Price = 100.0 // имитация цены
	// 	result.Status = "filled"
	// }

	// Сохраняем результат — безопасно
	te.mu.Lock()
	te.results[cmd] = &result
	te.mu.Unlock()

	// Можно также отправить событие в канал для внешнего слушателя (если нужно)
	// Например: te.events <- result
}

// GetResult — неблокирующий способ получить результат по OrderID
func (te *FetchDataExecutor) GetResult(cmd FetchDataCommand) (*FetchDataResult, bool) {
	te.mu.RLock()
	defer te.mu.RUnlock()
	res, ok := te.results[cmd]
	return res, ok
}

// GetAllResults — получить все результаты (для мониторинга)
func (te *FetchDataExecutor) GetAllResults() map[FetchDataCommand]*FetchDataResult {
	te.mu.RLock()
	defer te.mu.RUnlock()

	result := make(map[FetchDataCommand]*FetchDataResult, len(te.results))
	for k, v := range te.results {
		result[k] = v
	}
	return result
}

// Close — корректное завершение
func (te *FetchDataExecutor) Close() {
	te.cancel()
	<-te.done
}

// func NewDataCollector() *FetchDataExecutor {
// 	return &FetchDataExecutor{
// 		results:  make(map[FetchDataCommand][]byte),
// 		channels: make(map[FetchDataCommand]chan FetchDataResult),
// 	}
// }

// // StartWorker запускает goroutine для конкретной настройки
// func (dc *FetchDataExecutor) StartWorker(ctx context.Context, config FetchDataCommand) {
// 	ch := make(chan FetchDataResult, 1) // буферизованный канал, чтобы не блокировать горутину
// 	dc.channels[config] = ch

// 	go func() {
// 		defer close(ch)
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 				// Здесь имитируем получение данных с биржи
// 				// Замените на реальный вызов API
// 				data := []byte(fmt.Sprintf("data for %s/%s", config.Symbol, config.Interval))
// 				result := FetchDataResult{
// 					Config:    config,
// 					Data:      data,
// 					Timestamp: time.Now(),
// 				}

// 				select {
// 				case ch <- result:
// 					// Успешно отправили
// 				default:
// 					// Канал переполнен — пропускаем, если клиент не успел забрать
// 					// Можно логировать или увеличить буфер
// 				}

// 				time.Sleep(2 * time.Second) // имитация задержки получения данных
// 			}
// 		}
// 	}()
// }

// // CollectNonBlocking — основной цикл, который НЕ БЛОКИРУЕТСЯ
// func (dc *FetchDataExecutor) CollectNonBlocking() {
// 	for config, ch := range dc.channels {
// 		select {
// 		case result, ok := <-ch:
// 			if !ok {
// 				// Канал закрыт — можно удалить из коллекции, если нужно
// 				delete(dc.channels, config)
// 				fmt.Printf("Channel for %s/%s closed\n", config.Symbol, config.Interval)
// 				continue
// 			}

// 			// Сохраняем результат
// 			dc.mu.Lock()
// 			dc.results[config] = result.Data
// 			dc.mu.Unlock()

// 			fmt.Printf("✅ Got result for %s/%s at %v\n", config.Symbol, config.Interval, result.Timestamp)
// 		default:
// 			// Нет данных — ничего не делаем, продолжаем цикл
// 		}
// 	}
// }

// // GetResults — безопасный способ получить текущие результаты
// func (dc *FetchDataExecutor) GetResults() map[FetchDataCommand][]byte {
// 	dc.mu.RLock()
// 	defer dc.mu.RUnlock()

// 	result := make(map[FetchDataCommand][]byte, len(dc.results))
// 	for k, v := range dc.results {
// 		result[k] = v
// 	}
// 	return result
// }
