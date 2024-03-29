package binance

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mrwill84/goex"
	"github.com/mrwill84/goex/internal/logger"
)

type FuturesWs struct {
	base  *BinanceFutures
	fOnce sync.Once
	dOnce sync.Once

	wsBuilder *goex.WsBuilder
	f         *goex.WsConn
	d         *goex.WsConn

	depthCallFn  func(depth *goex.Depth)
	tickerCallFn func(ticker *goex.FutureTicker)
	tradeCalFn   func(trade *goex.Trade)
}

func NewFuturesWs() *FuturesWs {
	futuresWs := new(FuturesWs)

	futuresWs.wsBuilder = goex.NewWsBuilder().
		ProtoHandleFunc(futuresWs.handle).AutoReconnect()

	httpCli := &http.Client{
		Timeout: 10 * time.Second,
	}

	if os.Getenv("HTTPS_PROXY") != "" {
		httpCli = &http.Client{
			Transport: &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) {
					return url.Parse(os.Getenv("HTTPS_PROXY"))
				},
			},
			Timeout: 10 * time.Second,
		}
	}

	futuresWs.base = NewBinanceFutures(&goex.APIConfig{
		HttpClient: httpCli,
	})

	return futuresWs
}

func (s *FuturesWs) connectUsdtFutures() {
	s.fOnce.Do(func() {
		s.f = s.wsBuilder.WsUrl("wss://fstream.binance.com/ws").Build()
	})
}

func (s *FuturesWs) connectFutures() {
	s.dOnce.Do(func() {
		s.d = s.wsBuilder.WsUrl("wss://dstream.binance.com/ws").Build()
	})
}

func (s *FuturesWs) DepthCallback(f func(depth *goex.Depth)) {
	s.depthCallFn = f
}

func (s *FuturesWs) TickerCallback(f func(ticker *goex.FutureTicker)) {
	s.tickerCallFn = f
}

func (s *FuturesWs) TradeCallback(f func(trade *goex.Trade)) {
	s.tradeCalFn = f
}

func (s *FuturesWs) SubscribeDepth(pair goex.CurrencyPair, contractType string) error {
	switch contractType {
	case goex.SWAP_USDT_CONTRACT:
		s.connectUsdtFutures()
		return s.f.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{pair.AdaptUsdToUsdt().ToLower().ToSymbol("") + "@depth10@100ms"},
			Id:     1,
		})
	default:
		s.connectFutures()
		sym, _ := s.base.adaptToSymbol(pair.AdaptUsdtToUsd(), contractType)
		return s.d.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{strings.ToLower(sym) + "@depth20@100ms"},
			Id:     2,
		})
	}
	return errors.New("contract is error")
}

func (s *FuturesWs) SubscribeTicker(pair goex.CurrencyPair, contractType string) error {
	switch contractType {
	case goex.SWAP_USDT_CONTRACT:
		s.connectUsdtFutures()
		return s.f.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{pair.AdaptUsdToUsdt().ToLower().ToSymbol("") + "@miniTicker"},
			Id:     1,
		})
	default:
		s.connectFutures()
		sym, _ := s.base.adaptToSymbol(pair.AdaptUsdtToUsd(), contractType)
		return s.d.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{strings.ToLower(sym) + "@ticker"},
			Id:     2,
		})
	}
	return errors.New("contract is error")
}

func (s *FuturesWs) SubscribeTrade(pair goex.CurrencyPair, contractType string) error {
	///panic("implement me")
	switch contractType {
	case goex.SWAP_USDT_CONTRACT:
		s.connectUsdtFutures()
		return s.f.Subscribe(req{
			Method: "SUBSCRIBE",
			Params: []string{pair.AdaptUsdToUsdt().ToLower().ToSymbol("") + "@aggTrade"},
			Id:     1,
		})
	}
	return nil
}

