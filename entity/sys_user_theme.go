package entity

type SysUserTheme struct {
	Id      string `gorm:"primary_key" json:"id"`
	UserId  string `gorm:"column:user_id" json:"userId"`
	ThemeId string `gorm:"column:theme_id" json:"themeId"`
}

func (SysUserTheme) TableName() string {
	return "sys_user_theme"
}
