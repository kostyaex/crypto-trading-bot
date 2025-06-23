package trader

import (
	"crypto-trading-bot/internal/core/utils"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/service/clusters"
	"crypto-trading-bot/internal/service/marketdata/sources"
	"crypto-trading-bot/internal/service/series"
	"crypto-trading-bot/internal/trading/dispatcher"
	"fmt"
)

// Функция распределения данных по группам
func groupStrategiesBySymbolInterval(strategies []models.Strategy) map[string][]models.Strategy {
	grouped := make(map[string][]models.Strategy)
	for _, strategy := range strategies {
		settings, err := strategy.Settings()
		if err != nil {
			panic("Не удалось получить параметры стратегии")
		}

		key := fmt.Sprintf("%s|%s", settings.Symbol, settings.Interval)
		grouped[key] = append(grouped[key], strategy)
	}
	return grouped
}

func runStrategyForSource(
	strategy models.Strategy,
	source sources.MarketDataSource,
	dispatcher *dispatcher.Dispatcher,
) error {
	defer source.Close()

	strategySettings, err := strategy.Settings()
	if err != nil {
		return err
	}

	marketDataCh := source.GetMarketDataCh()

	// Разбиваем полученные торговые данные на интевалы по настройкам из стратегии
	intervalsCh := make(chan []*models.MarketData)
	go func() {
		utils.SplitChannelWithOverlap(marketDataCh, strategySettings.Waves.BlockSize, strategySettings.Waves.Overlap, intervalsCh)
		//close(intervalsCh)
	}()

	// Пример конфигурации
	config := map[string]interface{}{
		"type":         "simple",
		"value_factor": 1.0,
		"time_factor":  100.0,
	}

	builder, err := series.NewSeriesBuilder(config)
	if err != nil {
		panic(err)
	}

	var activeSeries []series.Series

	for interval := range intervalsCh {
		clusteredMd := clusters.ClusterMarketData(interval, "1h", 5)
		var points []series.Point
		for _, md := range clusteredMd {
			point := series.Point{
				Value:      md.ClusterPrice,
				Weight:     md.Volume,
				Time:       md.Timestamp,
				MarketData: md,
			}
			points = append(points, point)
		}

		activeSeries = builder.AddPoints(activeSeries, points)

		for _, sr := range activeSeries {
			dispatcher.Dispatch(&sr)
		}
	}

	return nil
}
