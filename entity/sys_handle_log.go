package entity

import (
	"time"
)

type SysHandleLog struct {
	Id         string    `gorm:"primary_key" json:"id"`
	ClientIp   string    `gorm:"column:client_ip" json:"clientIp"`
	Status     string    `gorm:"column:status" json:"status"`
	Method     string    `gorm:"column:method" json:"method"`
	Url        string    `gorm:"column:url" json:"url"`
	Content    string    `gorm:"column:content" json:"content"`
	UserId     string    `gorm:"column:user_id" json:"userId"`
	HandleTime time.Time `gorm:"column:handle_time" json:"handleTime"`
}

func (SysHandleLog) TableName() string {
	return "sys_handle_log"
}
