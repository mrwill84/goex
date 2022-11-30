package localexchange

import (
	"encoding/json"
	"testing"
)

var leconfig = LE{
	wsAddress: "ws://127.0.0.1:3001/ws",
}

func HHH(channel string, data json.RawMessage) error {
	return nil
}

func TestKuCoin_GetTicker(t *testing.T) {

	wx := NewLocalExchangeWs(&leconfig, HHH)
	ps := []map[string]string{}

	ps = append(ps, map[string]string{
		"channel": "tickers",
		"instId":  "BTC-USDT",
	})

	wx.Subscribe(
		map[string]interface{}{
			"op":   "subscribe",
			"args": ps,
		},
	)
}
