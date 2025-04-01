package models

import "time"

// интервал с указанием таймфрейма и начала соответствующего отрезка
type MarketInterval struct {
	TimeFrame string    `json:"time_frame"`
	Start     time.Time `json:"timestamp"`
}

// MarketData представляет модель данных о рынке.
type MarketData struct {
	Timestamp  time.Time `json:"timestamp"`
	Exchange   string    `json:"exchange"`
	Symbol     string    `json:"symbol"`
	TimeFrame  string    `json:"time_frame"`
	OpenPrice  float64   `json:"open_price"`
	ClosePrice float64   `json:"close_price"`
	Volume     float64   `json:"volume"`
	BuyVolume  float64   `json:"buy_volume"`
	SellVolume float64   `json:"sell_volume"`
}

type ClusterData struct {
	Timestamp    time.Time `json:"timestamp"`
	Symbol       string    `json:"symbol"`
	TimeFrame    string    `json:"time_frame"`
	IsBuySell    bool      `json:"is_buysell"`
	ClusterPrice float64   `json:"buy_cluster"`
	Volume       float64   `json:"volume"`
}
