package fixndp

const (
    // MsgPacketMiniSize fix通讯协议报文数据块最低长度
    // 8=FIX.1.1<SOH>9=?<SOH>35=?<SOH>49=?<SOH>56=?<SOH>34=?<SOH>52=?<SOH>10=XXX<SOH>
    MsgPacketMiniSize = 46

    // SOH delimiter is ASCII 1 (SOH) symbol.
    SOH = byte(0x01)

    // CheckSumMinLen 3 byte simple checksum. Always represented as 3 ASCII characters.
    // 10=XXX<SOH>
    CheckSumMinLen = 7

    // MaxReqMsgLen max len
    MaxReqMsgLen = uint64(4 * 1024)
)
