package filter

import (
	"mifanpark/service/auth"
	"mifanpark/utilities/route"
	"net/http"
	"strings"
)

type AuthFilter struct {
}

func (this *AuthFilter) staticFile(url string) bool {
	if strings.HasPrefix(url, "/static") {
		return true
	}
	return url == "/favicon.ico"
}

func (this *AuthFilter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if this.staticFile(r.URL.Path) || auth.Identify(w, r) {
		nw := route.NewResponse(w)
		next(nw, r)
	}
}

func init() {
	// 设置白名单，免认证请求
	auth.AddConnUrl("/")
	auth.AddConnUrl("/auth/login")
}
