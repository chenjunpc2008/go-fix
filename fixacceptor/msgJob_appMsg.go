package fixacceptor

import (
	"github.com/quickfixgo/quickfix"

	"github.com/chenjunpc2008/go/util/kchanthreadpool"
)

type appMsgJob struct {
    clientID   uint64
    clientIP   string
    clientAddr string
    fixMsg     *quickfix.Message
    acceptor   *Acceptor
}

func (job *appMsgJob) Do(threadID uint) {
    //
    if nil == job.fixMsg {
        job.acceptor.appHdl.OnErrorStr("nil *adminMsgJob")
        return
    }

    select {
    case <-job.acceptor.chExit:
        break

    default:
    }

    // 业务处理交给上层
    job.acceptor.appHdl.FromApp(job.fixMsg)
}

func genAppMsgProc(clientID uint64, clientIP string, clientAddr string, fixMsg *quickfix.Message, acceptor *Acceptor) {
    select {
    case <-acceptor.chExit:
        return

    default:
    }

    if nil == acceptor.adminMsgTPool {
        acceptor.appHdl.OnErrorStr("nil *adminMsgTPool")
        return
    }

    job := &appMsgJob{
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
        acceptor.appHdl.OnError("genAppMsgProc", err)
    }

    if busy {
        acceptor.appHdl.OnErrorStr("genAppMsgProc busy")
    }
}
