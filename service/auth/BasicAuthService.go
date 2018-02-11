package auth

import (
	"mifanpark/utilities/common"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/i18n"
	"mifanpark/utilities/jwt"
	"mifanpark/utilities/logger"
	"net/http"
	"sync"
	"text/template"
)

type BasicAuthFilter struct {
	// filterAuthUrl 存放不需要校验用户权限的路由信息
	filterAuthUrl map[string]bool
	alock         *sync.RWMutex

	// filterConnUrl 存放不需要校验连接的路由信息
	filterConnUrl map[string]bool
	clock         *sync.RWMutex
}

var BAFilter = NewBasicAuthFilter()

func NewBasicAuthFilter() *BasicAuthFilter {
	return &BasicAuthFilter{
		filterAuthUrl: make(map[string]bool),
		alock:         new(sync.RWMutex),
		filterConnUrl: make(map[string]bool),
		clock:         new(sync.RWMutex),
	}
}

func Identify(w http.ResponseWriter, r *http.Request) bool {
	return BAFilter.Identify(w, r)
}

// 校验用户权限信息
func (this *BasicAuthFilter) Identify(w http.ResponseWriter, r *http.Request) bool {
	// 校验连接信息
	this.clock.RLock()
	if _, yes := this.filterConnUrl[r.URL.Path]; !yes {
		if !jwt.ValidHttp(r) {
			this.clock.RUnlock()
			hz, _ := template.ParseFiles("./views/auth/disconnect.tpl")
			hz.Execute(w, nil)
			return false
		}
	} else {
		this.clock.RUnlock()
		return true
	}
	this.clock.RUnlock()

	// 校验授权信息
	this.alock.RLock()
	defer this.alock.RUnlock()
	if _, yes := this.filterAuthUrl[r.URL.Path]; !yes {
		flag := this.basicAuth(r)
		if !flag {
			logger.Info("您没有被授权访问：", r.URL)
			hret.Error(w, 403, i18n.NoAuth(r))
			return false
		}
	}
	return true
}

func (this *BasicAuthFilter) basicAuth(r *http.Request) bool {
	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		return false
	}
	if common.IsAdmin(jclaim.LoginUser.Id) {
		return true
	}
	method := r.Method
	if method == http.MethodPost && r.FormValue("_method") == http.MethodDelete {
		method = http.MethodDelete
	}
	logger.Debug("basicAuth,method is:", method, ",path is:", r.URL.Path, ",user is:", jclaim.LoginUser.Account)
	return RouteService.CheckUrlAuth(jclaim.LoginUser.Account, r.URL.Path, method)
}

func AddConnUrl(url string) {
	BAFilter.AddConnUrl(url)
}

func (this *BasicAuthFilter) AddConnUrl(url string) {
	this.clock.Lock()
	defer this.clock.Unlock()
	this.filterConnUrl[url] = true
}

func BasicAuth(r *http.Request) bool {
	return BAFilter.basicAuth(r)
}
