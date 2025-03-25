package services

import (
	"crypto-trading-bot/internal/exchange"
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/repositories"
	"crypto-trading-bot/internal/utils"
	"fmt"
	"strings"
)

type ExchangeService interface {
	LoadData() []*models.MarketData
}

type exchangeService struct {
	repo              *repositories.Repository
	exchanges         []exchange.Exchange
	logger            *utils.Logger
	marketDataService MarketDataService
}

func NewEchangeService(repo *repositories.Repository, logger *utils.Logger, exchanges []exchange.Exchange, marketDataService MarketDataService) ExchangeService {
	return &exchangeService{
		repo:              repo,
		logger:            logger,
		exchanges:         exchanges,
		marketDataService: marketDataService,
	}
}

func (s *exchangeService) LoadData() []*models.MarketData {
	s.logger.Infof("Starting data fetching")
	var allMarketData []*models.MarketData

	// перебираем все биржи и таблицу состояния данных - для каждой активной выполняем загрузку
	// В таблице состояний записаны данные для загрузки: пара (символ), интервал, дата актуальности с которой нужно продолжить загрузку

	statusList, err := s.marketDataService.GetMarketDataStatusList()
	if err != nil {
		s.logger.Errorf("Failed to GetMarketDataStatusList %v", err)
		return nil
	}

	for _, ex := range s.exchanges {

		for _, status := range statusList {

			if !status.Active || status.Exchange != strings.ToLower(ex.GetName()) {
				continue
			}

			s.logger.Infof("Fetching data from exchange: %s %s %v", ex.GetName(), status.Symbol, status.ActualTime)

			marketData, lastTime, err := ex.GetMarketData(status.Symbol, status.TimeFrame, status.ActualTime)
			if err != nil {
				s.logger.Errorf("Failed to fetch data from exchange %s: %v", ex.GetName(), err)
				continue
			}

			// if err := s.repo.MarketData.SaveMarketData(marketData); err != nil {
			// 	s.logger.Errorf("Failed to save market data for exchange %s: %v", ex.GetName(), err)
			// } else {
			// 	s.logger.Infof("Market data saved for exchange %s", ex.GetName())
			// }

			allMarketData = append(allMarketData, marketData...)

			// Сохранение рыночных данных в базу данных
			if err := s.marketDataService.SaveMarketData(marketData); err != nil {
				s.logger.Errorf("Failed to save market data: %v", err)

				status.Status = fmt.Sprintf("ОШИБКА: %v", err)
				if err := s.marketDataService.SaveMarketDataStatus(status); err != nil {
					s.logger.Errorf("Failed to save market data: %v", err)
					return nil
				}

				return nil
			}

			// Сохранение статуса загрузки данных
			status.ActualTime = lastTime
			status.Status = "OK"
			if err := s.marketDataService.SaveMarketDataStatus(status); err != nil {
				s.logger.Errorf("Failed to save market data: %v", err)
				return nil
			}
		}

	}

	s.logger.Infof("Data fetching task completed")

	return allMarketData
}
