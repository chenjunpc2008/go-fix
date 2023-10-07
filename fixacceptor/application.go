package fixacceptor

import (
    "github.com/quickfixgo/quickfix"
)

// ApplicationIF interface for upper FIX application
type ApplicationIF interface {
    OnError(msg string, err error)
    OnErrorStr(msg string)
    OnEvent(msg string)
    FromAdmin(fixMsg *quickfix.Message)
    FromApp(fixMsg *quickfix.Message)
    OnSendedData(fixMsg *quickfix.Message)
    OnDisconnected(serverIP string, serverPort uint16)
}
