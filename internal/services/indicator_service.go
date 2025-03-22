package services

import (
	"crypto-trading-bot/internal/repositories"
	"time"
)

type IndicatorService interface {
	SaveIndicator(symbol, indicatorName string, value float64, timestamp time.Time) error
}

type indicatorService struct {
	repo *repositories.Repository
}

func NewIndicatorService(repo *repositories.Repository) IndicatorService {
	return &indicatorService{repo: repo}
}

func (s *indicatorService) SaveIndicator(symbol, indicatorName string, value float64, timestamp time.Time) error {
	return s.repo.IndicatorRepository.SaveIndicator(symbol, indicatorName, value, timestamp)
}
