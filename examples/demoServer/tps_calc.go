package main

import (
    "fmt"
    "sync/atomic"
    "time"
)

var (
    GReqProcNum int64
)

func SvrTpsTest() {
    var (
        timeout   = time.Duration(10) * time.Second
        sTime     string
        startNums int64
        endNums   int64
        delta     int64
    )

    for {
        startNums = atomic.AddInt64(&GReqProcNum, 0)

        time.Sleep(timeout)

        sTime = time.Now().Format("2006-01-02 15:04:05")

        endNums = atomic.AddInt64(&GReqProcNum, 0)

        delta = endNums - startNums

        // 下单操作每个请求回复两条msg
        fmt.Printf("Until:%s, 10s total process orders:%d, tps:%d\n", sTime, delta, delta/10)
        fmt.Printf("total process orders:%d\n", GReqProcNum)

    }
}
