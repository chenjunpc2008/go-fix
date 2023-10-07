package fixacceptor

import (
    "fmt"
    "sync"
    "time"
)

/*
defines for client session control
*/

// client heartbeat
type clientHbSt struct {
    // LastRcvTime 上次接收报文时间
    LastRcvTime int64
    // LastSendTime 上次发送报文时间
    LastSendTime int64
}

// client info
type clientInfoSt struct {
    beginstring  string
    senderCompID string
    targetCompID string
    heartBtInt   uint16
}

// client session
type clientSessionSt struct {
    cid          uint64
    ip           string
    addr         string
    senderSeqnum uint64
    targetSeqnum uint64
    bLogon       bool
    bLogout      bool
    cliInfo      clientInfoSt
    heartbeat    clientHbSt
}

// client sessions ctrl
type clientSessionCtrlSt struct {
    // lock for below values
    lock sync.Mutex
    // client sessions key--clientID
    mapCIDCtrl map[uint64]*clientSessionSt
    // client sessions key--targetCompID
    mapSIDCtrl map[string]*clientSessionSt
}

//
type hbtUnitSt struct {
    cid          uint64
    beginstring  string
    senderCompID string
    targetCompID string
}

// client close type
type cliCloseType int

const (
    closeTypeUnknown    = cliCloseType(0)
    closeTypeHbtTimeOut = cliCloseType(1)
    closeTypeNotLogon   = cliCloseType(2)
    closeTypeLogout     = cliCloseType(3)
)

// client close unit
type cliCloseUnitSt struct {
    cid          uint64
    closeType    cliCloseType
    beginstring  string
    senderCompID string
    targetCompID string
}

func newClientSessionCtrlSt() *clientSessionCtrlSt {
    var obj = &clientSessionCtrlSt{}
    obj.mapCIDCtrl = make(map[uint64]*clientSessionSt, 0)
    obj.mapSIDCtrl = make(map[string]*clientSessionSt, 0)

    return obj
}

// new connection
func (ctrl *clientSessionCtrlSt) addNewConnection(clientID uint64, clientIP, clientAddr string) bool {
    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    var cli = &clientSessionSt{cid: clientID, ip: clientIP, addr: clientAddr}

    now := time.Now().Unix()

    cli.cliInfo.heartBtInt = CiDefaultHbtTime
    cli.heartbeat.LastRcvTime = now
    cli.heartbeat.LastSendTime = now

    ctrl.mapCIDCtrl[clientID] = cli

    return true
}

// connection disconnected
func (ctrl *clientSessionCtrlSt) disconnect(clientID uint64, clientIP, clientAddr string) error {
    const ftag = "clientSessionCtrlSt.disconnect()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    cli, ok := ctrl.mapCIDCtrl[clientID]
    if !ok {
        return fmt.Errorf("%s cli-ctrl don't have ID:%d", ftag, clientID)
    }

    delete(ctrl.mapCIDCtrl, clientID)

    if nil == cli {
        return fmt.Errorf("%s cli-ctrl nil *clientSessionSt ID:%d", ftag, clientID)
    }

    if 0 == len(cli.cliInfo.targetCompID) {
        return nil
    }

    delete(ctrl.mapSIDCtrl, cli.cliInfo.targetCompID)

    return nil
}

// record receive, update heartbeat
func (ctrl *clientSessionCtrlSt) recordRcv(clientID uint64, clientIP, clientAddr string) error {
    const ftag = "clientSessionCtrlSt.recordRcv()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    cli, ok := ctrl.mapCIDCtrl[clientID]
    if !ok {
        return fmt.Errorf("%s cli-ctrl don't have ID:%d", ftag, clientID)
    }

    if nil == cli {
        return fmt.Errorf("%s cli-ctrl nil *clientSessionSt ID:%d", ftag, clientID)
    }

    now := time.Now().Unix()

    cli.heartbeat.LastRcvTime = now

    return nil
}

