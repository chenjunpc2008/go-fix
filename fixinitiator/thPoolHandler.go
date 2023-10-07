package fixinitiator

import "fmt"

type threadPoolHandlerSt struct {
	initiator *Initiator
}

func newThreadPoolHandlerSt(intor *Initiator) *threadPoolHandlerSt {
	return &threadPoolHandlerSt{initiator: intor}
}

func (h *threadPoolHandlerSt) OnError(msg string) {
	const ftag = "threadPoolHandlerSt::OnError()"

	h.initiator.appHdl.OnErrorStr(ftag + " " + msg)
}

func (h *threadPoolHandlerSt) OnEvent(msg string) {
	const ftag = "threadPoolHandlerSt::OnEvent()"

	fmt.Println(ftag, msg)
}
