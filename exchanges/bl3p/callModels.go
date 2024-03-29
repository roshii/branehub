package bl3p

import "encoding/json"

//Bl3pResult | Main result struct
type Result struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data"`
}

//Orderbook | Orderbook call struct
type Orderbook struct {
	Bids []OrderbookItem
	Asks []OrderbookItem
}

//OrderbookItem | Orderbook item struct
type OrderbookItem struct {
	Count     int   `json:"count"`
	PriceInt  int64 `json:"price_int"`
	AmountInt int64 `json:"amount_int"`
}

//Orders | Order array struct
type Orders struct {
	Order []Order `json:"orders"`
}

//Trades | Trades array struct
type Trades struct {
	Trade []Trade `json:"trades"`
}

//Order | Order struct
type Order struct {
	OrderID             int64     `json:"order_id"`
	Label               string    `json:"label"`
	Currency            string    `json:"currency"`
	Item                string    `json:"item"`
	Type                string    `json:"type"`
	Status              string    `json:"status"`
	Date                int64     `json:"date"`
	Amount              AmountObj `json:"amount"`
	AmountExecuted      AmountObj `json:"amount_executed"`
	AmountFunds         AmountObj `json:"amount_funds"`
	AmountFundsExecuted AmountObj `json:"amount_funds_executed"`
	Price               AmountObj `json:"price"`
	TotalAmount         AmountObj `json:"total_amount"`
	TotalSpent          AmountObj `json:"total_spent"`
	TotalFee            AmountObj `json:"total_fee"`
	AvgCost             AmountObj `json:"avg_cost"`
	Trades              []Trade   `json:"trades"`
}

//Trade | Trade struct
type Trade struct {
	TradeID   int64     `json:"trade_id"`
	Date      int64     `json:"date"`
	Currency  string    `json:"currency"`
	Amount    AmountObj `json:"amount"`
	Price     AmountObj `json:"price"`
	AmountInt int64     `json:"amount_int"`
	Item      string    `json:"item"`
	PriceInt  int64     `json:"price_int"`
}

//AmountObj | AmountObj struct
type AmountObj struct {
	ValueInt     string `json:"value_int"`
	DisplayShort string `json:"display_short"`
	Display      string `json:"display"`
	Currency     string `json:"currency"`
	Value        string `json:"value"`
}

//DepositAddress | DepositAddress call struct
type DepositAddress struct {
	Address string `json:"address"`
}

//Transactions | WalletHistory call struct
type Transactions struct {
	Page         int64         `json:"page"`
	Records      int64         `json:"records"`
	MaxPpage     int64         `json:"max_page"`
	Transactions []Transaction `json:"transactions"`
}

//Transaction | Transaction struct
type Transaction struct {
	TransactionID int64     `json:"transaction_id"`
	Amount        AmountObj `json:"amount"`
	Date          int64     `json:"date"`
	DebitCredit   string    `json:"debit_credit"`
	Price         AmountObj `json:"price"`
	OrderID       int64     `json:"order_id"`
	Type          string    `json:"type"`
	Balance       AmountObj `json:"balance"`
	TradeID       int64     `json:"trade_id"`
	ContraAmount  AmountObj `json:"contra_amount"`
	Fee           AmountObj `json:"fee"`
}

//AddOrder | AddOrder call struct
type AddOrder struct {
	OrderID int64 `json:"order_id"`
}

//Wallet | Wallet struct
type Wallet struct {
	Balance   AmountObj `json:"balance"`
	Available AmountObj `json:"available"`
}

//Info | Account call struct
type Info struct {
	UserID  int64             `json:"user_id"`
	Wallets map[string]Wallet `json:"wallets"`
}

//Ticker | Ticker call struct
type rawTicker struct {
	Currency  string  `json:"currency"`
	Last      float32 `json:"last"`
	Bid       float32 `json:"bid"`
	Ask       float32 `json:"ask"`
	High      float32 `json:"high"`
	Low       float32 `json:"low"`
	Timestamp int32   `json:"timestamp"`
	Volume    volume  `json:"volume"`
}

//Ticker | Ticker call struct
type volume struct {
	Daily   float32 `json:"24h"`
	Monthly float32 `json:"30d"`
}

//Bl3p struct
type Bl3p struct {
	url     string
	version uint8
	pubkey  string
	privkey string
}

//Error struct
type Error struct {
	Result string `json:"result"`
	Data   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"data"`
}
