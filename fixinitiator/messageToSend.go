package fixinitiator

// message to send type
type msgToSendType int

const (
	mtsTypeNormal msgToSendType = 1 // normal message
	mtsTypeResend msgToSendType = 2 // resend message
)

type msgToSendSt struct {
	pkg     *MsgPkg
	msgType msgToSendType
}
