package entity

type SysResAttr struct {
	Id   string `gorm:"primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (SysResAttr) TableName() string {
	return "sys_res_attr"
}
