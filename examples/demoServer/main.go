package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "runtime"
    "time"

    "github.com/chenjunpc2008/go-fix/fixacceptor"
    "github.com/chenjunpc2008/go/util/onlinepprof"
    "github.com/chenjunpc2008/go/util/panic-catch/cpanic"
)

var (
    // panic
    gPanicFile   *os.File
    gPanicLogger *log.Logger

    gbPprofEnable bool
    gPprofPort    uint
)

func init() {
    flag.BoolVar(&gbPprofEnable, "pprof_enbale", true, "pprof enable")
    flag.UintVar(&gPprofPort, "pprof_port", 10005, "pprof listen port")
}

func main() {

    // 获取参数
    flag.Parse()

    runtime.GOMAXPROCS(20)

    var err error

    // 将 stderr 重定向到 f
    lNow := time.Now().Unix()
    sPFileName := fmt.Sprintf("log/%v_panic.log", lNow)
    gPanicFile, gPanicLogger, err = cpanic.NewPanicFile(sPFileName)
    if nil != err {
        log.Printf("cpanic.NewPanicFile %v\n", time.Now())
        return
    }

    var chExit = make(chan int)

    if gbPprofEnable {
        _, err := onlinepprof.StartOnlinePprof(uint16(gPprofPort), true)
        if nil != err {
            fmt.Printf("pprof, %v, %d, %v\n", gbPprofEnable, gPprofPort, err)
            return
        }
    }

    var setting = fixacceptor.DefaultConfig()
    setting.Port = 5001
    setting.SenderCompID = "NASDAQ"
    setting.FixDataDictPath = "./spec/FIX42.xml"

    app := new(myApp)

    acceptor := fixacceptor.NewAcceptor(setting, app)

    acceptor.Start()

    go SvrTpsTest()

    select {
    case <-chExit:
        return
    }
}
