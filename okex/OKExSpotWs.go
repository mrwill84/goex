package okex

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	. "github.com/mrwill84/goex"
)

type OKExV3SpotWs struct {
	base           *OKEx
	v3Ws           *OKExV3Ws
	tickerCallback func(*Ticker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade)
	klineCallback  func(*Kline, KlinePeriod)
}

func NewOKExSpotV3Ws(base *OKEx) *OKExV3SpotWs {
	okV3Ws := &OKExV3SpotWs{
		base: base,
	}
	okV3Ws.v3Ws = NewOKExV3Ws(base, okV3Ws.handle)
	return okV3Ws
}

func (okV3Ws *OKExV3SpotWs) TickerCallback(tickerCallback func(*Ticker)) {
	okV3Ws.tickerCallback = tickerCallback
}

func (okV3Ws *OKExV3SpotWs) DepthCallback(depthCallback func(*Depth)) {
	okV3Ws.depthCallback = depthCallback
}

func (okV3Ws *OKExV3SpotWs) TradeCallback(tradeCallback func(*Trade)) {
	okV3Ws.tradeCallback = tradeCallback
}

func (okV3Ws *OKExV3SpotWs) KLineCallback(klineCallback func(kline *Kline, period KlinePeriod)) {
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3SpotWs) SetCallbacks(tickerCallback func(*Ticker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade),
	klineCallback func(*Kline, KlinePeriod)) {
	okV3Ws.tickerCallback = tickerCallback
	okV3Ws.depthCallback = depthCallback
	okV3Ws.tradeCallback = tradeCallback
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3SpotWs) SubscribeDepth(currencyPair CurrencyPair) error {
	if okV3Ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf("spot/depth5:%s", currencyPair.ToSymbol("-"))}})
}

func (okV3Ws *OKExV3SpotWs) SubscribeTicker(currencyPair CurrencyPair) error {
	if okV3Ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf("spot/ticker:%s", currencyPair.ToSymbol("-"))}})
}

func (okV3Ws *OKExV3SpotWs) SubscribeTrade(currencyPair CurrencyPair) error {
	if okV3Ws.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}
	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf("spot/trade:%s", currencyPair.ToSymbol("-"))}})
}

func (okV3Ws *OKExV3SpotWs) SubscribeKline(currencyPair CurrencyPair, period int) error {
	if okV3Ws.klineCallback == nil {
		return errors.New("place set kline callback func")
	}

	seconds := adaptKLinePeriod(KlinePeriod(period))
	if seconds == -1 {
		return fmt.Errorf("unsupported kline period %d in okex", period)
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf("spot/candle%ds:%s", seconds, currencyPair.ToSymbol("-"))}})
}

func (okV3Ws *OKExV3SpotWs) getCurrencyPair(instrumentId string) CurrencyPair {
	return NewCurrencyPair3(instrumentId, "-")
}

func (okV3Ws *OKExV3SpotWs) handle(*wsResp) error {

	return fmt.Errorf("unknown websocket message:")
}

func (okV3Ws *OKExV3SpotWs) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}
