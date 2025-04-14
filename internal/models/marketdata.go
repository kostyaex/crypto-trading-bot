package models

import "time"

// интервал с указанием таймфрейма и начала соответствующего отрезка
type MarketDataInterval struct {
	TimeFrame        string
	Symbol           string
	Start            time.Time
	End              time.Time
	Records          []*MarketData
	PreviousInterval *MarketDataInterval
}

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
