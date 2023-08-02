package okex

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	. "github.com/mrwill84/goex"
	"github.com/mrwill84/goex/internal/logger"
)

type OKExV3FuturesWs struct {
	base           *OKEx
	v3Ws           *OKExV3Ws
	tickerCallback func(*FutureTicker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade, string)
	klineCallback  func(*FutureKline, int)
}

func NewOKExV3FuturesWs(base *OKEx) *OKExV3FuturesWs {
	okV3Ws := &OKExV3FuturesWs{
		base: base,
	}
	okV3Ws.v3Ws = NewOKExV3Ws(base, okV3Ws.handle)
	return okV3Ws
}

func (okV3Ws *OKExV3FuturesWs) TickerCallback(tickerCallback func(*FutureTicker)) {
	okV3Ws.tickerCallback = tickerCallback
}

func (okV3Ws *OKExV3FuturesWs) DepthCallback(depthCallback func(*Depth)) {
	okV3Ws.depthCallback = depthCallback
}

func (okV3Ws *OKExV3FuturesWs) TradeCallback(tradeCallback func(*Trade, string)) {
	okV3Ws.tradeCallback = tradeCallback
}

func (okV3Ws *OKExV3FuturesWs) KlineCallback(klineCallback func(*FutureKline, int)) {
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3FuturesWs) SetCallbacks(tickerCallback func(*FutureTicker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade, string),
	klineCallback func(*FutureKline, int)) {
	okV3Ws.tickerCallback = tickerCallback
	okV3Ws.depthCallback = depthCallback
	okV3Ws.tradeCallback = tradeCallback
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3FuturesWs) getChannelName(currencyPair CurrencyPair, contractType string) string {
	var (
		prefix      string
		contractId  string
		channelName string
	)

	if contractType == SWAP_CONTRACT {
		prefix = "swap"
		contractId = fmt.Sprintf("%s-SWAP", currencyPair.ToSymbol("-"))
	} else {
		prefix = "futures"
		contractId = okV3Ws.base.OKExFuture.GetFutureContractId(currencyPair, contractType)
		//	logger.Info("contractid=", contractId)
	}

	if contractId == "" {
		return ""
	}

	channelName = prefix + "/%s:" + contractId

	return channelName
}

func (okV3Ws *OKExV3FuturesWs) SubscribeDepth(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "books")}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeTicker(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "ticker")}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeTrade(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tradeCallback == nil {
		return errors.New("please set trade callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, "trade")}})
}

func (okV3Ws *OKExV3FuturesWs) SubscribeKline(currencyPair CurrencyPair, contractType string, period int) error {
	if okV3Ws.klineCallback == nil {
		return errors.New("place set kline callback func")
	}

	seconds := adaptKLinePeriod(KlinePeriod(period))
	if seconds == -1 {
		return fmt.Errorf("unsupported kline period %d in okex", period)
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": []string{fmt.Sprintf(chName, fmt.Sprintf("candle%ds", seconds))}})
}

func (okV3Ws *OKExV3FuturesWs) getContractAliasAndCurrencyPairFromInstrumentId(instrumentId string) (alias string, pair CurrencyPair) {
	if strings.HasSuffix(instrumentId, "SWAP") {
		ar := strings.Split(instrumentId, "-")
		return instrumentId, NewCurrencyPair2(fmt.Sprintf("%s_%s", ar[0], ar[1]))
	} else {
		contractInfo, err := okV3Ws.base.OKExFuture.GetContractInfo(instrumentId)
		if err != nil {
			logger.Error("instrument id invalid:", err)
			return "", UNKNOWN_PAIR
		}
		alias = contractInfo.Alias
		pair = NewCurrencyPair2(fmt.Sprintf("%s_%s", contractInfo.UnderlyingIndex, contractInfo.QuoteCurrency))
		return alias, pair
	}
}

func (okV3Ws *OKExV3FuturesWs) handle(*wsResp) error {

	return fmt.Errorf("unknown websocket message:")
}

func (okV3Ws *OKExV3FuturesWs) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}
