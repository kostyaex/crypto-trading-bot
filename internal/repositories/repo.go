package repositories

import (
	"crypto-trading-bot/internal/utils"
)

type Repository struct {
	db                     *DB
	logger                 *utils.Logger
	Strategy               StrategyRepository
	MarketData             MarketDataRepository
	MarketDataStatus       MarketDataStatusRepository
	IndicatorRepository    IndicatorRepository
	BehaviorTreeRepository BehaviorTreeRepository
}

func NewRepository(db *DB, logger *utils.Logger) *Repository {
	return &Repository{
		db:                     db,
		logger:                 logger,
		Strategy:               NewStrategyRepository(db, logger),
		MarketData:             NewMarketDataRepository(db, logger),
		MarketDataStatus:       NewMarketDataStatusRepository(db, logger),
		IndicatorRepository:    NewIndicatorRepository(db, logger),
		BehaviorTreeRepository: NewBehaviorTreeRepositoryRepository(db, logger),
	}
}
