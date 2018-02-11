package jwt

import (
	"mifanpark/entity"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// 自定义Claims类
type customClaims struct {
	*jwt.StandardClaims
	UserData *UserData
}

func newClaims(conf *Config, userData *UserData) *customClaims {
	c := &customClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + conf.duration,
			Issuer:    conf.owner,
		},
		UserData: userData,
	}
	return c
}

type UserData struct {
	// // 用户账号
	// UserId string
	// // 用户机构号
	// OrgId string
	//用户
	LoginUser entity.SysUser
	// 用户角色
	Authorities string `json:"authorities"`
}

func NewUserData() *UserData {
	return &UserData{
		Authorities: "",
	}
}

func (r *UserData) SetLoginUser(user entity.SysUser) *UserData {
	r.LoginUser = user
	return r
}

func (r *UserData) SetAuthorities(role string) *UserData {
	r.Authorities = role
	return r
}
