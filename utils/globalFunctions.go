package utils

import (
	"time"
)

/*func ByteToUint32(arr []byte) uint32 {
	//not optimized (i will not use it)
	return binary.LittleEndian.Uint32(arr)
}*/

const (
	epochDiff = 116444736000000000
)

func IntFTToDate(decimalFT int64) time.Time {
	if decimalFT == 0 {
		return time.Time{}
	}
	return IntUnixToDate((decimalFT - epochDiff) * 100)
}

func IntFTToUnix(decimalFT uint64) int64 {
	return (int64(decimalFT) - epochDiff) * 100
}

func IntFTToStringDate(decimalFT int64) string {
	return IntFTToDate(decimalFT).Format("2006-01-02 15:04:05")
}

func IntUnixToDate(decimalU int64) time.Time {
	if decimalU <= 0 {
		return time.Time{}
	}
	return time.Unix(0, decimalU)
}

func IntUnixToStringDate(UnixN int64) string {
	return IntUnixToDate(UnixN).Format("2006-01-02 15:04:05")
}
