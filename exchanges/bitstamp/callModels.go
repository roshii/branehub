package bitstamp

//Bitstamp struct
type Bitstamp struct {
	url     string
	pubkey  string
	privkey string
}

//Ticker | Ticker call struct
type rawTicker struct {
	High      string `json:"high"`
	Last      string `json:"last"`
	Timestamp string `json:"timestamp"`
	Bid       string `json:"bid"`
	VWAP      string `json:"vwap"`
	Volume    string `json:"volume"`
	Low       string `json:"low"`
	Ask       string `json:"ask"`
	Open      string `json:"open"`
}
