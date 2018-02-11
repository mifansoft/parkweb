package controllers

import (
	"encoding/json"
	"errors"
	"mifanpark/dto"
	"mifanpark/entity"
	"mifanpark/models"
	"mifanpark/service/auth"
	"mifanpark/utilities/common"
	"mifanpark/utilities/crypto/aes"
	"mifanpark/utilities/groupcache"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/i18n"
	"mifanpark/utilities/jwt"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/uuid"
	"mifanpark/utilities/validator"
	"net/http"
	"strings"
	"time"
)

type userModel struct {
	model models.UserModel
}

var UserCtl = &userModel{}

func init() {
	groupcache.RegisterStaticFile("MiFanParkUserPage", "./views/auth/user_page.tpl")
}

func (this *userModel) UserPage(w http.ResponseWriter, r *http.Request) {
	rst, err := groupcache.GetStaticFile("MiFanParkUserPage")
	if err != nil {
		logger.Error(err)
		hret.Error(w, 404, i18n.PageNotFound(r))
		return
	}
	hz, err := auth.ParseText(r, string(rst))
	if err != nil {
		logger.Error(err)
		hret.Error(w, 404, i18n.PageNotFound(r))
		return
	}
	hz.Execute(w, nil)
}

func (this *userModel) GetUserAll(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rst, err := this.model.GetUserAllDto()
	if err != nil {
		logger.Error(err)
		hret.Error(w, 410, i18n.Get(r, "error_user_query"), err)
		return
	}
	hret.Json(w, rst)
}

