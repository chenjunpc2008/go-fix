package fixndp

import (
    "fmt"
    "strconv"
    "strings"
)

// NetPacketSt net neet package
type NetPacketSt struct {
    MsgType    string // 消息类型
    MsgSeqNum  uint64 // 消息序号
    MsgBodyLen uint32 // 消息体长度
    RawBody    []byte
    RawData    []byte
}

/*
ToString debug string
*/
func (pack *NetPacketSt) ToString() string {

    var (
        sOut string
        b    strings.Builder
        sTmp string
    )

    b.WriteString("{MsgType:")
    b.WriteString(pack.MsgType)

    b.WriteString(", MsgSeqNum:")
    sTmp = strconv.FormatUint(pack.MsgSeqNum, 10)
    b.WriteString(sTmp)

    b.WriteString(", MsgBodyLen:")
    sTmp = strconv.FormatUint(uint64(pack.MsgBodyLen), 10)
    b.WriteString(sTmp)

    if nil == pack.RawData || 0 == len(pack.RawData) {
        b.WriteString(", RawData:{}")
    } else {
        sOut = fmt.Sprintf(", RawData:{%s}", string(pack.RawData))

        b.WriteString(sOut)
    }

    b.WriteString("}")

    sOut = b.String()
    return sOut
}
