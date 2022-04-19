package utils

import (
	"fmt"
	"testing"
)

func TestStrToUintSlice(t *testing.T) {
	res := StrToUintSlice("1,2,3")
	if len(res) != 3 {
		t.Fatal("wrong")
	}
	fmt.Println(res)
}
