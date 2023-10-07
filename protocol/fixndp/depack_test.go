package fixndp

import (
    "bytes"
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_Depack_logon_1(t *testing.T) {
    fixMsg := "8=FIX.4.29=6135=A34=149=TW52=20220110-08:39:59.49956=ISLD98=0108=3010=226"

    rawMsg := bytes.NewBufferString(fixMsg)

    afterDepackBuff, pPacks, ok, errstr := Depack(rawMsg.Bytes())
    assert.Equal(t, true, ok)
    assert.Equal(t, "", errstr)
    assert.Equal(t, 0, len(afterDepackBuff))

    if nil == pPacks {
        assert.Equal(t, 1, 0)
    }

    assert.Equal(t, 1, len(pPacks))

    v := pPacks[0]
    msg, ok := (v).(*NetPacketSt)
    assert.Equal(t, true, ok)

    assert.Equal(t, "A", msg.MsgType)
    assert.Equal(t, uint32(83), msg.MsgBodyLen)
    assert.Equal(t, fixMsg, string(msg.RawData))
}

func Test_Depack_hbt_2(t *testing.T) {
    fixMsg := "8=FIX.4.29=4935=034=249=TW52=20220110-08:40:29.51856=ISLD10=172"

    rawMsg := bytes.NewBufferString(fixMsg)

    afterDepackBuff, pPacks, ok, errstr := Depack(rawMsg.Bytes())
    assert.Equal(t, true, ok)
    assert.Equal(t, "", errstr)
    assert.Equal(t, 0, len(afterDepackBuff))

    if nil == pPacks {
        assert.Equal(t, 1, 0)
    }

    assert.Equal(t, 1, len(pPacks))

    v := pPacks[0]
    msg, ok := (v).(*NetPacketSt)
    assert.Equal(t, true, ok)

    assert.Equal(t, "0", msg.MsgType)
    assert.Equal(t, uint32(71), msg.MsgBodyLen)
    assert.Equal(t, fixMsg, string(msg.RawData))
}

func Test_Depack_missing_1(t *testing.T) {
    fixMsg := "8=FIX.4.29=6135=A34=149=TW52=20220110-08:39:59.49956=ISLD98=0108=3010=226"

    rawMsg := bytes.NewBufferString(fixMsg)

    afterDepackBuff, pPacks, ok, errstr := Depack(rawMsg.Bytes())
    assert.Equal(t, true, ok)
    assert.Equal(t, "", errstr)
    assert.Equal(t, 82, len(afterDepackBuff))

    if nil == pPacks {
        assert.Equal(t, 1, 0)
    }

    assert.Equal(t, 0, len(pPacks))
}

func Test_Depack_halfPacket_1(t *testing.T) {
    fixMsg := "8=FIX.4.29=6135=A34=149=TW52=20220110-08:39:59.49956=ISLD98=0108=3010=226"

    rawMsg := bytes.NewBufferString(fixMsg)

    afterDepackBuff, pPacks, ok, errstr := Depack(rawMsg.Bytes())
    assert.Equal(t, true, ok)
    assert.Equal(t, "", errstr)
    assert.Equal(t, 82, len(afterDepackBuff))

    if nil == pPacks {
        assert.Equal(t, 1, 0)
    }

    assert.Equal(t, 0, len(pPacks))

    // add missing piece
    fixMsg += ""
    rawMsg = bytes.NewBufferString(fixMsg)

    afterDepackBuff, pPacks, ok, errstr = Depack(rawMsg.Bytes())
    assert.Equal(t, true, ok)
    assert.Equal(t, "", errstr)
    assert.Equal(t, 0, len(afterDepackBuff))

    if nil == pPacks {
        assert.Equal(t, 1, 0)
    }

    assert.Equal(t, 1, len(pPacks))

    v := pPacks[0]
    msg, ok := (v).(*NetPacketSt)
    assert.Equal(t, true, ok)

    assert.Equal(t, "A", msg.MsgType)
    assert.Equal(t, uint32(83), msg.MsgBodyLen)
    assert.Equal(t, fixMsg, string(msg.RawData))
}
