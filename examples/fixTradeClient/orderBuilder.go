package main

import (
    "github.com/quickfixgo/quickfix"

    "github.com/chenjunpc2008/go-fix/common/fixtimeutil"
    "github.com/chenjunpc2008/go-fix/common/ibfixdefine"
)

const (
    cBeginString  = "FIX.4.2"
    cSenderCompID = "ETA"
    cTargetCompID = "NASDAQ"
)

func buildIBNewOrderSingle() *quickfix.Message {
    var order = quickfix.NewMessage()

    // header
    order.Header.SetString(ibfixdefine.Tag8, cBeginString)
    order.Header.SetString(ibfixdefine.Tag35, ibfixdefine.MsgTypeNewOrderSingle)
    order.Header.SetString(ibfixdefine.Tag49, cSenderCompID)
    order.Header.SetString(ibfixdefine.Tag56, cTargetCompID)

    // body
    order.Body.SetString(ibfixdefine.Tag1, "U901221")

    sCliOrdID := fixtimeutil.GetRandom16()
    order.Body.SetString(ibfixdefine.Tag11, sCliOrdID)

    order.Body.SetString(ibfixdefine.Tag55, "EUR")
    order.Body.SetString(ibfixdefine.Tag21, "2")
    order.Body.SetString(ibfixdefine.Tag54, "1")
    order.Body.SetString(ibfixdefine.Tag38, "100.00")
    order.Body.SetString(ibfixdefine.Tag40, "2")
    order.Body.SetString(ibfixdefine.Tag44, "18.80")

    return order
}
