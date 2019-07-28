package kraken

import "encoding/json"

//Result | Main result struct
type Result struct {
	Error  [1]string                  `json:"error"`
	Result map[string]json.RawMessage `json:"result"` // Rest of the fields should go here.
}

type rawTicker struct {
	Ask        [3]string `json:"a"`
	Bid        [3]string `json:"b"`
	Last       [2]string `json:"c"`
	Volume     [2]string `json:"v"`
	VWA        [2]string `json:"p"`
	TradeCount [2]int32  `json:"t"`
	Low        [2]string `json:"l"`
	High       [2]string `json:"h"`
	Open       string    `json:"o"`
}

type Ticker struct {
	Ask    float32
	Bid    float32
	Last   float32
	Volume float32
	VWA    float32
	Low    float32
	High   float32
}

//Kraken struct
type Kraken struct {
	url     string
	version uint8
	pubkey  string
	privkey string
	nonce   uint64
	otp     uint32
}

//Error struct
type Error struct {
	Data string
}
