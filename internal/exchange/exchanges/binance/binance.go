package binance

import (
	"crypto-trading-bot/internal/exchange"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type BinanceExchange struct {
	httpClient *http.Client
	wsConn     *websocket.Conn
	asyncMgr   *exchange.AsyncManager
	candlesSub map[string]func(exchange.Candle) // symbol_interval -> handler
	mu         sync.RWMutex
}

func NewBinanceExchange() *BinanceExchange {
	return &BinanceExchange{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		asyncMgr:   exchange.NewAsyncManager(5 * time.Minute),
		candlesSub: make(map[string]func(exchange.Candle)),
	}
}

// Пример асинхронной команды
func (b *BinanceExchange) FetchCandlesAsync(symbol, interval string, limit int) exchange.CommandID {
	cmdID := exchange.CommandID(fmt.Sprintf("candles_%s_%s_%d", symbol, interval, time.Now().UnixNano()))

	go func() {
		candles, err := b.fetchCandlesHTTP(symbol, interval, limit)
		if err != nil {
			b.asyncMgr.StoreResult(cmdID, err)
		} else {
			b.asyncMgr.StoreResult(cmdID, candles)
		}
	}()

	return cmdID
}

func (b *BinanceExchange) fetchCandlesHTTP(symbol, interval string, limit int) ([]exchange.Candle, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%d", symbol, interval, limit)
	resp, err := b.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawData [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		return nil, err
	}

	var candles []exchange.Candle
	for _, item := range rawData {
		candle := exchange.Candle{
			Symbol:    symbol,
			Timestamp: time.UnixMilli(int64(item[0].(float64))),
			Open:      toFloat64(item[1]),
			High:      toFloat64(item[2]),
			Low:       toFloat64(item[3]),
			Close:     toFloat64(item[4]),
			Volume:    toFloat64(item[5]),
		}
		candles = append(candles, candle)
	}
	return candles, nil
}

func toFloat64(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	default:
		return 0
	}
}

// Реализация GetResult
func (b *BinanceExchange) GetResult(cmdID exchange.CommandID) (interface{}, bool) {
	return b.asyncMgr.GetResult(cmdID)
}

// WebSocket подписка на свечи
func (b *BinanceExchange) SubscribeCandles(symbol, interval string, handler func(exchange.Candle)) {
	stream := fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
	wsURL := "wss://stream.binance.com:9443/ws/" + stream

	go b.connectWebSocket(wsURL, handler)
}

func (b *BinanceExchange) connectWebSocket(url string, handler func(exchange.Candle)) {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Printf("WebSocket error: %v", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		var msg struct {
			K struct {
				T int64  `json:"t"` // timestamp
				O string `json:"o"`
				H string `json:"h"`
				L string `json:"l"`
				C string `json:"c"`
				V string `json:"v"`
			} `json:"k"`
		}

		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		candle := exchange.Candle{
			Symbol:    "BTCUSDT", // можно парсить из stream name
			Timestamp: time.UnixMilli(msg.K.T),
			Open:      toFloat64(msg.K.O),
			High:      toFloat64(msg.K.H),
			Low:       toFloat64(msg.K.L),
			Close:     toFloat64(msg.K.C),
			Volume:    toFloat64(msg.K.V),
		}

		handler(candle)
	}
}

// Остальные методы PlaceOrderAsync, FetchBalanceAsync и т.д. — аналогично
