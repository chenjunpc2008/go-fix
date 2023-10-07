package fixtimeutil

import (
	"math/rand"
	"strconv"
	"time"
)

/*
GetRandom16 生成16位的随机数 10位时间戳+6位随机数

@return string
*/
func GetRandom16() string {
	const (
		MaxRandom = 999999
		Carry     = 1000000
	)

	var (
		timeNow int64
		randnum int64
		iRand16 int64
		sRes    string
	)
	timeNow = time.Now().Unix()
	randnum = rand.Int63n(MaxRandom)
	iRand16 = timeNow*Carry + randnum
	sRes = strconv.FormatInt(iRand16, 10)
	return sRes
}
