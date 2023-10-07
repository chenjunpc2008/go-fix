package fixndp

import (
    "bytes"
    "testing"

    "github.com/stretchr/testify/assert"
)

//go test -bench=. -run=none
func Benchmark_Depack(b *testing.B) {
    for n := 0; n < b.N; n++ {
        fixMsg := "8=FIX.4.29=6135=A34=149=TW52=20220110-08:39:59.49956=ISLD98=0108=3010=226"

        rawMsg := bytes.NewBufferString(fixMsg)

        _, _, ok, _ := Depack(rawMsg.Bytes())
        assert.Equal(b, true, ok)
    }
}
