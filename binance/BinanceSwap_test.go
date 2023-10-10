package binance

import (
	"fmt"
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
	t.Log(bs.PlaceFutureOrder("waht", goex.BTC_USDT, "", "8322", "0.01", "openlong", 0, 0))
}

func TestBinanceSwap_PlaceFutureOrder2(t *testing.T) {
	t.Log(bs.PlaceFutureOrder("wahtthefuck", goex.BTC_USDT,
		goex.SWAP_USDT_CONTRACT,
		"25999",
		"0.01",
		"openlong", 0, 1))
}

func TestBinanceSwap_GetFutureOrder(t *testing.T) {
	t.Log(bs.GetFutureOrder("197306870655", goex.BTC_USDT, goex.SWAP_USDT_CONTRACT))
}

func TestBinanceSwap_GetFutureOrderByCid(t *testing.T) {
	t.Log(bs.GetFutureOrderByCid("wahtthefuck", goex.BTC_USDT, goex.SWAP_USDT_CONTRACT))
}

func TestBinanceSwap_FutureCancelOrder(t *testing.T) {
	t.Log(bs.FutureCancelOrder(goex.BTC_USDT,
		goex.SWAP_USDT_CONTRACT,
		"197306870655"))
}

func TestBinanceSwap_GetFuturePosition(t *testing.T) {
	t.Log(bs.GetFuturePosition(goex.BTC_USDT, ""))
}

func TestBinanceIntegation(t *testing.T) {
	order, err := bs.PlaceFutureOrder("wahtthefuck", goex.BTC_USDT,
		goex.SWAP_USDT_CONTRACT,
		"25999",
		"0.01",
		"", 0, 1)
	fmt.Println(err)
	bs.FutureCancelOrder(goex.BTC_USDT,
		goex.SWAP_USDT_CONTRACT,
		order)
}

func TestBinanceIntegationCid(t *testing.T) {
	order, err := bs.PlaceFutureOrder("wahtthefuck34", goex.BTC_USDT,
		goex.SWAP_USDT_CONTRACT,
		"25999",
		"0.01",
		"openlong", 0, 1)
	fmt.Println(err)
	fmt.Println(bs.GetFutureOrderByCid("wahtthefuck34", goex.BTC_USDT, goex.SWAP_USDT_CONTRACT))

	bs.FutureCancelOrder(goex.BTC_USDT,
		goex.SWAP_USDT_CONTRACT,
		order)
	fmt.Println(bs.GetFutureOrderByCid("wahtthefuck34", goex.BTC_USDT, goex.SWAP_USDT_CONTRACT))

}
