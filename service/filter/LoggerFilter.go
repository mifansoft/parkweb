package filter

import (
	"mifanpark/service/auth"
	"net/http"
)

type LoggerFilter struct {
}

func (this *LoggerFilter) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
	go auth.WriteHandleLogs(w, r)
}

func init() {
	// 开启操作日志监听
	go auth.LogSync()
}
