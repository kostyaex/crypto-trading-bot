package components

import (
	"crypto-trading-bot/internal/engine"
	"testing"
	"time"
)

func TestDataSource_NextPosition(t *testing.T) {
	// Create test data with sequential timestamps
	t1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Minute)
	t3 := t2.Add(time.Minute)

	testData := []engine.MarketData{
		{Timestamp: t1},
		{Timestamp: t2},
		{Timestamp: t3},
	}

	ds := NewDataSource(testData)

	// Test case 1: Valid next timestamp
	next, ok := ds.NextPosition(t1)
	if !ok {
		t.Error("Expected valid next timestamp, got false")
	}
	if next != t2 {
		t.Errorf("Expected %v, got %v", t2, next)
	}

	// Test case 2: Last element
	next, ok = ds.NextPosition(t3)
	if ok {
		t.Error("Expected end of data, got true")
	}
	if !next.IsZero() {
		t.Errorf("Expected zero time, got %v", next)
	}

	// Test case 3: Middle element
	next, ok = ds.NextPosition(t2)
	if !ok {
		t.Error("Expected valid next timestamp, got false")
	}
	if next != t3 {
		t.Errorf("Expected %v, got %v", t3, next)
	}
}
