package entity

type SysResType struct {
	Id   string `gorm:"primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (SysResType) TableName() string {
	return "sys_res_type"
}
