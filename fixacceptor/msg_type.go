package fixacceptor

import (
    "bytes"
)

var (
    msgTypeHeartbeat     = []byte("0")
    msgTypeLogon         = []byte("A")
    msgTypeTestRequest   = []byte("1")
    msgTypeResendRequest = []byte("2")
    msgTypeReject        = []byte("3")
    msgTypeSequenceReset = []byte("4")
    msgTypeLogout        = []byte("5")
)

// isAdminMessageType returns true if the message type is a session level message.
func isAdminMessageType(m []byte) bool {
    switch {
    case bytes.Equal(msgTypeHeartbeat, m),
        bytes.Equal(msgTypeLogon, m),
        bytes.Equal(msgTypeTestRequest, m),
        bytes.Equal(msgTypeResendRequest, m),
        bytes.Equal(msgTypeReject, m),
        bytes.Equal(msgTypeSequenceReset, m),
        bytes.Equal(msgTypeLogout, m):
        return true
    }

    return false
}

var (
    msgTypeNews                      = []byte("B")
    msgTypeNewOrderSingle            = []byte("D")
    msgTypeOrderCancelRequest        = []byte("F")
    msgTypeOrderCancelReplaceRequest = []byte("G")
    msgTypeOrderStatusRequest        = []byte("H")
    msgTypeAllocationInstruction     = []byte("J")
    msgTypeAllocationInstructionAck  = []byte("P")
)

func isAppMessageType(m []byte) bool {

    switch {
    case bytes.Equal(msgTypeNews, m),
        bytes.Equal(msgTypeNewOrderSingle, m),
        bytes.Equal(msgTypeOrderCancelRequest, m),
        bytes.Equal(msgTypeOrderCancelReplaceRequest, m),
        bytes.Equal(msgTypeOrderStatusRequest, m),
        bytes.Equal(msgTypeAllocationInstruction, m),
        bytes.Equal(msgTypeAllocationInstructionAck, m):
        return true
    }

    return false
}
