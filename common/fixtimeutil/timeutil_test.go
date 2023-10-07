package fixtimeutil_test

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"

    "github.com/chenjunpc2008/go-fix/common/fixtimeutil"
)

func Test_GetFIXUTCTimestamp(t *testing.T) {
    t1 := fixtimeutil.GetFIXUTCTimestamp()

    time.Sleep(1 * time.Millisecond)

    t2 := fixtimeutil.GetFIXUTCTimestamp()

    assert.Equal(t, t1[0:18], t2[0:18])

    // now := time.Now()
    // fmt.Println(now)
    // fmt.Println(now.Unix())
    // fmt.Println(now.UnixMilli())
    // fmt.Println(now.UnixNano())
    // fmt.Println(time.Unix(now.Unix(), 0))
    // fmt.Println(time.Unix(0, now.UnixNano()))
}
