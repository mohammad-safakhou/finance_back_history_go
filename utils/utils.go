package utils

import (
	"github.com/mohammad-safakhou/finance_back_history_go/models"
	"time"
)

func ConvertTo30MinuteCandles(oneMinCandles []models.Candle) []models.Candle {
	var thirtyMinCandles []models.Candle
	var tempCandle models.Candle
	var tempCandleOpenSet bool

	for i, c := range oneMinCandles {
		if i == 0 || c.Time.Minute()%30 == 0 {
			if tempCandleOpenSet {
				thirtyMinCandles = append(thirtyMinCandles, tempCandle)
			}
			tempCandle = models.Candle{
				Time:  c.Time.Truncate(30 * time.Minute),
				Open:  c.Open,
				High:  c.High,
				Low:   c.Low,
				Close: c.Close,
			}
			tempCandleOpenSet = true
		} else {
			tempCandle.High = max(tempCandle.High, c.High)
			tempCandle.Low = min(tempCandle.Low, c.Low)
			tempCandle.Close = c.Close
		}
	}

	// Append the last 30-minute candle if not already appended
	if tempCandleOpenSet {
		thirtyMinCandles = append(thirtyMinCandles, tempCandle)
	}

	return thirtyMinCandles
}

func Find30MinCandleBasedOn1MinCandle(candles30 []models.Candle, lastCandle1Min models.Candle) (int, models.Candle) {
	for i, candle := range candles30 {
		if candle.Time == lastCandle1Min.Time {
			return i, candle
		}
	}
	return -1, models.Candle{}
}
