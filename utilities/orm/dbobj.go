package orm

import (
	"mifanpark/utilities/orm/dbhandle"

	"github.com/jinzhu/gorm"
)

var (
	Dbobj gorm.DB
)

func InitDB(dbDriver string) error {
	if Dbobj.DB() == nil {
		if val, ok := dbhandle.Adapter[dbDriver]; ok {
			Dbobj = val()
		}
	}
	return nil
}
