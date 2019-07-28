//branehub.go
package main

import (
	"branehub/exchanges/bl3p"
	"branehub/exchanges/kraken"
	_ "branehub/ui"
	"fmt"
)

var market string = "BTCEUR"

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

func main() {
	fmt.Println("*** BraneHub ***")

	market := "BTCEUR"

	bl3p := bl3p.NewBl3p("", "")
	ticker, _ := bl3p.GetTicker(market)
	bl3pTick := [2]float32{ticker.Volume, ticker.Last}
	fmt.Println("@BL3P Last: ", bl3pTick[1])

	kraken := kraken.NewKraken("", "")
	ticker, _ = kraken.GetTicker(market)
	krakenTick := [2]float32{ticker.Volume, ticker.Last}
	fmt.Println("@Kraken Last: ", krakenTick[1])

	average := vwap(bl3pTick, krakenTick)
	fmt.Println("BTC/EUR Average: ", average)

}
