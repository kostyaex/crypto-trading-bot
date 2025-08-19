package trader

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/service/series"
	"encoding/json"
	"os"
)

type BacktestContext struct {
	//StrategyParams map[string]interface{} `json:"strategy_params"`
	StrategySettings *models.StrategySettings
	//IntermediateResults []map[string]interface{} `json:"intermediate_results"`
	Timestamp           string               `json:"timestamp"`
	MarketData          []*models.MarketData `json:"marketdata"`
	ClusteredMarketData []*models.MarketData `json:"clustered_marketdata"`
	SeriesList          []*series.Series     `json:"series_list"`
}

// func (s *traderService) RunBacktesting(strategy *models.Strategy, startTime, endTime time.Time) {

// 	strategySettings, err := strategy.Settings()
// 	if err != nil {
// 		s.logger.Errorf("Ошибка создания новой стратегии %v", err)
// 		return
// 	}

// 	ctx := &BacktestContext{
// 		StrategySettings: strategySettings,
// 		Timestamp:        time.Now().Format("2006-01-02T15:04:05Z"),
// 	}

// 	// Получить данные за период
// 	marketData, err := s.marketDataService.GetMarketDataPeriod(strategySettings.Symbol, strategySettings.Interval, startTime, endTime)
// 	//	marketData, err := s.GetMarketData("BTCUSDT", 1000)
// 	if err != nil {
// 		s.logger.Errorf("Ошибка получения торговых данных: %s\n", err)
// 		return
// 	}

// 	source := sources.NewMockMarketDataSource(marketData)

// 	disp := dispatcher.NewDispatcher(
// 		&dispatcher.VolumeTrendRule{MinVolumeChangePercent: 10},
// 	)

// 	disp.Register(dispatcher.SignalBuy, &dispatcher.LoggerHandler{})
// 	disp.Register(dispatcher.SignalSell, &dispatcher.LoggerHandler{})
// 	//disp.Register(dispatcher.SignalHold, &dispatcher.LoggerHandler{})

// 	// Вызов тестируемой функции
// 	err = runStrategyForSource(*strategy, source, disp, ctx)
// 	if err != nil {
// 		s.logger.Errorf("Ошибка выполнения RunStrategyForSource %v", err)
// 		return
// 	}

// 	// Сохранение контекста в JSON
// 	filepath := filepath.Join(s.conf.Data.Dir, "backtests", strategySettings.Symbol+"_"+strategySettings.Interval+"_"+ctx.Timestamp+".json")
// 	if err := ctx.saveBacktestContext(filepath); err != nil {
// 		s.logger.Errorf("Ошибка выполнения saveBacktestContext %v", err)
// 	} else {
// 		s.logger.Printf("Сохоанен файл - %s", filepath)
// 	}
// }

// func (ctx *BacktestContext) collectMarketData(md *models.MarketData) {
// 	ctx.MarketData = append(ctx.MarketData, md)
// }

// func (ctx *BacktestContext) collectClusteredMarketData(md *models.MarketData) {
// 	ctx.ClusteredMarketData = append(ctx.ClusteredMarketData, md)
// }

func (ctx *BacktestContext) saveBacktestContext(filePath string) error {
	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
