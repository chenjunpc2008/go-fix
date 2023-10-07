package main

import (
    "fmt"
    "sync/atomic"

    "github.com/quickfixgo/quickfix"

    "github.com/chenjunpc2008/go-fix/fixinitiator"
)

type tradeCliApp struct {
}

var _ fixinitiator.ApplicationIF = (*tradeCliApp)(nil)

func (app *tradeCliApp) OnConnected(serverIP string, serverPort uint16) {

}

func (app *tradeCliApp) OnError(msg string, err error) {
    fmt.Println(msg, err)
}

func (app *tradeCliApp) OnErrorStr(msg string) {
    fmt.Println(msg)
}

func (app *tradeCliApp) OnEvent(msg string) {
    fmt.Println(msg)
}

func (app *tradeCliApp) OnAdminMsg(fixmsg *quickfix.Message) {
    fmt.Println(fixmsg)
}

func (app *tradeCliApp) OnAppMsg(fixmsg *quickfix.Message) {
    atomic.AddInt64(&GRspRcvNum, 1)
    fmt.Println(fixmsg)
}

func (app *tradeCliApp) OnDisconnected(serverIP string, serverPort uint16) {

}
func (app *tradeCliApp) OnSendedData(msg *fixinitiator.MsgPkg) {

}
