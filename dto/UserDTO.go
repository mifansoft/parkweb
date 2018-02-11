package dto

type UserDto struct {
	UserId     string `json:"userId"`
	Account    string `json:"account"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	OrgName    string `json:"orgName"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	CreateTime string `json:"createTime"`
	UpdateUser string `json:"updateUser"`
	UpdateTime string `json:"updateTime"`
}
