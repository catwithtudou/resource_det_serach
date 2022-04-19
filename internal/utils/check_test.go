package utils

import (
	"fmt"
	"testing"
)

func TestCheckPswd(t *testing.T) {
	a := "aaaaaaaaaaaaaaaaaaa"
	b := "aaaaaaaa949812478"
	c := "2131232143241321"
	d := "1111"
	e := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	fmt.Println(CheckPswd(a))
	fmt.Println(CheckPswd(b))
	fmt.Println(CheckPswd(c))
	fmt.Println(CheckPswd(d))
	fmt.Println(CheckPswd(e))
}

func TestCheckDocTypeStr(t *testing.T) {
	res, ok := CheckDocTypeStr("[1,2,3]")
	if !ok {
		t.Fatal("wrong")
	}
	fmt.Println(res)
	_, ok = CheckDocTypeStr("1,2,3")
	if ok {
		t.Fatal("wrong")
	}
	_, ok = CheckDocTypeStr("1,2,3.2,-1")
	if ok {
		t.Fatal("wrong")
	}
	_, ok = CheckDocTypeStr("[1,2,3.2]")
	if ok {
		t.Fatal("wrong")
	}
	_, ok = CheckDocTypeStr("[1,2,-1]")
	if ok {
		t.Fatal("wrong")
	}
}
