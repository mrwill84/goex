package okex

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mrwill84/goex/okex/bst"

	. "github.com/mrwill84/goex"
	"github.com/mrwill84/goex/internal/logger"
	cmap "github.com/orcaman/concurrent-map"
)

type OKExV3SwapWs struct {
	base           *OKEx
	v3Ws           *OKExV3Ws
	tickerCallback func(*FutureTicker)
	depthCallback  func(*Depth)
	tradeCallback  func(*Trade, string)
	klineCallback  func(*FutureKline, int)
}

type OkexSnapshot struct {
	Asks bst.BSTree
	Bids bst.BSTree
}

var container cmap.ConcurrentMap = cmap.New()

func NewOKExV3SwapWs(base *OKEx) *OKExV3SwapWs {
	okV3Ws := &OKExV3SwapWs{
		base: base,
	}
	okV3Ws.v3Ws = NewOKExV3Ws(base, okV3Ws.handle)
	return okV3Ws
}

func (okV3Ws *OKExV3SwapWs) TickerCallback(tickerCallback func(*FutureTicker)) {
	okV3Ws.tickerCallback = tickerCallback
}

func (okV3Ws *OKExV3SwapWs) DepthCallback(depthCallback func(*Depth)) {
	okV3Ws.depthCallback = depthCallback
}

func (okV3Ws *OKExV3SwapWs) TradeCallback(tradeCallback func(*Trade, string)) {
	okV3Ws.tradeCallback = tradeCallback
}

func (okV3Ws *OKExV3SwapWs) KlineCallback(klineCallback func(*FutureKline, int)) {
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3SwapWs) SetCallbacks(tickerCallback func(*FutureTicker),
	depthCallback func(*Depth),
	tradeCallback func(*Trade, string),
	klineCallback func(*FutureKline, int)) {
	okV3Ws.tickerCallback = tickerCallback
	okV3Ws.depthCallback = depthCallback
	okV3Ws.tradeCallback = tradeCallback
	okV3Ws.klineCallback = klineCallback
}

func (okV3Ws *OKExV3SwapWs) getChannelName(currencyPair CurrencyPair, contractType string) string {
	var (
		contractId string
	)

	if contractType == SWAP_CONTRACT {
		contractId = fmt.Sprintf("%s-SWAP", currencyPair.ToSymbol("-"))
	} else {

		contractId = okV3Ws.base.OKExFuture.GetFutureContractId(currencyPair, contractType)
		//	logger.Info("contractid=", contractId)
	}

	return contractId
}

func (okV3Ws *OKExV3SwapWs) SubscribeDepth(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.depthCallback == nil {
		return errors.New("please set depth callback func")
	}
	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}
	ps := []map[string]string{}

	ps = append(ps, map[string]string{
		"channel": "books",
		"instId":  chName,
	})

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": ps,
	})
}

func (okV3Ws *OKExV3SwapWs) SubscribeTicker(currencyPair CurrencyPair, contractType string) error {
	if okV3Ws.tickerCallback == nil {
		return errors.New("please set ticker callback func")
	}

	chName := okV3Ws.getChannelName(currencyPair, contractType)
	if chName == "" {
		return errors.New("subscribe error, get channel name fail")
	}
	ps := []map[string]string{}

	ps = append(ps, map[string]string{
		"channel": "tickers",
		"instId":  chName,
	})

	return okV3Ws.v3Ws.Subscribe(map[string]interface{}{
		"op":   "subscribe",
		"args": ps,
	})
}

