package main

import (
    "fmt"
    "github.com/quickfixgo/quickfix"
    "sync/atomic"
    "time"
)

func buildNewOrderSingleMsg() (newOrderMsg *quickfix.Message) {
    newOrderMsg = quickfix.NewMessage()

    newOrderMsg.Header.SetString(quickfix.Tag(8), "FIX.4.2")
    newOrderMsg.Header.SetString(quickfix.Tag(35), "D")
    newOrderMsg.Header.SetString(quickfix.Tag(49), "ETA")
    newOrderMsg.Header.SetString(quickfix.Tag(56), "NASDAQ")

    newOrderMsg.Body.SetString(quickfix.Tag(1), "U123456")
    newOrderMsg.Body.SetString(quickfix.Tag(11), "ClOrdID")
    newOrderMsg.Body.SetString(quickfix.Tag(55), "AAPL")
    newOrderMsg.Body.SetString(quickfix.Tag(21), "1")
    newOrderMsg.Body.SetString(quickfix.Tag(54), "1")
    newOrderMsg.Body.SetString(quickfix.Tag(38), "100")
    newOrderMsg.Body.SetString(quickfix.Tag(40), "1")
    newOrderMsg.Body.SetString(quickfix.Tag(100), "1")

    return
}

var (
    maxPerSecond int64
)

func orderNewSingle() {
    var (
        cnt          int64
        orderNewMsg  *quickfix.Message
    )
    orderNewMsg = buildNewOrderSingleMsg()

    for {
        for cnt = 0; cnt < maxPerSecond; cnt++ {
            quickfix.Send(orderNewMsg)
        }
        time.Sleep(1 * time.Second)
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
    }
}
