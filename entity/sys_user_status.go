package entity

type SysUserStatus struct {
	Id   string `gorm:"primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (SysUserStatus) TableName() string {
	return "sys_user_status"
}
