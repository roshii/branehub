//observables.go
package marketObservables

//Volume, Low and High are for the last 24h
type Ticker struct {
	Last   float32
	Bid    float32
	Ask    float32
	Volume float32
	Low    float32
	High   float32
}
