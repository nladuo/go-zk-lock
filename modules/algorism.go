//some algorism
package modules

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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

func GetLastNodeName(lockerName, basePath, prefix string) string {
	path := basePath + "/" + prefix
	num := getSerialNumber(lockerName, path)
	if num == err_num {
		panic(errors.New(ErrConvertMsg))
	}
	lastNumStr := fmt.Sprintf("%010d", num-1)
	return prefix + lastNumStr
}

func CheckOutTimeOut(data []byte, timeout time.Duration) bool {
	nowUnixTime := time.Now().Unix()
	//get the znode create time
	createUnixTime, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return true
	}
	timeoutUnixTime := int64(createUnixTime) + int64(timeout.Seconds())
	return timeoutUnixTime < nowUnixTime
}