// record send, update heartbeat
func (ctrl *clientSessionCtrlSt) recordSend(clientID uint64, clientIP, clientAddr string) error {
    const ftag = "clientSessionCtrlSt.recordSend()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    cli, ok := ctrl.mapCIDCtrl[clientID]
    if !ok {
        return fmt.Errorf("%s cli-ctrl don't have ID:%d", ftag, clientID)
    }

    if nil == cli {
        return fmt.Errorf("%s cli-ctrl nil *clientSessionSt ID:%d", ftag, clientID)
    }

    now := time.Now().Unix()

    cli.heartbeat.LastSendTime = now

    return nil
}

// get fix session logon status
func (ctrl *clientSessionCtrlSt) getLogonStatus(clientID uint64, clientIP, clientAddr string) (bool, error) {
    const ftag = "clientSessionCtrlSt.getLogonStatus()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    cli, ok := ctrl.mapCIDCtrl[clientID]
    if !ok {
        return false, fmt.Errorf("%s cli-ctrl don't have ID:%d", ftag, clientID)
    }

    return cli.bLogon, nil
}

// TargetCompIdHasLoggedOn 登录互踢，方法中不加锁，同步工作交给调用函数处理
func (ctrl *clientSessionCtrlSt) clientAlreadyLogonUnlock(sSenderCompID string) bool {
    if _, ok := ctrl.mapSIDCtrl[sSenderCompID]; ok {
        return true
    }
    return false
}

// fix session logon
func (ctrl *clientSessionCtrlSt) logon(clientID uint64, clientIP, clientAddr, sBeginstring, sSenderCompID, sTargetCompID string,
    iHeartBtInt uint16) (sRtSenderCompID, sRtTargetCompID string, iRtHeartBtInt uint16, logonRes bool, err error) {
    const ftag = "clientSessionCtrlSt.logon()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    cli, ok := ctrl.mapCIDCtrl[clientID]
    if !ok {
        err = fmt.Errorf("%s cli-ctrl don't have ID:%d", ftag, clientID)
        return
    }

    if ctrl.clientAlreadyLogonUnlock(sSenderCompID) {
        // 将bLogon设置为false，剩余的工作交给监控进程，监控进程会检测到client未登录，然后执行发送logout消息等操作
        cli.bLogon = false
        cli.bLogout = true
        err = fmt.Errorf("SenderCompID:%s already logged in", sSenderCompID)
        return
    }

    cli.bLogon = true
    cli.bLogout = false

    cli.cliInfo.beginstring = sBeginstring
    cli.cliInfo.senderCompID = sTargetCompID
    cli.cliInfo.targetCompID = sSenderCompID

    // 当报盘端发送Logon 消息中的HeartBtInt 取值属于[5,60]时，服务端
    // 返回原值，否则取边界值（5 或60）。
    if CiHbtMin > iHeartBtInt || CiHbtMax < iHeartBtInt {
        cli.cliInfo.heartBtInt = CiHbtMax
    } else {
        cli.cliInfo.heartBtInt = iHeartBtInt
    }

    return cli.cliInfo.senderCompID, cli.cliInfo.targetCompID, cli.cliInfo.heartBtInt, true, nil
}

// fix session logout
func (ctrl *clientSessionCtrlSt) logout(clientID uint64) error {
    const ftag = "clientSessionCtrlSt.logout()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    cli, ok := ctrl.mapCIDCtrl[clientID]
    if !ok {
        return fmt.Errorf("%s cli-ctrl don't have ID:%d", ftag, clientID)
    }

    cli.bLogout = true

    return nil
}

// set sender MsgSeqNum
func (ctrl *clientSessionCtrlSt) setSenderSeqnum(clientID, seqNo uint64) error {
    const ftag = "clientSessionCtrlSt.setSenderSeqnum()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    cli, ok := ctrl.mapCIDCtrl[clientID]
    if !ok {
        return fmt.Errorf("%s cli-ctrl don't have ID:%d", ftag, clientID)
    }

    cli.senderSeqnum = seqNo

    return nil
}

