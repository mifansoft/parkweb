package models

import (
	"mifanpark/entity"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
)

type ThemeResModel struct {
}

func (this *ThemeResModel) GetThemeRes(themeId string) ([]entity.SysThemeRes, error) {
	var rst []entity.SysThemeRes
	err := orm.Dbobj.Table("sys_theme_res").Joins("inner join sys_theme on sys_theme_res.theme_id = sys_theme.id").Where("sys_theme_res.theme_id = ?", themeId).
		Order("sys_theme_res.group_id,sys_theme_res.sort_id asc").Find(&rst).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return rst, err
}
