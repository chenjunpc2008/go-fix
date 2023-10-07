package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "runtime"
    "time"

    "github.com/chenjunpc2008/go-fix/fixinitiator"
    "github.com/chenjunpc2008/go/util/onlinepprof"
    "github.com/chenjunpc2008/go/util/panic-catch/cpanic"
)

var (
    // panic
    gPanicFile   *os.File
    gPanicLogger *log.Logger

    // exit channel
    gChExit chan int

    gFixCli *fixinitiator.Initiator
)

var (
    gbPprofEnable bool
    gPprofPort    uint
)

func init() {
    flag.BoolVar(&gbPprofEnable, "pprof_enbale", true, "pprof enable")
    flag.UintVar(&gPprofPort, "pprof_port", 10006, "pprof listen port")
    flag.IntVar(&cTradeNumsToSend, "max", 30000, "max req per second")
}

func main() {
    var (
        err error
    )

    // 获取参数
    flag.Parse()

    runtime.GOMAXPROCS(20)

    // 将 stderr 重定向到 f
    lNow := time.Now().Unix()
    sPFileName := fmt.Sprintf("log/%v_panic.log", lNow)
    gPanicFile, gPanicLogger, err = cpanic.NewPanicFile(sPFileName)
    if nil != err {
        log.Printf("cpanic.NewPanicFile %v\n", time.Now())
        return
    }

    gChExit = make(chan int)

    if gbPprofEnable {
        _, err = onlinepprof.StartOnlinePprof(uint16(gPprofPort), true)
        if nil != err {
            fmt.Printf("pprof, %v, %d, %v\n", gbPprofEnable, gPprofPort, err)
            return
        }
    }

    var setting = fixinitiator.DefaultConfig()
    setting.FixDataDictPath = "./spec/FIX42.xml"
    setting.RemoteIP = "127.0.0.1"
    setting.RemotePort = 5001

    setting.BeginString = "FIX.4.2"
    setting.SenderCompID = "qaib84"
    setting.TargetCompID = "IB"

    app := new(tradeCliApp)

    gFixCli = fixinitiator.NewInitiator(setting, app)

    var (
        senderStartSeqNum uint64 = 1
        targetStartSeqNum uint64 = 1
    )
    gFixCli.Start(senderStartSeqNum, targetStartSeqNum)

    select {
    case <-time.After(time.Duration(5) * time.Second):
        // wait for client to establish
    }

    // go tradeSender()
    // go packetRcvTpsCalc()

    select {
    case <-gChExit:
        return
    }
}
