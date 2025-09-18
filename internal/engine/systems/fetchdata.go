package systems

import (
	"context"
	"crypto-trading-bot/internal/engine/components"
	"crypto-trading-bot/internal/engine/ecsx"
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/exchange/exchanges/mockexchange"
	"fmt"
	"time"

	"github.com/andygeiss/ecs"
)

type fetchDataSystem struct {
	em         ecs.EntityManager
	ctx        context.Context
	exchange   exchange.Exchange
	subscribes map[string]int // map[Symbol+Interval]count - кол-во подписок на данный символ и интервал
}

func NewFetchDataSystem(ctx context.Context, em ecs.EntityManager) *fetchDataSystem {

	s := &fetchDataSystem{
		em:  em,
		ctx: ctx,
	}

	ex := mockexchange.NewMockExchange()
	ex.DelayMin = 10 * time.Millisecond
	ex.DelayMax = 50 * time.Millisecond

	s.exchange = ex

	return s
}

func (s *fetchDataSystem) Process(em ecs.EntityManager) (state int) {

	return ecs.StateEngineContinue
}

func (s *fetchDataSystem) Setup() {
	s.subscribes = make(map[string]int)
}

func (s *fetchDataSystem) Teardown() {
	//fmt.Println("[fetchDataSystem | Teardown]")
	for _, dataComp := range s.em.FilterByMask(components.MaskDatasource) {

		datasource := dataComp.Get(components.MaskDatasource).(*components.DataSource)

		s.subscribes[dataSourceKey(datasource)]--
		if s.subscribes[dataSourceKey(datasource)] == 0 {
			fmt.Printf("[fetchDataSystem | Teardown] Unsubscribe %s %s\n", datasource.Symbol, datasource.Interval)
		}

	}
}

func (s *fetchDataSystem) OnEntityAdded(entity *ecs.Entity, _components []ecs.Component) {

	datasource := entity.Get(components.MaskDatasource).(*components.DataSource)
	if s.subscribes[dataSourceKey(datasource)] == 0 {
		fmt.Printf("[fetchDataSystem] Subscribe %s %s\n", datasource.Symbol, datasource.Interval)
	} else {
		fmt.Printf("[fetchDataSystem] Alredy subscribed %s %s\n", datasource.Symbol, datasource.Interval)
	}
	s.subscribes[dataSourceKey(datasource)]++

}

func (s *fetchDataSystem) OnEntityRemoved(entity *ecs.Entity) {
	datasource := entity.Get(components.MaskDatasource).(*components.DataSource)
	if s.subscribes[dataSourceKey(datasource)] == 0 {
		fmt.Printf("[fetchDataSystem] Ошибка. Не было подписки на %s %s\n", datasource.Symbol, datasource.Interval)
	}

	s.subscribes[dataSourceKey(datasource)]--
	if s.subscribes[dataSourceKey(datasource)] == 0 {
		fmt.Printf("[fetchDataSystem] Unsubscribe %s %s\n", datasource.Symbol, datasource.Interval)
	} else {
		fmt.Printf("[fetchDataSystem] .. %s %s\n", datasource.Symbol, datasource.Interval)
	}
}

func dataSourceKey(datasource *components.DataSource) string {
	if datasource != nil {
		return datasource.Symbol + datasource.Interval
	} else {
		return ""
	}
}

// Проверка соответствия интерфейсу
var _ ecs.System = (*fetchDataSystem)(nil)
var _ ecsx.EntityLifecycleListener = (*fetchDataSystem)(nil)
