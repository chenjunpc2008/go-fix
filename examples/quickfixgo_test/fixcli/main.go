package main

import (
    "bytes"
    "flag"
    "fmt"
    "github.com/quickfixgo/quickfix"
    "io/ioutil"
    "os"
    "os/signal"
    "path"
    "sync/atomic"
    "syscall"
)

//tradeClient implements the quickfix.Application interface
type tradeClient struct {
}

func newTradeClient() *tradeClient {
	return &tradeClient{}
}

//OnCreate implemented as part of Application interface
func (e *tradeClient) OnCreate(sessionID quickfix.SessionID) {}

//OnLogon implemented as part of Application interface
func (e *tradeClient) OnLogon(sessionID quickfix.SessionID) {}

//OnLogout implemented as part of Application interface
func (e *tradeClient) OnLogout(sessionID quickfix.SessionID) {}

//FromAdmin implemented as part of Application interface
func (e *tradeClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return nil
}

//ToAdmin implemented as part of Application interface
func (e *tradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {}

//ToApp implemented as part of Application interface
func (e *tradeClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	//fmt.Printf("Sending %s\n", msg)
	return
}

//FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e *tradeClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	//fmt.Printf("FromApp: %s\n", msg.String())
    atomic.AddInt64(&GRspRcvNum, 1)
	return
}


func main() {
    flag.Int64Var(&maxPerSecond, "max", 100, "max req per second")
    flag.Parse()
    fmt.Println("maxPerSecond=", maxPerSecond)
	fmt.Println("fix client start...")

	cfgFileName := path.Join("./" + "fixcli.cfg")
	cfg, err := os.Open(cfgFileName)
	if nil != err {
		fmt.Printf("Error opening %v, %v\n", cfgFileName, err)
		return
	}
	defer cfg.Close()

	stringData, readErr := ioutil.ReadAll(cfg)
	if nil != readErr {
		fmt.Printf("Error reading cfg: %s,", readErr)
		return
	}

	appSettings, err := quickfix.ParseSettings(bytes.NewReader(stringData))
	if nil != err {
		fmt.Printf("Error reading cfg: %s,", err)
		return
	}

	logFactory, err := quickfix.NewFileLogFactory(appSettings)
	if nil != err {
		fmt.Printf("Error NewFileLogFactory: %s,", err)
		return
	}

	storFactory := quickfix.NewFileStoreFactory(appSettings)
	app := newTradeClient()

	initiator, err := quickfix.NewInitiator(app, storFactory, appSettings, logFactory)
	if nil != err {
		fmt.Printf("Unable to create Initiator: %s\n", err)
		return
	}

	err = initiator.Start()
	if err != nil {
		fmt.Printf("Unable to start Initiator: %s\n", err)
		return
	}

	go orderNewSingle()
	go packetRcvTpsCalc()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	initiator.Stop()

	return
}

