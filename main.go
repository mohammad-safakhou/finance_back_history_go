package main

import (
	"encoding/csv"
	"fmt"
	"github.com/mohammad-safakhou/finance_back_history_go/core"
	"github.com/mohammad-safakhou/finance_back_history_go/models"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//// Open the CSV file
	//file, err := os.Open("./hist_data_XAUUSD_M1_2009_2023/DAT_MT_XAUUSD_M1_2009.csv")
	//if err != nil {
	//	fmt.Println("Error opening file:", err)
	//	return
	//}
	//core.Start(file, "2009")
	//file.Close()
	//
	//file, err = os.Open("./hist_data_XAUUSD_M1_2009_2023/DAT_MT_XAUUSD_M1_2010.csv")
	//if err != nil {
	//	fmt.Println("Error opening file:", err)
	//	return
	//}
	//core.Start(file, "2010")
	//file.Close()

	dir := "./hist_data_XAUUSD_M1_2009_2023/"
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return
	}

	capital := 100000.
	initialCapital := capital
	percents := []float64{}
	for _, f := range files {
		if f.Name()[0:6] != "DAT_MT" {
			continue
		}
		file, err := os.Open(dir + f.Name())
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		candles := Candles(file)
		//if candles[0].Time.Year() != 2023 {
		//	continue
		//}

		//candles30 := convertTo30MinuteCandles(candles)

		tradeCounts := 0
		strategy := core.CandleCount{
			NumRedCandle:      1,
			NumGreenCandle:    1,
			StopLossPips:      80,
			TakeProfitPips:    2000,
			StopLossPercent:   .01,
			TakeProfitPercent: .30,
			TimeFrame:         30,
			Leverage:          20,
		}
		capital, tradeCounts = strategy.StrategyCandleCount(initialCapital, candles)
		percent := (float64(capital) * 100 / initialCapital) - 100
		fmt.Printf("Year %s Trade Counts: %d Exit Capital: %f (%f)\n", f.Name(), tradeCounts, capital, percent)
		percents = append(percents, percent)
		file.Close()
	}

	total := 0.
	for _, percent := range percents {
		fmt.Println(percent)
		total += percent
	}

	fmt.Println(total)

	//capital := 1000000.
	//for _, f := range files {
	//	file, err := os.Open(dir + f.Name())
	//	if err != nil {
	//		fmt.Println("Error opening file:", err)
	//		return
	//	}
	//	candles := Candles(file)
	//	if candles[0].Time.Year() != 2023 {
	//		continue
	//	}
	//
	//	//candles30 := convertTo30MinuteCandles(candles)
	//
	//	a := capital
	//	tradeCounts := 0
	//	capital, tradeCounts = core.StrategyV1(capital, candles)
	//	fmt.Printf("Year %s Trade Counts: %d Exit Capital: %f (%f)\n", f.Name(), tradeCounts, capital, (capital*100/a)-100)
	//	file.Close()
	//}
	//
	//fmt.Println("\n\n------------------------------\n")
	//
	//capital = 1000000.
	//for _, f := range files {
	//	file, err := os.Open(dir + f.Name())
	//	if err != nil {
	//		fmt.Println("Error opening file:", err)
	//		return
	//	}
	//	candles := Candles(file)
	//	candles30 := utils.ConvertTo30MinuteCandles(candles)
	//
	//	a := capital
	//	tradeCounts := 0
	//	capital, tradeCounts = core.StrategyV2(capital, candles30)
	//	fmt.Printf("Year %s Trade Counts: %d Exit Capital: %f (%f)\n", f.Name(), tradeCounts, capital, (capital*100/a)-100)
	//	file.Close()
	//}
}

func Candles(file *os.File) []models.Candle {
	// Read the CSV file
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return nil
	}

	// Process each record
	var candles []models.Candle
	for _, record := range records {
		candle, err := parseRecord(record)
		if err != nil {
			fmt.Println("Error parsing record:", err)
			continue
		}
		candles = append(candles, candle)
	}
	return candles
}

// parseRecord parses a CSV record into a Candle struct
func parseRecord(record []string) (models.Candle, error) {
	if len(record) < 6 {
		return models.Candle{}, fmt.Errorf("invalid record length")
	}

	// Parse DateTime
	dateTimeStr := record[0] + " " + strings.Split(strings.ReplaceAll(record[1], " ", ""), ">")[0]
	dateTime, err := time.Parse("2006.01.02 15:04", dateTimeStr)
	if err != nil {
		return models.Candle{}, fmt.Errorf("error parsing DateTime: %v", err)
	}

	// Parse Open
	open, err := strconv.ParseFloat(strings.ReplaceAll(record[2], " ", ""), 64)
	if err != nil {
		return models.Candle{}, fmt.Errorf("error parsing Open: %v", err)
	}

	// Parse High
	high, err := strconv.ParseFloat(strings.ReplaceAll(record[3], " ", ""), 64)
	if err != nil {
		return models.Candle{}, fmt.Errorf("error parsing High: %v", err)
	}

	// Parse Low
	low, err := strconv.ParseFloat(strings.ReplaceAll(record[4], " ", ""), 64)
	if err != nil {
		return models.Candle{}, fmt.Errorf("error parsing Low: %v", err)
	}

	// Parse Close
	close, err := strconv.ParseFloat(strings.ReplaceAll(record[5], " ", ""), 64)
	if err != nil {
		return models.Candle{}, fmt.Errorf("error parsing Close: %v", err)
	}

	// Parse Volume
	volume, err := strconv.ParseFloat(strings.ReplaceAll(record[6], " ", ""), 64)
	if err != nil {
		return models.Candle{}, fmt.Errorf("error parsing Volume: %v", err)
	}

	return models.Candle{
		Time:   dateTime,
		Open:   open,
		High:   high,
		Low:    low,
		Close:  close,
		Volume: volume,
	}, nil
}
