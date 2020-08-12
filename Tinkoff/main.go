package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Candle struct {
	ticker       string
	maxPrice     float64
	minPrice     float64
	maxPriceTime time.Time
	minPriceTime time.Time
}

func main() {
	candles, users := readFiles()
	fmt.Println(candles)
	fmt.Println(users)
}

//функция для чтения обоих csv файлов
func readFiles() ([][]string, [][]string) {
	candlesFile, _ := os.Open("candles_5m.csv")
	userFile, _ := os.Open("user_trades.csv")
	candlesReader := csv.NewReader(candlesFile)
	userReader := csv.NewReader(userFile)
	candles, _ := candlesReader.ReadAll()
	users, _ := userReader.ReadAll()
	return candles, users
}

func stringToFloat(a string) float64 {
	b, _ := strconv.ParseFloat(a, 64)
	return b
}
