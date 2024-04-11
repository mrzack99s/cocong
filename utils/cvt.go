package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func StringToInt64(s string) (v int64, err error) {
	iVal, err := strconv.Atoi(s)
	if err != nil {
		return
	}

	v = int64(iVal)
	return
}

func StringToUInt64(s string) (v uint64, err error) {
	iVal, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return
	}

	v = iVal
	return
}

func Int64ToString(s int64) (v string) {
	v = fmt.Sprintf("%d", s)
	return
}

func InterfaceToString(v interface{}) (s string, err error) {
	bV, err := json.Marshal(v)
	if err != nil {
		return
	}

	s = string(bV)
	return
}
