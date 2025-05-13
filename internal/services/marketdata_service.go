package services

import (
	"context"
	"crypto-trading-bot/internal/calc"
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/utils"
	"encoding/json"
	"fmt"
	"os"
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
	RunBacktesting(startTime, endTime time.Time) (result []BacktestingResult)
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
	lastTime        time.Time
}

func NewMarketDataService(repo *repositories.Repository, logger *utils.Logger, exchanges []exchange.Exchange, exchangeService ExchangeService) MarketDataService {
	return &marketDataService{
		repo:            repo,
		logger:          logger,
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
		points = append(points, calc.WeightedPoint{
			Value:  d.OpenPrice,
			Weight: d.BuyVolume,
		})
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

type BacktestingResult struct {
	Strategy     *models.Strategy `json:"-"`
	StrategyName string           `json:"strategy_name"`
	Log          string           `json:"log"`
}

func (s *marketDataService) RunBacktesting(startTime, endTime time.Time) (result []BacktestingResult) {

	results := make([]BacktestingResult, 0)

	// Получить данные за период
	// надо переделать на функцию получения данных за период
	symbol := "BTCUSDT"
	marketData, err := s.repo.MarketData.GetMarketDataPeriod(symbol, startTime, endTime)
	//	marketData, err := s.GetMarketData("BTCUSDT", 1000)
	if err != nil {
		s.logger.Errorf("Ошибка получения торговых данных: %s\n", err)
		return results
	}

	s.logger.Debugf("marketData: %d\n", len(marketData))

	// Передаем тестовые данные в канал биржевых данных
	marketDataCh := make(chan *models.MarketData)
	go func() {
		for _, marketDataItem := range marketData {
			marketDataCh <- marketDataItem
		}
		close(marketDataCh)
	}()

	//-------------------------------------------------------------------

	// Разбивка данных биржи по блокам (интервалам)

	intervalsCh := make(chan []*models.MarketData) //s.GetIntervals(marketDataCh)
	blockSize := 5                                 // 60 - для минутных данных биржи получается группировка по часам
	overlap := 4                                   // Наложение блоков между собой. 0 - без наложения. Должен быть меньше размера блока.

	go utils.SplitChannelWithOverlap(marketDataCh, blockSize, overlap, intervalsCh)

	waves := make([]*models.MarketWave, 0)
	//pricesToWave := make(map[float64]*models.MarketWave, 0)           // определение волны по цене
	waveToPrices := make(map[*models.MarketWave][]calc.WeightedPoint) // обратное соответствие

	// получаем интервалы из канала
	var previousInterval []*models.MarketData
	for interval := range intervalsCh {

		if len(previousInterval) == 0 {
			previousInterval = interval
			continue
		}

		intervalStart := interval[0].Timestamp

		s.logger.Debugf("Блок: %v - %d\n", intervalStart, len(interval)+len(previousInterval))

		// Подготовка данных для кластеризации
		// Веса для кластеризации - берем общий объем проаж и покупок
		points := make([]calc.WeightedPoint, len(interval))
		pointsPreviousInterval := make([]calc.WeightedPoint, len(previousInterval))
		for n, d := range interval {
			points[n] = calc.WeightedPoint{
				Value:  d.OpenPrice,
				Weight: d.Volume,
			}
		}
		for n, d := range previousInterval {
			pointsPreviousInterval[n] = calc.WeightedPoint{Value: d.OpenPrice, Weight: d.Volume}
		}

		// Кластеризация
		numClusters := 5
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
		clustersWithPricecFromPrevIntervals := make(map[float64]*models.MarketWave) // ключ - цена кластера
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

			s.logger.Debugf("clusterValue for wave %p - %s:\n", wave, wave.String())
			// Здесь собственно подбор волны для кластера
			for cluster, value := range clusterValue {
				fav := ""
				if cluster == maxCluster {
					clustersWithPricecFromPrevIntervals[cluster.Center] = wave // для этого кластера не будет создаваться волна

					md := interval[len(interval)-1]
					wave.Points = append(wave.Points, models.MarketWavePoint{
						Timestamp:  intervalStart,
						Price:      cluster.Center,
						Volume:     md.Volume,
						BuyVolume:  md.BuyVolume,
						SellVolume: md.SellVolume,
					})
					fav = "*"
				}
				s.logger.Debugf("%f - %10.2f %s\n", cluster.Center, value, fav)
			}
			s.logger.Debugln()
		}

		//clear(pricesToWave) // соответствие уже использовано, очищаем, чтобы заполнить данными текущего интервала
		clear(waveToPrices)

		// 2.1 Выбираем кластеры, которые не попали в clustersWithPricecFromPrevIntervals
		for _, cluster := range clusters {
			var wave *models.MarketWave
			if _, ok := clustersWithPricecFromPrevIntervals[cluster.Center]; ok {
				wave = clustersWithPricecFromPrevIntervals[cluster.Center]
			} else {

				// подсчитываем объемы
				// берем только последнюю запись интервала
				md := interval[len(interval)-1]

				wave = &models.MarketWave{
					Start: intervalStart,
					Stop:  intervalStart,
					//Symbol: interval.Symbol,
					Points: make([]models.MarketWavePoint, 1),
				}
				wave.Points[0] = models.MarketWavePoint{
					Timestamp:  intervalStart,
					Price:      cluster.Center,
					Volume:     md.Volume,
					BuyVolume:  md.BuyVolume,
					SellVolume: md.SellVolume,
				}
				waves = append(waves, wave)

			}

			// заполняем соответствие для следующего интервала
			for _, point := range points {
				// фильтр точек текущего блока (points) по выбранному кластеру
				if c, ok := pricesToCluster[point.Value]; !ok || c.Center != cluster.Center {
					continue
				}

				//pricesToWave[point.Value] = wave
				if waveToPrices[wave] == nil {
					waveToPrices[wave] = make([]calc.WeightedPoint, 0)
				}
				waveToPrices[wave] = append(waveToPrices[wave], point)
			}

		}

	}

	res := fmt.Sprintf("Waves: %d\n", len(waves))
	//printWavesStat(waves)
	saveWaves(waves)

	results = append(results, BacktestingResult{
		Log: res,
	})

	return results
}

// Выводим статистику по полученным волнам
func printWaves(waves []*models.MarketWave) {
	fmt.Printf("Waves: %d\n", len(waves))
	timeFormat := "02.01.2006 15:04:05"
	for _, wave := range waves {
		fmt.Printf("Wave: %p - %s - %s\n", wave, wave.Start.Format(timeFormat), wave.Stop.Format(timeFormat))
		for _, point := range wave.Points {
			fmt.Printf(" - %s %f\n", point.Timestamp.Format(timeFormat), point.Price)
		}
		fmt.Println()
	}
}

// Выводим статистику по полученным волнам
func printWavesStat(waves []*models.MarketWave) {

	var maxLen int // максимальная длина волны

	//timeFormat := "02.01.2006 15:04:05"

	for _, wave := range waves {
		waveLen := len(wave.Points)
		if waveLen > maxLen {
			maxLen = waveLen
		}
		//fmt.Printf("Wave: %p - %s - %s\n", wave, wave.Start.Format(timeFormat), wave.Stop.Format(timeFormat))
		// for _, point := range wave.Points {
		// 	fmt.Printf(" - %s %f\n", point.Timestamp.Format(timeFormat), point.Price)
		// }
		//fmt.Println()
	}

	fmt.Printf("Waves: %d\n", len(waves))
	fmt.Printf("maxLen: %d\n", maxLen)
	fmt.Println()
}

func saveWaves(waves []*models.MarketWave) {
	// Кодируем массив структур в JSON
	jsonData, err := json.MarshalIndent(waves, "", "	")
	if err != nil {
		fmt.Println("Ошибка кодирования в JSON:", err)
		return
	}

	// Записываем JSON в файл
	filename := "./data/waves.json"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Printf("Ошибка записи в файл:%s\n", err)
		return
	}

	fmt.Printf("JSON записан в файл %s\n", filename)
}
