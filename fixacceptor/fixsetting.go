package fixacceptor

const (
    // DefaultSendBuffSize default send buff size
    DefaultSendBuffSize = 1 * 10000
)

// Settings fix setting
type Settings struct {
    FixDataDictPath string

    Port uint16

    // tcp send buff size
    SendBuffsize int
    // receive callback if async or sync
    AsyncReceive bool

    // sended callback if async or sync
    AsyncSended bool

    BeginString  string
    SenderCompID string

    // ResetOnLogon Determines if sequence numbers should be reset when recieving a logon request.
    ResetOnLogon bool
}

func DefaultConfig() Settings {
    return Settings{
        SendBuffsize: DefaultSendBuffSize,
        AsyncReceive: true,
        AsyncSended:  true,
    }
}
