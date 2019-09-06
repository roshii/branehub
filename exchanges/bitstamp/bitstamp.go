package bitstamp

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

const bitstampURL string = "https://www.bitstamp.net"

//NewBitstamp | Returns new Bitstamp struct
func NewBitstamp(apiPubkey string, apiPrivkey string) *Bitstamp {
	a := Bitstamp{
		url:     bitstampURL,
		pubkey:  apiPubkey,
		privkey: apiPrivkey,
	}
	return &a
}

//requester | Creates the request to Bitstamp API
func (b Bitstamp) requester(call string, params map[string]string) ([]byte, error) {

	u, err := url.ParseRequestURI(b.url)
	if err != nil {
		return []byte{}, err
	}

	u.Path = "/api/v2/" + call
	apiCallURL := fmt.Sprintf("%v", u)

	//prepare and convert params into querystring
	data := url.Values{}
	if len(params) > 0 {
		for k, p := range params {
			data.Set(k, p)
		}
	}

	//create request
	client := &http.Client{}
	r, err := http.NewRequest("GET", apiCallURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return []byte{}, err
	}

	//do request
	res, err := client.Do(r)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != 200 {
		return []byte{}, fmt.Errorf("Request did return a HTTP status %v", res.StatusCode)
	}

	return ioutil.ReadAll(res.Body)
}

func (r rawTicker) ticker() marketObservables.Ticker {
	a, _ := strconv.ParseFloat(r.Ask, 32)
	b, _ := strconv.ParseFloat(r.Bid, 32)
	c, _ := strconv.ParseFloat(r.Last, 32)
	v, _ := strconv.ParseFloat(r.Volume, 32)
	l, _ := strconv.ParseFloat(r.Low, 32)
	h, _ := strconv.ParseFloat(r.High, 32)

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

//GetTicker returns a standard Ticker for `market`
func (b Bitstamp) GetTicker(market string) (marketObservables.Ticker, error) {

	call := "ticker/" + market

	contents, err := b.requester(call, nil)

	raw := rawTicker{}

	if err == nil {
		err = json.Unmarshal(contents, &raw)
		if err != nil {
			return marketObservables.Ticker{}, err
		}
	}

	return raw.ticker(), err
}

//ChannelTicker returns a standard Ticker to a channel
func (b Bitstamp) ChannelTicker(market string, c chan marketObservables.Ticker) {
	ticker, err := b.GetTicker(market)
	if err == nil {
		c <- ticker
	}
	close(c)
}
