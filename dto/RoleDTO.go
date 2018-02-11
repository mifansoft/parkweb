package dto

type RoleDto struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	CreateUser string `json:"createUser"`
	CreateTime string `json:"createTime"`
	UpdateUser string `json:"updateUser"`
	UpdateTime string `json:"updateTime"`
}
