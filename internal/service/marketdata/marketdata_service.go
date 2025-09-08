package marketdata

import (
	"context"
	"crypto-trading-bot/internal/config"
	"crypto-trading-bot/internal/logger"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/service/exchange"
	"crypto-trading-bot/internal/types"
	"fmt"
	"strings"
	"sync"
	"time"
)

type MarketDataService interface {
	SaveMarketData(data []*types.MarketData) error
	GetMarketData(symbol string, limit int) ([]*types.MarketData, error)
	GetMarketDataStatus(id int) (*types.MarketDataStatus, error)
	SaveMarketDataStatus(marketdatastatus *types.MarketDataStatus) error
	GetMarketDataStatusList() ([]*types.MarketDataStatus, error)
	RunSchudeler(ctx context.Context)
	GetMarketDataPeriod(symbol string, interval string, start time.Time, end time.Time) ([]*types.MarketData, error)
}

type marketDataService struct {
	conf            *config.Config
	repo            *repositories.Repository
	logger          *logger.Logger
	exchanges       []exchange.Exchange
	exchangeService exchange.ExchangeService
	mu              sync.Mutex
	lastTime        time.Time
}

func NewMarketDataService(conf *config.Config,
	repo *repositories.Repository,
	logger *logger.Logger,
	exchanges []exchange.Exchange,
	exchangeService exchange.ExchangeService) MarketDataService {

	return &marketDataService{
		conf:            conf,
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
				s.logger.Debug("Выполняем не чаще 1 раза в секунду\n")
				time.Sleep(1 * time.Second)
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
	s.logger.Debugf("Starting data fetching")

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

			s.logger.Debugf("Fetching data from exchange: %s %s %v", ex.GetName(), status.Symbol, startTime)

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

	s.logger.Debugf("Data fetching task completed")

	return
}

// SaveMarketData сохраняет рыночные данные.
func (s *marketDataService) SaveMarketData(data []*types.MarketData) error {
	if err := s.repo.MarketData.SaveMarketData(data); err != nil {
		s.logger.Errorf("Failed to save market data: %v", err)
		return err
	}
	return nil
}

func (s *marketDataService) GetMarketData(symbol string, limit int) ([]*types.MarketData, error) {
	return s.repo.MarketData.GetMarketData(symbol, limit)
}

// GetMarketDataStatus получает информацию о marketdatastatus по ID.
func (s *marketDataService) GetMarketDataStatus(id int) (*types.MarketDataStatus, error) {
	marketdatastatus, err := s.repo.MarketData.GetMarketDataStatus(id)
	if err != nil {
		s.logger.Errorf("Failed to get marketdatastatus: %v", err)
		return nil, err
	}
	return marketdatastatus, nil
}

// SaveMarketDataStatus сохраняет информацию о marketdatastatus.
func (s *marketDataService) SaveMarketDataStatus(marketdatastatus *types.MarketDataStatus) error {
	if err := s.repo.MarketData.SaveMarketDataStatus(marketdatastatus); err != nil {
		s.logger.Errorf("Failed to save marketdatastatus: %v", err)
		return err
	}
	return nil
}

// GetMarketDataStatus получает список marketdatastatus.
func (s *marketDataService) GetMarketDataStatusList() ([]*types.MarketDataStatus, error) {
	marketdatastatus, err := s.repo.MarketData.GetMarketDataStatusList()
	if err != nil {
		s.logger.Errorf("Failed to get marketdatastatus list: %v", err)
		return nil, err
	}
	return marketdatastatus, nil
}

func (s *marketDataService) GetMarketDataPeriod(symbol string, interval string, start time.Time, end time.Time) ([]*types.MarketData, error) {
	return s.repo.MarketData.GetMarketDataPeriod(symbol, interval, start, end)
}
