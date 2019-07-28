//branehub.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Brane's Hub!")
}

func TickerIndex0(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Ticker v0 endpoint")
}

func ShowTicker0(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	market := vars["market"]

	price := priceIndex{
		Market: market,
		VWAP:   BranePriceIndex(market),
	}

	json.NewEncoder(w).Encode(price)
}
