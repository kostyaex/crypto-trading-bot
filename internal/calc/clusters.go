package calc

import (
	"math"
	"math/rand"
)

type WeightedPoint struct {
	Value  float64
	Weight float64
	//Cluster int
}

type Cluster struct {
	Center float64
	Points []WeightedPoint
}

// определение центра точек (1 кластер)
func kMeansSingleCluster(points []WeightedPoint) Cluster {
	sumVal := 0.0
	totalWeight := 0.0
	for _, p := range points {
		sumVal += p.Value * p.Weight
		totalWeight += p.Weight
	}
	return Cluster{
		Center: sumVal / totalWeight,
		Points: points,
	}
}

func distance(p1, p2 float64) float64 {
	return math.Abs(p1 - p2)
}

// определение кластеров для точек в одномерном пространстве
func KMeansWeighted1D(points []WeightedPoint, k int, maxIterations int) []Cluster {
	// Инициализация центров
	centers := make([]float64, k)
	//rand.Seed(time.Now().UnixNano())
	for i := range centers {
		centers[i] = points[rand.Intn(len(points))].Value
	}

	clusters := make([]Cluster, k)
	converged := false
	iteration := 0

	for !converged && iteration < maxIterations {
		// Присвоение точек кластерам
		for i := range clusters {
			clusters[i].Points = []WeightedPoint{}
		}
		for _, p := range points {
			var minDist float64 = math.MaxFloat64
			var closest int
			for i, c := range centers {
				d := distance(p.Value, c)
				if d < minDist {
					minDist = d
					closest = i
				}
			}
			clusters[closest].Points = append(clusters[closest].Points, p)
		}

		// Обновление центров с учётом весов
		newCenters := make([]float64, k)
		converged = true
		for i := range clusters {
			if len(clusters[i].Points) == 0 {
				newCenters[i] = centers[i]
				continue
			}
			sumVal := 0.0
			totalWeight := 0.0
			for _, p := range clusters[i].Points {
				sumVal += p.Value * p.Weight
				totalWeight += p.Weight
			}
			newCenters[i] = sumVal / totalWeight
			if math.Abs(centers[i]-newCenters[i]) > 1e-6 {
				converged = false
			}
		}
		centers = newCenters
		iteration++
	}

	// Формирование финальных кластеров
	filterdClusters := make([]Cluster, 0)
	for i := range clusters {
		clusters[i].Center = centers[i]
		if math.IsNaN(clusters[i].Center) || clusters[i].Center == 0 {
			continue
		}
		filterdClusters = append(filterdClusters, clusters[i])
	}

	return filterdClusters
}