func (s *FuturesWs) handle(data []byte) error {
	var m = make(map[string]interface{}, 4)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	if e, ok := m["e"].(string); ok && e == "depthUpdate" {
		dep := s.depthHandle(m["b"].([]interface{}), m["a"].([]interface{}))
		dep.ContractType = "SWAP" // m["s"].(string)

		symbol, ok := m["s"].(string)

		if ok {
			dep.Pair = adaptSymbolToCurrencyPair(symbol)
		} else {
			dep.Pair = adaptSymbolToCurrencyPair(dep.ContractType) //usdt swap
		}
		dep.ContractId = symbol[:len(symbol)-4] + "-USDT-SWAP"
		dep.Timestamp = goex.ToInt64(m["T"])
		dep.Exchange = "BINANCE"
		s.depthCallFn(dep)

		return nil
	}

	if e, ok := m["e"].(string); ok && e == "24hrMiniTicker" {
		s.tickerCallFn(s.tickerHandle(m))
		return nil
	}
	//fmt.Println("m", m)
	if e, ok := m["e"].(string); ok && e == "aggTrade" {
		s.tradeCalFn(s.aggTradeHandle(m))
		return nil
	}

	logger.Warn("unknown ws response:", string(data))

	return nil
}

func (s *FuturesWs) depthHandle(bids []interface{}, asks []interface{}) *goex.Depth {
	var dep goex.Depth

	for _, item := range bids {
		bid := item.([]interface{})
		dep.BidList = append(dep.BidList,
			goex.DepthRecord{
				Price:  goex.ToFloat64(bid[0]),
				Amount: goex.ToFloat64(bid[1]),
			})
	}

	for _, item := range asks {
		ask := item.([]interface{})
		dep.AskList = append(dep.AskList, goex.DepthRecord{
			Price:  goex.ToFloat64(ask[0]),
			Amount: goex.ToFloat64(ask[1]),
		})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return &dep
}

func (s *FuturesWs) aggTradeHandle(m map[string]interface{}) *goex.Trade {
	var trade goex.Trade

	symbol, ok := m["s"].(string)

	if ok {
		trade.Pair = adaptSymbolToCurrencyPair(symbol)
	} else {
		trade.Pair = adaptSymbolToCurrencyPair(m["s"].(string)) //usdt swap
	}

	trade.ContractType = "SWAP"
	sym := m["s"].(string)
	trade.Date = goex.ToInt64(m["T"])
	trade.Price = goex.ToFloat64(m["p"])
	trade.Tid = goex.ToInt64(m["a"])
	trade.Amount = goex.ToFloat64(m["q"])
	trade.Type = goex.SELL
	if goex.ToBool(m["m"]) == false {
		trade.Type = goex.BUY
	}
	trade.Slots = goex.ToInt64(m["l"]) - goex.ToInt64(m["f"])
	trade.Exchange = "BINANCE"
	trade.ContractId = sym[:len(sym)-4] + "-USDT-SWAP"
	//fmt.Println("trade", trade)
	return &trade
}

func (s *FuturesWs) tickerHandle(m map[string]interface{}) *goex.FutureTicker {
	var ticker goex.FutureTicker
	ticker.Ticker = new(goex.Ticker)

	symbol, ok := m["s"].(string)
	if ok {
		ticker.Pair = adaptSymbolToCurrencyPair(symbol)
	} else {
		ticker.Pair = adaptSymbolToCurrencyPair(m["s"].(string)) //usdt swap
	}

	ticker.ContractType = "SWAP"
	sym := m["s"].(string)
	ticker.Date = goex.ToUint64(m["E"])
	ticker.High = goex.ToFloat64(m["h"])
	ticker.Low = goex.ToFloat64(m["l"])
	ticker.Last = goex.ToFloat64(m["c"])
	ticker.Vol = goex.ToFloat64(m["v"])
	ticker.Exchange = "BINANCE"
	ticker.ContractId = sym[:len(sym)-4] + "-USDT-SWAP"
	return &ticker
}
