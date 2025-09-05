package sources

import (
	"context"
	"log"
	"time"
)

type TickSource struct{}

func (s *TickSource) Run(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Println("Sampling tick")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
