package services

import (
	"testing"
	"time"
)

func Test_exchangeService_LoadData(t *testing.T) {

	setup := NewTestSetup()

	// Выбираем данные час назад - они должны быть

	exchange := setup.exchanges[0]
	symbol := "BTCUSDT"
	timeFrame := "5m"
	startTime := time.Now().Add(-time.Hour)

	marketData, lastTime, err := setup.exchangeService.LoadData(exchange, symbol, timeFrame, startTime)
	if err != nil {
		t.Errorf("exchangeService.LoadData() error = %v", err)
		return
	}

	if len(marketData) == 0 {
		t.Error("exchangeService.LoadData() не получены данные")
		return
	}

	t.Logf("Loaded data %v to %v", len(marketData), lastTime)

}

func Test_exchangeService_LoadData_Empty(t *testing.T) {
	// Здесь случай когда данных нет за указанное время. Время указываем на час вперед.

	setup := NewTestSetup()

	exchange := setup.exchanges[0]
	symbol := "BTCUSDT"
	timeFrame := "5m"
	startTime := time.Now().Add(time.Hour)

	marketData, lastTime, err := setup.exchangeService.LoadData(exchange, symbol, timeFrame, startTime)
	if err != nil {
		t.Errorf("exchangeService.LoadData() error = %v", err)
		return
	}

	if len(marketData) != 0 {
		t.Error("exchangeService.LoadData() не должно возвращать данные")
		return
	}

	t.Logf("Loaded data %v to %v", len(marketData), lastTime)

}
