package dispatcher

import "crypto-trading-bot/internal/service/series"

type VolumeTrendRule struct {
	MinVolumeChangePercent float64 // минимальный процент изменения объема
}

func (r *VolumeTrendRule) Evaluate(series *series.Series) (TradeSignal, bool) {
	if len(series.Points) < 2 {
		return TradeSignal{Type: SignalHold}, false
	}

	first := series.Points[0]
	last := series.Points[len(series.Points)-1]

	buyVolChange := (last.MarketData.BuyVolume - first.MarketData.BuyVolume) / first.MarketData.BuyVolume * 100
	sellVolChange := (last.MarketData.SellVolume - first.MarketData.SellVolume) / first.MarketData.SellVolume * 100

	signal := TradeSignal{
		Timestamp:  last.MarketData.Timestamp,
		Symbol:     last.MarketData.Symbol,
		Price:      last.MarketData.ClosePrice,
		Volume:     last.MarketData.Volume,
		BuyVolume:  last.MarketData.BuyVolume,
		SellVolume: last.MarketData.SellVolume,
		Series:     series,
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
