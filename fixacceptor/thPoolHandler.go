package fixacceptor

import "fmt"

type threadPoolHandlerSt struct {
    acceptor *Acceptor
}

func newThreadPoolHandlerSt(ac *Acceptor) *threadPoolHandlerSt {
    return &threadPoolHandlerSt{acceptor: ac}
}

func (h *threadPoolHandlerSt) OnError(msg string) {
    const ftag = "threadPoolHandlerSt::OnError()"

    h.acceptor.appHdl.OnErrorStr(ftag + " " + msg)
}

func (h *threadPoolHandlerSt) OnEvent(msg string) {
    const ftag = "threadPoolHandlerSt::OnEvent()"

    fmt.Println(ftag, msg)
}
