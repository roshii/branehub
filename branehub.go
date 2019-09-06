//branehub.go
package main

import (
	"log"
	"net/http"
)

var prices []internalIndex

func main() {

	prices = []internalIndex{}

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":80", router))
}
