package fixinitiator

import (
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

// client session ctrl
type clientSessionCtrlSt struct {
    lock sync.Mutex // lock for values below

    svrIP        string
    svrPort      uint16
    bIsConnected bool
    lDisconnTime int64

    nextSenderMsgSeqNum uint64
    nextTargetMsgSeqNum uint64
    bLogon              bool
    bLogout             bool
    cliInfo             clientInfoSt
    heartbeat           clientHbSt
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

func newClientSessionCtrlSt(nxtSenderMsgSeqNum, nxtTargetMsgSeqNum uint64) *clientSessionCtrlSt {
    var obj = &clientSessionCtrlSt{nextSenderMsgSeqNum: nxtSenderMsgSeqNum, nextTargetMsgSeqNum: nxtTargetMsgSeqNum}

    return obj
}

// new connection
func (ctrl *clientSessionCtrlSt) addNewConnection(serverIP string, serverPort uint16) bool {
    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    ctrl.svrIP = serverIP
    ctrl.svrPort = serverPort
    ctrl.bIsConnected = true

    ctrl.bLogon = false
    ctrl.bLogout = false

    now := time.Now().Unix()

    ctrl.cliInfo.heartBtInt = CiDefaultHbtTime
    ctrl.heartbeat.LastRcvTime = now
    ctrl.heartbeat.LastSendTime = now

    ctrl.lDisconnTime = now

    return true
}

// connection disconnected
func (ctrl *clientSessionCtrlSt) disconnect(serverIP string, serverPort uint16) {
    const ftag = "clientSessionCtrlSt.disconnect()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    ctrl.svrIP = ""
    ctrl.svrPort = 0
    ctrl.bIsConnected = false
    ctrl.lDisconnTime = time.Now().Unix()

    ctrl.bLogon = false
    ctrl.bLogout = false
}

// record receive, update heartbeat
func (ctrl *clientSessionCtrlSt) recordRcv(serverIP string, serverPort uint16) {
    // const ftag = "clientSessionCtrlSt.recordRcv()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    now := time.Now().Unix()

    ctrl.heartbeat.LastRcvTime = now
}

// record send, update heartbeat
func (ctrl *clientSessionCtrlSt) recordSend(serverIP string, serverPort uint16) {
    // const ftag = "clientSessionCtrlSt.recordSend()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    now := time.Now().Unix()

    ctrl.heartbeat.LastSendTime = now
}

// get fix session logon status
func (ctrl *clientSessionCtrlSt) getLogonStatus() bool {
    // const ftag = "clientSessionCtrlSt.getLogonStatus()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    return ctrl.bLogon
}

// fix session logon
func (ctrl *clientSessionCtrlSt) logon(serverIP string, serverPort uint16,
    sBeginstring string, sSenderCompID string, sTargetCompID string, iHeartBtInt uint16) {
    const ftag = "clientSessionCtrlSt.logon()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    ctrl.bLogon = true
    ctrl.bLogout = false

    ctrl.cliInfo.beginstring = sBeginstring
    ctrl.cliInfo.senderCompID = sTargetCompID
    ctrl.cliInfo.targetCompID = sSenderCompID

    ctrl.cliInfo.heartBtInt = iHeartBtInt

    // reset time
    ctrl.lDisconnTime = time.Now().Unix()
}

// fix session logout
func (ctrl *clientSessionCtrlSt) logout() {
    const ftag = "clientSessionCtrlSt.logout()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    ctrl.bLogout = true
}

// get fix session logout status
func (ctrl *clientSessionCtrlSt) getLogoutStatus() bool {
    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    return ctrl.bLogout
}

// set next SenderMsgSeqNum
func (ctrl *clientSessionCtrlSt) setNextSenderMsgSeqNum(newSeqNo uint64) {

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    if 0 == newSeqNo {
        ctrl.nextSenderMsgSeqNum = 1
        return
    }

    ctrl.nextSenderMsgSeqNum = newSeqNo
}

// get next SenderMsgSeqNum, willn't increase the actual number
func (ctrl *clientSessionCtrlSt) getNextSenderMsgSeqNum() uint64 {

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    return ctrl.nextSenderMsgSeqNum
}

// get next SenderMsgSeqNum, then increase
func (ctrl *clientSessionCtrlSt) getThenIncrNextSenderMsgSeqNum() (nxtSeqNo uint64) {
    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    nxtSeqNo = ctrl.nextSenderMsgSeqNum

    ctrl.nextSenderMsgSeqNum++

    return nxtSeqNo
}

// set next TargetMsgSeqNum
func (ctrl *clientSessionCtrlSt) setNextTargetMsgSeqNum(newSeqNo uint64) {

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    if 0 == newSeqNo {
        ctrl.nextTargetMsgSeqNum = 1
        return
    }

    ctrl.nextTargetMsgSeqNum = newSeqNo
}

// get next TargetMsgSeqNum, willn't increase the actual number
func (ctrl *clientSessionCtrlSt) getNextTargetMsgSeqNum() (nxtSeqNo uint64) {

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    return ctrl.nextTargetMsgSeqNum
}

// get next TargetMsgSeqNum, then increace
func (ctrl *clientSessionCtrlSt) getThenIncrNextTargetMsgSeqNum() (nxtSeqNo uint64) {
    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    nxtSeqNo = ctrl.nextTargetMsgSeqNum

    ctrl.nextTargetMsgSeqNum++

    return nxtSeqNo
}

// check heartbeat and clients to be closed
func (ctrl *clientSessionCtrlSt) checkClientHeartbeat() (connected bool,
    disconnTime int64,
    needsToSendHb bool,
    needsToSendTestRequest bool,
    needsToClose bool, closeType cliCloseType) {
    // const ftag = "clientSessionCtrlSt.checkClientHeartbeat()"

    // lock
    ctrl.lock.Lock()
    defer ctrl.lock.Unlock()

    lNow := time.Now().Unix()

    //fmt.Printf("%s clientID:%v, IP:%s, LinkID:%s, logon:%v, logout:%v, tgap:%d",
    //	ftag, v.cid, v.ip, v.linkID, v.bLogon, v.bLogout, lNow-v.heartbeat.LastRcvTime)

    if !ctrl.bIsConnected {
        return false, ctrl.lDisconnTime, false, false, false, closeTypeUnknown
    }

    connected = true

    if ctrl.bLogout {
        // 注销后
        if lNow-ctrl.heartbeat.LastRcvTime >= CiForceCloseTimeOut {
            // 未成功登录
            needsToClose = true
            closeType = closeTypeLogout

            return
        }
    }

    if !ctrl.bLogon {
        if lNow-ctrl.heartbeat.LastRcvTime >= CiForceCloseTimeOut {
            // 未成功登录
            needsToClose = true
            closeType = closeTypeNotLogon

            return
        }
    }

    if lNow-ctrl.heartbeat.LastRcvTime >= int64(ctrl.cliInfo.heartBtInt) {
        if lNow-ctrl.heartbeat.LastRcvTime >= int64(ctrl.cliInfo.heartBtInt*2) {
            // 若接收方在2个心跳间隔内未收到任何消息，则可以认为会话出现异常并立即关闭连接

            needsToClose = true
            closeType = closeTypeHbtTimeOut
        } else if lNow-ctrl.heartbeat.LastRcvTime >= int64(ctrl.cliInfo.heartBtInt)+1 &&
            lNow-ctrl.heartbeat.LastRcvTime <= int64(ctrl.cliInfo.heartBtInt)+3 {
            // 测试请求消息能强制对方发出心跳消息
            // 对测试请求的发送不希望额外的记录标识，故采用在某一时间窗口内发送
            needsToSendTestRequest = true
        }
    } else {
        if lNow-ctrl.heartbeat.LastSendTime >= int64(ctrl.cliInfo.heartBtInt) {

            needsToSendHb = true
        }
    }

    return
}
