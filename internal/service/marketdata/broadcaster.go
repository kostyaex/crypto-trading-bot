package marketdata

import (
	"crypto-trading-bot/internal/types"
	"sync"
)

// Мультикастер для распределения потока торговых данных по сгрупированным стратегиям.

type Broadcaster struct {
	subscribers []chan *types.MarketData
	source      <-chan *types.MarketData
	wg          sync.WaitGroup
}

func NewBroadcaster(source <-chan *types.MarketData) *Broadcaster {
	return &Broadcaster{
		subscribers: make([]chan *types.MarketData, 0),
		source:      source,
	}
}

func (b *Broadcaster) Subscribe() <-chan *types.MarketData {
	ch := make(chan *types.MarketData, 100)
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
