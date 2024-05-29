package core

import (
	"fmt"
	"github.com/mohammad-safakhou/finance_back_history_go/models"
)

const (
	inputGreenCandles = 4
	inputRedCandles   = 4
)

func StrategyV1(capital float64, candles []models.Candle) (float64, int) {
	var trades []models.Trade

	for i := inputGreenCandles; i < len(candles); i++ {
		if isGreenCandleConditionMet(candles, i) {
			if len(trades) != 0 && trades[len(trades)-1].Type == "long" {
				continue
			}
			if len(trades) != 0 && trades[len(trades)-1].IsOpen && trades[len(trades)-1].Type != "long" {
				trades[len(trades)-1].IsOpen = false
				trades[len(trades)-1].ExitTime = candles[i].Time
				trades[len(trades)-1].ExitPrice = candles[i].Close
			}
			trades = append(trades, models.Trade{IsOpen: true, EntryTime: candles[i].Time, EntryPrice: candles[i].Close, Type: "long"})
		} else if isRedCandleConditionMet(candles, i) {
			if len(trades) != 0 && trades[len(trades)-1].Type == "short" {
				continue
			}
			if len(trades) != 0 && trades[len(trades)-1].IsOpen && trades[len(trades)-1].Type != "short" {
				trades[len(trades)-1].IsOpen = false
				trades[len(trades)-1].ExitTime = candles[i].Time
				trades[len(trades)-1].ExitPrice = candles[i].Close
			}
			trades = append(trades, models.Trade{IsOpen: true, EntryTime: candles[i].Time, EntryPrice: candles[i].Close, Type: "short"})
		}
	}

	return calculateFinalCapital(trades, capital), len(trades)
}

func StrategyV2(capital float64, candles []models.Candle) (float64, int) {
	heikinCandles := computeHeikinAshi(candles)
	var trades []models.Trade

	for i := inputGreenCandles; i < len(candles); i++ {
		if isGreenCandleConditionMet(heikinCandles, i) {
			if len(trades) != 0 && trades[len(trades)-1].Type == "long" {
				continue
			}
			if len(trades) != 0 && trades[len(trades)-1].IsOpen && trades[len(trades)-1].Type != "long" {
				trades[len(trades)-1].IsOpen = false
				trades[len(trades)-1].ExitTime = candles[i].Time
				trades[len(trades)-1].ExitPrice = candles[i].Close
			}
			trades = append(trades, models.Trade{IsOpen: true, EntryTime: candles[i].Time, EntryPrice: candles[i].Close, Type: "long"})
		} else if isRedCandleConditionMet(heikinCandles, i) {
			if len(trades) != 0 && trades[len(trades)-1].Type == "short" {
				continue
			}
			if len(trades) != 0 && trades[len(trades)-1].IsOpen && trades[len(trades)-1].Type != "short" {
				trades[len(trades)-1].IsOpen = false
				trades[len(trades)-1].ExitTime = candles[i].Time
				trades[len(trades)-1].ExitPrice = candles[i].Close
			}
			trades = append(trades, models.Trade{IsOpen: true, EntryTime: candles[i].Time, EntryPrice: candles[i].Close, Type: "short"})
		}
	}

	return calculateFinalCapital(trades, capital), len(trades)
}

func calculateFinalCapital(trades []models.Trade, capital float64) float64 {
	c := true
	openTrades := 0
	for _, trade := range trades {
		if trade.IsOpen {
			openTrades++
		}
		if capital <= 0 && c {
			c = false
			fmt.Printf("Liquidity: %v\n", trade.EntryTime)
		}
		if !trade.IsOpen {
			if trade.Type == "long" {
				profit := (trade.ExitPrice - trade.EntryPrice) / trade.EntryPrice * capital
				capital += profit
			} else if trade.Type == "short" {
				profit := (trade.EntryPrice - trade.ExitPrice) / trade.EntryPrice * capital
				capital += profit
			}
		}
	}

	return capital
}

func isGreenCandleConditionMet(candles []models.Candle, index int) bool {
	for i := 0; i < inputGreenCandles; i++ {
		if candles[index-i].Close <= candles[index-i-1].Close {
			return false
		}
	}
	return true
}

func isRedCandleConditionMet(candles []models.Candle, index int) bool {
	for i := 0; i < inputRedCandles; i++ {
		if candles[index-i].Close >= candles[index-i-1].Close {
			return false
		}
	}
	return true
}

func computeHeikinAshi(candles []models.Candle) []models.Candle {
	var heikinAshiCandles []models.Candle
	var haClose, haOpen, haHigh, haLow float64

	for i, c := range candles {
		if i == 0 {
			haClose = (c.Open + c.High + c.Low + c.Close) / 4
			haOpen = (c.Open + c.Close) / 2
			haHigh = c.High
			haLow = c.Low
		} else {
			haClose = (c.Open + c.High + c.Low + c.Close) / 4
			haOpen = (heikinAshiCandles[i-1].Open + heikinAshiCandles[i-1].Close) / 2
			haHigh = max(c.High, max(haOpen, haClose))
			haLow = min(c.Low, min(haOpen, haClose))
		}
		heikinAshiCandles = append(heikinAshiCandles, models.Candle{
			Time:  c.Time,
			Open:  haOpen,
			High:  haHigh,
			Low:   haLow,
			Close: haClose,
		})
	}
	return heikinAshiCandles
}
