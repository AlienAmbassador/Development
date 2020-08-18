package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Candle - структура свечи
type Candle struct {
	name         string
	minPrice     float64
	maxPrice     float64
	income       float64
	minPriceTime time.Time
	maxPriceTime time.Time
}

// User - структура юзера
type User struct {
	id      string
	tickers map[string]map[string]float64
}

func main() {
	candles, err := readFile("candles_5m.csv")
	if err != nil {
		log.Fatalf("can't read candles from csv %s", err)
	}

	trades, err := readFile("user_trades.csv")
	if err != nil {
		log.Fatalf("can't read trades from csv %s", err)
	}

	maxRevenueMap, err := mapCandles(candles)
	if err != nil {
		log.Fatalf("can't make mapCandles map %s", err)
	}

	tradeInfo, err := UserDeals(trades)
	if err != nil {
		log.Fatalf("can't make tradeInfo map %s", err)
	}

	result := finish(maxRevenueMap, tradeInfo)
	writeFile(result, "result.csv")
}

func mapCandles(candles [][]string) (map[string]Candle, error) {
	const (
		name = iota
		times
		openPrice
		maxPrice
		minPrice
		closePrice
	)

	maxRevenueMap := make(map[string]Candle)

	for _, candle := range candles {
		t, err := time.Parse(time.RFC3339, candle[times])
		if err != nil {
			return nil, err
		}

		maxPrice, err := stringToFloat(candle[maxPrice])
		minPrice, err := stringToFloat(candle[minPrice])

		if val, ok := maxRevenueMap[candle[name]]; ok {
			if minPrice < val.minPrice {
				val.minPrice = minPrice
				val.minPriceTime = t
			}

			if maxPrice > val.maxPrice {
				val.maxPrice = maxPrice
				val.maxPriceTime = t
			}

			val.income = val.maxPrice - val.minPrice
			maxRevenueMap[candle[name]] = val
		} else {
			maxRevenueMap[candle[name]] = Candle{candle[name], minPrice, maxPrice, maxPrice - minPrice, t, t}
		}
	}

	return maxRevenueMap, nil
}

func UserDeals(trades [][]string) (map[string]User, error) {
	const (
		id = iota
		times
		ticker
		buyPrice
		salePrice
	)

	tradeInfo := make(map[string]User)

	for _, trade := range trades {
		if val, ok := tradeInfo[trade[id]]; ok {
			if _, ok := val.tickers[trade[ticker]]; ok {
				salePrice, err := stringToFloat(trade[salePrice])
				if err != nil {
					return tradeInfo, err
				}
				val.tickers[trade[ticker]]["salePrice"] = salePrice
				val.tickers[trade[ticker]]["income"] = val.tickers[trade[ticker]]["salePrice"] - val.tickers[trade[ticker]]["buyPrice"]
			} else {
				buyPrice, err := stringToFloat(trade[buyPrice])
				if err != nil {
					return tradeInfo, err
				}
				val.tickers[trade[ticker]] = map[string]float64{
					"buyPrice": buyPrice,
				}
			}
		} else {
			buyPrice, err := stringToFloat(trade[buyPrice])
			if err != nil {
				return tradeInfo, err
			}
			ticker := map[string]map[string]float64{
				trade[ticker]: {
					"buyPrice": buyPrice,
				},
			}
			tradeInfo[trade[id]] = User{trade[id], ticker}
		}
	}

	return tradeInfo, nil
}

//функция для структурирования информации
func finish(maxRevenueMap map[string]Candle, tradeInfo map[string]User) [][]string {
	var result [][]string
	const (
		layout = "2006-01-02T15:04:05Z07:00"
	)

	for deal := range tradeInfo {
		for key, value := range tradeInfo[deal].tickers {
			var a []string

			userID := deal
			userRevenue := value["income"]
			maxRevenue := maxRevenueMap[key].income
			diff := maxRevenue - userRevenue
			timeToSale := maxRevenueMap[key].maxPriceTime.Format(layout)
			timeToBuy := maxRevenueMap[key].minPriceTime.Format(layout)

			a = append(a, userID, key, fmt.Sprintf("%.2f", userRevenue), fmt.Sprintf("%.2f", maxRevenue), fmt.Sprintf("%.2f", diff), timeToSale, timeToBuy)
			result = append(result, a)
		}
	}

	return result
}

// функция для чтения csv файла
func readFile(filename string) ([][]string, error) {
	File, err := os.Open(filename)
	defer File.Close()
	if err != nil {
		return nil, err
	}

	Reader := csv.NewReader(File)

	data, err := Reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

//функция для записи csv файла
func writeFile(result [][]string, filename string) error {
	File, err := os.Create(filename)
	if err != nil {
		return err
	}

	Writer := csv.NewWriter(File)

	err = Writer.WriteAll(result)
	if err != nil {
		return err
	}

	return nil
}

// функция для конвертирования string в float64
func stringToFloat(a string) (float64, error) {
	b, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return b, err
	}
	return b, nil
}
