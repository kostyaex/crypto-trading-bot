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
	RunBacktesting(startTime, endTime time.Time) error
	Push(marketData *models.MarketData)
	Pull(timeLimit time.Time) []*models.MarketDataInterval
	GetIntervals(marketDataCh <-chan *models.MarketData) <-chan *models.MarketDataInterval
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
		timeFrame:       "5m", // таймфрем для группировки торговых данных для анализа
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
			// Проверяем время последнего запуска.
			if time.Now().Add(-1 * time.Second).Before(s.lastTime) {
				s.logger.Debug("Ждем как пройдет 300 секунда с последнего запуска\n")
				time.Sleep(300 * time.Second)
				continue
			}
			//s.logger.Debug("LoadData\n")
			s.LoadData()
			s.lastTime = time.Now()
		}
	}
}

// добавить загруженные данные. При этом происходит группировка по интервалам.
func (s *marketDataService) Push(marketData *models.MarketData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	startTime, endTime, _, _ := utils.GetIntervalBounds(marketData.Timestamp, s.timeFrame)
	if _, exists := s.intervals[startTime]; !exists {
		s.intervals[startTime] = &models.MarketDataInterval{
			Start:   startTime,
			End:     endTime,
			Records: make([]*models.MarketData, 0),
		}
		s.intervalsOrder = append(s.intervalsOrder, marketData.Timestamp)
	}
	s.intervals[startTime].Records = append(s.intervals[startTime].Records, marketData)

}

// Получить сгрупированные по интервалам данные
func (s *marketDataService) Pull(timeLimit time.Time) []*models.MarketDataInterval {
	s.mu.Lock()
	defer s.mu.Unlock()

	var completed []*models.MarketDataInterval
	var newOrder []time.Time

	// Итерация по сохранённому порядку
	for _, start := range s.intervalsOrder {
		group, exists := s.intervals[start]
		if exists && timeLimit.After(group.End) {
			completed = append(completed, group)
			delete(s.intervals, start)
		} else {
			newOrder = append(newOrder, start)
		}
	}

	s.intervalsOrder = newOrder // Обновление списка ключей

	return completed
}

// Разбивает поток биржевых данных на интервалы по timeFrame
func (s *marketDataService) GetIntervals(marketDataCh <-chan *models.MarketData) <-chan *models.MarketDataInterval {
	intervals := make(chan *models.MarketDataInterval)

	var timeMarker time.Time

	go func() {

		var currentInterval *models.MarketDataInterval

		for marketData := range marketDataCh {
			startTime, endTime, _, _ := utils.GetIntervalBounds(marketData.Timestamp, s.timeFrame)

			if startTime.After(timeMarker) {
				// вышли за границу текущего интервала

				// текущий интервал сбрасываем в канал
				if currentInterval != nil && currentInterval.PreviousInterval != nil {
					intervals <- currentInterval
				}

				// формируем новый интервал
				currentInterval = &models.MarketDataInterval{
					Start:            startTime,
					End:              endTime,
					Records:          make([]*models.MarketData, 0),
					PreviousInterval: currentInterval,
				}
			}

			timeMarker = startTime

			currentInterval.Records = append(currentInterval.Records, marketData)
		}

		if currentInterval != nil && currentInterval.PreviousInterval != nil {
			intervals <- currentInterval
		}

		close(intervals)
	}()

	return intervals
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

			if len(marketData) == 0 {
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

func (s *marketDataService) RunBacktesting(startTime, endTime time.Time) error {

	// Получить данные за период
	// надо переделать на функцию получения данных за период
	marketData, err := s.GetMarketData("BTCUSDT", 100)
	if err != nil {
		s.logger.Errorf("Ошибка получения торговых данных: %s\n", err)
		return err
	}

	fmt.Printf("marketData: %d\n", len(marketData))

	// Передаем тестовые данные в канал биржевых данных
	marketDataCh := make(chan *models.MarketData)
	go func() {
		for _, marketDataItem := range marketData {
			marketDataCh <- marketDataItem
		}
		close(marketDataCh)
	}()

	intervalsCh := s.GetIntervals(marketDataCh)

	// получаем интервалы из канала
	for interval := range intervalsCh {
		fmt.Printf("Interval: %v - %d\n", interval.Start, len(interval.Records))
	}

	return nil
}
