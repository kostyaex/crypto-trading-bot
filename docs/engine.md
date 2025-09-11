# Реализация движка

Предполагается, что в одном экземпляре engine будут выполнятся одновременно как реальные торговые сессии, так и бектестинг. Выполняться будут по разным парам и интервалам, всё это одновременно должно выполняться.

Торговые данные будут представляться компонентом datasource. Данные в этом компоненте будут подгружаться из базы и обновляться из биржи.

Торговые сессии (в том числе бектестинг) будут содержать компонент position, который позиционируется по времени на источнике данных. Позицию будет сдвигать  система MovementSystem.

```go
type DataSource struct {
	data           []engine.MarketData
	indexTimestamp map[time.Time]int // индекс записи по времени
}

func NewDataSource(data []engine.MarketData) *DataSource {
	s := &DataSource{
		data: data,
	}

	s.indexTimestamp = make(map[time.Time]int, len(data))

	for n, v := range data {
		s.indexTimestamp[v.Timestamp] = n
	}

	return s
}

func (c *DataSource) Mask() uint64 {
	return MaskDatasource
}

// Получить следующую отметку времени, если она есть
func (c *DataSource) NextPosition(timestamp time.Time) (time.Time, bool) {

	if timestamp.IsZero() {
		if len(c.data) > 0 {
			return c.data[0].Timestamp, true
		} else {
			return time.Time{}, false
		}
	}

	ind := c.indexTimestamp[timestamp]
	if ind < len(c.data)-1 {
		return c.data[ind+1].Timestamp, true
	} else {
		return time.Time{}, false
	}

}
```

```go
type Position struct {
	PrevTimestamp time.Time // Предыдущее время позиции
	Timestamp     time.Time // Текущее время позиции
}

func NewPosition(Timestamp time.Time) *Position {
	return &Position{
		Timestamp: Timestamp,
	}
}

func (c *Position) Mask() uint64 {
	return MaskPosition
}

func (c *Position) SetPosition(timestamp time.Time) {
	c.PrevTimestamp = c.Timestamp
	c.Timestamp = timestamp
}
```

```go
type movementSystem struct {
}

func NewMovementSystem() ecs.System {
	return &movementSystem{}
}

func (s *movementSystem) Process(em ecs.EntityManager) (state int) {

	// сдвигаем позиции в источнике данных

	for _, dataComp := range em.FilterByMask(components.MaskDatasource) {

		datasource := dataComp.Get(components.MaskDatasource).(*components.DataSource)

		for _, posComp := range em.FilterByMask(components.MaskPosition) {
			position := posComp.Get(components.MaskPosition).(*components.Position)

			if next, ok := datasource.NextPosition(position.Timestamp); ok {
				position.SetPosition(next)
				fmt.Printf("Позиция перемещена: %s\n", next.Format("2006-01-02 15:04:05"))
			}

		}
	}

	return ecs.StateEngineContinue
}

func (s *movementSystem) Setup() {}

func (s *movementSystem) Teardown() {}

// Проверка соответствия интерфейсу
var _ ecs.System = (*movementSystem)(nil)
```
