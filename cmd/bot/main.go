package main

import (
	"context"
	"crypto-trading-bot/internal/components/sources"
)

func main() {
	s := &sources.TickSource{}
	s.Run(context.Background())
}
