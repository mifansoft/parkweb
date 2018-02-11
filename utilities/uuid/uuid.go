package uuid

import (
	"mifanpark/utilities/logger"
	"strings"

	"github.com/satori/go.uuid"
)

// 采用随机数的方式生成uuid
func Random() string {
	uid, err := uuid.NewV4()
	if err != nil {
		logger.Error(err)
	}
	return strings.Replace(uid.String(), "-", "", -1)
}

// 采用随机数和sha1组合的方式生成uuid
func UUID() string {
	uid1, err := uuid.NewV1()
	if err != nil {
		logger.Error(err)
	}
	uid2, err := uuid.NewV4()
	if err != nil {
		logger.Error(err)
	}
	return strings.Replace(uuid.NewV5(uid2, uid1.String()).String(), "-", "", -1)
}
