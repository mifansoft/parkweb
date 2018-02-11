package models

import (
	"mifanpark/entity"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
)

type HandleLogModel struct {
}

func (this *HandleLogModel) SaveLog(log_buf []entity.SysHandleLog) {
	tx := orm.Dbobj.Begin()
	if tx.Error != nil {
		logger.Error(tx.Error)
	}
	for _, val := range log_buf {
		if err := orm.Dbobj.Create(&val).Error; err != nil {
			logger.Error(err)
			tx.Rollback()
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		logger.Error(tx.Error)
	}
}
