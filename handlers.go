package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Brane's Hub!")
	fmt.Fprintf(w, "Please refer to https://gitlab.com/braneproject/branehub/blob/master/README.md")
}

func TickerIndex0(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Ticker v0 endpoint")
	fmt.Fprintf(w, "Please refer to https://gitlab.com/braneproject/branehub/blob/master/README.md")
}

func ShowTicker0(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	market := vars["market"]

	// Look for market in price list and tag if old
	ref, old := -1, true
	for i, v := range prices {
		if v.market == market {
			ref = i
			fmt.Println("ref updated")
			if v.last+2 > time.Now().Unix() {
				old = false
			}
		}
	}

	// If market has no entry or is old -> add/update entry
	if ref == -1 || old == true {
		average := BranePriceIndex(market)
		freshIndex := internalIndex{
			market: market,
			vwap:   average,
			last:   time.Now().Unix(),
		}
		if ref == -1 {
			prices = append(prices, freshIndex)
			ref = len(prices) - 1
		} else if old == true {
			prices[ref] = freshIndex
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(prices[ref].externalize()); err != nil {
		panic(err)
	}
}
