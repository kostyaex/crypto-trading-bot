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

type WavesCollection struct {
	Wawes []*models.MarketWave
	//ConnectionPoints []calc.WeightedPoint // точки соприкосновения интервалов
	priceToWave map[float64]*models.MarketWave
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
		timeFrame:       "1h", // таймфрем для группировки торговых данных для анализа
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
	marketData, err := s.GetMarketData("BTCUSDT", 1000)
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

	//-------------------------------------------------------------------

	intervalsCh := s.GetIntervals(marketDataCh)

	waves := make([]*models.MarketWave, 0)
	//pricesToWave := make(map[float64]*models.MarketWave, 0)           // определение волны по цене
	waveToPrices := make(map[*models.MarketWave][]calc.WeightedPoint) // обратное соответствие

	// получаем интервалы из канала
	for interval := range intervalsCh {
		fmt.Printf("Interval: %v - %d\n", interval.Start, len(interval.Records)+len(interval.PreviousInterval.Records))

		// Подготовка данных для кластеризации
		// !! Здесь надо сворачивать объемы с одинаковой ценой в одну точку
		points := make([]calc.WeightedPoint, len(interval.Records))
		pointsPreviousInterval := make([]calc.WeightedPoint, len(interval.PreviousInterval.Records))
		for n, d := range interval.Records {
			points[n] = calc.WeightedPoint{Value: d.OpenPrice, Weight: d.BuyVolume}
		}
		for n, d := range interval.PreviousInterval.Records {
			pointsPreviousInterval[n] = calc.WeightedPoint{Value: d.OpenPrice, Weight: d.BuyVolume}
		}

		// Кластеризация
		numClusters := 3
		clusters := calc.KMeansWeighted1D(append(points, pointsPreviousInterval...), numClusters, 100)

		// Для понимания алгоритма:
		// - Новая волна формируется из точек только нового интервала (точки прошлого уже сформировали волны)
		// - в кластеры берём точки как прошлого интервала, так и нового
		// - какие то точки нового интервала попадут в один кластер с точками прошлого интервала, а какие-то - не попадут,
		//   т.е. будут кластеры только с точками нового интервала. По таким кластерам формируем новую волну.

		// 1. Нужно определть, какие кластеры только с ценами нового интервала, а какие содержат предыдущие цены

		// 1.1 для удобства получения кластеров - Формируем соответствие для определения по цен соответствующего кластера
		pricesToCluster := make(map[float64]*calc.Cluster)
		for _, cluster := range clusters {
			for _, point := range cluster.Points {
				pricesToCluster[point.Value] = &cluster
			}
		}

		// 1.2 Сначала определим какие кластеры содержат цены предыдущего интервала
		//     Заодно можем эти кластеры соотнести с волнами
		//     Ключевым здесь будет соответсвие - pricesToWave

		// Определим, в каких кластерах были цены из прошлого интервала.
		clustersWithPricecFromPrevIntervals := make(map[*calc.Cluster]bool) // в ресурсе храним объем
		//clusterToWave := make(map[string]float64)                              // на одну волну может попасть несколько кластеров, надо выбрать больший по объему

		// Обходим уже сформированные на прошлом шаге волны, перебираем цены в каждой
		for wave, points := range waveToPrices {
			// по каждому кластуре подсчитываем объемы из точек этого кластера, которые есть в волне
			clusterValue := make(map[*calc.Cluster]float64, 0)
			for _, point := range points {
				if cluster, ok := pricesToCluster[point.Value]; ok {
					clusterValue[cluster] += point.Weight
				}
			}

			// Теперь выбираем кластер с максимальным значением
			var maxCluster *calc.Cluster
			for key := range clusterValue {
				if maxCluster == nil || clusterValue[key] > clusterValue[maxCluster] {
					maxCluster = key
				}
			}

			fmt.Printf("clusterValue for wave %s:\n", wave.String())
			for cluster, value := range clusterValue {
				fav := ""
				if cluster == maxCluster {
					clustersWithPricecFromPrevIntervals[cluster] = true
					wave.Points = append(wave.Points, models.MarketWavePoint{
						Timestamp: interval.Start,
						Price:     cluster.Center,
					})
					fav = "*"
				}
				fmt.Printf("%f - %10.2f %s\n", cluster.Center, value, fav)
			}
			fmt.Println()
		}

		// for _, point := range pointsPreviousInterval {
		// 	if cluster, ok := pricesToCluster[point.Value]; ok {

		// 		// if wave, ok3 := pricesToWave[point.Value]; ok3 {
		// 		// 	key := fmt.Sprintf("%f-%s-%f", cluster.Center, wave.Start.Format("02.01.2006 15:04:05"), wave.Points[0].Price)
		// 		// 	clusterToWave[key] = clusterToWave[key] + point.Weight
		// 		// }

		// 		// if _, ok2 := clustersWithPricecFromPrevIntervals[cluster]; !ok2 { // делаем проверку, может уже добавляли по другой цене
		// 		// 	clustersWithPricecFromPrevIntervals[cluster] = point.Value

		// 		// }
		// 		clustersWithPricecFromPrevIntervals[cluster] = clustersWithPricecFromPrevIntervals[cluster] + point.Value
		// 	}
		// }
		// fmt.Printf("clustersWithPricecFromPrevIntervals:\n")
		// for key, value := range clustersWithPricecFromPrevIntervals {
		// 	fmt.Printf("%f - %10.2f\n", key.Center, value)
		// }
		// fmt.Println()

		// // находим волну и добавляем туда цену кластера
		// if wave, ok3 := pricesToWave[point.Value]; ok3 {
		// 	wave.Points = append(wave.Points, models.MarketWavePoint{
		// 		Timestamp: interval.Start,
		// 		Price:     cluster.Center,
		// 	})
		// 	wave.Stop = interval.Start
		// }

		//clear(pricesToWave) // соответствие уже использовано, очищаем, чтобы заполнить данными текущего интервала
		clear(waveToPrices)

		// 2.1 Выбираем кластеры, которые не попали в clustersWithPricecFromPrevIntervals
		clustersForNewWaves := make([]*calc.Cluster, 0)
		for _, cluster := range clusters {
			if _, ok := clustersWithPricecFromPrevIntervals[&cluster]; !ok {
				clustersForNewWaves = append(clustersForNewWaves, &cluster)
				wave := &models.MarketWave{
					Start:  interval.Start,
					Stop:   interval.Start,
					Symbol: interval.Symbol,
					Points: make([]models.MarketWavePoint, 1),
				}
				wave.Points[0] = models.MarketWavePoint{
					Timestamp: interval.Start,
					Price:     cluster.Center,
				}
				waves = append(waves, wave)

				// заполняем соответствие для следующего интервала
				for _, point := range cluster.Points {
					//pricesToWave[point.Value] = wave
					if waveToPrices[wave] == nil {
						waveToPrices[wave] = make([]calc.WeightedPoint, 0)
					}
					waveToPrices[wave] = append(waveToPrices[wave], point)
				}
			}
		}

		//controlCount := 0
		//for _, cluster := range clusters {
		// if math.IsNaN(cluster.Center) || cluster.Center == 0 {
		// 	continue
		// }

		// controlCount += len(cluster.Points)
		// fmt.Printf("Cluster: %-20f - %-3d\n", cluster.Center, len(cluster.Points))

		// for w, _ := range wavesForCluster {
		// 	fmt.Printf("Wave: %v - %f\n", w.Start, w.ClusterPrice)
		// }

		//

		// // Точки кластера должны распределить в существующие волны либо создать новую волну.

		// // определяем по точкам волны, в которых они были
		// wavesForCluster := make(map[*models.MarketWave]bool, 0)
		// for _, point := range cluster.Points {
		// 	if w, ok := priceToWave[point.Value]; ok {
		// 		wavesForCluster[w] = true
		// 	}
		// }

		// clear(priceToWave)

		// //

		// // запоминаем связку цены с волной для следующего шага

		// for _, point := range cluster.Points {
		// 	priceToWave[point.Value] = &wave
		// }

		//}

		//fmt.Printf("Контрольное количество: %d\n", controlCount)
	}

	fmt.Printf("Waves: %d\n", len(waves))
	timeFormat := "02.01.2006 15:04:05"
	for _, wave := range waves {
		fmt.Printf("Wave: %s - %s\n", wave.Start.Format(timeFormat), wave.Stop.Format(timeFormat))
		for _, point := range wave.Points {
			fmt.Printf(" - %s %f\n", point.Timestamp.Format(timeFormat), point.Price)
		}
	}

	return nil
}

