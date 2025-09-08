package sampling

import (
	"context"
	"crypto-trading-bot/pkg/pipeline"
)

type MockSampler struct{}

func (s *MockSampler) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {

	return p, nil
}
