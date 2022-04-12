package utils

import "encoding/json"

func JsonToString(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}
