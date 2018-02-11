package entity

type SysIndexPage struct {
	Id      string `gorm:"primary_key" json:"id"`
	ThemeId string `gorm:"column:theme_id" json:"themeId"`
	ResUrl  string `gorm:"column:res_url" json:"resUrl"`
}

func (SysIndexPage) TableName() string {
	return "sys_index_page"
}
