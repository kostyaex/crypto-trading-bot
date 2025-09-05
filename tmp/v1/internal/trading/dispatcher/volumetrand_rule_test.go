package dispatcher

import (
	"crypto-trading-bot/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_VolumeTrendRule_Buy(t *testing.T) {

	marketData := []types.MarketData{
		{Timestamp: time.Now(), ClosePrice: 30000, BuyVolume: 100, SellVolume: 50},
		{Timestamp: time.Now().Add(time.Minute), ClosePrice: 30100, BuyVolume: 200, SellVolume: 70},
	}

	seriesItem := types.Series{}
	for _, md := range marketData {
		seriesItem.Points = append(seriesItem.Points, types.Point{MarketData: &md})
	}

	rule := &VolumeTrendRule{MinVolumeChangePercent: 10}
	signal, ok := rule.Evaluate(&seriesItem)

	assert.True(t, ok)
	assert.Equal(t, SignalBuy, signal.Type)
	assert.Equal(t, 30100.0, signal.Price)
}

func Test_VolumeTrendRule_Sell(t *testing.T) {
	marketData := []types.MarketData{
		{Timestamp: time.Now(), ClosePrice: 30000, BuyVolume: 100, SellVolume: 50},
		{Timestamp: time.Now().Add(time.Minute), ClosePrice: 29900, BuyVolume: 80, SellVolume: 200},
	}

	seriesItem := types.Series{}
	for _, md := range marketData {
		seriesItem.Points = append(seriesItem.Points, types.Point{MarketData: &md})
	}

	rule := &VolumeTrendRule{MinVolumeChangePercent: 10}
	signal, ok := rule.Evaluate(&seriesItem)

	assert.True(t, ok)
	assert.Equal(t, SignalSell, signal.Type)
	assert.Equal(t, 29900.0, signal.Price)
}

func Test_VolumeTrendRule_Hold_NoSignificantChange(t *testing.T) {
	marketData := []types.MarketData{
		{Timestamp: time.Now(), ClosePrice: 30000, BuyVolume: 100, SellVolume: 50},
		{Timestamp: time.Now().Add(time.Minute), ClosePrice: 30050, BuyVolume: 110, SellVolume: 55},
	}

	seriesItem := types.Series{}
	for _, md := range marketData {
		seriesItem.Points = append(seriesItem.Points, types.Point{MarketData: &md})
	}

	rule := &VolumeTrendRule{MinVolumeChangePercent: 20}
	_, ok := rule.Evaluate(&seriesItem)

	assert.False(t, ok)
}
