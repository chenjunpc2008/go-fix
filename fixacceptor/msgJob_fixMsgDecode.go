package fixacceptor

import (
    "bytes"
    "fmt"

    "github.com/quickfixgo/quickfix"

    "github.com/chenjunpc2008/go-fix/protocol/fixndp"
    "github.com/chenjunpc2008/go/util/kchanthreadpool"
)

type fixMsgDecodeJob struct {
    clientID   uint64
    clientIP   string
    clientAddr string
    pPacks     []interface{}
    acceptor   *Acceptor
}

func (job *fixMsgDecodeJob) Do(threadID uint) {
    if nil == job.pPacks {
        job.acceptor.appHdl.OnErrorStr("nil pPacks")
        return
    }
    var (
        netPacket *fixndp.NetPacketSt
        ok        bool
        err       error
        fixMsg    *quickfix.Message
        msgType   string
        rawBody   *bytes.Buffer
    )

    for _, elem := range job.pPacks {
        netPacket, ok = (elem).(*fixndp.NetPacketSt)
        if !ok {
            job.acceptor.appHdl.OnErrorStr("fixMsgDecodeJob get invalid packet")
            continue
        }

        fixMsg = quickfix.NewMessage()
        rawBody = bytes.NewBuffer(netPacket.RawBody)
        err = quickfix.ParseMessageWithDataDictionary(fixMsg, rawBody, job.acceptor.fixDataModelDict, nil)
        if err != nil {
            job.acceptor.appHdl.OnErrorStr(fmt.Sprintf("Msg Parse Error: %v, %q", err.Error(), rawBody))
            continue
        }

        msgType, err = fixMsg.MsgType()
        if err != nil {
            job.acceptor.appHdl.OnError("get msgType from fixMsg failed ", err)
            continue
        }

        if isAdminMessageType([]byte(msgType)) {
            genAdminMsgProc(job.clientID, job.clientIP, job.clientAddr, fixMsg, job.acceptor)

        } else if isAppMessageType([]byte(msgType)) {
            genAppMsgProc(job.clientID, job.clientIP, job.clientAddr, fixMsg, job.acceptor)

        } else {
            job.acceptor.appHdl.OnErrorStr(fmt.Sprintf("invalid msg type:%s", string(msgType)))
        }
    }
}

// generate fix msg decode
func genJobFixMsgDecode(clientID uint64, clientIP string, clientAddr string, pPacks []interface{}, acceptor *Acceptor) {
    select {
    case <-acceptor.chExit:
        return

    default:
    }

    if nil == acceptor.appMsgTPool {
        acceptor.appHdl.OnErrorStr("nil *appMsgTPool")
        return
    }

    var (
        job = &fixMsgDecodeJob{
            clientID:   clientID,
            clientIP:   clientIP,
            clientAddr: clientAddr,
            pPacks:     pPacks,
            acceptor:   acceptor}
        taskHold = kchanthreadpool.NewTask()
    )

    taskHold.Data = job

    busy, err := acceptor.appMsgTPool.AddTaskByMini(taskHold)
    if nil != err {
        acceptor.appHdl.OnError("genJobFixMsgDecode", err)
    }

    if busy {
        acceptor.appHdl.OnErrorStr("genJobFixMsgDecode busy")
    }
}
