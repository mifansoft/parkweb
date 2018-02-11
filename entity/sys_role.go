package entity

import (
	"time"
)

type SysRole struct {
	Id           string    `gorm:"primary_key" json:"id"`
	Name         string    `gorm:"column:name" json:"name"`
	StatusId     string    `gorm:"column:status_id" json:"statusId"`
	CreateUserId string    `gorm:"column:create_user_id" json:"createUserId"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateUserId string    `gorm:"column:update_user_id" json:"updateUserId"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (SysRole) TableName() string {
	return "sys_role"
}
