package orm

import (
	"mifanpark/utilities/orm/dbhandle"
)

func init() {
	conf, err := dbhandle.GetConfig()
	if err != nil {
		panic("init database failed." + err.Error())
	}
	//驱动
	dbDriver, err := conf.Get("DB.Driver")
	if err != nil {
		panic("get database driver failed." + err.Error())
	}
	InitDB(dbDriver)
}
