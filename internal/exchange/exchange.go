package exchange

import (
	"crypto-trading-bot/internal/models"
	"crypto-trading-bot/internal/utils"
)

type Exchange interface {
	GetName() string
	GetMarketData() ([]*models.MarketData, error)
}

type exchangeBase struct {
	name   string
	logger *utils.Logger
}

func (e *exchangeBase) GetName() string {
	return e.name
}

func (e *exchangeBase) logInfo(msg string, args ...interface{}) {
	e.logger.Infof("%s: "+msg, append([]interface{}{e.name}, args...)...)
}

func (e *exchangeBase) logError(err error, msg string, args ...interface{}) {
	e.logger.Errorf("%s: "+msg, append([]interface{}{e.name}, args...)...)
}
