package service

import (
	"mifanpark/controllers"
	"mifanpark/utilities/route"
)

func Register() {
	//菜单管理
	route.Handler("GET", "/", controllers.IndexPage)
	route.Handler("POST", "/auth/login", controllers.LoginCtl.LoginSystem)
	route.Handler("GET", "/auth/home", controllers.LoginCtl.HomePage)
	route.Handler("GET", "/auth/menu/index", controllers.MenuCtl.IndexMenu)
	route.Handler("GET", "/auth/menu/entry", controllers.MenuCtl.MenuEntry)

	//组织机构管理
	route.Handler("GET", "/auth/org/page", controllers.OrgCtl.OrgPage)
	route.Handler("GET", "/auth/org/getall", controllers.OrgCtl.GetOrgAll)
	route.Handler("GET", "/auth/org/details", controllers.OrgCtl.GetDetails)
	route.Handler("POST", "/auth/org/addorg", controllers.OrgCtl.AddOrg)
	route.Handler("POST", "/auth/org/deleteorg", controllers.OrgCtl.DeleteOrg)
	route.Handler("PUT", "/auth/org/updateorg", controllers.OrgCtl.UpdateOrg)
	route.Handler("GET", "/auth/org/exportorg", controllers.OrgCtl.ExportOrg)
	route.Handler("POST", "/auth/org/importorg", controllers.OrgCtl.ImportOrg)

	//用户管理
	route.Handler("GET", "/auth/user/page", controllers.UserCtl.UserPage)
	route.Handler("GET", "/auth/user/getall", controllers.UserCtl.GetUserAll)
	route.Handler("POST", "/auth/user/adduser", controllers.UserCtl.AddUser)
	route.Handler("PUT", "/auth/user/updateuser", controllers.UserCtl.UpdateUser)
	route.Handler("POST", "/auth/user/deleteuser", controllers.UserCtl.DeleteUser)
	route.Handler("PUT", "/auth/user/modifystatus", controllers.UserCtl.UpdateStatus)
	route.Handler("PUT", "/auth/user/modifypassword", controllers.UserCtl.UpdatePassword)
	route.Handler("GET", "/auth/user/search", controllers.UserCtl.Search)

	//角色管理
	route.Handler("GET", "/auth/role/page", controllers.RoleCtl.RolePage)
	route.Handler("GET", "/auth/role/getall", controllers.RoleCtl.GetRoleAll)
	route.Handler("POST", "/auth/role/addrole", controllers.RoleCtl.AddRole)
	route.Handler("PUT", "/auth/role/updaterole", controllers.RoleCtl.UpdateRole)
	route.Handler("POST", "/auth/role/deleterole", controllers.RoleCtl.DeleteRole)
	route.Handler("GET", "/auth/role/details", controllers.RoleCtl.GetDetails)
	route.Handler("GET", "/auth/role/getroleother", controllers.RoleCtl.GetRoleOther)

	//角色资源
	route.Handler("GET", "/auth/role/res/getroleres", controllers.RoleResCtl.GetRoleRes)
	route.Handler("POST", "/auth/role/res/authorized", controllers.RoleResCtl.Authorized)
	route.Handler("POST", "/auth/role/res/unauthorized", controllers.RoleResCtl.UnAuthorized)

	//用户角色管理
	route.Handler("GET", "/auth/role/user/page", controllers.RoleUserCtl.RoleUserPage)
	route.Handler("GET", "/auth/role/user/relationpage", controllers.RoleUserCtl.RoleUserRelationPage)
	route.Handler("GET", "/auth/role/user/getrelation", controllers.RoleUserCtl.GetRoleUserRelation)
	route.Handler("GET", "/auth/role/user/getrelationother", controllers.RoleUserCtl.GetRoleUserRelationOther)
	route.Handler("POST", "/auth/role/user/relationuser", controllers.RoleUserCtl.RelationUser)
	route.Handler("POST", "/auth/role/user/unrelationuser", controllers.RoleUserCtl.UnRelationUser)
}
