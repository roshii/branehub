package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TickerIndex",
		"GET",
		"/0/ticker",
		TickerIndex0,
	},
	Route{
		"ShowTicker",
		"GET",
		"/0/ticker/{market}",
		ShowTicker0,
	},
}
