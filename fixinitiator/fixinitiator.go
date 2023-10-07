package fixinitiator

import (
    "errors"
    "fmt"
    "sync"
    "time"

    "github.com/quickfixgo/quickfix/datadictionary"

    "github.com/chenjunpc2008/go-tcp/tcpclient"
    "github.com/chenjunpc2008/go/util/kchanthreadpool"
)

/*
Message Recovery

Message gaps may occur during login or in the middle of a FIX session that require message
recovery. When this occurs, a Resend Request message can be sent requesting a range of
missing messages. The re-sender will then respond with a Sequence Reset message that has
tag 123 (GapFillFlag) set to “Y” and tag 43 (PossDupFlag) set to “Y”. Tag 36 (NewSeqNo) will be
set to the sequence number of the message to be redelivered next. Following this, the missing
Application level messages will be resent.

当收到Resend Request消息，initiator会优先发送回放的消息，做如下处理：
1. initiator进入保护状态，将待发送的新消息缓存，后续的新消息在回放结束前均进入缓存队列
2. initiator生成Sequence Reset - Reset数据，Reset数据将直接进入正常发送队列发送
4. initiator回放结束，结束保护状态，恢复正常消息发送，将待发送的新消息缓存放入正常的发送队列
*/

const (
    ciTaskPoolBuffSize = 2 * 10000
)

// Initiator initiator
type Initiator struct {
    setting Settings

    lock sync.Mutex // locker for below values

    fixDataDict *datadictionary.DataDictionary // fix dict DataDictionary

    client       *tcpclient.CtcpCli
    appHdl       ApplicationIF
    cliSessnCtrl *clientSessionCtrlSt
    chExit       chan int // exit channel

    // thread pool for process admin msg
    adminMsgTPool *kchanthreadpool.ThreadPool
    // thread pool for process app msg
    appMsgTPool *kchanthreadpool.ThreadPool

    // Holding patten, may be in a Resend Request message have been received, entering message recovery process
    // bHoldingPatten          bool
    // holdingPattenBeginTime  int64
    // messageCacheWhenHolding *list.List
}

// NewInitiator new
func NewInitiator(set Settings, app ApplicationIF) *Initiator {
    var intor = &Initiator{setting: set, appHdl: app}
    return intor
}

// GetNeededProccessNum get gorountine numbers needed
func GetNeededProccessNum() int {
    tcpCliTNum := 2
    tpoolNum := CiAdminThPoolNum + CiAppThPoolNum
    /*
       --heartbeatWatcher()
    */
    localDeamon := 2

    var total = tcpCliTNum + tpoolNum + localDeamon

    return total
}

// Start initial initiator client
func (intor *Initiator) Start(nxtSenderMsgSeqNum, nxtTargetMsgSeqNum uint64) error {
    // lock
    intor.lock.Lock()
    defer intor.lock.Unlock()

    var err error

    intor.fixDataDict, err = datadictionary.Parse(intor.setting.FixDataDictPath)
    if nil != err {
        return err
    }

    var cnf = tcpclient.DefaultConfig()
    cnf.SendBuffsize = intor.setting.SendBuffsize
    cnf.AsyncReceive = intor.setting.AsyncReceive

    hdler := &initiatorHandler{appHdl: intor.appHdl, initiator: intor}

    intor.client = tcpclient.New(hdler, cnf)
    intor.cliSessnCtrl = newClientSessionCtrlSt(nxtSenderMsgSeqNum, nxtTargetMsgSeqNum)
    intor.chExit = make(chan int)
    // intor.bHoldingPatten = false

    // thread pool
    tpHdler := newThreadPoolHandlerSt(intor)
    intor.adminMsgTPool, err = kchanthreadpool.NewThreadPool(CiAdminThPoolNum, ciTaskPoolBuffSize, tpHdler)
    if nil != err {
        return err
    }

    err = intor.adminMsgTPool.Start()
    if nil != err {
        return err
    }

    // thread pool
    intor.appMsgTPool, err = kchanthreadpool.NewThreadPool(CiAppThPoolNum, ciTaskPoolBuffSize, tpHdler)
    if nil != err {
        intor.adminMsgTPool.Stop()

        return err
    }

    err = intor.appMsgTPool.Start()
    if nil != err {
        intor.adminMsgTPool.Stop()

        return err
    }

    // if error in this time, we will leave it for auto-reconnect
    _ = intor.client.ConnectToServer(intor.setting.RemoteIP, intor.setting.RemotePort)

    // wait for client to establish
    time.Sleep(time.Duration(1) * time.Second)

    go heartbeatWatcher(intor.chExit, intor.cliSessnCtrl, intor)

    return nil
}

// Stop initial initiator client
func (intor *Initiator) Stop() {
    // lock
    intor.lock.Lock()
    defer intor.lock.Unlock()

    if nil != intor.client {
        close(intor.chExit)

        intor.adminMsgTPool.Stop()
        intor.appMsgTPool.Stop()

        intor.client.Close()

        intor.client = nil
        intor.cliSessnCtrl = nil
    }
}

// disconnect from server
func (intor *Initiator) disconnectFromSvr() {
    // lock
    intor.lock.Lock()
    defer intor.lock.Unlock()

    if nil != intor.client {
        intor.client.Close()
    }
}

