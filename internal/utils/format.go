package utils

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
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

func TimestampFormat(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
