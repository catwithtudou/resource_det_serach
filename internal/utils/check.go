package utils

import (
	"encoding/json"
	"github.com/dlclark/regexp2"
	"regexp"
	"resource_det_search/internal/constants"
	"strings"
)

func CheckEmail(email string) bool {
	reg := regexp.MustCompile("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$")
	return reg.MatchString(email)
}

func CheckRole(str string) bool {
	roles := []string{"学生", "教师"}
	return Contains(roles, str)
}

func CheckSex(str string) bool {
	roles := []string{"男", "女"}
	return Contains(roles, str)
}

func CheckPswd(str string) bool {
	//密码大于6位小于20位
	//必须包含字母和数字
	reg := regexp2.MustCompile("^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,20}$", 0)
	result, _ := reg.MatchString(str)
	return result
}

func CheckUserType(str string) bool {
	types := []string{"tag", "category"}
	// part 类型只能官方手动创建
	return Contains(types, str)
}

func CheckAllType(str string) bool {
	types := []string{"tag", "category", "part"}
	return Contains(types, str)
}

func CheckDocTypeStr(str string) ([]uint, bool) {
	if str == "" {
		return nil, true
	}

	var result []uint
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		return nil, false
	}
	return result, true
}

func CheckDocFileType(fileName string) (string, bool) {
	slices := strings.Split(fileName, ".")
	if len(slices) < 2 {
		return "", false
	}

	fileType := slices[len(slices)-1]

	// 直接识别：doc/docx、ppt/pptx、md、txt
	// OCR识别：jpg/jpeg、png、pdf
	if !DetOcrTypesContains(fileType) && !DetByteTypesContains(fileType) {
		return "", false
	}

	return fileType, true
}

func CheckDocFileSize(fileSize int64) bool {
	return fileSize < 1024*1024*20 && fileSize > 0
}

func DetByteTypesContains(fileType string) bool {
	fType := constants.DetByteType(fileType)
	for _, e := range constants.DetByteTypes {
		if fType == e {
			return true
		}
	}

	return false
}

func DetOcrTypesContains(fileType string) bool {
	fType := constants.DetOcrType(fileType)
	for _, e := range constants.DetOcrTypes {
		if fType == e {
			return true
		}
	}

	return false
}

func Contains(elems []string, elem string) bool {
	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}

func ContainsUint(elems []uint, elem uint) bool {
	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}
