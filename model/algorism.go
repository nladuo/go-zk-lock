//some algorism
package model

import (
	"log"
	"strconv"
	"strings"
)

const (
	sequence_bit = 10
)

func GetMinIndex(strs []string, prefix string) int {
	index := 0
	min := 999999999
	for i := 0; i < len(strs); i++ {
		strVal := strings.TrimPrefix(strs[i], prefix)
		num, err := strconv.Atoi(strVal)
		if err != nil {
			log.Println("fifo.getMinIndex() , conversion err", err.Error())
			panic(err)
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
	numStr := strings.TrimPrefix(lockerName, path)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		log.Println("fifo.getLastNodeName() , conversion err", err.Error())
		panic(err)
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
