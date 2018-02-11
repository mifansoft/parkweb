package crypto_test

import (
	"fmt"
	"mifanpark/utilities/crypto"
	"testing"
)

func TestSha1(t *testing.T) {
	fmt.Println(crypto.Sha1("hello world", "abc"))
	newSha1 := crypto.NewSHA1("mifanpark")
	fmt.Println(newSha1.Sha1("zhanwei", "wei"))
	fmt.Println(newSha1.Sha1("hello world", "abc"))
}

func TestNewSHA1(t *testing.T) {
	fmt.Println(crypto.Sha1("hello world", "abc"))

	newSha1 := crypto.NewSHA1("huang")
	fmt.Println(newSha1.Sha1("zhanwei", "wei"))
	fmt.Println(newSha1.Sha1("hello world", "abc"))
}
