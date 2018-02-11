package dto

type RoleUserDto struct {
	UserId   string `json:"userId"`
	Account  string `json:"account"`
	Name     string `json:"name"`
	RoleId   string `json:"roleId"`
	OrgCode  string `json:"orgCode"`
	OrgName  string `json:"orgName"`
	AuthUser string `json:"authUser"`
	AuthTime string `json:"authTime"`
}
