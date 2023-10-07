package fixinitiator

import (
	"fmt"

	"github.com/chenjunpc2008/go-fix/common/ibfixdefine"

	"github.com/quickfixgo/quickfix"
)

/*
Session Level Reject
*/
func procSessionLevelRejectReq(serverIP string, serverPort uint16,
    fixMsg *quickfix.Message, fInfo ibfixdefine.FixInfoLocal, initiator *Initiator) {
    // RefSeqNo
    iRefSeqNo, fixerr := fixMsg.Body.GetInt(ibfixdefine.Tag45)
    if nil != fixerr {
        initiator.appHdl.OnError("Sequence Reset Request get RefSeqNo failed", fixerr)
        return
    }

    // Text
    sText, _ := fixMsg.Body.GetString(ibfixdefine.Tag58)

    sErrMsg := fmt.Sprintf("receive reject, RefSeqNo:%d, text:%s", iRefSeqNo, sText)

    initiator.appHdl.OnErrorStr(sErrMsg)
}
