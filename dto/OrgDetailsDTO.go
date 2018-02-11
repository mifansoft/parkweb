package dto

type OrgDetailsDto struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	ParentName string `json:"parentName"`
	CreateTime string `json:"createTime"`
	CreateUser string `json:"createUser"`
	UpdateTime string `json:"updateTime"`
	UpdateUser string `json:"updateUser"`
}
