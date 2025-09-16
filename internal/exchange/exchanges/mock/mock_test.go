package mock

import (
	"crypto-trading-bot/internal/exchange"
	"testing"
	"time"
)

func TestMockExchange_FetchCandlesAsync(t *testing.T) {
	ex := NewMockExchange()
	ex.DelayMin = 10 * time.Millisecond
	ex.DelayMax = 50 * time.Millisecond

	cmdID := ex.FetchCandlesAsync("BTCUSDT", "1m", 5)

	// Сразу результат не готов
	if result, ok := ex.GetResult(cmdID); ok {
		t.Errorf("Expected no result immediately, got: %+v", result)
	}

	// Ждём немного
	time.Sleep(100 * time.Millisecond)

	// Теперь должен быть результат
	if result, ok := ex.GetResult(cmdID); !ok {
		t.Fatal("Expected result after delay, got nothing")
	} else {
		candles, ok := result.([]exchange.Candle)
		if !ok {
			t.Fatalf("Expected []Candle, got %T", result)
		}
		if len(candles) != 5 {
			t.Errorf("Expected 5 candles, got %d", len(candles))
		}
	}
}

func TestMockExchange_PlaceOrderAsync_ErrorSimulation(t *testing.T) {
	ex := NewMockExchange()
	ex.DelayMin = 10 * time.Millisecond
	ex.DelayMax = 30 * time.Millisecond
	ex.ErrRate = 1.0 // всегда ошибка

	cmdID := ex.PlaceOrderAsync(exchange.Order{
		Symbol: "BTCUSDT",
		Side:   "buy",
		Type:   "market",
		Amount: 0.01,
	})

	time.Sleep(50 * time.Millisecond)

	if result, ok := ex.GetResult(cmdID); !ok {
		t.Fatal("Expected error result, but got nothing")
	} else {
		if _, isErr := result.(error); !isErr {
			t.Fatalf("Expected error, got %+v", result)
		}
	}
}

func TestMockExchange_SubscribeCandles(t *testing.T) {
	ex := NewMockExchange()

	received := 0
	ex.SubscribeCandles("BTCUSDT", "1m", func(c exchange.Candle) {
		received++
	})

	// Дадим стриму поработать
	time.Sleep(3 * time.Second)
	ex.UnsubscribeCandles("BTCUSDT", "1m")

	if received == 0 {
		t.Error("Expected to receive at least one candle via WebSocket mock")
	}
	t.Logf("Received %d mock candles", received)
}
