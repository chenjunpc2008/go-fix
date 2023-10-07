package fixinitiator

import (
	"github.com/chenjunpc2008/go-fix/common/fixtimeutil"
	"github.com/chenjunpc2008/go-fix/common/ibfixdefine"

	"github.com/quickfixgo/quickfix"
)

// build logon message
func buildLogonMsg(sBeginString, sSenderCompID, sTargetCompID string,
    iHeartBtInt uint16, bResetSeqNum bool) (*quickfix.Message, error) {

    // Logon MsgType = A
    var logonMsg = quickfix.NewMessage()

    logonMsg.Header.SetString(ibfixdefine.Tag8, sBeginString)
    logonMsg.Header.SetString(ibfixdefine.Tag35, ibfixdefine.MsgTypeLogon)
    logonMsg.Header.SetString(ibfixdefine.Tag49, sSenderCompID)
    logonMsg.Header.SetString(ibfixdefine.Tag56, sTargetCompID)

    logonMsg.Body.SetInt(ibfixdefine.Tag98, 0)
    logonMsg.Body.SetInt(ibfixdefine.Tag108, int(iHeartBtInt))

    if bResetSeqNum {
        logonMsg.Body.SetBool(ibfixdefine.Tag141, true)
    }

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

    if 0 != len(sTestReqID) {
        hbtMsg.Body.SetString(ibfixdefine.Tag112, sTestReqID)
    }

    return hbtMsg, nil
}

// build Test Request messasge
func buildTestReqMsg(sBeginString, sSenderCompID, sTargetCompID string) (*quickfix.Message, error) {
    // Test Request MsgType = 1
    var testMsg = quickfix.NewMessage()

    testMsg.Header.SetString(ibfixdefine.Tag8, sBeginString)
    testMsg.Header.SetString(ibfixdefine.Tag35, ibfixdefine.MsgTypeTestRequest)
    testMsg.Header.SetString(ibfixdefine.Tag49, sSenderCompID)
    testMsg.Header.SetString(ibfixdefine.Tag56, sTargetCompID)

    sTimestamp := fixtimeutil.GetFIXUTCTimestamp()
    testMsg.Body.SetString(ibfixdefine.Tag112, sTimestamp)

    return testMsg, nil
}

// build Sequence Reset - Reset messasge
func buildSequenceResetResetMsg(sBeginString, sSenderCompID, sTargetCompID string, iNewSeqNo int) (*quickfix.Message, error) {
    // Sequence Reset and Gap Fill MsgType = 4
    var resetMsg = quickfix.NewMessage()

    resetMsg.Header.SetString(ibfixdefine.Tag8, sBeginString)
    resetMsg.Header.SetString(ibfixdefine.Tag35, ibfixdefine.MsgTypeSequenceReset)
    resetMsg.Header.SetString(ibfixdefine.Tag49, sSenderCompID)
    resetMsg.Header.SetString(ibfixdefine.Tag56, sTargetCompID)
    resetMsg.Header.SetInt(ibfixdefine.Tag34, 1) // when reset, MsgSeqNum=1

    sendTime := fixtimeutil.GetFIXUTCTimestamp()
    resetMsg.Header.SetString(ibfixdefine.Tag52, sendTime)

    resetMsg.Body.SetString(ibfixdefine.Tag123, "N")
    resetMsg.Body.SetInt(ibfixdefine.Tag36, iNewSeqNo)

    return resetMsg, nil
}

// build logout message
func buildLogoutMsg(sBeginString, sSenderCompID, sTargetCompID, logoutReason string) (*quickfix.Message, error) {
    // Logout MsgType = 5
    var loutMsg = quickfix.NewMessage()

    loutMsg.Header.SetString(ibfixdefine.Tag8, sBeginString)
    loutMsg.Header.SetString(ibfixdefine.Tag35, ibfixdefine.MsgTypeLogout)
    loutMsg.Header.SetString(ibfixdefine.Tag49, sSenderCompID)
    loutMsg.Header.SetString(ibfixdefine.Tag56, sTargetCompID)

    loutMsg.Header.SetString(ibfixdefine.Tag58, logoutReason)

    return loutMsg, nil
}
