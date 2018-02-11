package entity

import "time"

type SysRoleUser struct {
	Id           string    `gorm:"primary_key" json:"id"`
	RoleId       string    `gorm:"column:role_id" json:"roleId"`
	UserId       string    `gorm:"column:user_id" json:"userId"`
	CreateUserId string    `gorm:"column:create_user_id" json:"createUserId"`
	CreateTime   time.Time `grom:"column:create_time" json:"createTime"`
	UpdateUserId string    `gorm:"column:update_user_id" json:"updateUserId"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (SysRoleUser) TableName() string {
	return "sys_role_user"
}
