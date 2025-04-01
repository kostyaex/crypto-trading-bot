package services

import (
	"crypto-trading-bot/internal/calc"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/utils"
	"fmt"
	"time"
)

type MarketDataService interface {
	SaveMarketData(data []*models.MarketData) error
	GetMarketData(symbol string, limit int) ([]*models.MarketData, error)
	GetMarketDataStatus(id int) (*models.MarketDataStatus, error)
	SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error
	GetMarketDataStatusList() ([]*models.MarketDataStatus, error)
	ClusterMarketData(data []*models.MarketData, numClusters int) ([]*models.MarketData, error)
}

type marketDataService struct {
	repo   *repositories.Repository
	logger *utils.Logger
}

func NewMarketDataService(repo *repositories.Repository, logger *utils.Logger) MarketDataService {
	return &marketDataService{repo: repo, logger: logger}
}

// SaveMarketData сохраняет рыночные данные.
func (s *marketDataService) SaveMarketData(data []*models.MarketData) error {
	if err := s.repo.MarketData.SaveMarketData(data); err != nil {
		s.logger.Errorf("Failed to save market data: %v", err)
		return err
	}
	return nil
}

func (s *marketDataService) GetMarketData(symbol string, limit int) ([]*models.MarketData, error) {
	return s.repo.MarketData.GetMarketData(symbol, limit)
}

// GetMarketDataStatus получает информацию о marketdatastatus по ID.
func (s *marketDataService) GetMarketDataStatus(id int) (*models.MarketDataStatus, error) {
	marketdatastatus, err := s.repo.MarketData.GetMarketDataStatus(id)
	if err != nil {
		s.logger.Errorf("Failed to get marketdatastatus: %v", err)
		return nil, err
	}
	return marketdatastatus, nil
}

// SaveMarketDataStatus сохраняет информацию о marketdatastatus.
func (s *marketDataService) SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error {
	if err := s.repo.MarketData.SaveMarketDataStatus(marketdatastatus); err != nil {
		s.logger.Errorf("Failed to save marketdatastatus: %v", err)
		return err
	}
	return nil
}

// GetMarketDataStatus получает список marketdatastatus.
func (s *marketDataService) GetMarketDataStatusList() ([]*models.MarketDataStatus, error) {
	marketdatastatus, err := s.repo.MarketData.GetMarketDataStatusList()
	if err != nil {
		s.logger.Errorf("Failed to get marketdatastatus list: %v", err)
		return nil, err
	}
	return marketdatastatus, nil
}

//
// Кластеризация
//

