package fixinitiator

import (
    "fmt"
    "math"

    "github.com/quickfixgo/quickfix"

    "github.com/chenjunpc2008/go-fix/common/ibfixdefine"
    "github.com/chenjunpc2008/go/util/kchanthreadpool"
)

type adminMsgTask struct {
    serverIP   string
    serverPort uint16
    fixMsg     *quickfix.Message
    fInfo      ibfixdefine.FixInfoLocal
    initiator  *Initiator
}

func (task *adminMsgTask) Do(threadID uint) {
    //
    if nil == task.fixMsg {
        task.initiator.appHdl.OnErrorStr("nil *adminMsgTask")
        return
    }

    select {
    case <-task.initiator.chExit:
        break

    default:
    }

    var (
        sEventMsg string
    )

    switch task.fInfo.MsgType {
    case ibfixdefine.MsgTypeLogon:
        // 108 HeartBtIn
        iHeartBtInt, fixerr := task.fixMsg.Body.GetInt(ibfixdefine.Tag108)
        if nil != fixerr || iHeartBtInt > math.MaxUint16 {
            task.initiator.appHdl.OnError("get HeartBtIn", fixerr)

            // logout
            logoutmsg, _ := buildLogoutMsg(task.fInfo.BeginString,
                task.fInfo.TargetCompID,
                task.fInfo.SenderCompID,
                ibfixdefine.InvalidHbtInt)

            logoutMpkg := NewMsgPkg()
            logoutMpkg.Fixmsg = logoutmsg

            _, err := task.initiator.SendMsgToSvr(logoutMpkg)
            if nil != err {
                task.initiator.appHdl.OnError("send logout", err)
            }

            // set session control
            task.initiator.cliSessnCtrl.logout()

        } else {
            uiHbt := uint16(iHeartBtInt)

            task.initiator.cliSessnCtrl.logon(task.serverIP, task.serverPort,
                task.fInfo.BeginString, task.fInfo.TargetCompID, task.fInfo.SenderCompID,
                uiHbt)
        }

        // 141 ResetSeqNumFlag
        bResetSeqNumFlag, fixerr := task.fixMsg.Body.GetBool(ibfixdefine.Tag141)
        if nil != fixerr {
            if bResetSeqNumFlag {
                task.initiator.cliSessnCtrl.setNextSenderMsgSeqNum(1)
                task.initiator.cliSessnCtrl.setNextTargetMsgSeqNum(1)

                task.initiator.appHdl.OnEvent("ResetSeqNumFlag true")
            }
        }

    case ibfixdefine.MsgTypeHeartbeat:
        // check if this heartbeat is forced by a test request.
        // IBKR suggests the use of a timestamp in the TestReqID field as it is useful
        // to verify that the Heartbeat is the result of the Test Request and not as the result of a regular
        // timeout.
        sTestReqID, fixerr := task.fixMsg.Body.GetString(ibfixdefine.Tag112)
        if nil == fixerr {
            sEventMsg = "receive heartbeat forced by test request, TestReqID:" + sTestReqID
            task.initiator.appHdl.OnEvent(sEventMsg)
        }

    case ibfixdefine.MsgTypeTestRequest:
        // The test request message forces a heartbeat from the opposing application.
        sTestReqID, fixerr := task.fixMsg.Body.GetString(ibfixdefine.Tag112)
        if nil != fixerr {
            task.initiator.appHdl.OnError("get TestReqID failed", fixerr)
        }

        hbtReqMsg, _ := buildHbtMsg(task.initiator.setting.BeginString,
            task.initiator.setting.SenderCompID, task.initiator.setting.TargetCompID, sTestReqID)

        hbtMpkg := NewMsgPkg()
        hbtMpkg.Fixmsg = hbtReqMsg

        _, err := task.initiator.SendMsgToSvr(hbtMpkg)
        if nil != err {
            task.initiator.appHdl.OnError("send heartbeat forced by test request", err)
        }

    case ibfixdefine.MsgTypeResendRequest:
        // The Resend Request initiates the retransmission of messages.
        procResendRequest(task.serverIP, task.serverPort, task.fixMsg, task.fInfo, task.initiator)

    case ibfixdefine.MsgTypeSequenceReset:
        procSequenceResetReq(task.serverIP, task.serverPort, task.fixMsg, task.fInfo, task.initiator)

    case ibfixdefine.MsgTypeSessionLevelReject:
        procSessionLevelRejectReq(task.serverIP, task.serverPort, task.fixMsg, task.fInfo, task.initiator)

    case ibfixdefine.MsgTypeLogout:
        // 58 Text
        sText, fixerr := task.fixMsg.Body.GetString(ibfixdefine.Tag58)
        if nil != fixerr {
            sEventMsg = fmt.Sprintf("%s logout %s", task.fInfo.SenderCompID, sText)
        } else {
            sEventMsg = fmt.Sprintf("%s logout", task.fInfo.SenderCompID)
        }

        task.initiator.appHdl.OnEvent(sEventMsg)

        bIsLogouted := task.initiator.cliSessnCtrl.getLogoutStatus()
        if !bIsLogouted {
            // 被动收到logout

            // logout
            logoutmsg, _ := buildLogoutMsg(task.fInfo.BeginString, task.fInfo.TargetCompID, task.fInfo.SenderCompID, ibfixdefine.InvalidHbtInt)

            logoutMpkg := NewMsgPkg()
            logoutMpkg.Fixmsg = logoutmsg

            _, err := task.initiator.SendMsgToSvr(logoutMpkg)
            if nil != err {
                task.initiator.appHdl.OnError("send logout", err)
            }
        }

        // set session control
        task.initiator.cliSessnCtrl.logout()

    default:

    }
    task.initiator.appHdl.OnAdminMsg(task.fixMsg)
}

// generate admin fix msg process task
func genAdminMsgProc(serverIP string, serverPort uint16,
    fixMsg *quickfix.Message, fInfo ibfixdefine.FixInfoLocal, initiator *Initiator) {
    select {
    case <-initiator.chExit:
        return

    default:
    }

    if nil == initiator.adminMsgTPool {
        initiator.appHdl.OnErrorStr("nil *adminMsgTPool")
        return
    }

    var (
        job = &adminMsgTask{
            serverIP: serverIP, serverPort: serverPort,
            fixMsg: fixMsg, fInfo: fInfo, initiator: initiator}
        taskHold = kchanthreadpool.NewTask()
    )

    taskHold.Data = job

    busy, err := initiator.adminMsgTPool.AddTaskByMini(taskHold)
    if nil != err {
        initiator.appHdl.OnError("genAdminMsgProc", err)
    }

    if busy {
        initiator.appHdl.OnErrorStr("genAdminMsgProc busy")
    }
}
