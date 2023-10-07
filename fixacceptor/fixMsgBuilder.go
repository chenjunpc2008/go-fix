package fixacceptor

import (
    "github.com/chenjunpc2008/go-fix/common/ibfixdefine"

    "github.com/quickfixgo/quickfix"
)

// build logon message
func buildLogonMsg(sBeginString, sSenderCompID, sTargetCompID string, iHeartBtInt uint16) (*quickfix.Message, error) {
    // Logon MsgType = A
    var logonMsg = quickfix.NewMessage()

    logonMsg.Header.SetString(ibfixdefine.Tag8, sBeginString)
    logonMsg.Header.SetString(ibfixdefine.Tag35, ibfixdefine.MsgTypeLogon)
    logonMsg.Header.SetString(ibfixdefine.Tag49, sSenderCompID)
    logonMsg.Header.SetString(ibfixdefine.Tag56, sTargetCompID)

    logonMsg.Body.SetInt(ibfixdefine.Tag98, 0)
    logonMsg.Body.SetInt(ibfixdefine.Tag108, int(iHeartBtInt))

    return logonMsg, nil
}

// build Heartbeat message
func buildHbtMsg(sBeginString, sSenderCompID, sTargetCompID, sTestReqID string) (*quickfix.Message, error) {

    // Heartbeat MsgType = 0
    var hbtMsg = quickfix.NewMessage()

    hbtMsg.Header.SetString(ibfixdefine.Tag8, sBeginString)
    hbtMsg.Header.SetString(ibfixdefine.Tag35, ibfixdefine.MsgTypeHeartbeat)
    hbtMsg.Header.SetString(ibfixdefine.Tag49, sSenderCompID)
    hbtMsg.Header.SetString(ibfixdefine.Tag56, sTargetCompID)

    hbtMsg.Body.SetInt(ibfixdefine.Tag98, 0)
    hbtMsg.Body.SetInt(ibfixdefine.Tag108, int(30))

    if 0 != len(sTestReqID) {
        hbtMsg.Body.SetString(ibfixdefine.Tag112, sTestReqID)
    }

    return hbtMsg, nil
}

// build logout message
func buildLogoutMsg(sBeginString, sSenderCompID, sTargetCompID, logoutReason string) (*quickfix.Message, error) {
    // Logout MsgType = 5
    var hbtMsg = quickfix.NewMessage()

    hbtMsg.Header.SetString(ibfixdefine.Tag8, sBeginString)
    hbtMsg.Header.SetString(ibfixdefine.Tag35, ibfixdefine.MsgTypeLogout)
    hbtMsg.Header.SetString(ibfixdefine.Tag49, sSenderCompID)
    hbtMsg.Header.SetString(ibfixdefine.Tag56, sTargetCompID)

    hbtMsg.Header.SetString(ibfixdefine.Tag58, logoutReason)

    return hbtMsg, nil
}
