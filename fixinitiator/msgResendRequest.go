package fixinitiator

import (
    "fmt"

    "github.com/chenjunpc2008/go-fix/common/ibfixdefine"

    "github.com/quickfixgo/quickfix"
)

/*
Message Recovery

Message gaps may occur during login or in the middle of a FIX session that require message
recovery. When this occurs, a Resend Request message can be sent requesting a range of
missing messages. The re-sender will then respond with a Sequence Reset message that has
tag 123 (GapFillFlag) set to “Y” and tag 43 (PossDupFlag) set to “Y”. Tag 36 (NewSeqNo) will be
set to the sequence number of the message to be redelivered next. Following this, the missing
Application level messages will be resent.
*/

/*
Resend Request
The resend request is sent by the receiving application to initiate the retransmission of messages. This
function is utilized if a sequence number gap is detected, if the receiving application lost a message, or
as a function of the initialization process.

The resend request can be used to request a single message, a range of messages or all messages
subsequent to a particular message.

Note: the sending application may wish to consider the message type when resending messages; e.g. if
a new order is in the resend series and a significant time period has elapsed since its original inception,
the sender may not wish to retransmit the order given the potential for changed market conditions.
(The Sequence Reset-GapFill message is used to skip messages that a sender does not wish to resend.)

Note: it is imperative that the receiving application process messages in sequence order, e.g. if
message number 7 is missed and 8-9 received, the application should ignore 8 and 9 and ask for a
resend of 7-9, or, preferably, 7-0 (0 represents infinity). This latter approach is strongly recommended
to recover from out of sequence conditions as it allows for faster recovery in the presence of certain
race conditions when both sides are simultaneously attempting to recover a gap.

• To request a single message: BeginSeqNo = EndSeqNo

• To request a range of messages: BeginSeqNo = first message of range, EndSeqNo = last message
of range

• To request all messages subsequent to a particular message: BeginSeqNo = first message of range,
EndSeqNo = 0 (represents infinity) .
*/

