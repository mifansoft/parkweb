package auth

import "mifanpark/models"

var OrgService = &OrgServiceImpl{}

type OrgServiceImpl struct {
	user models.UserModel
}
