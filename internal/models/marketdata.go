package models

import "time"

// MarketData представляет модель данных о рынке.
type MarketData struct {
	ID         int       `json:"id"`
	Exchange   string    `json:"exchange"`
	Symbol     string    `json:"symbol"`
	OpenPrice  float64   `json:"open_price"`
	ClosePrice float64   `json:"close_price"`
	Volume     float64   `json:"volume"`
	BuyVolume  float64   `json:"buy_volume"`
	SellVolume float64   `json:"sell_volume"`
	TimeFrame  string    `json:"time_frame"`
	Timestamp  time.Time `json:"timestamp"`
}
