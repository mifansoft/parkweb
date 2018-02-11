package common_test

import (
	"fmt"
	"mifanpark/utilities/common"
	"testing"
)

func TestGetKey(t *testing.T) {
	s := common.JoinKey("hello world", "my name is", "c huang zhan wei")
	fmt.Println(s)
	fmt.Println(common.GetKey(s, 4))
	fmt.Println([]byte(s))
	fmt.Println(byte(0x1f))
}
