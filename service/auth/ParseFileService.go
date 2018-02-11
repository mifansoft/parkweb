package auth

import (
	"html/template"
	"io/ioutil"
	"mifanpark/models"
	"mifanpark/utilities/common"
	"mifanpark/utilities/jwt"
	"net/http"
)

var roleResModel = models.RoleResModel{}

func ParseText(r *http.Request, content string) (*template.Template, error) {
	claim, err := jwt.ParseHttp(r)
	if err != nil {
		return nil, err
	}
	return template.New("template").Funcs(template.FuncMap{"checkResIDAuth": func(args ...string) bool {
		if len(args) < 2 {
			return false
		}
		if common.IsAdmin(claim.LoginUser.Id) {
			return true
		}
		if args[0] == "2" {
			return roleResModel.CheckResIDAuth(claim.LoginUser.Id, args[1])
		} else {
			return false
		}
	}}).Parse(content)
}

func ParseFile(r *http.Request, filePath string) (*template.Template, error) {
	rst, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return ParseText(r, string(rst))
}
