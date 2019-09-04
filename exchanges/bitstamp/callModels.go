package bitstamp

//Bitstamp struct
type Bitstamp struct {
	url     string
	pubkey  string
	privkey string
}

//Ticker | Ticker call struct
type rawTicker struct {
	High      float32 `json:"high"`
	Last      float32 `json:"last"`
	Timestamp int32   `json:"timestamp"`
	Bid       float32 `json:"bid"`
	VWAP      float32 `json:"vwap"`
	Volume    float32 `json:"volume"`
	Low       float32 `json:"low"`
	Ask       float32 `json:"ask"`
	Open      float32 `json:"open"`
}
