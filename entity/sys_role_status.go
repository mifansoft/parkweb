package entity

type SysRoleStatus struct {
	Id   string `gorm:"primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (SysRoleStatus) TableName() string {
	return "sys_role_status"
}