func (okV3Ws *OKExV3SwapWs) SubscribeTrade(currencyPair CurrencyPair, contractType string) error {
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

func (okV3Ws *OKExV3SwapWs) SubscribeKline(currencyPair CurrencyPair, contractType string, period int) error {
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

func (okV3Ws *OKExV3SwapWs) getContractAliasAndCurrencyPairFromInstrumentId(instrumentId string) (alias string, pair CurrencyPair) {
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

func AcceptSnap(snapshot *OkexSnapshot, depth *depthResponse) {

	for _, itm := range depth.Asks {
		if itm[3] == "0" { //delete
			snapshot.Asks.Delete(itm[0])
		} else {
			snapshot.Asks.Upsert(itm[0], &DepthRecord{
				Price:  ToFloat64(itm[0]),
				Amount: ToFloat64(itm[1]),
				Slots:  ToInt64(itm[3]),
			})
		}
	}
	for _, itm := range depth.Bids {
		if itm[3] == "0" { //delete
			snapshot.Bids.Delete(itm[0])
		} else {
			snapshot.Bids.Upsert(itm[0], &DepthRecord{
				Price:  ToFloat64(itm[0]),
				Amount: ToFloat64(itm[1]),
				Slots:  ToInt64(itm[3]),
			})
		}
	}
}

func (okV3Ws *OKExV3SwapWs) handle(resp *wsResp) error {

	//fmt.Println("channel string, instId", channel, instId)
	var (
		err           error
		ch            string
		tickers       []tickerResponse
		depthResp     []depthResponse
		dep           Depth
		tradeResponse []struct {
			Side         string  `json:"side"`
			TradeId      int64   `json:"trade_id,string"`
			Price        float64 `json:"price,string"`
			Qty          float64 `json:"qty,string"`
			InstrumentId string  `json:"instrument_id"`
			Timestamp    string  `json:"timestamp"`
		}
		klineResponse []struct {
			Candle       []string `json:"candle"`
			InstrumentId string   `json:"instrument_id"`
		}
	)
	channel := resp.Arg.Channel
	data := resp.Data
	if strings.Contains(channel, "futures/candle") ||
		strings.Contains(channel, "swap/candle") {
		ch = "candle"
	} else {
		ch, err = okV3Ws.v3Ws.parseChannel(channel)
		if err != nil {
			logger.Errorf("[%s] parse channel err=%s ,  originChannel=%s", okV3Ws.base.GetExchangeName(), err, ch)
			return nil
		}
	}

	switch ch {
	case "tickers":
		err = json.Unmarshal(data, &tickers)
		if err != nil {
			return err
		}

		for _, t := range tickers {
			alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(t.InstrumentId)
			//date, _ := time.Parse(time.RFC3339, t.Timestamp)
			//fmt.Println(" t.Last,", t.Last)
			date, _ := strconv.ParseUint(t.Timestamp, 10, 64)
			okV3Ws.tickerCallback(&FutureTicker{
				Ticker: &Ticker{
					Pair: pair,
					Last: t.Last,
					Buy:  t.BestBid,
					Sell: t.BestAsk,
					Vol:  t.Volume24h,
					Date: date,
				},
				ContractId:   t.InstrumentId,
				ContractType: alias,
			})

		}
		return nil
	case "candle":
		err = json.Unmarshal(data, &klineResponse)
		if err != nil {
			return err
		}

		for _, t := range klineResponse {
			_, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(t.InstrumentId)
			ts, _ := time.Parse(time.RFC3339, t.Candle[0])
			//granularity := adaptKLinePeriod(KlinePeriod(period))
			okV3Ws.klineCallback(&FutureKline{
				Kline: &Kline{
					Pair:      pair,
					High:      ToFloat64(t.Candle[2]),
					Low:       ToFloat64(t.Candle[3]),
					Timestamp: ts.Unix(),
					Open:      ToFloat64(t.Candle[1]),
					Close:     ToFloat64(t.Candle[4]),
					Vol:       ToFloat64(t.Candle[5]),
				},
				Vol2: ToFloat64(t.Candle[6]),
			}, 1)
		}
		return nil
	case "books":
		err := json.Unmarshal(data, &depthResp)
		if err != nil {
			logger.Error(err)
			return err
		}
		if len(depthResp) == 0 {
			return nil
		}
		//alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(depthResp[0].InstrumentId)
		//dep.Pair = pair
		//channel string, instId
		instId := resp.Arg.InstId
		dep.ContractType = "SWAP" //alias
		dep.ContractId = instId   //depthResp[0].InstrumentId
		dep.Action = resp.Action
		ts := ToInt64(depthResp[0].Timestamp)
		dep.UTime = time.Unix(ts/1000, ts%1000) //time.Parse(time.RFC3339, depthResp[0].Timestamp)

		i, ok := container.Get(instId)
		var snapshot *OkexSnapshot
		if ok {
			if s, ok2 := i.(*OkexSnapshot); ok2 {
				snapshot = s
			}
		} else {
			snapshot = new(OkexSnapshot)
			container.Set(instId, snapshot)
		}

		AcceptSnap(snapshot, &depthResp[0])

		for i := range snapshot.Asks.Iter() {
			r := i.Val.(*DepthRecord)
			dep.AskList = append(dep.AskList, *r)
		}
		for i := range snapshot.Bids.RIter() {
			r := i.Val.(*DepthRecord)
			dep.BidList = append(dep.BidList, *r)
		}
		//call back func
		okV3Ws.depthCallback(&dep)
		return nil
	case "trade":
		err := json.Unmarshal(data, &tradeResponse)
		if err != nil {
			logger.Error("unmarshal error :", err)
			return err
		}

		for _, resp := range tradeResponse {
			alias, pair := okV3Ws.getContractAliasAndCurrencyPairFromInstrumentId(resp.InstrumentId)

			tradeSide := SELL
			switch resp.Side {
			case "buy":
				tradeSide = BUY
			}

			t, err := time.Parse(time.RFC3339, resp.Timestamp)
			if err != nil {
				logger.Warn("parse timestamp error:", err)
			}

			okV3Ws.tradeCallback(&Trade{
				Tid:    resp.TradeId,
				Type:   tradeSide,
				Amount: resp.Qty,
				Price:  resp.Price,
				Date:   t.Unix(),
				Pair:   pair,
			}, alias)
		}
		return nil
	}

	return fmt.Errorf("[%s] unknown websocket message: %s", ch, string(data))
}

func (okV3Ws *OKExV3SwapWs) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}
