package utils

import (
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"strings"
)

var (
	client  *resty.Client
	ocrLink string
)

func NewOcrClient(link string) {
	client = resty.New()
	ocrLink = link
}

func PostReqOcr(imgBase64 string) (string, error) {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"key":   []string{"image"},
			"value": []string{imgBase64},
		}).
		Post(ocrLink)
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
