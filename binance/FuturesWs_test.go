package binance

import (
	"fmt"
	"testing"
	"time"

	"github.com/mrwill84/goex"
)

var futuresWs *FuturesWs

func createFuturesWs() {
	//os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:1080")

	futuresWs = NewFuturesWs()

	futuresWs.DepthCallback(func(depth *goex.Depth) {
		fmt.Println(depth)
	})

	futuresWs.TickerCallback(func(ticker *goex.FutureTicker) {
		fmt.Println(ticker.Ticker, ticker.ContractType)
	})

	futuresWs.TradeCallback(func(trade *goex.Trade) {
		fmt.Println("slots", trade.Slots)
	})

}

func TestFuturesWs_DepthCallback(t *testing.T) {
	createFuturesWs()

	//futuresWs.SubscribeDepth(goex.LTC_USDT, goex.SWAP_USDT_CONTRACT)
	futuresWs.SubscribeDepth(goex.BTC_USDT, goex.SWAP_USDT_CONTRACT)
	//futuresWs.SubscribeDepth(goex.LTC_USDT, goex.QUARTER_CONTRACT)

	time.Sleep(30 * time.Second)
}

func TestFuturesWs_SubscribeTicker(t *testing.T) {
	createFuturesWs()

	//futuresWs.SubscribeTicker(goex.BTC_USDT, goex.SWAP_USDT_CONTRACT)
	futuresWs.SubscribeTicker(goex.BTC_USDT, goex.SWAP_CONTRACT)
	//futuresWs.SubscribeTicker(goex.BTC_USDT, goex.QUARTER_CONTRACT)

	time.Sleep(30 * time.Second)
}

func TestFuturesWs_TradeCallback(t *testing.T) {
	createFuturesWs()

	//futuresWs.SubscribeDepth(goex.LTC_USDT, goex.SWAP_USDT_CONTRACT)
	futuresWs.SubscribeTrade(goex.BTC_USDT, goex.SWAP_USDT_CONTRACT)
	//futuresWs.SubscribeDepth(goex.LTC_USDT, goex.QUARTER_CONTRACT)

	time.Sleep(30 * time.Second)
}
