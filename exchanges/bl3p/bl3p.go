package bl3p

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gitlab.com/braneproject/branehub/marketObservables"
)

const bl3pURL string = "https://api.bl3p.eu"
const bl3pVersion uint8 = 1

//NewBl3p | Returns new Bl3p struct
func NewBl3p(apiPubkey string, apiPrivkey string) *Bl3p {
	a := Bl3p{
		url:     bl3pURL,
		version: bl3pVersion,
		pubkey:  apiPubkey,
		privkey: apiPrivkey,
	}
	return &a
}

//Error | Extends default Error struct
func (e Error) Error() string {
	return fmt.Sprintf("Message: %v: Code: %v", e.Data.Message, e.Data.Code)
}

//requester | Creates the request to Bl3p API
func (b Bl3p) requester(call string, params map[string]string) (Result, error) {

	//create empty bl3presult
	result := Result{}

	//build url
	u, err := url.ParseRequestURI(b.url)

	//error handling
	if err != nil {
		return result, err
	}

	u.Path = "/" + fmt.Sprintf("%v", b.version) + "/" + call
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
		return result, fmt.Errorf("%s returned HTTP Status: %v", apiCallURL, res.StatusCode)
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
	if result.Result == "error" {
		blerr := Error{}
		json.Unmarshal(contents, &blerr)
		err = blerr
	}

	return result, err
}

//requester | Creates the request to Bl3p API
func (b Bl3p) tickerRequester(call string, params map[string]string) (rawTicker, error) {

	//create empty bl3presult
	result := rawTicker{}

	//build url
	u, err := url.ParseRequestURI(b.url)

	//error handling
	if err != nil {
		return result, err
	}

	u.Path = "/" + fmt.Sprintf("%v", b.version) + "/" + call
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
		return result, fmt.Errorf("%s returned HTTP Status: %v", apiCallURL, res.StatusCode)
	}

	//read request body
	contents, err := ioutil.ReadAll(res.Body)

	//parse json
	err = json.Unmarshal(contents, &result)

	//error handling
	if err != nil {
		return result, err
	}

	return result, err
}

func (r rawTicker) ticker() marketObservables.Ticker {
	t := marketObservables.Ticker{
		Ask:    r.Ask,
		Bid:    r.Bid,
		Last:   r.Last,
		Volume: r.Volume.Daily,
		Low:    r.Low,
		High:   r.High,
	}
	return t
}

//Public API

//GetTicker ...
func (b Bl3p) GetTicker(market string) (marketObservables.Ticker, error) {

	call := strings.ToUpper(market) + "/ticker"

	result, err := b.tickerRequester(call, nil)

	if err != nil {
		fmt.Println(err)
		return marketObservables.Ticker{}, err
	}

	return result.ticker(), err
}

//ChannelTicker returns a standard Ticker to a channel
func (b Bl3p) ChannelTicker(market string, c chan marketObservables.Ticker) {

	ticker, err := b.GetTicker(market)
	if err == nil {
		c <- ticker
	}
	close(c)
}

//GetOrderbook ...
func (b Bl3p) GetOrderbook(market string) (Orderbook, error) {

	call := strings.ToUpper(market) + "/orderbook"

	orderbook, err := b.requester(call, nil)

	result := Orderbook{}

	if err == nil {
		err = json.Unmarshal(orderbook.Data, &result)
	}

	return result, err
}

//Retrieve the last 1000 trades or the last 1000 trades after the specified tradeID
func (b Bl3p) GetLast1000Trades(market string, tradeID int) (Trades, error) {
	var trades Result
	var err error

	call := strings.ToUpper(market) + "/trades"

	if tradeID != 0 {
		params := map[string]string{"trade_id": strconv.FormatInt(int64(tradeID), 10)}
		trades, err = b.requester(call, params)
	} else {
		trades, err = b.requester(call, nil)
	}

	result := Trades{}

	if err == nil {
		err = json.Unmarshal(trades.Data, &result)
	}

	return result, err
}

//GetTradeHistory
//TODO Implement GetTradeHistory
