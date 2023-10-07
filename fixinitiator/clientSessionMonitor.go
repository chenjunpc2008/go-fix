package fixinitiator

import (
    "fmt"
    "time"

    "github.com/quickfixgo/quickfix"
)

/*
心跳及session监控
*/
func heartbeatWatcher(chExit <-chan int, cliSessnCtrl *clientSessionCtrlSt, initiator *Initiator) {
    const ftag = "heartbeatWatcher()"

    timeout := time.Duration(1) * time.Second

    var (
        connected    bool
        lDisconnTime int64

        heartbeatToSend bool
        testReqToSend   bool
        cliToClose      bool
        closeType       cliCloseType

        lNow int64

        // bHoldingPatten         bool
        // holdingPattenBeginTime int64

        err error
    )

    connected = true
    lDisconnTime = time.Now().Unix()

    for {
        select {
        case <-chExit:
            fmt.Printf("%s Exit\n", ftag)
            return

        case <-time.After(timeout):
        }

        connected, lDisconnTime, heartbeatToSend, testReqToSend, cliToClose, closeType = cliSessnCtrl.checkClientHeartbeat()

        if connected {
            // fmt.Printf("%s  connected:%v, lDisconnTime:%v, heartbeatToSend:%v, testReqToSend:%v, cliToClose:%v, closeType:%v\n", ftag,
            //     connected, lDisconnTime, heartbeatToSend, testReqToSend, cliToClose, closeType)

            if heartbeatToSend {
                sendHeartBeat(initiator)
            }

            if testReqToSend {
                sendTestRequest(initiator)
            }

            if cliToClose {
                closeClient(closeType, initiator)
            }

            lDisconnTime = time.Now().Unix()

            // check if in holding status
            // bHoldingPatten, holdingPattenBeginTime = initiator.getHoldingStatus()
            // if bHoldingPatten {
            //     if CiHoldingPattenMaxTime > lNow-holdingPattenBeginTime {
            //         // force break from holding patten
            //         initiator.leaveHoldingPatten()

            //         // report back
            //         sErrMsg := fmt.Sprintf("force break holding patten, holdingPattenBeginTime:%d, now:%d", holdingPattenBeginTime, lNow)
            //         initiator.appHdl.OnErrorStr(sErrMsg)
            //     }
            // }

        } else {
            // check if need reconnect
            lNow = time.Now().Unix()
            if CiReconnectTimeGap <= lNow-lDisconnTime {
                initiator.appHdl.OnEvent("trying to reconnect to server")

                // reconnect
                initiator.disconnectFromSvr()
                err = initiator.reconnectToSvr()
                if nil != err {
                    initiator.appHdl.OnError("reconnect to server", err)
                }
            }
        }
    }
}

// send heartbeat to client
func sendHeartBeat(initiator *Initiator) {
    const ftag = "sendHeartBeat()"

    // incase of closed
    select {
    case <-initiator.chExit:
        return

    default:
    }

    var (
        hbtMsg *quickfix.Message
        mpkg   = NewMsgPkg()
        err    error
    )

    // build heartbeat msg
    hbtMsg, _ = buildHbtMsg(initiator.setting.BeginString,
        initiator.setting.SenderCompID, initiator.setting.TargetCompID, "")

    mpkg.Fixmsg = hbtMsg

    // send to server
    _, err = initiator.SendMsgToSvr(mpkg)
    if nil != err {
        initiator.appHdl.OnError("send heartbeat", err)
    }
}

// send Test Request to client
func sendTestRequest(initiator *Initiator) {
    // const ftag = "sendTestRequest()"

    // in case of closed
    select {
    case <-initiator.chExit:
        return

    default:
    }

    var (
        testreqMsg *quickfix.Message
        mpkg       = NewMsgPkg()
        err        error
    )

    // build Request msg
    testreqMsg, _ = buildTestReqMsg(initiator.setting.BeginString,
        initiator.setting.SenderCompID, initiator.setting.TargetCompID)

    mpkg.Fixmsg = testreqMsg

    // send to server
    _, err = initiator.SendMsgToSvr(mpkg)
    if nil != err {
        initiator.appHdl.OnError("send Test Request", err)
    }
}

// close client
func closeClient(closeType cliCloseType, initiator *Initiator) {
    // const ftag = "closeClient"

    // in case of closed
    select {
    case <-initiator.chExit:
        return

    default:
    }

    var (
        err error
    )

    switch closeType {
    case closeTypeHbtTimeOut:
        // 心跳超时

        // logout
        logout1, _ := buildLogoutMsg(initiator.setting.BeginString,
            initiator.setting.SenderCompID, initiator.setting.TargetCompID,
            "hbt time out")

        mpkg1 := NewMsgPkg()
        mpkg1.Fixmsg = logout1

        _, err = initiator.SendMsgToSvr(mpkg1)
        if nil != err {
            initiator.appHdl.OnError("send logout", err)
        }

        // logout in system
        if nil != initiator.cliSessnCtrl {
            initiator.cliSessnCtrl.logout()
        }

        // 之后会由session扫描自动登出

    case closeTypeNotLogon:
        // 建立TCP 连接之后，超过5 秒未完成登录；
        // 在登录失败之后，未在5 秒内关闭连接；

        // logout
        logout2, _ := buildLogoutMsg(initiator.setting.BeginString,
            initiator.setting.SenderCompID, initiator.setting.TargetCompID,
            "not log on")

        mpkg2 := NewMsgPkg()
        mpkg2.Fixmsg = logout2

        _, err = initiator.SendMsgToSvr(mpkg2)
        if nil != err {
            initiator.appHdl.OnError("send logout", err)
        }

        // logout in system
        if nil != initiator.cliSessnCtrl {
            initiator.cliSessnCtrl.logout()
        }

        // 之后会由session扫描自动登出

    case closeTypeLogout:
        // 已登出
        initiator.disconnectFromSvr()

    default:
        initiator.disconnectFromSvr()
    }

    return
}
