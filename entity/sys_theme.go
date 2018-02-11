package entity

type SysTheme struct {
	Id   string `gorm:"primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (SysTheme) TableName() string {
	return "sys_theme"
}
