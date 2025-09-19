package mockexchange

import (
	"fmt"
	"testing"
	"time"
)

func TestMockExchange_FetchCandlesAsync(t *testing.T) {
	ex := NewMockExchange()
	ex.DelayMin = 10 * time.Millisecond
	ex.DelayMax = 50 * time.Millisecond

	cmdID := ex.FetchCandlesAsync("BTCUSDT", "1m", 5)

	// Сразу результат не готов
	if result, ok, _ := ex.PopCandle(cmdID); ok {
		t.Errorf("Expected no result immediately, got: %+v", result)
	}

	// Ждём немного
	time.Sleep(100 * time.Millisecond)

	// Теперь должен быть результат
	for i := 0; i < 5; i++ {
		candle, ok, _ := ex.PopCandle(cmdID)
		if !ok {
			t.Fatal("Пустой или не полный результат")
		}
		fmt.Printf("%v\n", candle)
	}

}

// func TestMockExchange_PlaceOrderAsync_ErrorSimulation(t *testing.T) {
// 	ex := NewMockExchange()
// 	ex.DelayMin = 10 * time.Millisecond
// 	ex.DelayMax = 30 * time.Millisecond
// 	ex.ErrRate = 1.0 // всегда ошибка

// 	cmdID := ex.PlaceOrderAsync(exchange.Order{
// 		Symbol: "BTCUSDT",
// 		Side:   "buy",
// 		Type:   "market",
// 		Amount: 0.01,
// 	})

// 	time.Sleep(50 * time.Millisecond)

// 	if result, ok := ex.GetResult(cmdID); !ok {
// 		t.Fatal("Expected error result, but got nothing")
// 	} else {
// 		if _, isErr := result.(error); !isErr {
// 			t.Fatalf("Expected error, got %+v", result)
// 		}
// 	}
// }

func TestMockExchange_SubscribeCandles(t *testing.T) {
	ex := NewMockExchange()

	cmdID := ex.SubscribeCandles("BTCUSDT", "1s")

	num := 3
	// Дадим стриму поработать
	time.Sleep(time.Duration(num) * time.Second)

	result, ok, err := ex.PopCandle(cmdID)

	if err != nil {
		t.Fatalf("%v", err)
	}

	// Теперь должен быть результат
	if !ok {
		t.Fatal("Expected result after delay, got nothing")
	}

	fmt.Printf("%v\n", result)

	ex.UnsubscribeCandles("BTCUSDT", "1m")

}
