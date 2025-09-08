package processing

import (
	"context"
	"crypto-trading-bot/pkg/pipeline"
	"fmt"
)

type LoggerSink struct{}

func (l *LoggerSink) Consume(_ context.Context, payload pipeline.Payload) error {
	tradingPayload, ok := payload.(*TradingPayload)
	if !ok {
		return fmt.Errorf("invalid payload type: %T", payload)
	}
	fmt.Printf("symbol: %s, interval: %s\n", tradingPayload.Symbol, tradingPayload.Interval)
	return nil
}
