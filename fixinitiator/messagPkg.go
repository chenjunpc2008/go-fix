package fixinitiator

import "github.com/quickfixgo/quickfix"

/*
MsgPkg outgoing send message package
*/
type MsgPkg struct {
	Appendix interface{}
	Fixmsg   *quickfix.Message
}

// NewMsgPkg new
func NewMsgPkg() *MsgPkg {
	return &MsgPkg{}
}
