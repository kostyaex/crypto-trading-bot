package systems

import (
	"context"
	"crypto-trading-bot/internal/engine/components"
	"crypto-trading-bot/internal/engine/ecsx"
	"crypto-trading-bot/internal/exchange"
	"fmt"
	"log"

	"github.com/andygeiss/ecs"
)

type fetchDataSystem struct {
	em          ecs.EntityManager
	ctx         context.Context
	exchange    exchange.Exchange
	subscribes  map[string]int // map[Symbol+Interval]count - кол-во подписок на данный символ и интервал
	commandsIds map[exchange.CommandID]struct{}
}

func NewFetchDataSystem(
	ctx context.Context,
	em ecs.EntityManager,
	ex exchange.Exchange,
) *fetchDataSystem {

	s := &fetchDataSystem{
		em:          em,
		ctx:         ctx,
		exchange:    ex,
		commandsIds: make(map[exchange.CommandID]struct{}),
	}

	return s
}

func (s *fetchDataSystem) Process(em ecs.EntityManager) (state int) {

	for _, dataComp := range s.em.FilterByMask(components.MaskDatasource) {

		datasource := dataComp.Get(components.MaskDatasource).(*components.DataSource)

		if datasource.CmdId != "" {
			// получаем данные
			if c, ok, err := s.exchange.PopCandle(datasource.CmdId); ok {

				if err != nil {
					log.Printf("❌ Ошибка при получении свечей: %v", err)

				} else if ok {

					log.Printf("✅ Получено %s: %s %s O=%.2f H=%.2f L=%.2f C=%.2f V=%.2f",
						c.Timestamp.Format("15:04:05"), c.Symbol, c.Interval, c.Open, c.High, c.Low, c.Close, c.Volume)

				}
			}
		}

	}

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
		cmdId := s.exchange.SubscribeCandles(datasource.Symbol, datasource.Interval)
		s.commandsIds[cmdId] = struct{}{} // запопинаем идентификатор команды для последующей обработки
		datasource.CmdId = cmdId

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
