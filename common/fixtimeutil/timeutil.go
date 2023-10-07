package fixtimeutil

import (
    "fmt"
    "time"
)

/*
GetFIXUTCTimestamp UTC Time Stamp
format: yyyymmdd-hh:mm:ss.sss
*/
func GetFIXUTCTimestamp() string {

    now := time.Now()
    year, month, day := now.UTC().Date()

    hour, minute, second := now.UTC().Clock()

    millisecond := now.UTC().UnixMilli() % 1000

    var sUTCTimestamp string = fmt.Sprintf("%04d%02d%02d-%02d:%02d:%02d.%03d",
        year, month, day, hour, minute, second, millisecond)

    return sUTCTimestamp
}
