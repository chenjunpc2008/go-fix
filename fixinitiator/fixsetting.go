package fixinitiator

const (
    // DefaultSendBuffSize default send buff size
    DefaultSendBuffSize = 1 * 10000
)

// Settings fix setting
type Settings struct {
    FixDataDictPath string

    RemoteIP   string
    RemotePort uint16

    // tcp send buff size
    SendBuffsize int
    // receive callback if async or sync
    AsyncReceive bool

    BeginString  string
    SenderCompID string
    TargetCompID string

    // configuration
    // ResetOnLogon Determines if sequence numbers should be reset when receiving a logon request.
    ResetOnLogon bool
}

func DefaultConfig() Settings {
    return Settings{
        SendBuffsize: DefaultSendBuffSize,
    }
}
