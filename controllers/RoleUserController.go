package controllers

import (
	"encoding/json"
	"errors"
	"mifanpark/dto"
	"mifanpark/models"
	"mifanpark/service/auth"
	"mifanpark/utilities/groupcache"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/i18n"
	"mifanpark/utilities/jwt"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/validator"
	"net/http"
)

type roleUserModel struct {
	roleUserModel *models.RoleUserModel
	roleModel     *models.RoleModel
}

var RoleUserCtl = &roleUserModel{
	roleUserModel: new(models.RoleUserModel),
	roleModel:     new(models.RoleModel),
}

func init() {
	groupcache.RegisterStaticFile("MiFanParkRoleUserPage", "./views/auth/role_user_page.tpl")
}

func (this *roleUserModel) RoleUserPage(w http.ResponseWriter, r *http.Request) {
	rst, err := groupcache.GetStaticFile("MiFanParkRoleUserPage")
	if err != nil {
		logger.Error(err)
		hret.Error(w, 404, i18n.Get(r, "page_not_exist"))
		return
	}
	hz, err := auth.ParseText(r, string(rst))
	if err != nil {
		logger.Error(err)
		hret.Error(w, 404, i18n.Get(r, "as_of_date_page_not_exist"))
		return
	}
	hz.Execute(w, nil)
}

func (this *roleUserModel) RoleUserRelationPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	roleId := r.FormValue("roleId")

	if validator.IsEmpty(roleId) {
		hret.Error(w, 421, i18n.Get(r, "role_id_is_empty"), errors.New("role_id_is_empty"))
		return
	}

	rst, err := this.roleModel.GetRoleById(roleId)
	if err != nil {
		hret.Error(w, 419, i18n.Get(r, "error_role_resource_query"))
		return
	}
	file, _ := auth.ParseFile(r, "./views/auth/role_user_relation.tpl")
	file.Execute(w, rst)
}

func (this *roleUserModel) GetRoleUserRelation(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roleId := r.FormValue("roleId")
	if validator.IsEmpty(roleId) {
		hret.Error(w, 421, i18n.Get(r, "role_id_is_empty"), errors.New("role_id_is_empty"))
		return
	}
	rst, err := this.roleUserModel.GetRoleUserRelationByRoleId(roleId)
	if err != nil {
		hret.Error(w, 423, err.Error())
		return
	}
	hret.Json(w, rst)
}

func (this *roleUserModel) GetRoleUserRelationOther(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roleId := r.FormValue("roleId")
	if validator.IsEmpty(roleId) {
		hret.Error(w, 421, i18n.Get(r, "role_id_is_empty"), errors.New("role_id_is_empty"))
		return
	}
	rst, err := this.roleUserModel.GetRoleUserRelationOtherByRoleId(roleId)
	if err != nil {
		hret.Error(w, 423, err.Error())
		return
	}
	hret.Json(w, rst)
}

func (this *roleUserModel) RelationUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var rst []dto.RoleUserDto
	err := json.Unmarshal([]byte(r.FormValue("dataJson")), &rst)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, i18n.Get(r, "error_unmarsh_json"), err)
		return
	}
	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}
	err = this.roleUserModel.RelationUser(rst, jclaim.LoginUser.Id)
	if err != nil {
		hret.Error(w, 419, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *roleUserModel) UnRelationUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var rst []dto.RoleUserDto
	err := json.Unmarshal([]byte(r.FormValue("dataJson")), &rst)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, i18n.Get(r, "error_unmarsh_json"), err)
		return
	}
	err = this.roleUserModel.UnRelationUser(rst)
	if err != nil {
		hret.Error(w, 419, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}
