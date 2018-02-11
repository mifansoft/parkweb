package hret_test

import (
	"fmt"
	"mifanpark/utilities/hret"
	"testing"
)

func TestHttpPanic(t *testing.T) {
	defer hret.RecvPanic(func() {
		fmt.Println("捕获异常")
	})
	panic("抛出异常")
}
