package okex

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/mrwill84/goex"
	"github.com/mrwill84/goex/internal/logger"
)

func init() {
	logger.SetLevel(logger.DEBUG)
}

func TestNewOKExV3SwapWs(t *testing.T) {
	//os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:1080")
	ok := NewOKEx(&goex.APIConfig{
		HttpClient: http.DefaultClient,
	})
	ok.OKExV3SwapWs.TickerCallback(func(ticker *goex.FutureTicker) {
		t.Log(ticker.Ticker, ticker.ContractType)
	})
	ok.OKExV3SwapWs.DepthCallback(func(depth *goex.Depth) {
		t.Log(depth)
	})
	ok.OKExV3SwapWs.TradeCallback(func(trade *goex.Trade) {
		fmt.Println(trade)
		//t.Log(trade)
	})
	ok.OKExV3SwapWs.SubscribeTrade(goex.BTC_USDT, goex.SWAP_CONTRACT)
	//ok.OKExV3SwapWs.SubscribeTicker(goex.BTC_USDT, goex.SWAP_CONTRACT)
	time.Sleep(1 * time.Minute)
}
