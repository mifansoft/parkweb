package entity

import (
	"time"
)

type SysRoleRes struct {
	Id           string    `gorm:"primary_key" json:"id"`
	RoleId       string    `gorm:"column:role_id" json:"roleId"`
	ResId        string    `gorm:"column:res_id" json:"resId"`
	CreateUserId string    `gorm:"column:create_user_id" json:"createUserId"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateUserId string    `gorm:"column:update_user_id" json:"updateUserId"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (SysRoleRes) TableName() string {
	return "sys_role_res"
}
