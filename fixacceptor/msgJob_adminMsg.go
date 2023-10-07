package fixacceptor

import (
    "github.com/quickfixgo/quickfix"

    "github.com/chenjunpc2008/go-fix/common/ibfixdefine"
    "github.com/chenjunpc2008/go/util/kchanthreadpool"
)

type adminMsgJob struct {
    clientID   uint64
    clientIP   string
    clientAddr string
    fixMsg     *quickfix.Message
    acceptor   *Acceptor
}

func (job *adminMsgJob) Do(threadID uint) {
    if nil == job.fixMsg {
        job.acceptor.appHdl.OnErrorStr("nil *adminMsgJob")
        return
    }

    select {
    case <-job.acceptor.chExit:
        break

    default:
    }

    msgType, fixErr := job.fixMsg.MsgType()
    if nil != fixErr {
        job.acceptor.appHdl.OnError("get MsgType", fixErr)
        return
    }
    switch msgType {
    case ibfixdefine.MsgTypeLogon:
        job.procLogon()

    case ibfixdefine.MsgTypeLogout:
        job.procLogout()

    case ibfixdefine.MsgTypeHeartbeat:
        // Required when the heartbeat is the result of a Test Request message.
        sTestReqID, fixerr := job.fixMsg.Body.GetString(ibfixdefine.Tag112)
        if nil == fixerr {
            sEventMsg := "receive heartbeat forced by test request, TestReqID:" + sTestReqID
            job.acceptor.appHdl.OnEvent(sEventMsg)
        }

    case ibfixdefine.MsgTypeTestRequest:
        job.procTestRequest()

    case ibfixdefine.MsgTypeResendRequest:

    default:

    }

    job.acceptor.appHdl.FromAdmin(job.fixMsg)
}

func (job *adminMsgJob) procLogon() (err error) {
    var (
        sBeginString  quickfix.FIXString
        sSenderCompID quickfix.FIXString
        sTargetCompID quickfix.FIXString
        iHeartBtInt   quickfix.FIXInt
    )
    err = job.fixMsg.Header.GetField(quickfix.Tag(8), &sBeginString)
    if err != nil {
        job.acceptor.appHdl.OnError("get beginString err ", err)
        return
    }
    err = job.fixMsg.Header.GetField(quickfix.Tag(49), &sSenderCompID)
    if err != nil {
        job.acceptor.appHdl.OnError("get senderCompID err ", err)
        return
    }
    err = job.fixMsg.Header.GetField(quickfix.Tag(56), &sTargetCompID)
    if err != nil {
        job.acceptor.appHdl.OnError("get targetCompID err ", err)
        return
    }
    err = job.fixMsg.Body.GetField(quickfix.Tag(108), &iHeartBtInt)
    if err != nil {
        job.acceptor.appHdl.OnError("procLogin get iHeartBtInt err ", err)
        return
    }

    var (
        sRtSenderCompID string
        sRtTargetCompID string
        iRtHeartBtInt   uint16
        fixMsg          *quickfix.Message
    )
    sRtSenderCompID, sRtTargetCompID, iRtHeartBtInt, _, err = job.acceptor.cliSessnCtrl.logon(job.clientID, job.clientIP,
        job.clientAddr, string(sBeginString), string(sSenderCompID), string(sTargetCompID), uint16(iHeartBtInt))
    if err != nil {
        job.acceptor.appHdl.OnError("logon err ", err)
        return err
    }

    fixMsg, _ = buildLogonMsg(string(sBeginString), sRtSenderCompID, sRtTargetCompID, iRtHeartBtInt)

    err = job.acceptor.SendToClient(job.clientID, fixMsg)
    if nil != err {
        job.acceptor.appHdl.OnError("send logon err:", err)
        job.acceptor.cliSessnCtrl.logout(job.clientID)
    }

    job.acceptor.cliSessnCtrl.mapSIDCtrl[sSenderCompID.String()] = job.acceptor.cliSessnCtrl.mapCIDCtrl[job.clientID]

    return nil
}

func (job *adminMsgJob) procLogout() {

}

func (job *adminMsgJob) procTestRequest() error {
    var (
        clientSenderCompID quickfix.FIXString
    )
    // The test request message forces a heartbeat from the opposing application.
    sTestReqID, fixerr := job.fixMsg.Body.GetString(ibfixdefine.Tag112)
    if nil != fixerr {
        job.acceptor.appHdl.OnError("get TestReqID failed", fixerr)
    }

    job.fixMsg.Header.GetField(quickfix.Tag(49), &clientSenderCompID)

    targetCompID := clientSenderCompID
    hbtReqMsg, _ := buildHbtMsg(job.acceptor.setting.BeginString,
        job.acceptor.setting.SenderCompID, targetCompID.String(), sTestReqID)

    err := job.acceptor.SendToClient(job.clientID, hbtReqMsg)
    if nil != err {
        job.acceptor.appHdl.OnError("send heartbeat forced by test request", err)
    }
    return nil
}

func genAdminMsgProc(clientID uint64, clientIP string, clientAddr string, fixMsg *quickfix.Message, acceptor *Acceptor) {
    select {
    case <-acceptor.chExit:
        return

    default:
    }

    if nil == acceptor.appMsgTPool {
        acceptor.appHdl.OnErrorStr("nil *adminMsgTPool")
        return
    }

    job := &adminMsgJob{
        clientID:   clientID,
        clientIP:   clientIP,
        clientAddr: clientAddr,
        fixMsg:     fixMsg,
        acceptor:   acceptor,
    }
    taskHold := kchanthreadpool.NewTask()
    taskHold.Data = job

    busy, err := acceptor.appMsgTPool.AddTaskByMini(taskHold)
    if nil != err {
        acceptor.appHdl.OnError("genAdminMsgProc", err)
    }

    if busy {
        acceptor.appHdl.OnErrorStr("genAdminMsgProc busy")
    }
}
