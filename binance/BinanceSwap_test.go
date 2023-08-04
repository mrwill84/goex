package binance

import (
	"net"
	"net/http"
	"testing"
	"time"

	goex "github.com/mrwill84/goex"
)

var bs = NewBinanceSwap(&goex.APIConfig{
	HttpClient: &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
		},
		Timeout: 10 * time.Second,
	},
	ApiKey:       "",
	ApiSecretKey: "",
})

func TestBinanceSwap_Ping(t *testing.T) {
	bs.Ping()
}

func TestBinanceSwap_GetFutureDepth(t *testing.T) {
	t.Log(bs.GetFutureDepth(goex.BTC_USDT, "", 1))
}

func TestBinanceSwap_GetFutureIndex(t *testing.T) {
	t.Log(bs.GetFutureIndex(goex.BTC_USDT))
}

func TestBinanceSwap_GetKlineRecords(t *testing.T) {
	//kline, err := bs.GetKlineRecords("", goex.BTC_USDT, goex.KLINE_PERIOD_4H, 1, 0)
	//t.Log(err, kline[0].Kline)
}

func TestBinanceSwap_GetTrades(t *testing.T) {
	t.Log(bs.GetTrades("", goex.BTC_USDT, 0))
}

func TestBinanceSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(bs.GetFutureUserinfo())
}

func TestBinanceSwap_PlaceFutureOrder(t *testing.T) {
	t.Log(bs.PlaceFutureOrder(goex.BTC_USDT, "", "8322", "0.01", goex.OPEN_BUY, 0, 0))
}

func TestBinanceSwap_PlaceFutureOrder2(t *testing.T) {
	t.Log(bs.PlaceFutureOrder(goex.BTC_USDT,
		goex.SWAP_USDT_CONTRACT,
		"28000",
		"0.01",
		goex.OPEN_BUY, 0, 100))
}

func TestBinanceSwap_GetFutureOrder(t *testing.T) {
	t.Log(bs.GetFutureOrder("177597851702", goex.BTC_USDT, goex.SWAP_USDT_CONTRACT))
}

func TestBinanceSwap_FutureCancelOrder(t *testing.T) {
	t.Log(bs.FutureCancelOrder(goex.BTC_USDT,
		goex.SWAP_USDT_CONTRACT,
		"177597851702"))
}

func TestBinanceSwap_GetFuturePosition(t *testing.T) {
	t.Log(bs.GetFuturePosition(goex.BTC_USDT, ""))
}
