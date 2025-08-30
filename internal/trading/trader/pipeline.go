package trader

import (
	"context"
	"crypto-trading-bot/internal/core/config"
	"crypto-trading-bot/internal/core/logger"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/service/marketdata/sources"
	"crypto-trading-bot/internal/service/series"
	"crypto-trading-bot/internal/trading/dispatcher"
	"path/filepath"
	"time"
)

type Pipeline struct {
	Conf          *config.Config
	Logger        *logger.Logger
	Mode          string // "live" or "backtest"
	DataSource    sources.MarketDataSource
	Strategy      models.Strategy
	SeriesBuilder *series.SeriesBuilder
	Dispatcher    *dispatcher.Dispatcher
}

// Обновить правила
func (p *Pipeline) UpdateDispatcher(Dispatcher *dispatcher.Dispatcher) {
	p.Dispatcher = Dispatcher
}

func (p *Pipeline) Run(ctx context.Context) error {

	//var currentTime int64
	var backtestConext *BacktestContext

	defer p.DataSource.Close()

	strategySettings := p.Strategy.Settings

	if p.Mode == "backtest" {
		backtestConext = &BacktestContext{
			StrategySettings: &strategySettings,
			Timestamp:        time.Now().Format("2006-01-02T15:04:05Z"),
		}
	}

	aggregator := marketdata.NewAggregator(strategySettings.Cluster.Block)

	// ---------------------------------------------------------------------

	// Выбираем данные из источника (через канал)
	for marketData := range p.DataSource.GetMarketDataCh() {

		// // Этап 1. Получение данных по текущему интервалу (время, цена)
		// // Также получение текущего PNL по открытым позициям
		// if p.Mode == "live" {
		// 	currentTime = time.Now().Unix()
		// } else {
		// 	currentTime = marketData.Timestamp.Unix()

		// }

		// p.Logger.Debugf("CURRENT TIME: %s\n", time.Unix(currentTime, 0).Format(utils.TimeFormatForHuman))

		// ---------------------------------------------------------------------

		// Этап 1. Группировка интервалов в массив для последующей кластеризации
		aggregator.Add(marketData)
		if !aggregator.IsReady() {
			continue
		}

		//p.Logger.Debugf("Aggregated data: %d\n", len(aggregator.GetAggregatedData()))

		// ---------------------------------------------------------------------

		// Этап 2. Кластеризация
		// здесь сворачиваем данные в кластеры. Т.е. к примеру данные за секундный интервал в 5 минутный, получим столько значений, сколько указано количество кластеров.
		clusteredMd := aggregator.GetClusteredData(strategySettings.Cluster.NumClusters, strategySettings.Cluster.Interval)

		// // Собираем данные для статистики
		// for _, md := range clusteredMd {
		// 	backtestContext.collectClusteredMarketData(md)
		// }

		// ---------------------------------------------------------------------

		// Этап 3. Построение серий
		p.SeriesBuilder.AddClusteredData(clusteredMd)
		activeSeries := p.SeriesBuilder.GetActiveSeries()

		// ---------------------------------------------------------------------

		// Этап 4. Обработка серий

		for _, sr := range activeSeries {
			p.Dispatcher.Dispatch(sr)
		}
	}

	// ---------------------------------------------------------------------

	if p.Mode == "backtest" {

		backtestConext.SeriesList = p.SeriesBuilder.GetSeries()

		// Сбор метрик
		metrics := series.CollectMetrics(backtestConext.SeriesList)
		metrics.Print()

		// Сохранение контекста в JSON
		filepath := filepath.Join(p.Conf.Data.Dir, "backtests", strategySettings.Symbol+"_"+strategySettings.Interval+"_"+backtestConext.Timestamp+".json")
		if err := backtestConext.saveBacktestContext(filepath); err != nil {
			p.Logger.Errorf("Ошибка выполнения saveBacktestContext %v", err)
		} else {
			p.Logger.Printf("Сохранен файл - %s", filepath)
		}
	}

	return nil
}
