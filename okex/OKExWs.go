package okex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nntaoli-project/goex/internal/logger"

	. "github.com/nntaoli-project/goex"
)

type wsResp struct {
	Arg struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId"`
	} `json:"arg"`
	Event string `json:"event"`
	Data  json.RawMessage
}

type OKExV3Ws struct {
	base *OKEx
	*WsBuilder
	once       *sync.Once
	WsConn     *WsConn
	respHandle func(channel string, data json.RawMessage) error
}

type OkexLoginArg struct {
	APIKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

type OKExLogin struct {
	OP   string         `json:"op"`
	Args []OkexLoginArg `json:"args"`
}

func NewOKExV3Ws(base *OKEx, handle func(channel string, data json.RawMessage) error) *OKExV3Ws {
	okV3Ws := &OKExV3Ws{
		once:       new(sync.Once),
		base:       base,
		respHandle: handle,
	}
	okV3Ws.WsBuilder = NewWsBuilder().
		WsUrl("wss://ws.okx.com:8443/ws/v5/public").
		ReconnectInterval(time.Second).
		AutoReconnect().
		Heartbeat(func() []byte { return []byte("ping") }, 28*time.Second).
		DecompressFunc(FlateDecompress).
		ProtoHandleFunc(okV3Ws.handle)
		/*.
		ConnectSuccessAfterSendMessage(func() (msg []byte) {
			sign, timestamp := base.doParamSign("GET", "/users/self/verify", "")

			loginMsg := OKExLogin{
				OP: "login",
				Args: []OkexLoginArg{
					{
						APIKey:     base.config.ApiKey,
						Passphrase: basre.config.ApiPassphrase,
						Timestamp:  timestamp,
						Sign:       sign,
					},
				},
			}
			bLoginMsg, err := json.Marshal(loginMsg)
			if err != nil {
				return msg
			}
			return bLoginMsg
		})*/
	return okV3Ws
}

func (okV3Ws *OKExV3Ws) clearChan(c chan wsResp) {
	for {
		if len(c) > 0 {
			<-c
		} else {
			break
		}
	}
}

func (okV3Ws *OKExV3Ws) getTablePrefix(currencyPair CurrencyPair, contractType string) string {
	if contractType == SWAP_CONTRACT {
		return "swap"
	}
	return "futures"
}

func (okV3Ws *OKExV3Ws) ConnectWs() {
	okV3Ws.once.Do(func() {
		okV3Ws.WsConn = okV3Ws.WsBuilder.Build()
	})
}

func (okV3Ws *OKExV3Ws) parseChannel(channel string) (string, error) {
	return channel, nil
}

func (okV3Ws *OKExV3Ws) getKlinePeriodFormChannel(channel string) int {
	metas := strings.Split(channel, ":")
	if len(metas) != 2 {
		return 0
	}
	i, _ := strconv.ParseInt(metas[1], 10, 64)
	return int(i)
}

func (okV3Ws *OKExV3Ws) handle(msg []byte) error {

	logger.Info("[ws] [response] ", string(msg))
	if string(msg) == "pong" {
		return nil
	}

	var wsResp wsResp
	err := json.Unmarshal(msg, &wsResp)
	if err != nil {
		logger.Error(err)
		return err
	}

	if wsResp.Event != "" {
		switch wsResp.Event {
		case "subscribe":
			logger.Info("subscribed:", wsResp.Arg.Channel)
			return nil
		case "error":
			logger.Errorf("fuck?", string(msg))
		default:
			logger.Info(string(msg))
		}
		return fmt.Errorf("unknown websocket message: %v", wsResp)
	}

	err = okV3Ws.respHandle(wsResp.Arg.Channel, wsResp.Data)
	if err != nil {
		logger.Error("handle ws data error:", err)
	}
	return err

	return fmt.Errorf("unknown websocket message: %v", wsResp)
}

func (okV3Ws *OKExV3Ws) Subscribe(sub map[string]interface{}) error {
	okV3Ws.ConnectWs()
	logger.Info("[ws] [response] ", sub)
	return okV3Ws.WsConn.Subscribe(sub)
}