// get sender MsgSeqNum
func (ctrl *clientSessionCtrlSt) getSenderSeqnum(clientID uint64) (seqNo uint64, err error) {
    const ftag = "clientSessionCtrlSt.getSenderSeqnum()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    cli, ok := ctrl.mapCIDCtrl[clientID]
    if !ok {
        return 0, fmt.Errorf("%s cli-ctrl don't have ID:%d", ftag, clientID)
    }

    cli.senderSeqnum++
    seqNo = cli.senderSeqnum

    return seqNo, nil
}

// check heartbeat and clients to be closed
func (ctrl *clientSessionCtrlSt) checkClientHeartbeat() (needsToSendHb []hbtUnitSt,
    needsToClose []cliCloseUnitSt) {
    // const ftag = "clientSessionCtrlSt.checkClientHeartbeat()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    lNow := time.Now().Unix()

    needsToSendHb = make([]hbtUnitSt, 0, 200)
    needsToClose = make([]cliCloseUnitSt, 0, 200)

    for _, v := range ctrl.mapCIDCtrl {
        if nil == v {
            continue
        }

        //fmt.Printf("%s clientID:%v, IP:%s, LinkID:%s, logon:%v, logout:%v, tgap:%d",
        //	ftag, v.cid, v.ip, v.linkID, v.bLogon, v.bLogout, lNow-v.heartbeat.LastRcvTime)

        if v.bLogout {
            // 注销后
            if lNow-v.heartbeat.LastRcvTime >= CiForceCloseTimeOut {
                // 未成功登录
                var (
                    singleCli cliCloseUnitSt
                )

                singleCli.cid = v.cid
                singleCli.closeType = closeTypeLogout

                singleCli.beginstring = v.cliInfo.beginstring
                singleCli.senderCompID = v.cliInfo.senderCompID
                singleCli.targetCompID = v.cliInfo.targetCompID

                needsToClose = append(needsToClose, singleCli)
                continue
            }
        }

        if !v.bLogon {
            if lNow-v.heartbeat.LastRcvTime >= CiForceCloseTimeOut {
                // 未成功登录
                var (
                    singleCli cliCloseUnitSt
                )

                singleCli.cid = v.cid
                singleCli.closeType = closeTypeNotLogon

                singleCli.beginstring = v.cliInfo.beginstring
                singleCli.senderCompID = v.cliInfo.senderCompID
                singleCli.targetCompID = v.cliInfo.targetCompID

                needsToClose = append(needsToClose, singleCli)
                continue
            }
        }

        if lNow-v.heartbeat.LastRcvTime >= int64(v.cliInfo.heartBtInt*2) {
            // 若接收方在2个心跳间隔内未收到任何消息，则可以认为会话出现异常并立即关闭连接

            var (
                singleCli cliCloseUnitSt
            )

            singleCli.cid = v.cid
            singleCli.closeType = closeTypeHbtTimeOut

            singleCli.beginstring = v.cliInfo.beginstring
            singleCli.senderCompID = v.cliInfo.senderCompID
            singleCli.targetCompID = v.cliInfo.targetCompID

            needsToClose = append(needsToClose, singleCli)
        } else {
            if lNow-v.heartbeat.LastSendTime >= int64(v.cliInfo.heartBtInt) {

                var (
                    singleCli hbtUnitSt
                )

                singleCli.cid = v.cid

                singleCli.beginstring = v.cliInfo.beginstring
                singleCli.senderCompID = v.cliInfo.senderCompID
                singleCli.targetCompID = v.cliInfo.targetCompID

                needsToSendHb = append(needsToSendHb, singleCli)
            }
        }
    }

    return
}

// GetClientIdByCompId get client id by targetCompID
// return negative integer if not find
func (ctrl *clientSessionCtrlSt) getClientIdByCompId(targetCompID string) (clientID uint64, err error) {
    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    sessionCtrlInfo, ok := ctrl.mapSIDCtrl[targetCompID]
    if !ok {
        err = fmt.Errorf("no corresponding clientID of targetCompID:%s", targetCompID)
        return
    }

    return sessionCtrlInfo.cid, nil
}
