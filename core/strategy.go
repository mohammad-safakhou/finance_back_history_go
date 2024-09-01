package core

import (
	"fmt"
	"github.com/mohammad-safakhou/finance_back_history_go/models"
	"github.com/mohammad-safakhou/finance_back_history_go/utils"
)

const (
	inputGreenCandles = 4
	inputRedCandles   = 4
)

func StrategyV1(capital float64, candles []models.Candle) (float64, int) {
	var trades []models.Trade

	for i := 10; i < len(candles); i++ {
		if isRedCandleConditionMet(candles, i, inputRedCandles) {
			if len(trades) != 0 && trades[len(trades)-1].Type == "long" {
				continue
			}
			if len(trades) != 0 && trades[len(trades)-1].IsOpen && trades[len(trades)-1].Type != "long" {
				trades[len(trades)-1].IsOpen = false
				trades[len(trades)-1].ExitTime = candles[i].Time
				trades[len(trades)-1].ExitPrice = candles[i].Close
			}
			trades = append(trades, models.Trade{IsOpen: true, EntryTime: candles[i].Time, EntryPrice: candles[i].Close, Type: "long"})
		} else if isGreenCandleConditionMet(candles, i, inputGreenCandles) {
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

	for i := 10; i < len(candles); i++ {
		if isGreenCandleConditionMet(heikinCandles, i, inputGreenCandles) {
			if len(trades) != 0 && trades[len(trades)-1].Type == "long" {
				continue
			}
			if len(trades) != 0 && trades[len(trades)-1].IsOpen && trades[len(trades)-1].Type != "long" {
				trades[len(trades)-1].IsOpen = false
				trades[len(trades)-1].ExitTime = candles[i].Time
				trades[len(trades)-1].ExitPrice = candles[i].Close
			}
			trades = append(trades, models.Trade{IsOpen: true, EntryTime: candles[i].Time, EntryPrice: candles[i].Close, Type: "long"})
		} else if isRedCandleConditionMet(heikinCandles, i, inputRedCandles) {
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

type CandleCount struct {
	NumRedCandle      int
	NumGreenCandle    int
	StopLossPips      float64
	TakeProfitPips    float64
	StopLossPercent   float64
	TakeProfitPercent float64
	TimeFrame         int
	Leverage          int
}

func (c CandleCount) TouchTP(entryPrice float64, price float64, positionType string) bool {
	if positionType == "long" {
		if price > entryPrice*(c.TakeProfitPercent+100)/100 {
			return true
		}
		//if price > entryPrice+(c.TakeProfitPips/1000) {
		//	return true
		//}
	} else {
		if price < entryPrice*(100-c.TakeProfitPercent)/100 {
			return true
		}
		//if price < entryPrice-(c.TakeProfitPips/1000) {
		//	return true
		//}
	}
	return false
}

func (c CandleCount) TouchSL(entryPrice float64, price float64, positionType string) bool {
	if positionType == "long" {
		if price < entryPrice*(100-c.StopLossPercent)/100 {
			return true
		}
		//if price < entryPrice-(c.StopLossPips/1000) {
		//	return true
		//}
	} else {
		if price > entryPrice*(c.StopLossPercent-100)/100 {
			return true
		}
		//if price > entryPrice+(c.StopLossPips/1000) {
		//	return true
		//}
	}
	return false
}

func (c CandleCount) StrategyCandleCount(capital float64, candles []models.Candle) (float64, int) {
	candles30 := utils.ConvertTo30MinuteCandles(candles)
	var trades []models.Trade

	for i := 10; i < len(candles); i++ {
		if len(trades) != 0 && trades[len(trades)-1].IsOpen {
			if c.TouchTP(trades[len(trades)-1].EntryPrice, candles[i].Close, trades[len(trades)-1].Type) ||
				c.TouchSL(trades[len(trades)-1].EntryPrice, candles[i].Close, trades[len(trades)-1].Type) {
				trades[len(trades)-1].IsOpen = false
				trades[len(trades)-1].ExitTime = candles[i].Time
				trades[len(trades)-1].ExitPrice = candles[i].Close
			}
			continue
		}
		if candles[i].Time.Minute()%30 == 0 {
			lastIndex, _ := utils.Find30MinCandleBasedOn1MinCandle(candles30, candles[i])
			if lastIndex == -1 {
				panic("ohoh")
			}
			if isRedCandleConditionMet(candles30, lastIndex, c.NumRedCandle) {
				trades = append(trades, models.Trade{IsOpen: true, EntryTime: candles[i].Time, EntryPrice: candles[i].Close, Type: "short"})
			} else if isGreenCandleConditionMet(candles30, lastIndex, c.NumGreenCandle) {
				trades = append(trades, models.Trade{IsOpen: true, EntryTime: candles[i].Time, EntryPrice: candles[i].Close, Type: "long"})
			}
		}
	}

	return calculateFinalCapital(trades, capital), len(trades)
}

type capitalStruct struct {
	time    string
	capital float64
}

func calculateFinalCapital(trades []models.Trade, capital float64) float64 {
	initialCap := capital
	monthCaps := []capitalStruct{}
	start := trades[0].ExitTime
	monthCaps = append(monthCaps, capitalStruct{start.Format("2006-01"), capital})
	wins := 0
	losses := 0
	profits := []float64{}
	for _, trade := range trades {
		if trade.IsOpen {
			continue
		}
		if capital <= 0 {
			fmt.Printf("Liquidity: %v\n", trade.EntryTime)
		}
		if !trade.IsOpen {
			if trade.ExitTime.Month().String() != start.Month().String() {
				found := false
				for i, monthCap := range monthCaps {
					if monthCap.time == start.Format("2006-01") {
						monthCaps[i].capital = capital
						found = true
						break
					}
				}
				if !found {
					monthCaps = append(monthCaps, capitalStruct{start.Format("2006-01"), capital})
				}
				start = trade.ExitTime
			}
			profit := 0.
			if trade.Type == "long" {
				profit = (trade.ExitPrice - trade.EntryPrice) / trade.EntryPrice * capital

			} else if trade.Type == "short" {
				profit = (trade.EntryPrice - trade.ExitPrice) / trade.EntryPrice * capital
			}
			// 2001.06.07,23:32,1965.000100,1965.000100,1965.000100,1965.000100,0
			if int64(profit) == 362157907 {
				fmt.Println(profit)
			}
			if profit > 0 {
				wins++
			} else {
				losses++
			}
			profits = append(profits, profit)
			capital += profit
		}
	}

	a := initialCap
	for _, monthCap := range monthCaps {
		fmt.Printf("Monthly Report %s: Exit Capital: %f (%f) \n", monthCap.time, monthCap.capital, (monthCap.capital*100/a)-100)
		a = monthCap.capital
	}

	maxProfit := -100000000000000.
	minProfit := 100000000000000.
	for _, profit := range profits {
		if profit > maxProfit {
			maxProfit = profit
		}
		if profit < minProfit {
			minProfit = profit
		}
	}
	fmt.Printf("Win: %d, Loss: %d, Win/Loss: %f, Max Profit: %d, Min Profit: %d\n", wins, losses, float64(wins)*100/float64(losses), int64(maxProfit), int64(minProfit))

	return capital
}

func isGreenCandleConditionMet(candles []models.Candle, index int, count int) bool {
	for i := 0; i < count; i++ {
		if candles[index-i].Close <= candles[index-i-1].Close {
			return false
		}
	}
	return true
}

func isRedCandleConditionMet(candles []models.Candle, index int, count int) bool {
	for i := 0; i < count; i++ {
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
