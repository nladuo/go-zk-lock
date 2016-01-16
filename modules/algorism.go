//some algorism
package modules

import (
	"errors"
	"strconv"
	"strings"
)

const (
	sequence_bit  = 10
	err_num       = -1
	ErrConvertMsg = "Number conversion err"
)

func getSequentialNumber(str, prefix string) int {

	numStr := strings.TrimPrefix(str, prefix)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return err_num
	}
	return num
}

func GetMinIndex(strs []string, prefix string) int {
	index := 0
	min := 999999999
	for i := 0; i < len(strs); i++ {
		num := getSequentialNumber(strs[i], prefix)
		if num == err_num {
			continue
		}

		if num < min {
			min = num
			index = i
		}

	}
	return index
}

func GetStrsSequenceLessThanLocker(strs []string, basePath, prefix, lockername string) []string {
	resultStrs := []string{}
	resultStrs = append(resultStrs, strings.TrimPrefix(lockername, basePath+"/"))
	lockerSeq := getSequentialNumber(lockername, basePath+"/"+prefix)
	if lockerSeq == err_num {
		return resultStrs
	}
	for i := 0; i < len(strs); i++ {
		seq := getSequentialNumber(strs[i], prefix)
		if seq == err_num {
			continue
		}
		if seq < lockerSeq {
			resultStrs = append(resultStrs, strs[i])
		}
	}
	return resultStrs

}

func GetLastNodeName(lockerName, basePath, prefix string) string {
	path := basePath + "/" + prefix
	num := getSequentialNumber(lockerName, path)
	if num == err_num {
		panic(errors.New(ErrConvertMsg))
	}
	lastNumStr := strconv.Itoa(num - 1)
	numBit := 1
	for i := num; i > 0; i /= 10 {
		numBit++
	}
	for i := 0; i <= sequence_bit-numBit; i++ {
		lastNumStr = "0" + lastNumStr
	}
	return prefix + lastNumStr
}
