package entity

import (
	"time"
)

type SysUser struct {
	Id            string    `gorm:"primary_key" json:"id"`
	Account       string    `gorm:"column:account" json:"account"`
	Password      string    `gorm:"column:password" json:"-"`
	Name          string    `gorm:"column:name" json:"name"`
	Phone         string    `gorm:"column:phone" json:"phone"`
	Email         string    `gorm:"column:email" json:"email"`
	StatusId      string    `gorm:"status_id" json:"statusId"`
	OrgId         string    `gorm:"column:org_id" json:"orgId"`
	LoginIp       string    `gorm:"column:login_ip" json:"loginIp"`
	LastLoginTime string    `gorm:"column:last_login_time" json:"lastLoginTime"`
	ErrorCount    int       `gorm:"column:error_count" json:"errorCount"`
	CreateTime    time.Time `gorm:"column:create_time" json:"createTime"`
	UpdateUserId  string    `gorm:"column:update_user_id" json:"updateUserId"`
	UpdateTime    time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
