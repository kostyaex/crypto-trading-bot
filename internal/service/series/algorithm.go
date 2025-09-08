package series

import "crypto-trading-bot/internal/types"

type AlgorithmType string

const (
	SimpleAlgorithmType   AlgorithmType = "simple"
	WindowedAlgorithmType AlgorithmType = "windowed"
)

type Algorithm interface {
	AddPoints(activeSeries []*types.Series, newPoints []types.Point) []*types.Series

	// Возвращает имя алгоритма (для логирования)
	Name() string
}
