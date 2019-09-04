package main

import (
	"gitlab.com/braneproject/branehub/exchanges/bitstamp"
	"gitlab.com/braneproject/branehub/exchanges/bl3p"
	"gitlab.com/braneproject/branehub/exchanges/kraken"
)

type priceIndex struct {
	Market string  `json:"market"`
	VWAP   float32 `json:"vwap"`
}

// vwap calculates the volume-weighted average price
// formula: sum(num shares * share price)/(total shares)
// input: [[volume, price], ...]
func vwap(positions ...[2]float32) float32 {
	var sum, total float32
	for _, x := range positions {
		sum += x[0] * x[1]
		total += x[0]
	}

	return sum / total
}

// BranePriceIndex returns the volume weighted average for `market`
func BranePriceIndex(market string) float32 {

	if market != "BTCEUR" {
		return 0
	}

	bl3p := bl3p.NewBl3p("", "")
	ticker, _ := bl3p.GetTicker(market)
	bl3pTick := [2]float32{ticker.Volume, ticker.Last}
	// fmt.Println("@BL3P Last: ", bl3pTick[1])

	kraken := kraken.NewKraken("", "")
	ticker, _ = kraken.GetTicker(market)
	krakenTick := [2]float32{ticker.Volume, ticker.Last}
	// fmt.Println("@Kraken Last: ", krakenTick[1])

	bitstamp := bitstamp.NewBitstamp("", "")
	ticker, _ = bitstamp.GetTicker(market)
	bitstampTick := [2]float32{ticker.Volume, ticker.Last}
	// fmt.Println("@Bitstamp Last: ", bitstampTick[1])

	average := vwap(bl3pTick, krakenTick, bitstampTick)
	// fmt.Println("BTC/EUR Average: ", average)

	return average

}
