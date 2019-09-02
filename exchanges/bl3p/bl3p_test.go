package bl3p

import (
	"encoding/json"
	"testing"
)

func TestTicker(t *testing.T) {
	contents := `{"currency":"BTC","last":8988,"bid":8966.02,"ask":8988,"high":8990.68,"low":8697.37999,"timestamp":1567425456,"volume":{"24h":52.3852199,"30d":2669.61017749}}`
	result := rawTicker{}
	json.Unmarshal([]byte(contents), &result)
	got := result.ticker()
	if got.Bid != 8966.02 {
		t.Errorf("Bid = %.2f; want 8966.02", got)
	}
}
