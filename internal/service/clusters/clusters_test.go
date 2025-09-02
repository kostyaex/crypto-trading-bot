package clusters

import (
	"crypto-trading-bot/pkg/types"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKMeansWeighted1D(t *testing.T) {

	points := []WeightedPoint{
		{Value: 1, Weight: 3},
		{Value: 2, Weight: 2},
		{Value: 3, Weight: 1},
		{Value: 10, Weight: 4},
		{Value: 11, Weight: 3},
		{Value: 12, Weight: 2},
	}

	clusters := KMeansWeighted1D(points, 2, 100)
	for i, c := range clusters {
		fmt.Printf("Кластер %d:\n", i+1)
		fmt.Printf("Центр: %.2f\n", c.Center)
		fmt.Println("Точки:")
		for _, p := range c.Points {
			fmt.Printf("  %.2f вес: %.2f\n", p.Value, p.Weight)
		}
		fmt.Println()
	}
}

func Test_ClusterMarketData(t *testing.T) {

	// Подготовка тестовых данных
	now, _ := time.Parse(time.RFC3339, "2025-01-05T00:00:00Z") //time.Now()
	testData := []*types.MarketData{
		{Timestamp: now, TimeFrame: "1m", ClosePrice: 100, Volume: 10, BuyVolume: 6, SellVolume: 4},
		{Timestamp: now.Add(time.Minute), TimeFrame: "1m", ClosePrice: 105, Volume: 20, BuyVolume: 22, SellVolume: 18},
		{Timestamp: now.Add(2 * time.Minute), TimeFrame: "1m", ClosePrice: 103, Volume: 15, BuyVolume: 7, SellVolume: 8},
		{Timestamp: now.Add(3 * time.Minute), TimeFrame: "1m", ClosePrice: 107, Volume: 25, BuyVolume: 10, SellVolume: 15},
		{Timestamp: now.Add(4 * time.Minute), TimeFrame: "1m", ClosePrice: 108, Volume: 30, BuyVolume: 14, SellVolume: 16},
	}

	results := ClusterMarketData(testData, "1h", 3)

	fmt.Printf("%s\n", types.MarketDataToString(testData))
	fmt.Printf("%s\n", types.MarketDataToString(results))

	assert.NotEmpty(t, results)
	// assert.Contains(t, results[0].Log, "Waves:")
}
