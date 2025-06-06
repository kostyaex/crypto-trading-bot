package calc

import (
	"fmt"
	"testing"
)

func TestKMeansWeighted1D(t *testing.T) {

	points := []WeightedPoint{
		{Value: 1, Weight: 3},
		{Value: 2, Weight: 2},
		{Value: 3, Weight: 1},
		{Value: 10, Weight: 4},
		{Value: 11, Weight: 3},
		{Value: 12, Weight: 2},
	}

	clusters := KMeansWeighted1D(points, 2, 100)
	for i, c := range clusters {
		fmt.Printf("Кластер %d:\n", i+1)
		fmt.Printf("Центр: %.2f\n", c.Center)
		fmt.Println("Точки:")
		for _, p := range c.Points {
			fmt.Printf("  %.2f вес: %.2f\n", p.Value, p.Weight)
		}
		fmt.Println()
	}
}
