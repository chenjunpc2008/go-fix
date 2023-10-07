package main

import (
    "bytes"
    "fmt"
    "log"

    "github.com/quickfixgo/quickfix"
    "github.com/quickfixgo/quickfix/datadictionary"
)

func main() {
    var (
        err              error
        fixDataModelDict *datadictionary.DataDictionary
    )

    FixDataDictPath := "./spec/FIX42.xml"

    fixDataModelDict, err = datadictionary.Parse(FixDataDictPath)
    if err != nil {
        log.Panicln(err)
    }

    fixMsg := quickfix.NewMessage()
    rawBody := logonMsg()
    err = quickfix.ParseMessageWithDataDictionary(fixMsg, rawBody, fixDataModelDict, nil)
    if err != nil {
        log.Printf(fmt.Sprintf("Msg Parse Error: %v, %q\n", err.Error(), rawBody))
    }
    log.Printf("parse fix msg success: %v", fixMsg)

    return
}

func logonMsg() *bytes.Buffer {
    fixMsg := "8=FIX.4.29=3635=A34=149=TW56=ISLD98=0108=3010=226"
    return bytes.NewBufferString(fixMsg)
}
