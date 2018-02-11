package entity

import (
	"time"
)

type SysThemeRes struct {
	Id           string    `gorm:"primary_key" json:"id"`
	ResId        string    `gorm:"column:res_id" json:"resId"`
	ThemeId      string    `gorm:"column:theme_id" json:"themeId"`
	ResUrl       string    `gorm:"column:res_url" json:"resUrl"`
	OpenTypeId   string    `gorm:"column:res_open_type_id" json:"resOpenTypeId"`
	BgColor      string    `gorm:"column:bg_color" json:"bgColor"`
	ResClass     string    `gorm:"column:res_class" json:"resClass"`
	GroupId      string    `gorm:"column:group_id" json:"groupId"`
	ResImg       string    `gorm:"column:res_img" json:"resImg"`
	SortId       string    `gorm:"column:sort_id" json:"sortId"`
	NewIFrame    string    `gorm:"column:new_iframe" json:"newIFrame"`
	CreateUserId string    `gorm:"column:create_user_id" json:"createUserId"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateUserId string    `gorm:"column:update_user_id" json:"updateUserId"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (SysThemeRes) TableName() string {
	return "sys_theme_res"
}
