package main

import (
	"fmt"

	"gitlab.com/braneproject/branehub/exchanges/bitstamp"
	"gitlab.com/braneproject/branehub/exchanges/bl3p"
	"gitlab.com/braneproject/branehub/exchanges/kraken"
	"gitlab.com/braneproject/branehub/marketObservables"
)

type priceIndex struct {
	Market string  `json:"market"`
	VWAP   float32 `json:"vwap"`
}

type internalIndex struct {
	market string
	vwap   float32
	last   int64
}

func (i internalIndex) externalize() priceIndex {
	return priceIndex{
		Market: i.market,
		VWAP:   i.vwap,
	}
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
	if total == 0 {
		return 0
	}
	return sum / total
}

// BranePriceIndex returns the volume weighted average for `market`
func BranePriceIndex(market string) float32 {

	bitstamp := bitstamp.NewBitstamp("", "")
	bitstampChannel := make(chan marketObservables.Ticker)
	go bitstamp.ChannelTicker(market, bitstampChannel)

	kraken := kraken.NewKraken("", "")
	krakenChannel := make(chan marketObservables.Ticker)
	go kraken.ChannelTicker(market, krakenChannel)

	var bitstampTicker, krakenTicker, bl3pTicker marketObservables.Ticker

	if market == "BTCEUR" || market == "LTCEUR" {
		bl3pChannel := make(chan marketObservables.Ticker)
		bl3p := bl3p.NewBl3p("", "")
		go bl3p.ChannelTicker(market, bl3pChannel)
		bitstampTicker, krakenTicker, bl3pTicker = <-bitstampChannel, <-krakenChannel, <-bl3pChannel
	} else {
		bitstampTicker, krakenTicker = <-bitstampChannel, <-krakenChannel
		bl3pTicker = marketObservables.Ticker{}
	}

	bitstampTick := [2]float32{bitstampTicker.Volume, bitstampTicker.Last}
	krakenTick := [2]float32{krakenTicker.Volume, krakenTicker.Last}
	bl3pTick := [2]float32{bl3pTicker.Volume, bl3pTicker.Last}

	average := vwap(bl3pTick, krakenTick, bitstampTick)

	fmt.Println("@Bitstamp Last: ", bitstampTick[1])
	fmt.Println("@Kraken Last: ", krakenTick[1])
	fmt.Println("@BL3P Last: ", bl3pTick[1])
	fmt.Println("Average: ", average)

	return average

}
