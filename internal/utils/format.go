package utils

import (
	"encoding/json"
	"strconv"
)

func JsonToString(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

func StringToInt64(str string) int64 {
	result, _ := strconv.Atoi(str)
	return int64(result)
}
