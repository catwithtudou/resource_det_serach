package utils

import (
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"strings"
)

var client *resty.Client

func init() {
	client = resty.New()
}

func PostReqOcr(imgBase64 string) (string, error) {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"key":   []string{"image"},
			"value": []string{imgBase64},
		}).
		Post("http://127.0.0.1:9299/ocr/prediction")
	if err != nil {
		return "", err
	}

	respValue := gjson.Get(resp.String(), "value").String()
	respValue = strings.TrimLeft(respValue, "[")
	respValue = strings.TrimRight(respValue, "]")
	respValue = strings.ReplaceAll(respValue, "'", "")
	respValue = strings.ReplaceAll(respValue, ",", " ")

	return DelPunctuation(respValue), nil
}
