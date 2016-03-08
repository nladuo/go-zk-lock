//some algorism
package modules

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	err_num       = -1
	ErrConvertMsg = "Number conversion err"
)

func getSerialNumber(path, prefix string) int {

	numStr := strings.TrimPrefix(path, prefix)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return err_num
	}
	return num
}

func GetMinSerialNumber(children []string, prefix string) int {
	index := 0
	min := 999999999
	for i := 0; i < len(children); i++ {
		num := getSerialNumber(children[i], prefix)
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

func GetPathListSerialNumberLessThanLocker(children []string, basePath, prefix, lockername string) []string {
	resultPathList := []string{}
	resultPathList = append(resultPathList, strings.TrimPrefix(lockername, basePath+"/"))
	lockerSeq := getSerialNumber(lockername, basePath+"/"+prefix)
	if lockerSeq == err_num {
		return resultPathList
	}
	for i := 0; i < len(children); i++ {
		seq := getSerialNumber(children[i], prefix)
		if seq == err_num {
			continue
		}
		if seq < lockerSeq {
			resultPathList = append(resultPathList, basePath+"/"+children[i])
		}
	}
	return resultPathList

}

func GetLastNodeName(lockerName, basePath, prefix string) string {
	path := basePath + "/" + prefix
	num := getSerialNumber(lockerName, path)
	if num == err_num {
		panic(errors.New(ErrConvertMsg))
	}
	lastNumStr := fmt.Sprintf("%010d", num-1)
	return prefix + lastNumStr
}
