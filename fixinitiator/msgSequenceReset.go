package fixinitiator

import (
	"fmt"

	"github.com/quickfixgo/quickfix"

	"github.com/chenjunpc2008/go-fix/common/ibfixdefine"
)

/*
Sequence Reset and Gap Fill
*/
func procSequenceResetReq(serverIP string, serverPort uint16,
    fixMsg *quickfix.Message, fInfo ibfixdefine.FixInfoLocal, initiator *Initiator) {
    // GapFillFlag
    bGapFillFlag, fixerr := fixMsg.Body.GetBool(ibfixdefine.Tag123)
    if nil != fixerr {
        initiator.appHdl.OnError("Sequence Reset Request get GapFillFlag failed", fixerr)
        return
    }

    // NewSeqNo
    iNewSeqNo, fixerr := fixMsg.Body.GetInt(ibfixdefine.Tag36)
    if nil != fixerr {
        initiator.appHdl.OnError("Resend Request get NewSeqNo failed", fixerr)
        return
    }

    uiNewSeqNo := uint64(iNewSeqNo)

    if bGapFillFlag {
        // Y = Gap Fill message, MsgSeqNum field valid
        return
    }

    // N = Sequence Reset, ignore MsgSeqNum
    nextTargetMsgSeqNum := initiator.cliSessnCtrl.getNextTargetMsgSeqNum()

    if uiNewSeqNo > nextTargetMsgSeqNum {
        initiator.cliSessnCtrl.setNextTargetMsgSeqNum(uiNewSeqNo)
    } else if uiNewSeqNo < nextTargetMsgSeqNum {
        // TODO: need to reject
        sErrMsg := fmt.Sprintf("Resend Request NewSeqNo:%d, nextTargetMsgSeqNum:%d", uiNewSeqNo, nextTargetMsgSeqNum)
        initiator.appHdl.OnErrorStr(sErrMsg)
    }
}
