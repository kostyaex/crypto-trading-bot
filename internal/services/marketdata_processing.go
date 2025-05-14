package services

import (
	"crypto-trading-bot/internal/calc"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
	"fmt"
	"sync"
)

// Источник данных (база данных или загруженные с биржи данные)
type MarketDataSource interface {
	GetMarketDataCh() <-chan *models.MarketData
	Close()
}

// ---------------------------------------------------------------------------
// источник исторических данных из БД
type HistoricalSource struct {
	data []*models.MarketData
}

func NewHistoricalSource(data []*models.MarketData) *HistoricalSource {
	return &HistoricalSource{data: data}
}

func (h *HistoricalSource) GetMarketDataCh() <-chan *models.MarketData {
	ch := make(chan *models.MarketData)
	go func() {
		for _, item := range h.data {
			ch <- item
		}
		close(ch)
	}()
	return ch
}

func (h *HistoricalSource) Close() {}

// ---------------------------------------------------------------------------
// источник реальных данных, загруженных с биржи

// type LiveMarketSource struct {
// 	wsClient *WebSocketClient // Предположим, есть клиент биржи
// }

// func NewLiveMarketSource(client *WebSocketClient) *LiveMarketSource {
// 	return &LiveMarketSource{wsClient: client}
// }

// func (l *LiveMarketSource) GetMarketDataCh() <-chan *models.MarketData {
// 	return l.wsClient.SubscribeToTrades()
// }

// func (l *LiveMarketSource) Close() {
// 	l.wsClient.Unsubscribe()
// }

// ---------------------------------------------------------------------------
// Мультикастер для распределения потока торговых данных по сгрупированным стратегиям.

type Broadcaster struct {
	subscribers []chan *models.MarketData
	source      <-chan *models.MarketData
	wg          sync.WaitGroup
}

func NewBroadcaster(source <-chan *models.MarketData) *Broadcaster {
	return &Broadcaster{
		subscribers: make([]chan *models.MarketData, 0),
		source:      source,
	}
}

func (b *Broadcaster) Subscribe() <-chan *models.MarketData {
	ch := make(chan *models.MarketData, 100)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *Broadcaster) Start() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for item := range b.source {
			for _, ch := range b.subscribers {
				ch <- item
			}
		}

		// Закрываем все подписки после окончания потока
		for _, ch := range b.subscribers {
			close(ch)
		}
	}()
}

func (b *Broadcaster) Wait() {
	b.wg.Wait()
}

// ---------------------------------------------------------------------------

// Функция распределения данных по группам
func groupStrategiesBySymbolInterval(strategies []models.Strategy) map[string][]models.Strategy {
	grouped := make(map[string][]models.Strategy)
	for _, strategy := range strategies {
		settings, err := strategy.Settings()
		if err != nil {
			panic("Не удалось получить параметры стратегии")
		}

		key := fmt.Sprintf("%s|%s", settings.Symbol, settings.Interval)
		grouped[key] = append(grouped[key], strategy)
	}
	return grouped
}

// ---------------------------------------------------------------------------

func (s *marketDataService) RunStrategyForSource(strategy models.Strategy, source MarketDataSource) error {
	defer source.Close()

	strategySettings, err := strategy.Settings()
	if err != nil {
		return err
	}

	marketDataCh := source.GetMarketDataCh()

	// Разбиваем полученные торговые данные на интевалы по настройкам из стратегии
	intervalsCh := make(chan []*models.MarketData)
	go func() {
		utils.SplitChannelWithOverlap(marketDataCh, strategySettings.Waves.BlockSize, strategySettings.Waves.Overlap, intervalsCh)
		//close(intervalsCh)
	}()

	// результирующий набор волн
	waves := make([]*models.MarketWave, 0)
	// временное соответствие для определения пересекающихся точек между разными интервалами
	waveToPrices := make(map[*models.MarketWave][]calc.WeightedPoint)

	var previousInterval []*models.MarketData
	for interval := range intervalsCh {
		if len(previousInterval) == 0 {
			previousInterval = interval
			continue
		}

		newWaves := analyzeInterval(interval, previousInterval, strategySettings.Waves.NumClusters, waveToPrices)
		waves = append(waves, newWaves...)
		previousInterval = interval

		// Генерация сигналов
		//for _, wave := range newWaves {
		//signal := TradeSignal{
		//Symbol:     strategy.Symbol,
		//Price:      wave.Points[0].Price,
		//Volume:     wave.Points[0].Volume,
		//Timestamp:  wave.Start,
		//ActionType: determineActionType(wave),
		//}
		//dispatcher.Dispatch(signal)
		//}
	}

	if s.conf.Backtesting.WavesDumpDir != "" {
		filename := fmt.Sprintf("%s/waves_%s.json", s.conf.Backtesting.WavesDumpDir, strategy.Name)
		saveWaves(waves, filename)
	}

	return nil
}

