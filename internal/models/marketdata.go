package models

import (
	"crypto-trading-bot/internal/calc"
	"fmt"
	"time"
)

// MarketDataStatus представляет модель MarketDataStatus.
type MarketDataStatus struct {
	ID         int       `db:"id"`
	Exchange   string    `db:"exchange"`
	Symbol     string    `db:"symbol"`
	TimeFrame  string    `db:"time_frame"`
	Active     bool      `db:"active"`
	ActualTime time.Time `db:"actual_time"`
	Status     string    `db:"status"`
}

// MarketData представляет модель данных о рынке.
type MarketData struct {
	Timestamp  time.Time `db:"timestamp"`
	Exchange   string    `db:"exchange"`
	Symbol     string    `db:"symbol"`
	TimeFrame  string    `db:"time_frame"`
	OpenPrice  float64   `db:"open_price"`
	ClosePrice float64   `db:"close_price"`
	Volume     float64   `db:"volume"`
	BuyVolume  float64   `db:"buy_volume"`
	SellVolume float64   `db:"sell_volume"`
}

type ClusterData struct {
	Timestamp    time.Time `db:"timestamp"`
	Symbol       string    `db:"symbol"`
	TimeFrame    string    `db:"time_frame"`
	IsBuySell    bool      `db:"is_buysell"`
	ClusterPrice float64   `db:"buy_cluster"`
	Volume       float64   `db:"volume"`
}

type MarketWave struct {
	Start  time.Time `json:"start"`
	Stop   time.Time `json:"stop"`
	Symbol string    `json:"symbol"`
	//IsBuySell        bool      `db:"is_buysell"`
	//ClusterPrice float64
	Points []MarketWavePoint `json:"points"`
}

type MarketWavePoint struct {
	Timestamp     time.Time            `json:"timestamp"`
	Price         float64              `json:"price"`
	Volume        float64              `json:"volume"`
	BuyVolume     float64              `json:"buyvolume"`
	SellVolume    float64              `json:"sellvolume"`
	ClusterPoints []calc.WeightedPoint `json:"cluster_points"`
}

func (wave *MarketWave) String() string {
	if len(wave.Points) == 0 {
		return fmt.Sprintf("%s", wave.Start.Format("02.01.2006 15:04:05"))
	}

	return fmt.Sprintf("%s-%f", wave.Start.Format("02.01.2006 15:04:05"), wave.Points[0].Price)
}
