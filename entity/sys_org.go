package entity

import (
	"time"
)

type SysOrg struct {
	Id           string    `gorm:"primary_key" json:"id"`
	Code         string    `gorm:"column:code" json:"code"`
	Name         string    `gorm:"column:name" json:"name"`
	ParentId     string    `gorm:"column:parent_id" json:"parentId"`
	CreateUserId string    `gorm:"column:create_user_id" json:"createUserId"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateUserId string    `gorm:"column:update_user_id" json:"updateUserId"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (SysOrg) TableName() string {
	return "sys_org"
}
