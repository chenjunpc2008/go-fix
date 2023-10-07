package fixinitiator

import "github.com/quickfixgo/quickfix"

// ApplicationIF interface for upper FIX application
type ApplicationIF interface {
    OnError(msg string, err error)
    OnErrorStr(msg string)
    OnEvent(msg string)
    OnConnected(serverIP string, serverPort uint16)    // connected event
    OnDisconnected(serverIP string, serverPort uint16) // disconnected event
    OnAdminMsg(fixmsg *quickfix.Message)               // receive admin message
    OnAppMsg(fixmsg *quickfix.Message)                 // receive application message
    OnSendedData(msg *MsgPkg)                          // message already sent
}
