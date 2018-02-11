package models

import (
	"mifanpark/entity"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
)

type UserThemeModel struct {
}

func (this *UserThemeModel) GetUserTheme(userId string) (entity.SysUserTheme, error) {
	var rst entity.SysUserTheme
	err := orm.Dbobj.Table("sys_user_theme").Joins("left join sys_theme on sys_user_theme.theme_id = sys_theme.id").
		Where("sys_user_theme.user_id = ?", userId).Find(&rst).Error
	if err != nil {
		logger.Error(err)
		return rst, err
	}
	return rst, err
}
