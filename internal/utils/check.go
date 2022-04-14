package utils

import (
	"github.com/dlclark/regexp2"
	"regexp"
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

func CheckType(str string) bool {
	types := []string{"tag", "category"}
	// part 类型只能官方手动创建
	return Contains(types, str)
}

func Contains(elems []string, elem string) bool {
	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}
