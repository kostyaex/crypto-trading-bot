package dispatcher

import (
	"crypto-trading-bot/internal/service/series"
	"fmt"
)

type VolumeTrendRule struct {
	MinVolumeChangePercent float64 // минимальный процент изменения объема
}

func (r *VolumeTrendRule) Evaluate(series *series.Series) (TradeSignal, bool) {
	if len(series.Points) < 2 {
		return TradeSignal{Type: SignalHold}, false
	}

	// берем последние 2 значения в серии
	prev := series.Points[len(series.Points)-2]
	last := series.Points[len(series.Points)-1]

	buyVolChange := (last.MarketData.BuyVolume - prev.MarketData.BuyVolume) / prev.MarketData.BuyVolume * 100
	sellVolChange := (last.MarketData.SellVolume - prev.MarketData.SellVolume) / prev.MarketData.SellVolume * 100

	signal := TradeSignal{
		Timestamp:   last.MarketData.Timestamp,
		Symbol:      last.MarketData.Symbol,
		Price:       last.MarketData.ClusterPrice,
		Volume:      last.MarketData.Volume,
		BuyVolume:   last.MarketData.BuyVolume,
		SellVolume:  last.MarketData.SellVolume,
		Series:      series,
		Description: "volumetrandrule | " + series.String(),
	}

	// Проверяем условия
	if buyVolChange > r.MinVolumeChangePercent && buyVolChange > sellVolChange {
		signal.Type = SignalBuy
		return signal, true
	}

	if sellVolChange > r.MinVolumeChangePercent && sellVolChange > buyVolChange {
		signal.Type = SignalSell
		return signal, true
	}

	return TradeSignal{Type: SignalHold}, false
}

func NewVolumeTrendRule(config map[string]interface{}) (SignalRule, error) {
	min_volume_change_percent, ok := config["min_volume_change_percent"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'min_volume_change_percent'")
	}
	return &VolumeTrendRule{MinVolumeChangePercent: min_volume_change_percent}, nil
}
