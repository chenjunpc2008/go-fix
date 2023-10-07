package fixinitiator

import (
    "errors"

    "github.com/chenjunpc2008/go-fix/common/fixtimeutil"
    "github.com/chenjunpc2008/go-fix/common/ibfixdefine"
    "github.com/chenjunpc2008/go-fix/protocol/fixndp"
)

type initiatorHandler struct {
    appHdl    ApplicationIF
    initiator *Initiator
}

// Pack pack message into the []byte to be written
func (hdler *initiatorHandler) Pack(msg interface{}) ([]byte, error) {

    toSendMsg, ok := (msg).(*msgToSendSt)
    if !ok {
        hdler.appHdl.OnErrorStr("pack assert failed")
        return []byte{}, nil
    }

    var fixMsg = toSendMsg.pkg.Fixmsg

    switch toSendMsg.msgType {
    case mtsTypeNormal:
        // normal message
        // add MsgSeqNum and SendingTime
        seqnum := hdler.initiator.cliSessnCtrl.getThenIncrNextSenderMsgSeqNum()
        sendTime := fixtimeutil.GetFIXUTCTimestamp()

        fixMsg.Header.SetInt(ibfixdefine.Tag34, int(seqnum))
        fixMsg.Header.SetString(ibfixdefine.Tag52, sendTime)

    case mtsTypeResend:
        // message resend
    }

    sData := fixMsg.String()
    byData := []byte(sData)

    return byData, nil
}

// Depack depack the message packages from read []byte
func (hdler *initiatorHandler) Depack(rawData []byte) ([]byte, []interface{}) {
    var (
        dataRemain []byte
        pakgs      []interface{}
        ok         bool
        errMsg     string
    )

    dataRemain, pakgs, ok, errMsg = fixndp.Depack(rawData)
    if !ok {
        // 解包过程中有错误发送
        hdler.appHdl.OnError(errMsg, errors.New("depack failed"))

        if nil != hdler.initiator {
            // need to close this connection
        }
    }

    return dataRemain, pakgs
}

// OnNewConnection new connections event
func (hdler *initiatorHandler) OnNewConnection(serverIP string, serverPort uint16) {
    // const ftag = "initiatorHandler.OnNewConnection()"

    clisc := hdler.initiator.getCliSessnCtrl()
    if nil == clisc {
        go hdler.appHdl.OnErrorStr("initiatorHandler.OnNewConnection() get nil cliSessnCtrl")
        return
    }

    clisc.addNewConnection(serverIP, serverPort)

    go hdler.appHdl.OnConnected(serverIP, serverPort)

    responeMsg, _ := buildLogonMsg(hdler.initiator.setting.BeginString,
        hdler.initiator.setting.SenderCompID, hdler.initiator.setting.TargetCompID,
        30, hdler.initiator.setting.ResetOnLogon)

    mpkg := NewMsgPkg()
    mpkg.Fixmsg = responeMsg

    ok, err := hdler.initiator.SendMsgToSvr(mpkg)
    if nil != err {
        hdler.OnError("send logon", err)
    } else if !ok {
        hdler.OnErrorStr("send logon not ok")
    }
}

// disconnected event
func (hdler *initiatorHandler) OnDisconnected(serverIP string, serverPort uint16) {
    clisc := hdler.initiator.getCliSessnCtrl()
    if nil != clisc {
        clisc.disconnect(serverIP, serverPort)
    }

    go hdler.appHdl.OnDisconnected(serverIP, serverPort)
}

// receive data event
func (hdler *initiatorHandler) OnReceiveData(serverIP string, serverPort uint16, pPacks []interface{}) {
    // const ftag = "initiatorHandler.OnReceiveData()"

    // fmt.Printf("%s , %s:%d\n", ftag, serverIP, serverPort)
    clisc := hdler.initiator.getCliSessnCtrl()
    if nil != clisc {
        clisc.recordRcv(serverIP, serverPort)
    }

    //
    genTaskFixMsgDecode(serverIP, serverPort, pPacks, hdler.initiator)
}

// data already sended event
func (hdler *initiatorHandler) OnSendedData(serverIP string, serverPort uint16, msg interface{}, bysSended []byte, len int) {
    clisc := hdler.initiator.getCliSessnCtrl()
    if nil != clisc {
        clisc.recordSend(serverIP, serverPort)
    }

    toSendMsg, ok := (msg).(*msgToSendSt)
    if !ok {
        hdler.appHdl.OnErrorStr("OnSendedData assert failed")
        return
    }

    hdler.appHdl.OnSendedData(toSendMsg.pkg)
}

// event
func (hdler *initiatorHandler) OnEvent(msg string) {
    hdler.appHdl.OnEvent(msg)
}

// error
func (hdler *initiatorHandler) OnError(msg string, err error) {
    hdler.appHdl.OnError(msg, err)
}

// error
func (hdler *initiatorHandler) OnErrorStr(msg string) {
    hdler.appHdl.OnErrorStr(msg)
}
