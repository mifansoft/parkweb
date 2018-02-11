package entity

import (
	"time"
)

type SysRes struct {
	Id           string    `gorm:"primary_key" json:"id"`
	Name         string    `gorm:"column:name" json:"name"`
	ResAttrId    string    `gorm:"column:res_attr_id" json:"resAttrId"`
	ParentId     string    `gorm:"column:parent_id" json:"parentId"`
	ResTypeId    string    `gorm:"column:res_type_id" json:"resTypeId"`
	Flag         int       `gorm:"column:flag" json:"flag"`
	Method       string    `gorm:"column:method" json:"method"`
	CreateUserId string    `gorm:"column:create_user_id" json:"createUserId"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateUserId string    `gorm:"column:update_user_id" json:"updateUserId"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (SysRes) TableName() string {
	return "sys_res"
}
