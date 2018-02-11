package controllers

import (
	"encoding/json"
	"mifanpark/models"
	"mifanpark/service/auth"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/i18n"
	"mifanpark/utilities/jwt"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/validator"
	"net/http"
)

type roleResModel struct {
	roleModel    *models.RoleModel
	roleResModel *models.RoleResModel
}

var RoleResCtl = &roleResModel{
	roleModel:    new(models.RoleModel),
	roleResModel: new(models.RoleResModel),
}

func (this *roleResModel) GetRoleRes(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !auth.BasicAuth(r) {
		hret.Error(w, 403, i18n.NoAuth(r))
		return
	}
	roleId := r.FormValue("roleId")
	typeId := r.FormValue("typeId")

	if typeId == "0" {
		rst, err := this.roleResModel.GetRoleResByRoleId(roleId)
		if err != nil {
			hret.Error(w, 419, i18n.Get(r, "error_role_get_resource"))
			return
		}
		hret.Json(w, rst)
	} else if typeId == "1" {
		rst, err := this.roleResModel.UnGetRoleResByRoleId(roleId)
		if err != nil {
			hret.Error(w, 419, i18n.Get(r, "error_role_get_resource"))
			return
		}
		hret.Json(w, rst)
	}
}

func (this *roleResModel) Authorized(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var resList []string
	err := json.Unmarshal([]byte(r.FormValue("dataJson")), &resList)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 422, err.Error())
		return
	}
	roleId := r.FormValue("roleId")
	if len(resList) == 0 {
		hret.Error(w, 421, i18n.Get(r, "error_resource_res_id"))
		return
	}

	if validator.IsEmpty(roleId) {
		hret.Error(w, 421, i18n.Get(r, "error_role_id_format"))
		return
	}

	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}

	err = this.roleResModel.Authorized(roleId, jclaim.LoginUser.Id, resList)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 419, i18n.Get(r, "error_role_add_resource_failed"))
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *roleResModel) UnAuthorized(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var resList []string
	err := json.Unmarshal([]byte(r.FormValue("dataJson")), &resList)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 422, err.Error())
		return
	}
	roleId := r.FormValue("roleId")

	if len(resList) == 0 {
		hret.Error(w, 421, i18n.Get(r, "error_resource_res_id"))
		return
	}

	if validator.IsEmpty(roleId) {
		hret.Error(w, 421, i18n.Get(r, "error_role_id_format"))
		return
	}

	err = this.roleResModel.UnAuthorized(roleId, resList)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 419, i18n.Get(r, "error_role_add_resource_failed"))
		return
	}
	hret.Success(w, i18n.Success(r))
}
