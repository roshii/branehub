package kraken

import (
	"encoding/json"
	"testing"
)

func TestBtc2xbt(t *testing.T) {
	tests := [4]string{"abcbtc", "btcuhg", "btcbtcbc", "aabtczz"}
	results := [4]string{"abcxbt", "xbtuhg", "xbtxbtbc", "aaxbtzz"}
	var result string
	for i, v := range tests {
		result = btc2xbt([]rune(v))
		if results[i] != result {
			t.Errorf("%s did not replace BTC by XBT in %s", result, results[i])
		}
	}

}

func TestTicker(t *testing.T) {
	contents := `{"error":[],"result":{"XXBTZEUR":{"a":["9060.80000","12","12.000"],"b":["9058.70000","1","1.000"],"c":["9058.80000","0.00200000"],"v":["2071.77990748","3456.27572256"],"p":["8960.80495","8912.04839"],"t":[8853,16464],"l":["8882.60000","8700.00000"],"h":["9078.00000","9078.00000"],"o":"8887.80000"}}}`
	rawResult := Result{}
	json.Unmarshal([]byte(contents), &rawResult)
	rawTicker := rawTicker{}
	json.Unmarshal(rawResult.Result["XXBTZEUR"], &rawTicker)
	got := rawTicker.ticker()
	if got.Bid != 9058.7000 {
		t.Errorf("Bid = %.2f; want 9058.70", got)
	}
}
