package fixacceptor

import (
    "errors"
    "fmt"

    "github.com/chenjunpc2008/go-fix/common/fixtimeutil"
    "github.com/chenjunpc2008/go-fix/common/ibfixdefine"
    "github.com/chenjunpc2008/go-fix/protocol/fixndp"
    "github.com/quickfixgo/quickfix"
)

type acceptorHandler struct {
    appHdl   ApplicationIF
    acceptor *Acceptor
}

// pack message into the []byte to be written
func (hdler *acceptorHandler) Pack(clientID uint64, cliIP string, cliAddr string, msg interface{}) ([]byte, error) {
    fixMsg, ok := (msg).(*quickfix.Message)
    if !ok {
        hdler.appHdl.OnErrorStr("invalid msg type")
    }
    // add MsgSeqNum and SendingTime
    seqnum, _ := hdler.acceptor.cliSessnCtrl.getSenderSeqnum(clientID)
    sendTime := fixtimeutil.GetFIXUTCTimestamp()

    fixMsg.Header.SetInt(ibfixdefine.Tag34, int(seqnum))
    fixMsg.Header.SetString(ibfixdefine.Tag52, sendTime)

    sData := fixMsg.String()
    byData := []byte(sData)

    return byData, nil
}

// depack the message packages from read []byte
func (hdler *acceptorHandler) Depack(clientID uint64, cliIP string, cliAddr string, rawData []byte) ([]byte, []interface{}) {
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

        if nil != hdler.acceptor {
            // need to close this connection
        }
    }

    return dataRemain, pakgs
}

// new connections event
func (hdler *acceptorHandler) OnNewConnection(clientID uint64, clientIP string, clientAddr string) {
    const ftag = "acceptorHandler.OnNewConnection()"

    ok := hdler.acceptor.cliSessnCtrl.addNewConnection(clientID, clientIP, clientAddr)
    if !ok {
        sErrMsg := fmt.Sprintf("%s addNewConnection failed, clientID:%d, clientIP:%s, clientAddr:%s",
            ftag, clientID, clientIP, clientAddr)
        hdler.appHdl.OnErrorStr(sErrMsg)
    }
}

// disconnected event
func (hdler *acceptorHandler) OnDisconnected(clientID uint64, clientIP string, clientAddr string) {
    err := hdler.acceptor.cliSessnCtrl.disconnect(clientID, clientIP, clientAddr)
    if nil != err {
        hdler.appHdl.OnError("", err)
    }
}

// receive data event
func (hdler *acceptorHandler) OnReceiveData(clientID uint64, clientIP string, clientAddr string, pPacks []interface{}) {
    err := hdler.acceptor.cliSessnCtrl.recordRcv(clientID, clientIP, clientAddr)
    if nil != err {
        hdler.appHdl.OnError("", err)
    }

    //
    genJobFixMsgDecode(clientID, clientIP, clientAddr, pPacks, hdler.acceptor)
}

// data already sended event
func (hdler *acceptorHandler) OnSendedData(clientID uint64, clientIP string, clientAddr string, msg interface{}, bysSended []byte, len int) {
    err := hdler.acceptor.cliSessnCtrl.recordSend(clientID, clientIP, clientAddr)
    if nil != err {
        hdler.appHdl.OnError("", err)
    }

    fixMsg, ok := msg.(*quickfix.Message)
    if !ok {
        return
    }
    hdler.appHdl.OnSendedData(fixMsg)
}

// event
func (hdler *acceptorHandler) OnEvent(msg string) {
    hdler.appHdl.OnEvent(msg)
}

// error
func (hdler *acceptorHandler) OnError(msg string, err error) {
    hdler.appHdl.OnError(msg, err)
}

// error
func (hdler *acceptorHandler) OnCliError(clientID uint64, clientIP string, clientAddr string, msg string, err error) {
    sErrMsg := fmt.Sprintf("cid:%d, ip:%s, addr:%s, %s", clientID, clientIP, clientAddr, msg)
    hdler.appHdl.OnError(sErrMsg, err)
}

// error
func (hdler *acceptorHandler) OnCliErrorStr(clientID uint64, clientIP string, clientAddr string, msg string) {
    sErrMsg := fmt.Sprintf("cid:%d, ip:%s, addr:%s, %s", clientID, clientIP, clientAddr, msg)
    hdler.appHdl.OnErrorStr(sErrMsg)
}
