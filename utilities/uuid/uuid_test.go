package uuid_test

import (
	"fmt"
	"testing"

	"mifanpark/utilities/uuid"
)

func TestUuid(t *testing.T) {
	fmt.Println(uuid.Random())
	fmt.Println(uuid.Random())
	fmt.Println(uuid.Random())
	fmt.Println(uuid.Random())

	fmt.Println(uuid.UUID())
	fmt.Println(uuid.UUID())
	fmt.Println(uuid.UUID())
}
