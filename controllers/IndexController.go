package controllers

import (
	"mifanpark/utilities/groupcache"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/i18n"
	"mifanpark/utilities/logger"
	"net/http"
)

func init() {
	groupcache.RegisterStaticFile("MiFanParkIndexPage", "./views/auth/login.tpl")
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	rst, err := groupcache.GetStaticFile("MiFanParkIndexPage")
	if err != nil {
		logger.Error(err)
		hret.Error(w, 404, i18n.PageNotFound(r))
		return
	}
	w.Write(rst)
}
