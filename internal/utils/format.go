package utils

import (
	"encoding/json"
	"strconv"
	"strings"
)

func JsonToString(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

func StringToInt64(str string) int64 {
	result, _ := strconv.Atoi(str)
	return int64(result)
}

// StrToUintSlice example "1,2,3" to [1,2,3]
func StrToUintSlice(str string) []uint {
	slice := strings.Split(str, ",")
	result := make([]uint, len(slice))
	for i := 0; i < len(slice); i++ {
		v, _ := strconv.Atoi(slice[i])
		result[i] = uint(v)
	}
	return result
}
