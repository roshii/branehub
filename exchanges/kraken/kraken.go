package kraken

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
)

const krakenUrl string = "https://api.kraken.com"
const krakenVersion uint8 = 0

func (r rawTicker) ticker() Ticker {
	a, _ := strconv.ParseFloat(r.Ask[0], 32)
	b, _ := strconv.ParseFloat(r.Bid[0], 32)
	c, _ := strconv.ParseFloat(r.Last[0], 32)
	v, _ := strconv.ParseFloat(r.Volume[1], 32)
	p, _ := strconv.ParseFloat(r.VWA[1], 32)
	l, _ := strconv.ParseFloat(r.Low[1], 32)
	h, _ := strconv.ParseFloat(r.High[1], 32)

	t := Ticker{
		Ask:    float32(a),
		Bid:    float32(b),
		Last:   float32(c),
		Volume: float32(v),
		VWA:    float32(p),
		Low:    float32(l),
		High:   float32(h),
	}
	return t
}

//NewKraken | Returns new Kraken struct
func NewKraken(apiPubkey string, apiPrivkey string) *Kraken {
	a := Kraken{
		url:     krakenUrl,
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
func (b Kraken) requester(call string, query string, params map[string]string) (Result, error) {

	//create empty Result
	result := Result{}

	//build url
	u, err := url.ParseRequestURI(b.url)

	//error handling
	if err != nil {
		return result, err
	}

	u.Path = "/" + fmt.Sprintf("%v", b.version) + "/" + call
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

	//request body
	body := []byte(call + string(0) + data.Encode())

	//decode privkey
	base64Decode := make([]byte, base64.StdEncoding.DecodedLen(len(b.privkey)))
	l, err := base64.StdEncoding.Decode(base64Decode, []byte(b.privkey))

	//error handling
	if err != nil {
		return result, err
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

//Public API

//Get market ticker
func (b Kraken) GetTicker(market string) (Ticker, error) {

	call := "public/Ticker"
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

	ticker, err := b.requester(call, query, nil)

	rawResult := rawTicker{}

	market = "X" + market[0:3] + "Z" + market[3:6]

	if err == nil {
		err = json.Unmarshal(ticker.Result[market], &rawResult)
	}

	result := rawResult.ticker()

	return result, err
}
