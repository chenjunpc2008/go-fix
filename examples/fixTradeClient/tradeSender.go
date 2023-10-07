package main

import (
    "fmt"
    "sync/atomic"
    "time"

    "github.com/quickfixgo/quickfix"

    "github.com/chenjunpc2008/go-fix/fixinitiator"
)

var (
    cTradeNumsToSend = 1000
)

func tradeSender() {

    var (
        timeOut  = time.Duration(1) * time.Second
        fixOrder *quickfix.Message
        mpkg     = fixinitiator.NewMsgPkg()
        i        int
    )

    for i = 0; i < 100; i++ {
        select {
        case <-gChExit:
            return

        case <-time.After(timeOut):
            // sleep for one
        }

        for i := 0; i < cTradeNumsToSend; i++ {
            fixOrder = buildIBNewOrderSingle()
            mpkg.Fixmsg = fixOrder
            gFixCli.SendMsgToSvr(mpkg)
        }
    }

}

var (
    GRspRcvNum int64
)

func packetRcvTpsCalc() {
    var (
        timeout   = time.Duration(10) * time.Second
        sTime     string
        startNums int64
        endNums   int64
        delta     int64
    )

    for {
        startNums = atomic.AddInt64(&GRspRcvNum, 0)
        select {
        case <-time.After(timeout):
        }

        sTime = time.Now().Format("2006-01-02 15:04:05")

        endNums = atomic.AddInt64(&GRspRcvNum, 0)

        delta = endNums - startNums

        // 下单操作每个请求回复两条msg
        fmt.Printf("Until:%s, 10s total receive response:%d, tps:%d\n", sTime, delta, delta/10)
        fmt.Printf("total receive response:%d\n", GRspRcvNum)
    }
}