// // Обработать точки нового шага.
// func (wavesCollection *WavesCollection) commitPoints(points []calc.WeightedPoint, startTime time.Time) {

// 	if wavesCollection.Wawes == nil {
// 		wavesCollection.Wawes = make([]*models.MarketWave, 0)
// 		wavesCollection.priceToWave = make(map[float64]*models.MarketWave, 0)
// 	}

// 	var wave *models.MarketWave

// 	// Если хоть одна цена кластера попадает в существующую волну - относим кластер к этой волне
// 	for _, point := range points {
// 		if w, ok := wavesCollection.priceToWave[point.Value]; !ok {
// 			wavesCollection.priceToWave[point.Value] = &models.MarketWave{
// 				Start: startTime,
// 				//Symbol:       interval.Symbol,
// 			}
// 		}
// 	}

// 	// предыдущие соответствия не нужны - очищаем. Ниже заполним.
// 	//clear(wavesCollection.priceToWave)

// 	// Формируем новую волну
// 	if wave != nil {
// 		wave =
// 		wavesCollection.Wawes = append(wavesCollection.Wawes, wave)
// 	}

// 	// заполняем соответствия цен и волн для следующего шага
// 	for _, point := range cluster.Points {
// 		if w, ok := wavesCollection.priceToWave[point.Value]; !ok {
// 			wavesCollection.priceToWave[point.Value] = w
// 		}
// 	}

// 	for _, point := range cluster.Points {
// 		//priceToWave[point.Value] = &wave
// 		wavesCollection.ConnectionPoints = append(wavesCollection.ConnectionPoints, point)
// 	}

// 	// controlCount += len(cluster.Points)
// 	// fmt.Printf("Cluster: %f - %d\n", cluster.Center, len(cluster.Points))
// 	// for w, _ := range wavesForCluster {
// 	// 	fmt.Printf("Wave: %v - %f\n", w.Start, w.ClusterPrice)
// 	// }
// }
