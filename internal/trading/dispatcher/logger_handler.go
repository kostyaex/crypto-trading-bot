package dispatcher

import (
	"fmt"
	"time"
)

type LoggerHandler struct {
}

func (l *LoggerHandler) Handle(signal TradeSignal) {
	fmt.Printf("[%s] %s @ %.2f | %.2f %s | %s\n",
		signal.Timestamp.Format(time.RFC3339),
		signal.Type,
		signal.Price,
		signal.Volume,
		signal.Symbol,
		signal.Description)
}
