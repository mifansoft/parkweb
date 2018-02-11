package dbhandle

import (
	"mifanpark/utilities/config"
	"mifanpark/utilities/logger"
	"sync"

	"github.com/jinzhu/gorm"
)

type instance func() gorm.DB

var (
	dbLock  = new(sync.RWMutex)
	Adapter = make(map[string]instance)
)

func GetConfig() (config.Handle, error) {
	return config.Load("conf/app.conf", config.INI)
}

func Register(dsn string, f instance) {
	dbLock.Lock()
	defer dbLock.Unlock()
	if f == nil {
		logger.Error("sql:Register driver is nil")
	}
	if _, dup := Adapter[dsn]; dup {
		logger.Error("reregister driver. dsn is: ", dsn)
	}
	Adapter[dsn] = f
}
