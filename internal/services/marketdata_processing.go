package services

import (
	"crypto-trading-bot/internal/models"
	"fmt"
	"sync"
)

// Источник данных (база данных или загруженные с биржи данные)
type MarketDataSource interface {
	GetMarketDataCh() <-chan *models.MarketData
	Close()
}

// ---------------------------------------------------------------------------
// источник исторических данных из БД
type HistoricalSource struct {
	data []*models.MarketData
}

func NewHistoricalSource(data []*models.MarketData) *HistoricalSource {
	return &HistoricalSource{data: data}
}

func (h *HistoricalSource) GetMarketDataCh() <-chan *models.MarketData {
	ch := make(chan *models.MarketData)
	go func() {
		for _, item := range h.data {
			ch <- item
		}
		close(ch)
	}()
	return ch
}

func (h *HistoricalSource) Close() {}

// ---------------------------------------------------------------------------
// источник реальных данных, загруженных с биржи

// type LiveMarketSource struct {
// 	wsClient *WebSocketClient // Предположим, есть клиент биржи
// }

// func NewLiveMarketSource(client *WebSocketClient) *LiveMarketSource {
// 	return &LiveMarketSource{wsClient: client}
// }

// func (l *LiveMarketSource) GetMarketDataCh() <-chan *models.MarketData {
// 	return l.wsClient.SubscribeToTrades()
// }

// func (l *LiveMarketSource) Close() {
// 	l.wsClient.Unsubscribe()
// }

// ---------------------------------------------------------------------------
// Мультикастер для распределения потока торговых данных по сгрупированным стратегиям.

type Broadcaster struct {
	subscribers []chan *models.MarketData
	source      <-chan *models.MarketData
	wg          sync.WaitGroup
}

func NewBroadcaster(source <-chan *models.MarketData) *Broadcaster {
	return &Broadcaster{
		subscribers: make([]chan *models.MarketData, 0),
		source:      source,
	}
}

func (b *Broadcaster) Subscribe() <-chan *models.MarketData {
	ch := make(chan *models.MarketData, 100)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *Broadcaster) Start() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for item := range b.source {
			for _, ch := range b.subscribers {
				ch <- item
			}
		}

		// Закрываем все подписки после окончания потока
		for _, ch := range b.subscribers {
			close(ch)
		}
	}()
}

func (b *Broadcaster) Wait() {
	b.wg.Wait()
}

// ---------------------------------------------------------------------------

// Функция распределения данных по группам
func groupStrategiesBySymbolInterval(strategies []models.Strategy) map[string][]models.Strategy {
	grouped := make(map[string][]models.Strategy)
	for _, strategy := range strategies {
		settings, err := strategy.Settings()
		if err != nil {
			panic("Не удалось получить параметры стратегии")
		}

		key := fmt.Sprintf("%s|%s", settings.Symbol, settings.Interval)
		grouped[key] = append(grouped[key], strategy)
	}
	return grouped
}
