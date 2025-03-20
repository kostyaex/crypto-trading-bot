package analysis

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"errors"
	"math"
)

// CalculateRSI вычисляет индикатор RSI для заданного символа
func CalculateRSI(repo *repositories.Repository, symbol string, period int) (float64, error) {
	marketData, err := repo.GetMarketData(symbol, period)
	if err != nil {
		return 0, err
	}

	if len(marketData) < period {
		return 0, errors.New("not enough data points to calculate RSI")
	}

	gains := make([]float64, period-1)
	losses := make([]float64, period-1)

	for i := 1; i < period; i++ {
		change := marketData[i].Price - marketData[i-1].Price
		if change > 0 {
			gains[i-1] = change
		} else {
			losses[i-1] = math.Abs(change)
		}
	}

	avgGain := average(gains)
	avgLoss := average(losses)

	if avgLoss == 0 {
		return 100, nil
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi, nil
}

// CalculateMACD вычисляет индикатор MACD для заданного символа
func CalculateMACD(repo *repositories.Repository, symbol string, fastPeriod, slowPeriod, signalPeriod int) (float64, float64, error) {
	marketData, err := repo.GetMarketData(symbol, slowPeriod+signalPeriod-1)
	if err != nil {
		return 0, 0, err
	}

	if len(marketData) < slowPeriod+signalPeriod-1 {
		return 0, 0, errors.New("not enough data points to calculate MACD")
	}

	fastEMA, err := calculateEMA(marketData, fastPeriod)
	if err != nil {
		return 0, 0, err
	}

	slowEMA, err := calculateEMA(marketData, slowPeriod)
	if err != nil {
		return 0, 0, err
	}

	macdValues := make([]float64, len(slowEMA))
	for i := range slowEMA {
		macdValues[i] = fastEMA[i] - slowEMA[i]
	}

	// Создаем срез данных для вычисления сигнальной линии
	macdMarketData := make([]*models.MarketData, len(macdValues))
	for i, macdValue := range macdValues {
		macdMarketData[i] = &models.MarketData{
			Symbol:    symbol,
			Price:     macdValue,
			Timestamp: marketData[i+fastPeriod-1].Timestamp,
		}
	}

	signalLineValues, err := calculateEMAMarketData(macdMarketData, signalPeriod)
	if err != nil {
		return 0, 0, err
	}

	if len(signalLineValues) == 0 {
		return 0, 0, errors.New("failed to calculate signal line")
	}

	macdLine := macdValues[len(macdValues)-1]
	signalLine := signalLineValues[len(signalLineValues)-1]

	return macdLine, signalLine, nil
}

// calculateEMA вычисляет экспоненциальное скользящее среднее
func calculateEMA(marketData []*models.MarketData, period int) ([]float64, error) {
	if len(marketData) < period {
		return nil, errors.New("not enough data points to calculate EMA")
	}

	ema := make([]float64, len(marketData))
	multiplier := 2.0 / float64(period+1)

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += marketData[i].Price
	}
	ema[period-1] = sum / float64(period)

	for i := period; i < len(marketData); i++ {
		ema[i] = (marketData[i].Price-ema[i-1])*multiplier + ema[i-1]
	}

	return ema, nil
}

// calculateEMAMarketData вычисляет экспоненциальное скользящее среднее для среза данных
func calculateEMAMarketData(marketData []*models.MarketData, period int) ([]float64, error) {
	if len(marketData) < period {
		return nil, errors.New("not enough data points to calculate EMA")
	}

	ema := make([]float64, len(marketData))
	multiplier := 2.0 / float64(period+1)

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += marketData[i].Price
	}
	ema[period-1] = sum / float64(period)

	for i := period; i < len(marketData); i++ {
		ema[i] = (marketData[i].Price-ema[i-1])*multiplier + ema[i-1]
	}

	return ema, nil
}

// average вычисляет среднее значение массива чисел
func average(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
