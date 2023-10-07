package fixinitiator

import (
    "github.com/quickfixgo/quickfix"

    "github.com/chenjunpc2008/go-fix/common/ibfixdefine"
    "github.com/chenjunpc2008/go/util/kchanthreadpool"
)

type appMsgTask struct {
    serverIP   string
    serverPort uint16
    fixMsg     *quickfix.Message
    fInfo      ibfixdefine.FixInfoLocal
    initiator  *Initiator
}

func (task *appMsgTask) Do(threadID uint) {
    //
    if nil == task.fixMsg {
        task.initiator.appHdl.OnErrorStr("nil *appMsgTask")
        return
    }

    select {
    case <-task.initiator.chExit:
        break

    default:
    }

    task.initiator.appHdl.OnAppMsg(task.fixMsg)
}

// generate app fix msg process task
func genAppMsgProc(serverIP string, serverPort uint16,
    fixMsg *quickfix.Message, fInfo ibfixdefine.FixInfoLocal, initiator *Initiator) {
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
        job = &appMsgTask{
            serverIP: serverIP, serverPort: serverPort,
            fixMsg: fixMsg, fInfo: fInfo, initiator: initiator}
        taskHold = kchanthreadpool.NewTask()
    )

    taskHold.Data = job

    busy, err := initiator.appMsgTPool.AddTaskByMini(taskHold)
    if nil != err {
        initiator.appHdl.OnError("genAppMsgProc", err)
    }

    if busy {
        initiator.appHdl.OnErrorStr("genAppMsgProc busy")
    }
}
