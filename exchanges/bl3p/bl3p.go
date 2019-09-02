package bl3p

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

//Get market ticker
func (b Bl3p) GetTicker(market string) (marketObservables.Ticker, error) {

	call := market + "/ticker"

	result, err := b.tickerRequester(call, nil)

	return result.ticker(), err
}

//Retrieve the orderbook
func (b Bl3p) GetOrderbook(market string) (Orderbook, error) {

	call := market + "/orderbook"

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

	call := market + "/trades"

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

//AddOrder | Add new order to the orderbook
func (b Bl3p) AddOrder(orderType string, orderAmount int, orderPrice int) (interface{}, error) {

	price := strconv.FormatInt(int64(orderPrice), 10)
	amount := strconv.FormatInt(int64(orderAmount), 10)

	params := map[string]string{"type": orderType, "amount_int": amount, "price_int": price, "fee_currency": "BTC"}

	addOrder, err := b.requester("BTCEUR/money/order/add", params)

	result := AddOrder{}

	if err == nil {
		err = json.Unmarshal(addOrder.Data, &result)
	}

	return result, err
}

//WalletHistory | Retrieve your account transaction history
func (b Bl3p) WalletHistory(currency string) (Transactions, error) {

	params := map[string]string{"currency": currency, "recs_per_page": "25"}

	transactions, err := b.requester("GENMKT/money/wallet/history", params)

	result := Transactions{}

	if err == nil {
		err = json.Unmarshal(transactions.Data, &result)
	}

	return result, err
}

//CancelOrder | Cancel an open order
func (b Bl3p) CancelOrder(orderID int) (Result, error) {

	params := map[string]string{"order_id": strconv.FormatInt(int64(orderID), 10)}

	result, err := b.requester("BTCEUR/money/order/cancel", params)

	return result, err
}

//OrderInfo | Retrieve information about an order
func (b Bl3p) OrderInfo(orderID int) (Order, error) {

	params := map[string]string{"order_id": strconv.FormatInt(int64(orderID), 10)}

	order, err := b.requester("BTCEUR/money/order/result", params)

	result := Order{}

	if err == nil {
		err = json.Unmarshal(order.Data, &result)
	}

	return result, err
}

//GetAllActiveOrders | Retrieve all your open orders
func (b Bl3p) GetAllActiveOrders() (Orders, error) {

	allActiveOrders, err := b.requester("BTCEUR/money/orders", nil)

	result := Orders{}

	if err == nil {
		err = json.Unmarshal(allActiveOrders.Data, &result)
	}

	return result, err
}

//GetNewDepositAddress | Create a new bitcoin deposit address
func (b Bl3p) GetNewDepositAddress() (DepositAddress, error) {

	depositAddress, err := b.requester("BTCEUR/money/new_deposit_address", nil)

	result := DepositAddress{}

	if err == nil {
		err = json.Unmarshal(depositAddress.Data, &result)
	}

	return result, err
}

//GetLastDepositAddress | Retrieve the last created bitcoin deposit address
func (b Bl3p) GetLastDepositAddress() (DepositAddress, error) {

	depositAddress, err := b.requester("BTCEUR/money/deposit_address", nil)

	result := DepositAddress{}

	if err == nil {
		err = json.Unmarshal(depositAddress.Data, &result)
	}

	return result, err
}

//GetInfo | Get account info
func (b Bl3p) GetInfo() (Info, error) {

	info, err := b.requester("GENMKT/money/info", nil)

	result := Info{}

	if err == nil {
		err = json.Unmarshal(info.Data, &result)
	}

	return result, err
}
