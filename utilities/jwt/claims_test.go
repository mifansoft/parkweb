package jwt_test

import (
	"mifanpark/utilities/jwt"
	"testing"

	"fmt"
)

func TestToken(t *testing.T) {
	token, err := jwt.GenToken(jwt.NewUserData().SetUserId("huang").SetOrgunitId("100"))
	fmt.Println(token, err)
	claims, err := jwt.ParseToken(token)
	fmt.Println(claims, err)

	flag := jwt.ValidToken(token)
	fmt.Println("true or false:", flag)

	handle := jwt.NewHandle(nil)
	flag = handle.ValidToken(token)
	fmt.Println("true or false:", flag)
	token, err = handle.GenToken(jwt.NewUserData().SetUserId("huang").SetOrgunitId("100"))
	fmt.Println(token, err)
	claims, err = jwt.ParseToken(token)
	fmt.Println(claims, err)

	flag = jwt.ValidToken(token)
	fmt.Println("true or false:", flag)
}