/*
Sequence Reset (Gap Fill) -
The sequence reset message is used by the sending application to reset the incoming sequence number
on the opposing side. This message has two modes: “Sequence Reset-Gap Fill” when GapFillFlag is
‘Y’ and “Sequence Reset-Reset” when GapFillFlag is N or not present. The “Sequence Reset-Reset”
mode should ONLY be used to recover from a disaster situation which cannot be otherwise recovered
via “Gap Fill” mode. The sequence reset message can be used in the following situations:

• During normal resend processing, the sending application may choose not to send a message (e.g.
an aged order). The Sequence Reset – Gap Fill is used to mark the place of that message.

• During normal resend processing, a number of administrative messages are not resent, the
Sequence Reset – Gap Fill message is used to fill the sequence gap created.

• In the event of an application failure, it may be necessary to force synchronization of sequence
numbers on the sending and receiving sides via the use of Sequence Reset - Reset

The sending application will initiate the sequence reset. The message in all situations specifies
NewSeqNo to reset as the value of the next sequence number immediately
following the messages and/or sequence numbers being skipped.

If the GapFillFlag field is not present (or set to N), it can be assumed that the purpose of the sequence
reset message is to recover from an out-of-sequence condition. The MsgSeqNum in the header should
be ignored (i.e. the receipt of a Sequence Reset - Reset message with an out of sequence MsgSeqNum
should not generate resend requests). Sequence Reset – Reset should NOT be used as a normal
response to a Resend Request (use Sequence Reset – Gap Fill). The Sequence Reset – Reset should
ONLY be used to recover from a disaster situation which cannot be recovered via the use of
Sequence Reset – Gap Fill. Note that the use of Sequence Reset – Reset may result in the possibility
of lost messages

If the GapFillFlag field is present (and equal to Y), the MsgSeqNum should conform to standard
message sequencing rules (i.e. the MsgSeqNum of the Sequence Reset-GapFill message should
represent the beginning MsgSeqNum in the GapFill range because the remote side is expecting that
next message).

The sequence reset can only increase the sequence number. If a sequence reset is received attempting
to decrease the next expected sequence number the message should be rejected and treated as a serious
error. It is possible to have multiple ResendRequests issued in a row (i.e. 5 to 10 followed by 5 to 11).
If sequence number 8, 10, and 11 represent application messages while the 5-7 and 9 represent
administrative messages, the series of messages as result of the Resend Request may appear as
SeqReset-GapFill with NewSeqNo of 8, message 8, SeqReset-GapFill with NewSeqNo of 10, and
message 10. This could then followed by SeqReset-GapFill with NewSeqNo of 8, message 8,
SeqReset-GapFill with NewSeqNo of 10, message 10, and message 11. One must be careful to ignore
the duplicate SeqReset-GapFill which is attempting to lower the next expected sequence number. This
can be detected by checking to see if its MsgSeqNum is less than expected. If so, the SeqReset-GapFill
is a duplicate and should be discarded.
*/
func procResendRequest(serverIP string, serverPort uint16,
    fixMsg *quickfix.Message, fInfo ibfixdefine.FixInfoLocal, initiator *Initiator) {

    // BeginSeqNo
    iBeginSeqNo, fixerr := fixMsg.Body.GetInt(ibfixdefine.Tag7)
    if nil != fixerr {
        initiator.appHdl.OnError("Resend Request get BeginSeqNo failed", fixerr)
        return
    }

    // EndSeqNo
    iEndSeqNo, fixerr := fixMsg.Body.GetInt(ibfixdefine.Tag16)
    if nil != fixerr {
        initiator.appHdl.OnError("Resend Request get EndSeqNo failed", fixerr)
        return
    }

    sEventMsg := fmt.Sprintf("Resend Request BeginSeqNo:%d, EndSeqNo:%d", iBeginSeqNo, iEndSeqNo)
    initiator.appHdl.OnEvent(sEventMsg)

    // 进入resend流程
    var err error
    // err = initiator.enterHoldingPatten()
    // if nil != err {
    //     initiator.appHdl.OnError("enterHoldingPatten failed", err)
    //     return
    // }

    if 0 == iEndSeqNo {
        // use Sequence Reset - Reset
        iNxtExptOut := initiator.cliSessnCtrl.getNextSenderMsgSeqNum()

        resetMsg, _ := buildSequenceResetResetMsg(initiator.setting.BeginString,
            initiator.setting.SenderCompID, initiator.setting.TargetCompID, int(iNxtExptOut))

        resetMpkg := NewMsgPkg()
        resetMpkg.Fixmsg = resetMsg

        err = initiator.SendPriorMsgToSvr(resetMpkg)
        if nil != err {
            initiator.appHdl.OnError("SendPriorMsgToSvr failed", err)
        }
    } else {
        // use Sequence Reset - Gap Fill
        // check if we need to push SenderMsgSeqNum higher
        iNxtExptOut := initiator.cliSessnCtrl.getNextSenderMsgSeqNum()

        if iNxtExptOut < uint64(iEndSeqNo) {
            // asked above local
            initiator.cliSessnCtrl.setNextSenderMsgSeqNum(uint64(iEndSeqNo + 1))
            sEventMsg = fmt.Sprintf("move local SenderMsgSeqNum from %d to %d", iNxtExptOut, iEndSeqNo+1)
            initiator.appHdl.OnEvent(sEventMsg)

            iNxtExptOut = uint64(iEndSeqNo)
        }

        resetMsg, _ := buildSequenceResetResetMsg(initiator.setting.BeginString,
            initiator.setting.SenderCompID, initiator.setting.TargetCompID, int(iNxtExptOut))

        resetMpkg := NewMsgPkg()
        resetMpkg.Fixmsg = resetMsg

        err = initiator.SendPriorMsgToSvr(resetMpkg)
        if nil != err {
            initiator.appHdl.OnError("SendPriorMsgToSvr failed", err)
        }
    }

    //
    // err = initiator.leaveHoldingPatten()
    // if nil != err {
    //     initiator.appHdl.OnError("leaveHoldingPatten failed", fixerr)
    //     return
    // }
}
