package kraken

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"gitlab.com/braneproject/branehub/marketObservables"
)

const krakenURL string = "https://api.kraken.com"
const krakenVersion uint8 = 0

func (r rawTicker) ticker() marketObservables.Ticker {
	a, _ := strconv.ParseFloat(r.Ask[0], 32)
	b, _ := strconv.ParseFloat(r.Bid[0], 32)
	c, _ := strconv.ParseFloat(r.Last[0], 32)
	v, _ := strconv.ParseFloat(r.Volume[1], 32)
	l, _ := strconv.ParseFloat(r.Low[1], 32)
	h, _ := strconv.ParseFloat(r.High[1], 32)

	t := marketObservables.Ticker{
		Last:   float32(c),
		Bid:    float32(b),
		Ask:    float32(a),
		Volume: float32(v),
		Low:    float32(l),
		High:   float32(h),
	}
	return t
}

//NewKraken | Returns new Kraken struct
func NewKraken(apiPubkey string, apiPrivkey string) *Kraken {
	a := Kraken{
		url:     krakenURL,
		version: krakenVersion,
		pubkey:  apiPubkey,
		privkey: apiPrivkey,
	}
	return &a
}

//Error | Extends default Error struct
func (e Error) Error() string {
	return fmt.Sprintf("Message: %v", e.Data)
}

//requester | Creates the request to Kraken API
func (k Kraken) requester(call string, query string, params map[string]string) (Result, error) {

	//create empty Result
	result := Result{}

	//build url
	u, err := url.ParseRequestURI(k.url)

	//error handling
	if err != nil {
		return result, err
	}

	u.Path = "/" + fmt.Sprintf("%v", k.version) + "/" + call
	u.RawQuery = query

	apiCallURL := fmt.Sprintf("%v", u)

	//prepare params
	data := url.Values{}

	//convert params into querystring
	if len(params) > 0 {
		for k, p := range params {
			data.Set(k, p)
		}
	}

	//create request
	client := &http.Client{}
	r, err := http.NewRequest("GET", apiCallURL, bytes.NewBufferString(data.Encode()))

	//error handling
	if err != nil {
		return result, err
	}

	//do request
	res, err := client.Do(r)

	//error handling
	if err != nil {
		return result, err
	}

	//error handling
	if res.StatusCode != 200 {
		return result, fmt.Errorf("request didn't return a HTTP Status 200 but HTTP Status: %v", res.StatusCode)
	}

	//read request body
	contents, err := ioutil.ReadAll(res.Body)

	//parse json
	err = json.Unmarshal(contents, &result)

	//error handling
	if err != nil {
		return result, err
	}

	//handle Result error
	if result.Error[0] != "" {
		blerr := Error{result.Error[0]}
		json.Unmarshal(contents, &blerr)
		err = blerr
	}

	return result, err
}

//Remove "BTC" and replace with "XBT"
func btc2xbt(input []rune) string {
	btc, xbt := []rune("BTC"), []rune("XBT")
	for i := 0; i < (len(input) - 2); i++ {
		if input[i] == btc[0] && input[i+1] == btc[1] && input[i+2] == btc[2] {
			input[i], input[i+1], input[i+2] = xbt[0], xbt[1], xbt[2]
		}
	}
	return string(input)
}

//Public API

//GetTicker ...
func (k Kraken) GetTicker(market string) (marketObservables.Ticker, error) {

	call := "public/Ticker"
	market = btc2xbt([]rune(market))
	query := "pair=" + market

	// Removed multi market ticker
	// for i, v := range market {
	// 	if i == len(market)-1 {
	// 		call += v
	// 		break
	// 	} else {
	// 		call += v + ","
	// 	}
	// }

	ticker, err := k.requester(call, query, nil)

	rawResult := rawTicker{}

	market = "X" + market[0:3] + "Z" + market[3:6]

	if err == nil {
		err = json.Unmarshal(ticker.Result[market], &rawResult)
	}

	result := rawResult.ticker()

	return result, err
}

//ChannelTicker returns a standard Ticker to a channel
func (k Kraken) ChannelTicker(market string, c chan marketObservables.Ticker) {
	ticker, err := k.GetTicker(market)
	if err == nil {
		c <- ticker
	}
	close(c)
}