// ClusterMarketData кластеризует рыночные данные по объемам покупок и продаж по ценам.
func (s *marketDataService) ClusterMarketData(data []*models.MarketData, numClusters int) ([]*models.MarketData, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data to cluster")
	}

	// // Подготовка данных для кластеризации
	// points := mat.NewDense(len(data), 2, nil)
	// for i, d := range data {
	// 	points.Set(i, 0, d.OpenPrice)
	// 	points.Set(i, 1, d.BuyVolume)
	// }

	// // Кластеризация объемов покупок
	// buyClusters, err := kMeansCluster(points, numClusters)
	// if err != nil {
	// 	s.logger.Errorf("Failed to cluster buy volumes: %v", err)
	// 	return nil, err
	// }

	// // Подготовка данных для кластеризации объемов продаж
	// for i, d := range data {
	// 	points.Set(i, 1, d.SellVolume)
	// }

	// // Кластеризация объемов продаж
	// sellClusters, err := kMeansCluster(points, numClusters)
	// if err != nil {
	// 	s.logger.Errorf("Failed to cluster sell volumes: %v", err)
	// 	return nil, err
	// }

	// Подготовка данных для кластеризации
	points := make([]calc.WeightedPoint, len(data))
	for _, d := range data {
		points = append(points, calc.WeightedPoint{Value: d.OpenPrice, Weight: d.BuyVolume})
	}

	// Кластеризация объемов покупок
	buyClusters := calc.KMeansWeighted1D(points, numClusters, 100)

	points = make([]calc.WeightedPoint, len(data))
	for _, d := range data {
		points = append(points, calc.WeightedPoint{Value: d.OpenPrice, Weight: d.SellVolume})
	}

	// Кластеризация объемов продаж
	sellClusters := calc.KMeansWeighted1D(points, numClusters, 100)

	buyClusterData := make([]*models.ClusterData, len(buyClusters))
	sellClusterData := make([]*models.ClusterData, len(sellClusters))

	// заполнение кластеров в данные биржи
	for _, d := range data {
		for _, cluster := range buyClusters {
			for _, point := range cluster.Points {
				if d.OpenPrice == point.Value {
					buyClusterData = append(buyClusterData, &models.ClusterData{
						Timestamp:    d.Timestamp,
						Symbol:       d.Symbol,
						TimeFrame:    d.TimeFrame,
						IsBuySell:    true,
						ClusterPrice: cluster.Center,
						Volume:       d.BuyVolume,
					})
				}
			}
		}
		for _, cluster := range sellClusters {
			for _, point := range cluster.Points {
				if d.OpenPrice == point.Value {
					sellClusterData = append(sellClusterData, &models.ClusterData{
						Timestamp:    d.Timestamp,
						Symbol:       d.Symbol,
						TimeFrame:    d.TimeFrame,
						IsBuySell:    false,
						ClusterPrice: cluster.Center,
						Volume:       d.BuyVolume,
					})
				}
			}
		}
	}

	s.repo.MarketData.SaveClusterData(buyClusterData)
	s.repo.MarketData.SaveClusterData(sellClusterData)

	return data, nil
}

//------------------------------------------------------------------

// Возвращает начало временного интервала для заданного времени и типа интервала
func GetIntervalBounds(t time.Time, interval string) (start, end, nextStart time.Time, err error) {

	// Рассчитываем начало интервала
	if interval == "1M" {
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		nextStart = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
		end = nextStart.Add(-time.Nanosecond) // Конец интервала
		return
	} else if interval == "1d" {
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
		nextDay := t.Add(24 * time.Hour)
		nextStart = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, t.Location())
		return
	} else if interval == "1w" {
		start, end, nextStart = getWeekBounds(t)
		return
	}

	// Парсинг интервала
	duration, err := parseInterval(interval)
	if err != nil {
		return
	}
	start = t.Truncate(duration)
	end = start.Add(duration).Add(-time.Nanosecond) // Конец интервала
	nextStart = start.Add(duration)

	return
}

// определить границы недели
func getWeekBounds(t time.Time) (start, end, nextStart time.Time) {
	year, week := t.ISOWeek()
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
	daysSinceEpoch := int(firstDay.Weekday()) - 1
	if daysSinceEpoch < 0 {
		daysSinceEpoch += 7
	}
	start = firstDay.AddDate(0, 0, (week-1)*7-daysSinceEpoch)
	end = start.Add(7 * 24 * time.Hour).Add(-time.Nanosecond)
	nextStart = start.Add(7 * 24 * time.Hour)
	return
}

// Парсинг строки интервала в time.Duration
func parseInterval(interval string) (time.Duration, error) {
	// Словарь интервалов
	intervals := map[string]time.Duration{
		"1s":  time.Second,
		"1m":  time.Minute,
		"3m":  3 * time.Minute,
		"5m":  5 * time.Minute,
		"15m": 15 * time.Minute,
		"30m": 30 * time.Minute,
		"1h":  time.Hour,
		"4h":  4 * time.Hour,
		"6h":  6 * time.Hour,
		"8h":  8 * time.Hour,
		"12h": 12 * time.Hour,
		"1d":  24 * time.Hour,
		"3d":  72 * time.Hour,
		"1w":  168 * time.Hour,
		//"1M":  30 * 24 * time.Hour, // Примерное значение для месяца
	}

	// Проверка допустимых значений
	if duration, ok := intervals[interval]; ok {
		return duration, nil
	}
	return 0, fmt.Errorf("недопустимый интервал: %s", interval)
}
