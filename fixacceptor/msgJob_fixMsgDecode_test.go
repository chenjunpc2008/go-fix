package fixacceptor

import (
    "bytes"
    "fmt"
    "testing"

    "github.com/quickfixgo/quickfix"
    "github.com/quickfixgo/quickfix/datadictionary"
    "github.com/stretchr/testify/assert"
)

type myApp struct {
}

func (a *myApp) OnError(msg string, err error) {
    fmt.Printf(msg, err)
}

func (a *myApp) OnErrorStr(msg string) {
    fmt.Printf(msg)
}

func Test_fixMsgDecodeJob_Do(t *testing.T) {
    testMsg := "8=FIX.4.29=6135=A34=149=TW52=20220110-08:39:59.49956=ISLD98=0108=3010=226"
    rawMsg := bytes.NewBufferString(testMsg).Bytes()
    fixMsg := quickfix.NewMessage()
    rawBody := bytes.NewBuffer(rawMsg)
    fixDataModelDict, err := datadictionary.Parse("../../common/spec/FIX42.xml")
    assert.Equal(t, nil, err)
    err = quickfix.ParseMessageWithDataDictionary(fixMsg, rawBody, fixDataModelDict, nil)
    assert.Equal(t, nil, err)
    msgType, err := fixMsg.MsgType()
    assert.Equal(t, "A", msgType)
    var beginString quickfix.FIXString
    fixMsg.Header.GetField(quickfix.Tag(8), &beginString)
    assert.Equal(t, quickfix.FIXString("FIX.4.2"), beginString)
}
