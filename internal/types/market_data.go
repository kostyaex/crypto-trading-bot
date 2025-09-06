package types

import (
	"fmt"
	"strings"
	"time"
)

const TimeFormat string = "02.01.2006 15:04:05"

// ЗДЕСЬ НЕ ДОЛЖНО БЫТЬ ДИРЕКТИВ db!!

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
	Timestamp    time.Time `db:"timestamp"`
	Exchange     string    `db:"exchange"`
	Symbol       string    `db:"symbol"`
	TimeFrame    string    `db:"time_frame"`
	OpenPrice    float64   `db:"open_price"`
	HightPrice   float64   `db:"hight_price"`
	LowPrice     float64   `db:"low_price"`
	ClosePrice   float64   `db:"close_price"`
	ClusterPrice float64   `db:"cluster_price"`
	Volume       float64   `db:"volume"`
	BuyVolume    float64   `db:"buy_volume"`
	SellVolume   float64   `db:"sell_volume"`
}

func (md MarketData) String() string {
	return fmt.Sprintf(
		"%s %s %8.2f %8.2f %8.2f %8.2f %8.2f",
		md.Timestamp.Format(TimeFormat),
		md.TimeFrame,
		md.OpenPrice,
		md.ClosePrice,
		md.ClusterPrice,
		md.BuyVolume,
		md.SellVolume)
}

func MarketDataToString(md []*MarketData) string {
	var s strings.Builder

	for _, md1 := range md {
		s.WriteString(md1.String() + "\n")
	}

	return s.String()
}
