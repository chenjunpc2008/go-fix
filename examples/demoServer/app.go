package main

import (
    "fmt"

    "github.com/quickfixgo/quickfix"

    "github.com/chenjunpc2008/go-fix/fixacceptor"
)

type myApp struct {
}

var _ fixacceptor.ApplicationIF = (*myApp)(nil)

func (a *myApp) OnError(msg string, err error) {
    fmt.Printf(msg, err)
}

func (a *myApp) OnErrorStr(msg string) {
    fmt.Printf(msg)
}

func (app *myApp) OnEvent(msg string) {
    fmt.Println(msg)
}

func (app *myApp) FromAdmin(fixMsg *quickfix.Message) {
    fmt.Println("FromAdmin:", fixMsg)
}

func (app *myApp) FromApp(fixMsg *quickfix.Message) {
    fmt.Println("FromApp:", fixMsg)
}

func (app *myApp) OnSendedData(fixMsg *quickfix.Message) {
    fmt.Println("OnSent:", fixMsg)
}

func (app *myApp) OnDisconnected(serverIP string, serverPort uint16) {
    fmt.Printf("OnDisconnected:%s:%d/n", serverIP, serverPort)
}
