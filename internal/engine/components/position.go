// позиция трейдера в торговых данных
package components

import "time"

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
