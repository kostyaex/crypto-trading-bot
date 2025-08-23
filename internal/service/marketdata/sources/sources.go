package sources

import (
	"context"
	"crypto-trading-bot/internal/models"
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
	ctx  context.Context
}

func NewHistoricalSource(data []*models.MarketData, ctx context.Context) *HistoricalSource {
	return &HistoricalSource{
		data: data,
		ctx:  ctx,
	}
}

func (h *HistoricalSource) GetMarketDataCh() <-chan *models.MarketData {
	ch := make(chan *models.MarketData)
	go func() {
		defer close(ch)

		for _, item := range h.data {

			// помещаем в выходной канал данные с проверкой контекста на завершение выполнения
			select {
			case ch <- item:
			case <-h.ctx.Done():
				return
			}

		}

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
