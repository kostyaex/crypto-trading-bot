package trader

import (
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/core/utils"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/service/clusters"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/service/marketdata/sources"
	"crypto-trading-bot/internal/service/series"
	"crypto-trading-bot/internal/trading/dispatcher"
	"fmt"
	"time"
)

type TraderService interface {
	RunBacktesting(strategy *models.Strategy, startTime, endTime time.Time)
}

type traderService struct {
	conf              *config.Config
	logger            *logger.Logger
	marketDataService marketdata.MarketDataService
}

func NewTraderService(conf *config.Config, logger *logger.Logger, marketDataService marketdata.MarketDataService) TraderService {
	return &traderService{
		conf:              conf,
		logger:            logger,
		marketDataService: marketDataService,
	}
}

type MockMarketDataSource struct {
	data []*models.MarketData
}

func NewMockMarketDataSource(data []*models.MarketData) *MockMarketDataSource {
	return &MockMarketDataSource{data: data}
}

func (m *MockMarketDataSource) GetMarketDataCh() <-chan *models.MarketData {
	ch := make(chan *models.MarketData)
	go func() {
		for _, item := range m.data {
			ch <- item
		}
		close(ch)
	}()
	return ch
}

func (m *MockMarketDataSource) Close() {}

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
	backtestContext *BacktestContext,
) error {
	defer source.Close()

	strategySettings, err := strategy.Settings()
	if err != nil {
		return err
	}

	//marketDataCh := source.GetMarketDataCh()

	// дублируем канал для обработки и сбора статистики
	broadcaster := marketdata.NewBroadcaster(source.GetMarketDataCh())
	broadcaster.Start()

	marketDataCh1 := broadcaster.Subscribe()
	marketDataCh2 := broadcaster.Subscribe()

	// Разбиваем полученные торговые данные на интевалы по настройкам из стратегии
	intervalsCh := make(chan []*models.MarketData)
	go func() {
		utils.SplitChannelWithOverlap(marketDataCh1, strategySettings.Cluster.Block, 0, intervalsCh)
		//close(intervalsCh)
	}()

	go func() {
		for md := range marketDataCh2 {
			backtestContext.collectMarketData(md)
		}
	}()

	//broadcaster.Wait()

	builder, err := series.NewSeriesBuilder(strategy.SeriesBuilderConfig)
	if err != nil {
		panic(err)
	}

	var activeSeries []series.Series

	for interval := range intervalsCh {
		clusteredMd := clusters.ClusterMarketData(interval, strategySettings.Cluster.Interval, strategySettings.Cluster.NumClusters)
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

		// здесь надо отфильтровать серии, выбрать только серии которые были обновлены на этой итерации
		for _, sr := range activeSeries {
			// проверяем последнюю точку серии
			last := sr.Last()
			if last == nil || last.Time.Before(clusteredMd[0].Timestamp) {
				continue
			}

			dispatcher.Dispatch(&sr)
		}
	}

	//filename := fmt.Sprintf("/home/kostya/projects/crypto-trading-bot/data/series/series_%s.json", time.Now().Format("2006-01-02_15-04-05"))
	//series.SaveSeries(activeSeries, filename)

	// Сбор метрик
	metrics := series.CollectMetrics(activeSeries)
	metrics.Print()

	return nil
}
