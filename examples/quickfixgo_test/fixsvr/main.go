package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
    "runtime"
    "sync/atomic"
	"syscall"
	"time"

	"github.com/quickfixgo/quickfix"
)

type executor struct {
	orderID int64
	execID  int64
}

func newExecutor() *executor {
	return &executor{}
}


//quickfix.Application interface
func (e executor) OnCreate(sessionID quickfix.SessionID)                           {}
func (e executor) OnLogon(sessionID quickfix.SessionID)                            {}
func (e executor) OnLogout(sessionID quickfix.SessionID)                           {}
func (e executor) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID)     {}
func (e executor) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error { return nil }
func (e executor) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

//Use Message Cracker on Incoming Application Messages
func (e *executor) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {

	//fmt.Printf("%s\n", msg.String())

	atomic.AddInt64(&GReqProcNum, 1)

	return nil
}

func main() {
    runtime.GOMAXPROCS(20)
	fmt.Println("fix server")

	cfgFileName := path.Join("./" + "fixsvr.cfg")
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
	app := newExecutor()

	acceptor, err := quickfix.NewAcceptor(app, storFactory, appSettings, logFactory)
	if nil != err {
		fmt.Printf("Unable to create Acceptor: %s\n", err)
		return
	}

	err = acceptor.Start()
	if nil != err {
		fmt.Printf("Unable to start Acceptor: %s\n", err)
		return
	}

	go SvrTpsTest()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	acceptor.Stop()

	return
}

var (
    GReqProcNum int64
)

func SvrTpsTest() {
    var (
        timeout   = time.Duration(10) * time.Second
        sTime     string
        startNums int64
        endNums   int64
        delta     int64
    )

    for {
        startNums = atomic.AddInt64(&GReqProcNum, 0)
        select {
        case <-time.After(timeout):
        }

        sTime = time.Now().Format("2006-01-02 15:04:05")

        endNums = atomic.AddInt64(&GReqProcNum, 0)

        delta = endNums - startNums

        // 下单操作每个请求回复两条msg
        fmt.Printf("Until:%s, 10s total process orders:%d, tps:%d\n", sTime, delta, delta/10)
    }
}
