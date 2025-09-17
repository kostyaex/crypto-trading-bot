package systems

import (
	"context"
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/exchange/exchanges/mockexchange"
	"time"

	"github.com/andygeiss/ecs"
)

type fetchDataSystem struct {
	ctx      context.Context
	exchange exchange.Exchange
}

func NewFetchDataSystem(ctx context.Context) ecs.System {

	s := &fetchDataSystem{
		ctx: ctx,
	}

	ex := mockexchange.NewMockExchange()
	ex.DelayMin = 10 * time.Millisecond
	ex.DelayMax = 50 * time.Millisecond

	s.exchange = ex

	return s
}

func (s *fetchDataSystem) Process(em ecs.EntityManager) (state int) {

	// выбираем все компоненты datasource, из них выбираем по каким символам и интервалам нужны данные.
	// Группируем и отправляем команды для выборки

	// grouped := make(map[string]exchange.FetchDataCommand)

	// for _, dataComp := range em.FilterByMask(components.MaskDatasource) {

	// 	datasource := dataComp.Get(components.MaskDatasource).(*components.DataSource)

	// 	key := datasource.Symbol + datasource.Interval
	// 	grouped[key] = exchange.FetchDataCommand{Symbol: datasource.Symbol, Interval: datasource.Interval}

	// }

	return ecs.StateEngineContinue
}

func (s *fetchDataSystem) Setup() {}

func (s *fetchDataSystem) Teardown() {}

// Проверка соответствия интерфейсу
var _ ecs.System = (*fetchDataSystem)(nil)
