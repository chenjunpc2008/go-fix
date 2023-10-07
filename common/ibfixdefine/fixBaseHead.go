package ibfixdefine

// FixInfoLocal fix base values for local use
type FixInfoLocal struct {
    BeginString  string // FIX.4.1, FIX.4.2 or FIX.4.3 are supported in this tag
    MsgType      string
    SenderCompID string // By default, this is the IBKR username.
    TargetCompID string // Default is “IB.”
    MsgSeqNum    uint64
}
