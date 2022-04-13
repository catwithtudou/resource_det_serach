package utils

import (
	"io/ioutil"
	"mime/multipart"
)

func MultipartFileHeaderToBytes(fileData *multipart.FileHeader) ([]byte, error) {
	src, err := fileData.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return data, nil
}