func analyzeInterval(
	interval, previousInterval []*models.MarketData,
	numClusters int,
	waveToPrices map[*models.MarketWave][]calc.WeightedPoint,
) []*models.MarketWave {

	points := make([]calc.WeightedPoint, len(interval))
	pointsPrevious := make([]calc.WeightedPoint, len(previousInterval))

	for i, d := range interval {
		points[i] = calc.WeightedPoint{Value: d.OpenPrice, Weight: d.Volume}
	}
	for i, d := range previousInterval {
		pointsPrevious[i] = calc.WeightedPoint{Value: d.OpenPrice, Weight: d.Volume}
	}

	clusters := calc.KMeansWeighted1D(append(points, pointsPrevious...), numClusters, 100)

	// соответствие - по цене получить кластер в который она попала
	pricesToCluster := make(map[float64]*calc.Cluster)
	for _, cluster := range clusters {
		for _, p := range cluster.Points {
			pricesToCluster[p.Value] = &cluster
		}
	}

	clustersWithWave := make(map[float64]*models.MarketWave)
	for wave, pointsInWave := range waveToPrices {
		maxCluster := calc.FindMaxVolumeCluster(pointsInWave, pricesToCluster)
		if maxCluster != nil {
			clustersWithWave[maxCluster.Center] = wave
			updateWavePoints(wave, interval, maxCluster, pricesToCluster)
		}
	}

	clear(waveToPrices)

	var newWaves []*models.MarketWave
	for _, cluster := range clusters {
		if _, ok := clustersWithWave[cluster.Center]; ok {
			continue
		}

		md := interval[len(interval)-1] // просто последняя запись торговых данных, чтобы определить некоторые параметры
		newWave := &models.MarketWave{
			Start: interval[0].Timestamp,
			Stop:  md.Timestamp,
			Points: []models.MarketWavePoint{{
				Timestamp:     md.Timestamp,
				Price:         cluster.Center,
				Volume:        md.Volume,
				BuyVolume:     md.BuyVolume,
				SellVolume:    md.SellVolume,
				ClusterPoints: cluster.Points,
			}},
		}

		newWaves = append(newWaves, newWave)
		waveToPrices[newWave] = filterPointsByCluster(&cluster, points, pricesToCluster)
	}

	return newWaves
}

// Отфильтровать переданные точки (points) по указанному кластеру (cluster)
// cluster - один из кластеров
// points - часть точек, по которым опредялись кластеры
// pricesToCluster - служебное соответствие - по цене получить кластер в который она попала
func filterPointsByCluster(cluster *calc.Cluster, points []calc.WeightedPoint, pricesToCluster map[float64]*calc.Cluster) []calc.WeightedPoint {
	var res []calc.WeightedPoint
	for _, p := range points {
		c, ok := pricesToCluster[p.Value]
		if ok && c.Center == cluster.Center {
			res = append(res, p)
		}
	}
	return res
}

func updateWavePoints(wave *models.MarketWave, interval []*models.MarketData, cluster *calc.Cluster, pricesToCluster map[float64]*calc.Cluster) {
	md := interval[len(interval)-1]
	wave.Points = append(wave.Points, models.MarketWavePoint{
		Timestamp:     md.Timestamp,
		Price:         cluster.Center,
		Volume:        md.Volume,
		BuyVolume:     md.BuyVolume,
		SellVolume:    md.SellVolume,
		ClusterPoints: cluster.Points,
	})
	wave.Stop = md.Timestamp
}
