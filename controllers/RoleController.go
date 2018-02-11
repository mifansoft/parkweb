package controllers

import (
	"encoding/json"
	"errors"
	"mifanpark/dto"
	"mifanpark/entity"
	"mifanpark/models"
	"mifanpark/service/auth"
	"mifanpark/utilities/groupcache"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/i18n"
	"mifanpark/utilities/jwt"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/uuid"
	"mifanpark/utilities/validator"
	"net/http"
	"time"
)

type roleModel struct {
	model *models.RoleModel
}

var RoleCtl = &roleModel{
	model: new(models.RoleModel),
}

func init() {
	groupcache.RegisterStaticFile("MiFanParkRolePage", "./views/auth/role_page.tpl")
}

func (this *roleModel) RolePage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !auth.BasicAuth(r) {
		hret.Error(w, 403, i18n.NoAuth(r))
		return
	}
	rst, err := groupcache.GetStaticFile("MiFanParkRolePage")
	if err != nil {
		hret.Error(w, 403, i18n.PageNotFound(r))
		return
	}
	hz, err := auth.ParseText(r, string(rst))
	if err != nil {
		hret.Error(w, 404, i18n.PageNotFound(r))
		return
	}
	hz.Execute(w, nil)
}

func (this *roleModel) GetRoleAll(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rst, err := this.model.GetRoleAll()
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, "error_role_query"), err)
		return
	}
	hret.Json(w, rst)
}

func (this *roleModel) AddRole(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	statusId := r.FormValue("status")

	if validator.IsEmpty(name) {
		hret.Error(w, 421, i18n.Get(r, "name_is_empty"), errors.New("name_is_empty"))
		return
	}

	if validator.IsEmpty(statusId) {
		hret.Error(w, 421, i18n.Get(r, "name_is_empty"), errors.New("name_is_empty"))
		return
	}
	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}
	role := new(entity.SysRole)
	role.Id = uuid.Random()
	role.Name = name
	role.StatusId = statusId
	role.CreateTime = time.Now()
	role.CreateUserId = jclaim.LoginUser.Id
	err = this.model.AddRole(*role)
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *roleModel) UpdateRole(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	name := r.FormValue("name")
	statusId := r.FormValue("status")

	if validator.IsEmpty(id) {
		hret.Error(w, 421, i18n.Get(r, "id_is_empty"), errors.New("id_is_empty"))
		return
	}

	if validator.IsEmpty(name) {
		hret.Error(w, 421, i18n.Get(r, "name_is_empty"), errors.New("name_is_empty"))
		return
	}

	if validator.IsEmpty(statusId) {
		hret.Error(w, 421, i18n.Get(r, "name_is_empty"), errors.New("name_is_empty"))
		return
	}

	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}

	role := new(entity.SysRole)
	role.Id = id
	role.Name = name
	role.StatusId = statusId
	role.UpdateTime = time.Now()
	role.UpdateUserId = jclaim.LoginUser.Id
	err = this.model.UpdateRole(*role)
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *roleModel) DeleteRole(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var roleDto []dto.RoleDto
	err := json.Unmarshal([]byte(r.FormValue("dataJson")), &roleDto)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, i18n.Get(r, "error_role_json"))
		return
	}

	for _, val := range roleDto {
		if val.Id == "0" {
			hret.Error(w, 403, i18n.Get(r, "error_role_forbid_delete_admin"))
			return
		}
	}

	err = this.model.DeleteRole(roleDto)
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *roleModel) GetDetails(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	roleId := r.FormValue("id")

	if validator.IsEmpty(roleId) {
		hret.Error(w, 421, i18n.Get(r, "role_id_is_empty"), errors.New("role_id_is_empty"))
		return
	}

	rst, err := this.model.GetRoleById(roleId)
	if err != nil {
		hret.Error(w, 419, i18n.Get(r, "error_role_resource_query"))
		return
	}

	file, err := auth.ParseFile(r, "./views/auth/role_res_page.tpl")
	if err != nil {
		logger.Error(err)
		hret.Error(w, 423, i18n.Get(r, "error_role_resource_query"))
		return
	}
	file.Execute(w, rst)
}

func (this *roleModel) GetRoleOther(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := r.FormValue("userId")
	if validator.IsEmpty(userId) {
		hret.Error(w, 419, i18n.Get(r, "user_id_is_empty"))
		return
	}

	rst, err := this.model.GetRoleOther(userId)
	if err != nil {
		hret.Error(w, 419, i18n.Get(r, "error_user_role_un_auth"), err)
		return
	}
	hret.Json(w, rst)
}