func (this *userModel) AddUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := new(entity.SysUser)
	user.Account = r.FormValue("account")
	pwd := r.FormValue("password")
	confirmPwd := r.FormValue("confirmPassword")
	user.Email = r.FormValue("email")
	user.Name = r.FormValue("name")
	user.Phone = r.FormValue("phone")
	user.StatusId = r.FormValue("status")
	user.OrgId = r.FormValue("orgId")

	if validator.IsEmpty(user.Account) {
		hret.Error(w, 421, i18n.Get(r, "account_is_empty"), errors.New("account_is_empty"))
		return
	}
	if validator.IsEmpty(pwd) {
		hret.Error(w, 421, i18n.Get(r, "password_is_empty"), errors.New("password_is_empty"))
		return
	}
	if validator.IsEmpty(confirmPwd) {
		hret.Error(w, 421, i18n.Get(r, "confirm_password_is_empty"), errors.New("confirm_password_is_empty"))
		return
	}
	if pwd != confirmPwd {
		hret.Error(w, 421, i18n.Get(r, "error_password_confirm_failed"), errors.New("error_password_confirm_failed"))
		return
	}
	if len(strings.TrimSpace(confirmPwd)) < 6 || len(strings.TrimSpace(confirmPwd)) > 30 {
		hret.Error(w, 421, i18n.Get(r, "error_password_short"), errors.New("error_password_short"))
		return
	}
	if validator.IsEmpty(user.Phone) {
		hret.Error(w, 421, i18n.Get(r, "phone_is_empty"), errors.New("phone_is_empty"))
		return
	}
	if !validator.IsMobilePhone(user.Phone) {
		hret.Error(w, 421, i18n.Get(r, "phone_number_error"), errors.New("phone_number_error"))
		return
	}
	if !validator.IsEmpty(user.Email) {
		if !validator.IsEmail(user.Email) {
			hret.Error(w, 421, i18n.Get(r, "email_error"), errors.New("email_error"))
			return
		}
	}
	encryptPwd, err := aes.Encrypt(confirmPwd)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}
	user.Password = encryptPwd
	if validator.IsEmpty(user.Name) {
		hret.Error(w, 421, i18n.Get(r, "name_is_empty"), errors.New("name_is_empty"))
		return
	}
	if validator.IsEmpty(user.StatusId) {
		hret.Error(w, 421, i18n.Get(r, "status_is_empty"), errors.New("status_is_empty"))
		return
	}
	if validator.IsEmpty(user.OrgId) {
		hret.Error(w, 421, i18n.Get(r, "org_is_empty"), errors.New("org_is_empty"))
		return
	}
	user.Id = uuid.Random()
	user.CreateTime = time.Now()
	err = this.model.AddUser(*user)
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *userModel) UpdateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	account := r.FormValue("account")
	name := r.FormValue("name")
	orgId := r.FormValue("orgId")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	if validator.IsEmpty(account) {
		hret.Error(w, 421, i18n.Get(r, "account_is_empty"), errors.New("account_is_empty"))
		return
	}

	if validator.IsEmpty(orgId) {
		hret.Error(w, 421, i18n.Get(r, "org_is_empty"), errors.New("org_is_empty"))
		return
	}

	if validator.IsEmpty(name) {
		hret.Error(w, 421, i18n.Get(r, "name_is_empty"), errors.New("name_is_empty"))
		return
	}

	if validator.IsEmpty(phone) {
		hret.Error(w, 421, i18n.Get(r, "phone_is_empty"), errors.New("name_is_empty"))
		return
	}
	if !validator.IsMobilePhone(phone) {
		hret.Error(w, 421, i18n.Get(r, "phone_number_error"), errors.New("phone_number_error"))
		return
	}
	if !validator.IsEmpty(email) {
		if !validator.IsEmail(email) {
			hret.Error(w, 421, i18n.Get(r, "email_error"), errors.New("email_error"))
			return
		}
	}
	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}
	user := new(entity.SysUser)
	user.Account = account
	user.OrgId = orgId
	user.Name = name
	user.Email = email
	user.Phone = phone
	user.UpdateTime = time.Now()
	user.UpdateUserId = jclaim.LoginUser.Id
	err = this.model.UpdateUser(*user)
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *userModel) DeleteUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var rst []dto.UserDto
	err := json.Unmarshal([]byte(r.FormValue("dataJson")), &rst)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, i18n.Get(r, "error_user_json"))
		return
	}
	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}

	for _, val := range rst {
		if common.IsAdmin(val.Account) {
			hret.Error(w, 403, i18n.Get(r, "error_user_forbid_delete_admin"))
			return
		}
		if val.Account == jclaim.LoginUser.Account {
			hret.Error(w, 403, i18n.Get(r, "error_user_forbid_yourself"))
			return
		}
	}
	err = this.model.DeleteUser(rst)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 419, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *userModel) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	account := r.FormValue("account")
	newPwd := r.FormValue("newPassword")
	confirmPwd := r.FormValue("confirmPassword")

	if validator.IsEmpty(account) {
		hret.Error(w, 421, i18n.Get(r, "account_is_empty"), errors.New("account_is_empty"))
		return
	}

	if validator.IsEmpty(newPwd) {
		hret.Error(w, 421, i18n.Get(r, "new_password_is_empty"), errors.New("new_password_is_empty"))
		return
	}

	if validator.IsEmpty(confirmPwd) {
		hret.Error(w, 421, i18n.Get(r, "confirm_password_is_empty"), errors.New("confirm_password_is_empty"))
		return
	}

	if len(strings.TrimSpace(confirmPwd)) < 6 || len(strings.TrimSpace(confirmPwd)) > 30 {
		hret.Error(w, 421, i18n.Get(r, "error_password_short"), errors.New("error_password_short"))
		return
	}

	if newPwd != confirmPwd {
		hret.Error(w, 421, i18n.Get(r, "error_password_confirm_failed"), errors.New("error_password_confirm_failed"))
		return
	}

	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}
	err = this.model.UpdatePassword(account, confirmPwd, jclaim.LoginUser.Id)
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *userModel) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	account := r.FormValue("account")
	status := r.FormValue("status")

	if validator.IsEmpty(account) {
		hret.Error(w, 421, i18n.Get(r, "account_is_empty"), errors.New("account_is_empty"))
		return
	}

	if validator.IsEmpty(status) {
		hret.Error(w, 421, i18n.Get(r, "status_is_empty"), errors.New("status_is_empty"))
		return
	}

	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}

	err = this.model.UpdateStatus(account, status, jclaim.LoginUser.Id)
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}

	hret.Success(w, i18n.Success(r))
}

func (this *userModel) Search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	orgId := r.FormValue("orgId")
	statusId := r.FormValue("status")
	rst, err := this.model.Search(orgId, statusId)
	if err != nil {
		hret.Error(w, 419, i18n.Get(r, "error_user_query"), err)
		return
	}
	hret.Json(w, rst)
}
