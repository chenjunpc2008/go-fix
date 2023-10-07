package fixacceptor

import (
    "time"

    "github.com/quickfixgo/quickfix"
)

/*
心跳及session监控
*/
func heartbeatWatcher(chExit <-chan int, cliSessnCtrl *clientSessionCtrlSt, acceptor *Acceptor) {
    const ftag = "heartbeatWatcher()"

    timeout := time.Duration(1) * time.Second

    var (
        heartbeatToSend []hbtUnitSt
        cliToClose      []cliCloseUnitSt
    )

    for {
        select {
        case <-chExit:
            return

        case <-time.After(timeout):
        }

        heartbeatToSend, cliToClose = cliSessnCtrl.checkClientHeartbeat()

        // fmt.Printf("%s hbToSend=%v, cliToClose=%v",ftag, heartbeatToSend, cliToClose)

        go sendHeartBeats(heartbeatToSend, acceptor)
        go closeClients(cliToClose, acceptor)
    }
}

// send heartbeats to clients
func sendHeartBeats(clients []hbtUnitSt, acceptor *Acceptor) {
    const ftag = "sendHeartBeats()"

    if 0 == len(clients) {
        return
    }

    // incase of closed
    select {
    case <-acceptor.chExit:
        return

    default:
    }

    var (
        cli    hbtUnitSt
        hbtMsg *quickfix.Message
        err    error
    )
    for _, cli = range clients {
        // build heartbeat msg
        hbtMsg, _ = buildHbtMsg(cli.beginstring, cli.senderCompID, cli.targetCompID, "")

        //send to client
        err = acceptor.SendToClient(cli.cid, hbtMsg)
        if nil != err {
            acceptor.appHdl.OnError("send heartbeat", err)
        }
    }
}

// close clients
func closeClients(clients []cliCloseUnitSt, acceptor *Acceptor) {
    const ftag = "closeClients"

    if 0 == len(clients) {
        return
    }

    // in case of closed
    select {
    case <-acceptor.chExit:
        return

    default:
    }

    var (
        err error
    )

    for _, cli := range clients {

        switch cli.closeType {
        case closeTypeHbtTimeOut:
            // 心跳超时

            // logout
            logout1, _ := buildLogoutMsg(cli.beginstring, cli.senderCompID, cli.targetCompID, "hbt time out")

            err = acceptor.SendToClient(cli.cid, logout1)
            if nil != err {
                acceptor.appHdl.OnError("send logout", err)
            }

            // logout in system
            if nil != acceptor.cliSessnCtrl {
                acceptor.cliSessnCtrl.logout(cli.cid)
            }

            // 之后会由session扫描自动登出

        case closeTypeNotLogon:
            // 建立TCP 连接之后，超过5 秒未完成登录；
            // 在登录失败之后，未在5 秒内关闭连接；

            // logout
            logout2, _ := buildLogoutMsg(cli.beginstring, cli.senderCompID, cli.targetCompID, "not log on")

            err = acceptor.SendToClient(cli.cid, logout2)
            if nil != err {
                acceptor.appHdl.OnError("send logout", err)
            }

            // logout in system
            if nil != acceptor.cliSessnCtrl {
                acceptor.cliSessnCtrl.logout(cli.cid)
            }

            // 之后会由session扫描自动登出

        case closeTypeLogout:
            // 已登出
            acceptor.closeCli(cli.cid, "logout")

        default:
            acceptor.closeCli(cli.cid, "unknown")
        }
    }

    return
}
