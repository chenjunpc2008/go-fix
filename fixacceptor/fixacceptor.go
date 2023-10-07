package fixacceptor

import (
    "fmt"

    "github.com/quickfixgo/quickfix/datadictionary"

    "github.com/chenjunpc2008/go-tcp/tcpserver"
    "github.com/chenjunpc2008/go/util/kchanthreadpool"
)

const (
    ciTaskPoolBuffSize = 2 * 10000
)

// Acceptor acceptor
type Acceptor struct {
    setting      Settings
    server       *tcpserver.Ctcpsvr
    appHdl       ApplicationIF
    cliSessnCtrl *clientSessionCtrlSt
    chExit       chan int // exit channel

    // thread pool for process admin msg
    adminMsgTPool *kchanthreadpool.ThreadPool
    // thread pool for process app msg
    appMsgTPool *kchanthreadpool.ThreadPool
    // dict data model for parsing fix42 message
    fixDataModelDict *datadictionary.DataDictionary
}

// NewAcceptor new
func NewAcceptor(set Settings, app ApplicationIF) *Acceptor {

    var ac = &Acceptor{setting: set, appHdl: app}

    return ac
}

// GetNeededProcessNum get gorountine numbers needed
func GetNeededProcessNum() int {
    return 6
}

// Start start acceptor server
func (ac *Acceptor) Start() error {

    var (
        err error
        cnf = tcpserver.DefaultConfig()
    )

    ac.fixDataModelDict, err = datadictionary.Parse(ac.setting.FixDataDictPath)
    if err != nil {
        return err
    }

    hdler := &acceptorHandler{appHdl: ac.appHdl, acceptor: ac}

    cnf.SendBuffsize = ac.setting.SendBuffsize
    cnf.AsyncReceive = ac.setting.AsyncReceive
    cnf.RequireSendedCb = true
    cnf.AsyncSended = ac.setting.AsyncSended

    ac.server = tcpserver.NewTCPSvr(hdler, cnf)
    ac.cliSessnCtrl = newClientSessionCtrlSt()
    ac.chExit = make(chan int)

    // thread pool
    tpHdler := newThreadPoolHandlerSt(ac)
    ac.adminMsgTPool, err = kchanthreadpool.NewThreadPool(1, ciTaskPoolBuffSize, tpHdler)
    if nil != err {
        return err
    }

    err = ac.adminMsgTPool.Start()
    if nil != err {
        return err
    }

    // thread pool
    ac.appMsgTPool, err = kchanthreadpool.NewThreadPool(4, ciTaskPoolBuffSize, tpHdler)
    if nil != err {
        return err
    }

    err = ac.appMsgTPool.Start()
    if nil != err {
        return err
    }

    // tcp server
    // tcp server
    err = ac.server.StartServer(ac.setting.Port)
    if nil != err {
        return err
    }

    go heartbeatWatcher(ac.chExit, ac.cliSessnCtrl, ac)

    return nil
}

// Stop stop acceptor server
func (ac *Acceptor) Stop() {
    if nil != ac.server {
        close(ac.chExit)

        ac.server.StopServer()

        ac.adminMsgTPool.Stop()
        ac.appMsgTPool.Stop()

        ac.server = nil
        ac.cliSessnCtrl = nil
    }
}

// close client
func (ac *Acceptor) closeCli(clientID uint64, reason string) {
    // const ftag = "Acceptor.closeCli()"

    // incase of closed
    select {
    case <-ac.chExit:
        return

    default:
    }

    if nil != ac.server {
        ac.server.CloseClient(clientID, reason)
    }
}

// SendToClient send to client
func (ac *Acceptor) SendToClient(clientID uint64, msg interface{}) error {
    // const ftag = "Acceptor.SendToClient()"

    // incase of closed
    select {
    case <-ac.chExit:
        return fmt.Errorf("acceptor closed")

    default:
    }

    if nil == ac.server {
        return fmt.Errorf("acceptor un-initial")
    }

    busy, err := ac.server.SendToClient(clientID, msg)
    if err != nil {
        return err
    }

    if busy {
        err = fmt.Errorf("acceptor busy")
        return err
    }

    return nil
}

func (ac *Acceptor) GetClientIdByCompId(targetCompID string) (clientID uint64, err error) {
    return ac.cliSessnCtrl.getClientIdByCompId(targetCompID)
}
