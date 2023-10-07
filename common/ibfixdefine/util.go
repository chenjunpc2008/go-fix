package ibfixdefine

import (
    "github.com/quickfixgo/quickfix"
)

/*
IsAdminMsg check if it's admin message
@return bool : true -- it is admin message
               false -- not admin message
*/
func IsAdminMsg(fixmsg *quickfix.Message) bool {
    if nil == fixmsg {
        return false
    }

    msgType, fixErr := fixmsg.Header.GetString(Tag35)
    if nil != fixErr {
        return false
    }

    switch msgType {
    case MsgTypeLogon,
        MsgTypeHeartbeat,
        MsgTypeTestRequest,
        MsgTypeResendRequest,
        MsgTypeSequenceReset,
        MsgTypeSessionLevelReject,
        MsgTypeLogout:
        return true

    default:
        return false
    }
}
