package services

import (
	"crypto-trading-bot/internal/repositories"
)

type BehaviorTreeService interface {
	GetBehaviorTreeState(strategyID int) (map[string]interface{}, error)
	SaveBehaviorTreeState(strategyID int, state map[string]interface{}) error
}

type behaviorTreeService struct {
	repo *repositories.Repository
}

func NewBehaviorTree(repo *repositories.Repository) BehaviorTreeService {
	return &behaviorTreeService{repo: repo}
}

func (s *behaviorTreeService) GetBehaviorTreeState(strategyID int) (map[string]interface{}, error) {
	return s.repo.BehaviorTreeRepository.GetBehaviorTreeState(strategyID)
}

func (s *behaviorTreeService) SaveBehaviorTreeState(strategyID int, state map[string]interface{}) error {
	return s.repo.BehaviorTreeRepository.SaveBehaviorTreeState(strategyID, state)
}
