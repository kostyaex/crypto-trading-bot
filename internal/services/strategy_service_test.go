package services

import (
	"fmt"
	"testing"
)

func Test_strategyService_GetActiveStrategies(t *testing.T) {
	setup := NewTestSetup()

	strategies, err := setup.strategyService.GetActiveStrategies()
	if err != nil {
		t.Errorf("strategyService.GetActiveStrategies() error = %v", err)
		return
	}

	for _, strategy := range strategies {
		settings, err := strategy.Settings()
		if err != nil {
			t.Errorf("strategyService.GetActiveStrategies() error = %v", err)
			return
		}

		fmt.Printf("Strategy: %s %v\n", strategy.Name, settings)
	}
}