// connect to server
func (intor *Initiator) reconnectToSvr() error {
    // lock
    intor.lock.Lock()
    defer intor.lock.Unlock()

    select {
    case <-intor.chExit:
        return fmt.Errorf("initiator closed")

    default:
    }

    if nil == intor.client {
        return fmt.Errorf("initiator un-initial")
    }

    var cnf = tcpclient.DefaultConfig()
    cnf.SendBuffsize = intor.setting.SendBuffsize
    cnf.AsyncReceive = intor.setting.AsyncReceive

    hdler := &initiatorHandler{appHdl: intor.appHdl, initiator: intor}

    intor.client = tcpclient.New(hdler, cnf)

    // reset time, or the reconnector will try too soon
    intor.cliSessnCtrl.lDisconnTime = time.Now().Unix()

    err := intor.client.ConnectToServer(intor.setting.RemoteIP, intor.setting.RemotePort)
    if nil != err {
        return err
    }

    return nil
}

// func (intor *Initiator) getHoldingStatus() (bHoldingPatten bool, holdingPattenBeginTime int64) {
//     // lock
//     intor.lock.Lock()
//     defer intor.lock.Unlock()

//     return intor.bHoldingPatten, intor.holdingPattenBeginTime
// }

func (intor *Initiator) getCliSessnCtrl() *clientSessionCtrlSt {
    // lock
    intor.lock.Lock()
    defer intor.lock.Unlock()

    select {
    case <-intor.chExit:
        return nil

    default:
    }

    if nil == intor.cliSessnCtrl {
        return nil
    }

    return intor.cliSessnCtrl
}

/*
SendMsgToSvr send to server

@return ready bool : initiator is ready to send message, but the condition is ok, if you try later it will succeed
@return err error : initiator have error
*/
func (intor *Initiator) SendMsgToSvr(msg *MsgPkg) (ready bool, err error) {
    // lock
    intor.lock.Lock()
    defer intor.lock.Unlock()

    select {
    case <-intor.chExit:
        return false, fmt.Errorf("initiator closed")

    default:
    }

    if nil == intor.client {
        return false, fmt.Errorf("initiator un-initial")
    }

    // if intor.bHoldingPatten {
    //     // in resending status, hold message send first
    //     sMsg := fmt.Sprintf("in holding patten, delay send msg:%s", msg.Fixmsg.String())
    //     intor.appHdl.OnEvent(sMsg)
    //     return false, nil
    // }

    var (
        localmsg = &msgToSendSt{}
        busy     bool
    )

    localmsg.pkg = msg
    localmsg.msgType = mtsTypeNormal

    busy, err = intor.client.SendToServer(localmsg)
    if nil != err {
        return false, err
    }

    if busy {
        err = errors.New("busy")
        return false, err
    }

    return true, nil
}

// SendPriorMsgToSvr send to server
func (intor *Initiator) SendPriorMsgToSvr(msg *MsgPkg) error {
    // lock
    intor.lock.Lock()
    defer intor.lock.Unlock()

    select {
    case <-intor.chExit:
        return fmt.Errorf("initiator closed")

    default:
    }

    if nil == intor.client {
        return fmt.Errorf("initiator un-initial")
    }

    var (
        err      error
        localmsg = &msgToSendSt{}
        busy     bool
    )

    localmsg.pkg = msg
    localmsg.msgType = mtsTypeResend

    // err = intor.client.SendPriorToServer(localmsg)
    busy, err = intor.client.SendToServer(localmsg)
    if nil != err {
        return err
    }

    if busy {
        err = errors.New("busy")
        return err
    }

    return nil
}

// SendPriorMsgesToSvr send to server
func (intor *Initiator) SendPriorMsgesToSvr(msgs []*MsgPkg) error {
    // lock
    intor.lock.Lock()
    defer intor.lock.Unlock()

    select {
    case <-intor.chExit:
        return fmt.Errorf("initiator closed")

    default:
    }

    if nil == intor.client {
        return fmt.Errorf("initiator un-initial")
    }

    var (
        err       error
        localmsg  *msgToSendSt
        localmsgs = make([]interface{}, 0)
        busy      bool
    )

    for _, v := range msgs {
        localmsg = new(msgToSendSt)
        localmsg.pkg = v
        localmsg.msgType = mtsTypeResend

        localmsgs = append(localmsgs, localmsg)
    }

    // err = intor.client.SendPrioresToServer(localmsgs)
    busy, err = intor.client.SendToServer(localmsg)
    if nil != err {
        return err
    }

    if busy {
        err = errors.New("busy")
        return err
    }

    return nil
}

// // enter holding patten
// func (intor *Initiator) enterHoldingPatten() error {
//     // lock
//     intor.lock.Lock()
//     defer intor.lock.Unlock()

//     select {
//     case <-intor.chExit:
//         return fmt.Errorf("initiator closed")

//     default:
//     }

//     if nil == intor.client {
//         return fmt.Errorf("initiator un-initial")
//     }

//     if intor.bHoldingPatten {
//         return nil
//     }

//     intor.bHoldingPatten = true
//     intor.holdingPattenBeginTime = time.Now().Unix()

//     buff := intor.client.DumpSendBuffer()

//     if nil == buff {
//         intor.messageCacheWhenHolding = list.New()
//     } else {
//         intor.messageCacheWhenHolding = buff
//     }

//     return nil
// }

// // leave holding patten
// func (intor *Initiator) leaveHoldingPatten() error {
//     // lock
//     intor.lock.Lock()
//     defer intor.lock.Unlock()

//     select {
//     case <-intor.chExit:
//         return fmt.Errorf("initiator closed")

//     default:
//     }

//     if nil == intor.client {
//         return fmt.Errorf("initiator un-initial")
//     }

//     if !intor.bHoldingPatten {
//         return nil
//     }

//     intor.bHoldingPatten = false

//     if nil != intor.messageCacheWhenHolding {
//         err := intor.client.SendBuffToServer(intor.messageCacheWhenHolding)
//         if nil != err {
//             intor.appHdl.OnError("leaveHoldingPatten", err)
//         }

//         intor.messageCacheWhenHolding = nil
//     }

//     return nil
// }
