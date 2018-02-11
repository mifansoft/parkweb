package controllers

import (
	"encoding/json"
	"html/template"
	"mifanpark/dto"
	"mifanpark/entity"
	"mifanpark/models"
	"mifanpark/utilities/crypto/aes"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/i18n"
	"mifanpark/utilities/jwt"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/route"
	"mifanpark/utilities/validator"
	"net/http"
	"strconv"
)

type loginModel struct {
	model *models.LoginModel
}

var LoginCtl = &loginModel{
	model: new(models.LoginModel),
}

/*
登录系统
*/
func (this *loginModel) LoginSystem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form
	account := form.Get("account")
	password := form.Get("password")
	duration := form.Get("duration")

	if validator.IsEmpty(account) {
		rdto := dto.AuthDto{
			Account: account,
			Code:    422,
			Msg:     "鉴权失败，账号为空",
		}
		result(w, rdto)
		return
	}

	if validator.IsEmpty(password) {
		rdto := dto.AuthDto{
			Account: account,
			Code:    423,
			Msg:     "鉴权失败，密码为空",
		}
		result(w, rdto)
		return
	}
	user, qdto := this.checkUserInfo(dto.AuthDto{
		Account:  account,
		Password: password,
		Duration: duration,
		LoginIp:  route.RequestIP(r),
		Code:     404,
		Msg:      "",
	})

	if qdto.Success && qdto.Code == 200 {
		if validator.IsEmpty(user.OrgId) {
			logger.Error(account, "用户没有指定组织机构")
			qdto.Code = 427
			qdto.Msg = "获取用户所在组织机构失败"
			result(w, qdto)
			return
		}
		et, err := strconv.ParseInt(duration, 10, 64)
		if err != nil {
			et = 17280
		}

		token, _ := jwt.GenToken(jwt.NewUserData().SetLoginUser(user))
		cookie := http.Cookie{Name: "Authorization", Value: token, Path: "/", MaxAge: int(et)}
		http.SetCookie(w, &cookie)
		qdto.Token = token
		result(w, qdto)
		return
	}
	result(w, qdto)
	return
}

/*
检查用户信息
*/
func (this *loginModel) checkUserInfo(qdto dto.AuthDto) (entity.SysUser, dto.AuthDto) {
	var user entity.SysUser
	//加密用户信息
	pd, err := aes.Encrypt(qdto.Password)
	if err != nil {
		logger.Error(err)
		qdto.Code = 434
		qdto.Msg = "加密用户信息失败"
		return user, qdto
	}
	qdto.Password = pd
	return this.model.Login(qdto)
}

/*
页面输出
*/
func result(response http.ResponseWriter, rdto dto.AuthDto) {
	msg, err := json.Marshal(rdto)
	if err != nil {
		response.WriteHeader(http.StatusExpectationFailed)
		response.Write([]byte(`{username:` + rdto.Account + `,Code:"431",Msg:"format json type info failed."}`))
		return
	}
	response.WriteHeader(rdto.Code)
	response.Header().Set("Authorization", rdto.Token)
	response.Write(msg)
}

/*
后台首页
*/
func (this *loginModel) HomePage(w http.ResponseWriter, r *http.Request) {
	defer hret.RecvPanic(func() {
		logger.Error("Get Home Page Failure.")
		http.Redirect(w, r, "/", 302)
	})

	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		http.Redirect(w, r, "/", 302)
		return
	}

	url := this.model.GetDefaultTheme(jclaim.LoginUser.Id)
	h, err := template.ParseFiles(url)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, i18n.Get(r, "error_get_login_page"), err)
		return
	}
	h.Execute(w, jclaim.LoginUser.Id)
}
