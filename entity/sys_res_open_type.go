package entity

type SysResOpenType struct {
	Id   string `gorm:"primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (SysResOpenType) TableName() string {
	return "sys_res_open_type"
}
