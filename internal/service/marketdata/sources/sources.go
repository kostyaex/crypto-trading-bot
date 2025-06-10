package sources

import "crypto-trading-bot/internal/models"

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
