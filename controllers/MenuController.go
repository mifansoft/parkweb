package controllers

import (
	"mifanpark/models"
	"mifanpark/utilities/crypto"
	"mifanpark/utilities/groupcache"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/i18n"
	"mifanpark/utilities/jwt"
	"mifanpark/utilities/logger"
	"net/http"
)

type menuModel struct {
	model *models.MenuModel
}

var MenuCtl = &menuModel{
	model: new(models.MenuModel),
}

func (this *menuModel) IndexMenu(w http.ResponseWriter, r *http.Request) {
	defer hret.RecvPanic()
	r.ParseForm()
	form := r.Form
	resId := form.Get("res_id")
	typeId := form.Get("type_id")

	claim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}

	ojs, err := this.model.GetHomePageMenus(resId, typeId, claim.LoginUser.Id)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, i18n.Get(r, "error_query_menu"))
		return
	}
	w.Write(ojs)
}

func (this *menuModel) MenuEntry(w http.ResponseWriter, r *http.Request) {
	defer hret.RecvPanic()
	r.ParseForm()
	themeResId := r.FormValue("theme_id")

	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}
	url, err := this.model.GetResUrl(jclaim.LoginUser.Id, themeResId)
	if err != nil {
		logger.Error(err)
		w.Write([]byte(url))
		return
	}

	key := crypto.Sha1(themeResId, jclaim.LoginUser.Id, url)
	if !groupcache.FileIsExist(key) {
		groupcache.RegisterStaticFile(key, url)
	}

	tpl, err := groupcache.GetStaticFile(key)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 404, i18n.PageNotFound(r))
		return
	}
	w.Write(tpl)
}
