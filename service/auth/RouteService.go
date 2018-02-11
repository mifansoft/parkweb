package auth

import "mifanpark/models"

var RouteService = &RouteServiceImpl{}

type RouteServiceImpl struct {
	roleResModel models.RoleResModel
}

func (this *RouteServiceImpl) CheckUrlAuth(userId string, url, method string) bool {
	return this.roleResModel.CheckUrlAuth(userId, url, method)
}
