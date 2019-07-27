//branehub.go
package main

import (
	"branehub/exchanges/bl3p"
	"branehub/exchanges/kraken"
	_ "branehub/ui"
	"fmt"
)

var market string = "BTCEUR"

func main() {
	market := "BTCEUR"
	fmt.Println("*** BraneHub ***")
	ex1 := bl3p.NewBl3p("", "")
	r1, _ := ex1.GetTicker(market)

	fmt.Println("@BL3P Last: ", r1.Last)

	market = "XBTEUR"
	ex2 := kraken.NewKraken("", "")
	r2, _ := ex2.GetTicker(market)

	fmt.Println("@Kraken Last: ", r2.Last)

	average := (r1.Last + r2.Last) / 2
	fmt.Println("BTC/EUR Average: ", average)
}
