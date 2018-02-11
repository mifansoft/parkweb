package common_test

import (
	"fmt"
	"mifanpark/utilities/common"
	"testing"
)

func TestDateFormat(t *testing.T) {
	fmt.Println(common.DateFormat("2017-12-06", "YYYY-MM-DD"))
	fmt.Println(common.DateFormat("2017-12-06 11:03:21", "YYYY-MM-DD HH:MM:SS"))
	fmt.Println(common.DateFormat("2017-12-06 21:03:21", "YYYY-MM-DD HH24:MM:SS"))
}
