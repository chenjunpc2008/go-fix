# go-fix
The Financial Information Exchange (FIX) Protocol Engine


## fixacceptor
Go FIX server

## fixinitiator
Go FIX client

# Usage

## fixacceptor
---
for example: ```example/demoServer```

1. Create a struct for server event call back, and put your own message process in the callback functions.
    ```go
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
    ```
2. Use the server and go
    ```go
    var setting = fixacceptor.DefaultConfig()
    setting.Port = 5001
    setting.SenderCompID = "NASDAQ"
    setting.FixDataDictPath = "./spec/FIX42.xml"

    app := new(myApp)

    acceptor := fixacceptor.NewAcceptor(setting, app)

    acceptor.Start()
    ```


## fixinitiator
---
for example: ```example/fixTradeClient```

1. Create a struct for server event call back, and put your own message process in the callback functions.
    ```go
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
    ```
2. Use the client and go
    ```go
    var setting = fixinitiator.DefaultConfig()
    setting.FixDataDictPath = "./spec/FIX42.xml"
    setting.RemoteIP = "127.0.0.1"
    setting.RemotePort = 5001

    setting.BeginString = "FIX.4.2"
    setting.SenderCompID = "qaib84"
    setting.TargetCompID = "IB"

    app := new(tradeCliApp)

    gFixCli = fixinitiator.NewInitiator(setting, app)

    var (
        senderStartSeqNum uint64 = 1
        targetStartSeqNum uint64 = 1
    )
    gFixCli.Start(senderStartSeqNum, targetStartSeqNum)
    ```
3. Send a new trade order
   ```go
    fixOrder = buildIBNewOrderSingle()
    mpkg.Fixmsg = fixOrder
    gFixCli.SendMsgToSvr(mpkg)
   ```
