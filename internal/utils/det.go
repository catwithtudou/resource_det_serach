package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/presentation"
	"github.com/unidoc/unioffice/spreadsheet"
	"image/jpeg"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

func NewUnidoc(key string) {
	err := license.SetMeteredKey(key)
	if err != nil {
		panic(err)
	}
}

func DetDocxByUnidoc(fileBytes []byte) (string, error) {
	tFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("resource_det_search_docx_%d-*.docx", time.Now().UnixNano()))
	if err != nil {
		return "", err
	}
	defer os.Remove(tFile.Name())

	_, err = tFile.Write(fileBytes)
	if err != nil {
		return "", err
	}

	doc, err := document.Open(tFile.Name())
	if err != nil {
		return "", err
	}
	defer doc.Close()

	extracted := doc.ExtractText()

	return DelPunctuation(extracted.Text()), nil
}

func DetPptxByUnidoc(fileBytes []byte) (string, error) {
	tFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("resource_det_search_pptx_%d-*.pptx", time.Now().UnixNano()))
	if err != nil {
		return "", err
	}
	defer os.Remove(tFile.Name())

	_, err = tFile.Write(fileBytes)
	if err != nil {
		return "", err
	}

	ppt, err := presentation.Open(tFile.Name())
	if err != nil {
		return "", err
	}
	defer ppt.Close()

	extracted := ppt.ExtractText()

	return DelPunctuation(extracted.Text()), nil
}

func DetXlsxByUnidoc(fileBytes []byte) (string, error) {
	tFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("resource_det_search_xlsx_%d-*.xlsx", time.Now().UnixNano()))
	if err != nil {
		return "", err
	}
	defer os.Remove(tFile.Name())

	_, err = tFile.Write(fileBytes)
	if err != nil {
		return "", err
	}

	xlsx, err := spreadsheet.Open(tFile.Name())
	if err != nil {
		return "", err
	}
	defer xlsx.Close()

	extracted := xlsx.ExtractText()

	return DelPunctuation(extracted.Text()), nil
}

func DetMd(fileBytes []byte) (string, error) {
	unsafe := blackfriday.MarkdownCommon(fileBytes)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return DelMdTags(string(html)), nil
}

func DetPdf(fileBytes []byte) (string, error) {
	tFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("resource_det_search_pdf_%d-*.xlsx", time.Now().UnixNano()))
	if err != nil {
		return "", err
	}
	defer os.Remove(tFile.Name())

	_, err = tFile.Write(fileBytes)
	if err != nil {
		return "", err
	}

	pdf, err := fitz.New(tFile.Name())
	if err != nil {
		return "", err
	}

	result := ""

	// TODO:??????PDF??????????????????????????????
	for n := 0; n < 3; n++ {
		img, err := pdf.Image(n)
		if err != nil {
			panic(err)
		}

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			panic(err)
		}

		imageBase64Val := base64.StdEncoding.EncodeToString(buf.Bytes())
		val, err := PostReqOcr(imageBase64Val)
		if err != nil {
			return "", err
		}

		result += val
	}

	return result, nil
}

func DetImg(fileBytes []byte) (string, error) {
	imageBase64Val := base64.StdEncoding.EncodeToString(fileBytes)
	val, err := PostReqOcr(imageBase64Val)
	if err != nil {
		return "", err
	}
	return val, nil
}

func DelPunctuation(p string) string {
	reg := regexp.MustCompile("[^0-9A-Za-z\u4e00-\u9fa5]")
	result := reg.ReplaceAllString(p, " ")
	sReg := regexp.MustCompile("\\s+")
	result = sReg.ReplaceAllString(result, " ")
	return result
}

func DelMdTags(p string) string {
	reg := regexp.MustCompile("<[^>]*>")
	result := reg.ReplaceAllString(p, "")
	nReg := regexp.MustCompile("\\n")
	result = nReg.ReplaceAllString(result, " ")
	return DelPunctuation(result)
}
