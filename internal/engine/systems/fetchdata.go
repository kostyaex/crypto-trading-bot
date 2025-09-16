package systems

import (
	"context"
	"crypto-trading-bot/internal/engine/components"
	"crypto-trading-bot/internal/engine/service/exchange"

	"github.com/andygeiss/ecs"
)

type fetchDataSystem struct {
	ctx      context.Context
	executor *exchange.FetchDataExecutor
}

func NewFetchDataSystem(ctx context.Context) ecs.System {

	s := &fetchDataSystem{
		ctx: ctx,
	}

	// подготавливаем исполнитель
	s.executor = exchange.NewFetchDataExecutor(ctx)
	// for _, cmd := range fetchDataCommands {
	// 	s.executor.SubmitCommand(cmd)
	// 	fmt.Printf("📤 Submitted: %v\n", cmd)
	// }
	return s
}

func (s *fetchDataSystem) Process(em ecs.EntityManager) (state int) {

	// выбираем все компоненты datasource, из них выбираем по каким символам и интервалам нужны данные.
	// Группируем и отправляем команды для выборки

	grouped := make(map[string]exchange.FetchDataCommand)

	for _, dataComp := range em.FilterByMask(components.MaskDatasource) {

		datasource := dataComp.Get(components.MaskDatasource).(*components.DataSource)

		key := datasource.Symbol + datasource.Interval
		grouped[key] = exchange.FetchDataCommand{Symbol: datasource.Symbol, Interval: datasource.Interval}

	}

	return ecs.StateEngineContinue
}

func (s *fetchDataSystem) Setup() {}

func (s *fetchDataSystem) Teardown() {}

// Проверка соответствия интерфейсу
var _ ecs.System = (*fetchDataSystem)(nil)
