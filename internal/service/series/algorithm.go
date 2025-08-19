package series

type AlgorithmType string

const (
	SimpleAlgorithmType   AlgorithmType = "simple"
	WindowedAlgorithmType AlgorithmType = "windowed"
)

type Algorithm interface {
	AddPoints(activeSeries []*Series, newPoints []Point) []*Series

	// Возвращает имя алгоритма (для логирования)
	Name() string
}
