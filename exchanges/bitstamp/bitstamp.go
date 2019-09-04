package bitstamp

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
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

	//build url
	u, err := url.ParseRequestURI(b.url)

	//error handling
	if err != nil {
		return []byte{}, err
	}

	u.Path = "/" + call
	apiCallURL := fmt.Sprintf("%v", u)

	//prepare params
	data := url.Values{}

	//convert params into querystring
	if len(params) > 0 {
		for k, p := range params {
			data.Set(k, p)
		}
	}

	fmt.Println(apiCallURL)

	//create request
	client := &http.Client{}
	r, err := http.NewRequest("GET", apiCallURL, bytes.NewBufferString(data.Encode()))

	//error handling
	if err != nil {
		return []byte{}, err
	}

	//request body
	body := []byte(call + string(0) + data.Encode())

	//decode privkey
	base64Decode := make([]byte, base64.StdEncoding.DecodedLen(len(b.privkey)))
	l, err := base64.StdEncoding.Decode(base64Decode, []byte(b.privkey))

	//error handling
	if err != nil {
		return []byte{}, err
	}

	decodedPrivkey := []byte(base64Decode[:l])

	//sign
	h := hmac.New(sha512.New, decodedPrivkey)
	h.Write(body)
	sign := h.Sum(nil)

	//encode signature
	encodedSign := string(base64.StdEncoding.EncodeToString([]byte(sign)))

	//add headers for authentication
	r.Header.Add("Rest-Key", b.pubkey)
	r.Header.Add("Rest-Sign", encodedSign)

	//do request
	res, err := client.Do(r)

	//error handling
	if err != nil {
		return []byte{}, err
	}

	//error handling
	if res.StatusCode != 200 {
		return []byte{}, fmt.Errorf("request didn't return a HTTP Status 200 but HTTP Status: %v", res.StatusCode)
	}

	//read request body
	contents, err := ioutil.ReadAll(res.Body)

	return contents, err
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

	call := "api/v2/ticker/" + market

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
