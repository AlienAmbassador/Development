package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Candle - структура свечи
type Candle struct {
	ticker       string
	maxPrice     float64
	minPrice     float64
	income       float64
	maxPriceTime time.Time
	minPriceTime time.Time
}

type User struct{
	id string
	tickers map[string]map[string]float64
}

func main() {
	candles := readFile("candles_5m.cvs")
	users := readFile("user_trades.csv")
	fmt.Println(candles)
	fmt.Println(users)
}
// todo lalalal
func mapCandles(candles [][]string) map[string]Candle {
	maxRevenueMap := make(map[string]Candle)
	for _, candle := range candles {
		t, _ := time.Parse(time.RFC3339, candle[1])

		maxPrice := stringToFloat(candle[3])
		minPrice := stringToFloat(candle[4])

		if note, ok := maxRevenueMap[candle[0]]; ok {
			if minPrice < note.minPrice {
				note.minPrice = minPrice
				note.minPriceTime = t
			}
			if maxPrice < note.maxPrice {
				note.maxPrice = maxPrice
				note.maxPriceTime = t
			}
			note.income = note.maxPrice - note.maxPrice
			maxRevenueMap[candle[0]] = note

		} else {
			maxRevenueMap[candle[0]] = Candle{candle[0], maxPrice, minPrice, maxPrice - minPrice, t, t}
		}
	}
	return maxRevenueMap
}

func UserDeals(users[][]string)map[string]{
	userInfo := make(map[string]User)

	for _, user := range users{
		if note, ok := UserDeals[info[user[0]]]; ok {
			salePrice := stringToFloat(user[4])
			note.tickers[user[2]]["salePrice"]=salePrice
			note.tickers[user[2]]["income"]=note.tickers[user[2]]["salePrice"]-note.tickers[user[2]]["buyPrice"]

		}else {

		}
	}
}

// функция для чтения csv файла
func readFile(filename string) [][]string {
	File, _ := os.Open(filename)
	Reader := csv.NewReader(File)
	data, _ := Reader.ReadAll()
	return data
}

// функция для конвертирования string в float64
func stringToFloat(a string) float64 {
	b, _ := strconv.ParseFloat(a, 64)
	return b
}
