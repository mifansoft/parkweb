package auth

import (
	"encoding/json"
	"mifanpark/entity"
	"mifanpark/models"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/jwt"
	"mifanpark/utilities/route"
	"mifanpark/utilities/uuid"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var log_buf = make(chan entity.SysHandleLog, 40960)
var handleModel = &models.HandleLogModel{}

func LogSync() {
	var buf []entity.SysHandleLog
	for {
		select {
		case <-time.After(time.Second * 10):
			// sync handle logs to database per 5 second.
			if len(buf) == 0 {
				continue
			}
			go saveLog(buf)
			buf = make([]entity.SysHandleLog, 0)
		case val, ok := <-log_buf:
			if ok {
				buf = append(buf, val)
				if len(buf) > 1000 {
					go saveLog(buf)
					buf = make([]entity.SysHandleLog, 0)
				}
			}
		}
	}
}

func saveLog(log_buf []entity.SysHandleLog) {
	handleModel.SaveLog(log_buf)
}

func WriteHandleLogs(w http.ResponseWriter, r *http.Request) {
	defer hret.RecvPanic()
	if nw, ok := w.(*route.Response); ok {
		var handleLog entity.SysHandleLog
		handleLog.ClientIp = route.RequestIP(r)
		handleLog.Content = formencode(r.Form)
		handleLog.Id = uuid.Random()
		handleLog.Method = r.Method
		handleLog.Status = strconv.Itoa(nw.Status)
		handleLog.Url = r.URL.Path
		handleLog.HandleTime = time.Now()
		jclaim, err := jwt.ParseHttp(r)
		if err != nil {
			handleLog.UserId = handleLog.ClientIp
		} else {
			handleLog.UserId = jclaim.LoginUser.Id
		}
		log_buf <- handleLog
	}
}

func formencode(form url.Values) string {
	rst := make(map[string]string)
	for key, val := range form {
		if key == "_" {
			continue
		}
		rst[key] = val[0]
	}
	str, _ := json.Marshal(rst)
	return string(str)
}
