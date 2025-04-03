package services

import (
	"context"
	"crypto-trading-bot/internal/calc"
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/utils"
	"fmt"
	"strings"
	"sync"
	"time"
)

type MarketDataService interface {
	SaveMarketData(data []*models.MarketData) error
	GetMarketData(symbol string, limit int) ([]*models.MarketData, error)
	GetMarketDataStatus(id int) (*models.MarketDataStatus, error)
	SaveMarketDataStatus(marketdatastatus *models.MarketDataStatus) error
	GetMarketDataStatusList() ([]*models.MarketDataStatus, error)
	ClusterMarketData(data []*models.MarketData, numClusters int) ([]*models.MarketData, error)
	RunSchudeler(ctx context.Context)
}

type marketDataService struct {
	repo            *repositories.Repository
	logger          *utils.Logger
	exchanges       []exchange.Exchange
	exchangeService ExchangeService
	mu              sync.Mutex
	timeFrame       string                                   // интервал для группировки данных
	intervalsOrder  []time.Time                              // Список ключей в порядке добавления
	intervals       map[time.Time]*models.MarketDataInterval // соответствие для хранения загруженных интервалов
	lastTime        time.Time
}

func NewMarketDataService(repo *repositories.Repository, logger *utils.Logger, exchanges []exchange.Exchange, exchangeService ExchangeService) MarketDataService {
	return &marketDataService{
		repo:            repo,
		logger:          logger,
		timeFrame:       "5m",
		intervalsOrder:  make([]time.Time, 0),
		intervals:       make(map[time.Time]*models.MarketDataInterval),
		exchanges:       exchanges,
		exchangeService: exchangeService,
	}
}

// Запустить регламентную загрузку.
func (s *marketDataService) RunSchudeler(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			// Получен сигнал завершения
			return
		default:
			// Проверяем время последнего запуска. Если не прошло 5 секунд - ждем 1 секунду
			if time.Now().Add(-5 * time.Second).Before(s.lastTime) {
				s.logger.Debug("Ждем как пройдет 5 секунд с последнего запуска\n")
				time.Sleep(1 * time.Second)
				continue
			}
			//s.logger.Debug("LoadData\n")
			s.LoadData()
			s.lastTime = time.Now()
		}
	}
}

// добавить загруженные данные. При этом происходит группировка по интервалам.
func (s *marketDataService) Push(marketData models.MarketData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	startTime, endTime, _, _ := GetIntervalBounds(marketData.Timestamp, s.timeFrame)
	if _, exists := s.intervals[startTime]; !exists {
		s.intervals[startTime] = &models.MarketDataInterval{
			Start:   startTime,
			End:     endTime,
			Records: make([]models.MarketData, 0),
		}
		s.intervalsOrder = append(s.intervalsOrder, marketData.Timestamp)
	}
	s.intervals[startTime].Records = append(s.intervals[startTime].Records, marketData)

}

func (s *marketDataService) Pull(timeLimit time.Time) []models.MarketDataInterval {
	s.mu.Lock()
	defer s.mu.Unlock()

	var completed []models.MarketDataInterval
	var newOrder []time.Time

	// Итерация по сохранённому порядку
	for _, start := range s.intervalsOrder {
		group, exists := s.intervals[start]
		if exists && timeLimit.After(group.End) {
			completed = append(completed, *group)
			delete(s.intervals, start)
		} else {
			newOrder = append(newOrder, start)
		}
	}

	s.intervalsOrder = newOrder // Обновление списка ключей

	return completed
}

// Загружает пары по интервалам указанным в таблице marketdataStatuses.
// Возвращает данными свернутыми интервалами: сами интервалы массивом (чтобы сохранить порядок) и соответствие массивов по интервалам.
// Интервалы должны быть полными, т.е. по всем биржам в этом интервале данные должны быть загружены полностью, другими словами уже есть данные за следующий интервал.
func (s *marketDataService) LoadData() {
	s.logger.Infof("Starting data fetching")

	// перебираем все биржи и таблицу состояния данных - для каждой активной выполняем загрузку
	// В таблице состояний записаны данные для загрузки: пара (символ), интервал, дата актуальности с которой нужно продолжить загрузку

	statusList, err := s.GetMarketDataStatusList()
	if err != nil {
		s.logger.Errorf("Failed to GetMarketDataStatusList %v", err)
		return
	}

	for _, ex := range s.exchanges {

		for _, status := range statusList {

			if !status.Active || status.Exchange != strings.ToLower(ex.GetName()) {
				continue
			}

			// получаем начало интервала, чтобы получить интервал полностью
			//startTime, _, _, _ := GetIntervalBounds(status.ActualTime, status.TimeFrame)
			startTime := status.ActualTime

			s.logger.Infof("Fetching data from exchange: %s %s %v", ex.GetName(), status.Symbol, startTime)

			marketData, lastTime, err := s.exchangeService.LoadData(ex, status.Symbol, status.TimeFrame, startTime)
			if err != nil {
				s.logger.Errorf("Failed to fetch data from exchange %s: %v", ex.GetName(), err)
				continue
			}

			// Сохранение рыночных данных в базу данных
			if err := s.SaveMarketData(marketData); err != nil {
				s.logger.Errorf("Failed to save market data: %v", err)

				status.Status = fmt.Sprintf("ОШИБКА: %v", err)
				if err := s.SaveMarketDataStatus(status); err != nil {
					s.logger.Errorf("Failed to save market data: %v", err)
					return
				}

				return
			}

			// Сохранение статуса загрузки данных
			status.ActualTime = lastTime
			status.Status = "OK"
			if err := s.SaveMarketDataStatus(status); err != nil {
				s.logger.Errorf("Failed to save market data: %v", err)
				return
			}
		}

	}

	s.logger.Infof("Data fetching task completed")

	return
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
