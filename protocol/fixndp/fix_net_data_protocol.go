package fixndp

import (
    "fmt"
    "math"
    "strconv"
)

/*
package for fix net data pack and depack protocol
*/

/*
Depack depack the message packages from read []byte

@param rawData []byte : raw datas from net

@return []byte : half package data remained after depack
@return []interface{} : list of mesage packages
@return bool : is depack process succee
			true -- success
			flase -- have error, see errmsg for detail
@return string : error string
*/
func Depack(rawData []byte) (rawDataLocal []byte, outPacks []interface{}, bOk bool, sErrMsg string) {
    // const ftag = "fix_protocol.Depack()"
    rawDataLocal = rawData

    if MsgPacketMiniSize > len(rawData) {
        // 剩余的数据包不够包最小长度，还需继续缓存
        return rawData, nil, true, ""
    }

    // pos
    var (
        beginStringBond = []byte{'8', '='}
        bodyLengthBond  = []byte{'9', '='}
        msgTypeBond     = []byte{'3', '5', '='}
        checkSumBond    = []byte{'1', '0', '='}
    )

    var (
        err                 error
        i                   uint64
        uiBufferSize        uint64
        uiCurrentPos        uint64
        byTmp               []byte
        iMsgBodyLen         uint64
        sMsgType            string
        iBodyLengthBeginPos uint64
        iLenNeeded          uint64
        iMsgEndPos          uint64
    )

    uiCurrentPos = 0

    outPacks = make([]interface{}, 0)
    sErrMsg = ""
    bOk = true

    for {
        uiCurrentPos = 0

        uiBufferSize = uint64(len(rawDataLocal))
        if MsgPacketMiniSize > uiBufferSize {
            return rawDataLocal, outPacks, true, ""
        }

        /*
        	消息长度通过BodyLength域记录，表示BodyLength域值之后第一个域界定符<SOH>（不包括）与CheckSum域号前的最后一个域界定符<SOH>（包括）之间的字符个数。
        	8=FIX.1.1<SOH>9=?<SOH>
        */
        // 匹配包头
        // 8 BeginString
        for i = 0; i < 2; i++ {
            if beginStringBond[i] != rawDataLocal[uiCurrentPos+i] {
                // 从头开始的数据不是包头，需抛出错误
                sErrMsg = fmt.Sprintf("e-1 discard:%v", rawDataLocal[0:])
                rawDataLocal = make([]byte, 0)
                return rawDataLocal, outPacks, false, sErrMsg
            }
        }

        // forward
        uiCurrentPos += i

        // see to SOH
        uiCurrentPos = seekToSOH(uiCurrentPos, uiBufferSize, rawDataLocal)

        uiCurrentPos++
        if uiCurrentPos >= uiBufferSize {
            // 从头开始的数据不是包头，需抛出错误
            sErrMsg = fmt.Sprintf("e-2 discard:%v", rawDataLocal[0:])
            rawDataLocal = make([]byte, 0)
            return rawDataLocal, outPacks, false, sErrMsg
        }

        //
        // 9 BodyLength 消息体长度 N9
        for i = 0; i < 2; i++ {
            if uiCurrentPos+i >= uiBufferSize {
                // 越界，需抛出错误
                sErrMsg = fmt.Sprintf("e-3 discard:%v", rawDataLocal[0:])
                rawDataLocal = make([]byte, 0)
                return rawDataLocal, outPacks, false, sErrMsg
            }

            if bodyLengthBond[i] != rawDataLocal[uiCurrentPos+i] {
                // 从头开始的数据不是包头，需抛出错误
                sErrMsg = fmt.Sprintf("e-4 discard:%v", rawDataLocal[0:])
                rawDataLocal = make([]byte, 0)
                return rawDataLocal, outPacks, false, sErrMsg
            }
        }

        // forward
        uiCurrentPos += i

        // see to SOH
        i = seekToSOH(uiCurrentPos, uiBufferSize, rawDataLocal)

        /*
           消息长度通过BodyLength域记录，表示BodyLength域值之后第一个域界定符<SOH>（不包括）与CheckSum域号前的最后一个域界定符<SOH>（包括）之间的字符个数。
        */
        iBodyLengthBeginPos = i + 1

        if iBodyLengthBeginPos >= uiBufferSize {
            // 越界，需抛出错误
            sErrMsg = fmt.Sprintf("e-5 discard:%v", rawDataLocal[0:])
            rawDataLocal = make([]byte, 0)
            return rawDataLocal, outPacks, false, sErrMsg
        }

        byTmp = rawDataLocal[uiCurrentPos:i]

        iMsgBodyLen, err = strconv.ParseUint(string(byTmp), 10, 64)
        if nil != err {
            // 验证错误，清理错误的包头
            sErrMsg = fmt.Sprintf("error when get BodyLength, discard:%v", rawDataLocal[0:])
            rawDataLocal = make([]byte, 0)
            return rawDataLocal, outPacks, false, sErrMsg
        }

        if MaxReqMsgLen < iMsgBodyLen {
            sErrMsg = fmt.Sprintf("illegal BodyLength, discard:%v", rawDataLocal[0:])
            rawDataLocal = make([]byte, 0)
            return rawDataLocal, outPacks, false, sErrMsg
        }

        iLenNeeded = iBodyLengthBeginPos + iMsgBodyLen + CheckSumMinLen
        if iLenNeeded > uiBufferSize {
            // 剩余的数据端不够包头+包身+包尾数据长度，还需继续缓存
            return rawDataLocal, outPacks, true, ""
        }

        // forward
        uiCurrentPos = i + 1

        // 35 MsgType 消息类型 C16
        for i = 0; i < 3; i++ {
            if uiCurrentPos+i >= uiBufferSize {
                // 越界，需抛出错误
                sErrMsg = fmt.Sprintf("e-6 discard:%v", rawDataLocal[0:])
                rawDataLocal = make([]byte, 0)
                return rawDataLocal, outPacks, false, sErrMsg
            }

            if msgTypeBond[i] != rawDataLocal[uiCurrentPos+i] {
                sErrMsg = fmt.Sprintf("illegal MsgType, discard:%v", rawDataLocal[0:])
                rawDataLocal = make([]byte, 0)
                return rawDataLocal, outPacks, false, sErrMsg
            }
        }

        // forward
        uiCurrentPos += i

        // see to SOH
        i = seekToSOH(uiCurrentPos, uiBufferSize, rawDataLocal)

        if i >= uiBufferSize {
            // 越界，需抛出错误
            sErrMsg = fmt.Sprintf("e-7 discard:%v", rawDataLocal[0:])
            rawDataLocal = make([]byte, 0)
            return rawDataLocal, outPacks, false, sErrMsg
        }

        byTmp = rawDataLocal[uiCurrentPos:i]
        sMsgType = string(byTmp)

        // forward to the checksum
        uiCurrentPos = iBodyLengthBeginPos + iMsgBodyLen

        // CheckSum 10=XXX<SOH>
        for i = 0; i < 3; i++ {
            if uiCurrentPos+i >= uiBufferSize {
                // 越界，需抛出错误
                sErrMsg = fmt.Sprintf("e-8 discard:%v", rawDataLocal[0:])
                rawDataLocal = make([]byte, 0)
                return rawDataLocal, outPacks, false, sErrMsg
            }

            if checkSumBond[i] != rawDataLocal[uiCurrentPos+i] {
                // 校验错误，需抛出错误
                sErrMsg = fmt.Sprintf("e-8 discard:%v", rawDataLocal[0:])
                rawDataLocal = make([]byte, 0)
                return rawDataLocal, outPacks, false, sErrMsg
            }
        }

        // see to SOH
        i = seekToSOH(uiCurrentPos, uiBufferSize, rawDataLocal)

        if i >= uiBufferSize {
            // 越界，需抛出错误
            sErrMsg = fmt.Sprintf("e-9 discard:%v", rawDataLocal[0:])
            rawDataLocal = make([]byte, 0)
            return rawDataLocal, outPacks, false, sErrMsg
        }

        //
        iMsgEndPos = i + 1

        if iMsgEndPos > math.MaxUint32 {
            // 数据过长，需抛出错误
            sErrMsg = fmt.Sprintf("too long msg:%v, discard:%v", iMsgEndPos, rawDataLocal[0:])
            rawDataLocal = make([]byte, 0)
            return rawDataLocal, outPacks, false, sErrMsg
        }

        onePack := new(NetPacketSt)
        onePack.MsgType = sMsgType
        onePack.MsgBodyLen = uint32(iMsgEndPos)

        onePack.RawData = make([]byte, iMsgEndPos, iMsgEndPos*2+1)
        copy(onePack.RawData, rawDataLocal[0:iMsgEndPos])

        onePack.RawBody = onePack.RawData

        outPacks = append(outPacks, onePack)

        // forward
        uiCurrentPos = iMsgEndPos

        // 清除已接受的数据包
        if uiCurrentPos == uiBufferSize {
            rawDataLocal = []byte{}
        } else {
            rawDataLocal = rawDataLocal[uiCurrentPos:]
        }
    }
}

func seekToSOH(currPos uint64, bufferSize uint64, data []byte) uint64 {
    //
    var i uint64

    // seek to SOH
    for i = currPos; i < bufferSize; i++ {
        if SOH == data[i] {
            break
        }
    }

    return i
}
