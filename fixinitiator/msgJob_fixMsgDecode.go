package fixinitiator

import (
    "bytes"
    "fmt"

    "github.com/quickfixgo/quickfix"

    "github.com/chenjunpc2008/go-fix/common/ibfixdefine"
    "github.com/chenjunpc2008/go-fix/protocol/fixndp"
    "github.com/chenjunpc2008/go/util/kchanthreadpool"
)

type fixMsgDecodeTask struct {
    serverIP   string
    serverPort uint16
    pPacks     []interface{}
    initiator  *Initiator
}

func (task *fixMsgDecodeTask) Do(threadID uint) {
    //
    if nil == task.pPacks {
        task.initiator.appHdl.OnErrorStr("nil pPacks")
        return
    }

    var (
        rawPack *fixndp.NetPacketSt
        ok      bool
        fixMsg  *quickfix.Message
        err     error
        rawMsg  *bytes.Buffer
        sErrMsg string
        fInfo   ibfixdefine.FixInfoLocal
    )

    for _, elem := range task.pPacks {
        rawPack, ok = (elem).(*fixndp.NetPacketSt)
        if !ok {
            task.initiator.appHdl.OnErrorStr("type assert not *fixndp.NetPacketSt")
            continue
        }

        fixMsg = quickfix.NewMessage()
        rawMsg = bytes.NewBuffer(rawPack.RawData)
        err = quickfix.ParseMessageWithDataDictionary(fixMsg, rawMsg, nil, task.initiator.fixDataDict)
        if nil != err {
            sErrMsg = fmt.Sprintf("parse msg failed, %s", string(rawPack.RawData))
            task.initiator.appHdl.OnErrorStr(sErrMsg)
            continue
        }

        // prevalid message
        fInfo, sErrMsg, err = fixMsgPreValid(fixMsg)
        if nil != err {
            task.initiator.appHdl.OnError(sErrMsg, err)
            continue
        }

        switch fInfo.MsgType {
        case ibfixdefine.MsgTypeLogon,
            ibfixdefine.MsgTypeHeartbeat,
            ibfixdefine.MsgTypeTestRequest,
            ibfixdefine.MsgTypeResendRequest,
            ibfixdefine.MsgTypeSequenceReset,
            ibfixdefine.MsgTypeSessionLevelReject,
            ibfixdefine.MsgTypeLogout:
            genAdminMsgProc(task.serverIP, task.serverPort,
                fixMsg, fInfo, task.initiator)

        default:
            genAppMsgProc(task.serverIP, task.serverPort,
                fixMsg, fInfo, task.initiator)
        }

    }

}

func fixMsgPreValid(fixMsg *quickfix.Message) (fInfo ibfixdefine.FixInfoLocal, errMsg string, errRet error) {

    var (
        fixErr quickfix.MessageRejectError
    )

    // 35 MsgType
    fInfo.MsgType, fixErr = fixMsg.Header.GetString(ibfixdefine.Tag35)
    if nil != fixErr {
        return fInfo, "get MsgType", fixErr
    }

    // 8 BeginString
    fInfo.BeginString, fixErr = fixMsg.Header.GetString(ibfixdefine.Tag8)
    if nil != fixErr {
        return fInfo, "get BeginString", fixErr
    }

    // 49 SenderCompID
    fInfo.SenderCompID, fixErr = fixMsg.Header.GetString(ibfixdefine.Tag49)
    if nil != fixErr {
        return fInfo, "get SenderCompID", fixErr
    }

    // 56 TargetCompID Default is “IB.”
    fInfo.TargetCompID, fixErr = fixMsg.Header.GetString(ibfixdefine.Tag56)
    if nil != fixErr {
        return fInfo, "get TargetCompID", fixErr
    }

    var (
        iMsgSeqNum int
    )

    // 34 MsgSeqNum
    iMsgSeqNum, fixErr = fixMsg.Header.GetInt(ibfixdefine.Tag34)
    if nil != fixErr {
        return fInfo, "get MsgSeqNum", fixErr
    }

    fInfo.MsgSeqNum = uint64(iMsgSeqNum)

    return fInfo, "", nil
}

// generate fix msg decode
func genTaskFixMsgDecode(serverIP string, serverPort uint16, pPacks []interface{}, initiator *Initiator) {
    select {
    case <-initiator.chExit:
        return

    default:
    }

    if nil == initiator.appMsgTPool {
        initiator.appHdl.OnErrorStr("nil *appMsgTPool")
        return
    }

    var (
        task = &fixMsgDecodeTask{
            serverIP: serverIP, serverPort: serverPort,
            pPacks: pPacks, initiator: initiator}
        taskHold = kchanthreadpool.NewTask()
    )

    taskHold.Data = task

    busy, err := initiator.appMsgTPool.AddTaskByMini(taskHold)
    if nil != err {
        initiator.appHdl.OnError("genTaskFixMsgDecode", err)
    }

    if busy {
        initiator.appHdl.OnErrorStr("genTaskFixMsgDecode busy")
    }
}
