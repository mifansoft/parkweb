package models

import (
	"mifanpark/dto"
	"mifanpark/entity"
	"mifanpark/utilities/common"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
	"mifanpark/utilities/uuid"
	"mifanpark/utilities/validator"
	"time"
)

type RoleUserModel struct {
	userModel UserModel
}

func (this *RoleUserModel) GetRoleByUserId(userId string) ([]entity.SysRole, error) {
	var rst []entity.SysRole
	err := orm.Dbobj.Table("sys_role").Joins("left join sys_role_user on sys_role.id = sys_role_user.role_id").
		Joins("left join sys_role_status on sys_role.status_id = sys_role_status.id").
		Where("sys_user_role.user_id = ?", userId).Where("sys_role_status.id = ?", "0").Find(&rst).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return rst, err
}

func (this *RoleUserModel) GetRoleUserRelationByRoleId(roleId string) ([]dto.RoleUserDto, error) {
	var roleUserDto []dto.RoleUserDto
	var roleUserList []entity.SysRoleUser
	err := orm.Dbobj.Joins("inner join sys_user on sys_user.id = sys_role_user.user_id").
		Joins("inner join sys_org on sys_org.id = sys_user.org_id").Where("role_id = ?", roleId).Find(&roleUserList).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for _, val := range roleUserList {
		dt := new(dto.RoleUserDto)
		var user entity.SysUser
		err := orm.Dbobj.Where("id = ?", val.UserId).First(&user).Error
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		dt.Account = user.Account
		dt.UserId = user.Id
		if !val.CreateTime.IsZero() {
			dt.AuthTime = common.DateToString(val.CreateTime)
		}
		if !validator.IsEmpty(val.CreateUserId) {
			err := orm.Dbobj.Table("sys_user").Select("name").Where("id = ?", val.CreateUserId).Row().Scan(&dt.AuthUser)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
		}
		dt.Name = user.Name
		var org entity.SysOrg
		err = orm.Dbobj.Where("id = ?", user.OrgId).First(&org).Error
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		dt.OrgCode = org.Code
		dt.OrgName = org.Name
		roleUserDto = append(roleUserDto, *dt)
	}
	return roleUserDto, nil
}

func (this *RoleUserModel) GetRoleUserRelationOtherByRoleId(roleId string) ([]dto.RoleUserDto, error) {
	var roleUserDto []dto.RoleUserDto
	var userList []entity.SysUser
	err := orm.Dbobj.Where("not exists(select 1 from sys_role_user where sys_user.id = sys_role_user.user_id and sys_role_user.role_id = ?)", roleId).
		Find(&userList).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for _, val := range userList {
		dt := new(dto.RoleUserDto)
		dt.UserId = val.Id
		dt.Account = val.Account
		dt.Name = val.Name
		var org entity.SysOrg
		err := orm.Dbobj.Where("id = ?", val.OrgId).First(&org).Error
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		dt.OrgCode = org.Code
		dt.OrgName = org.Name
		roleUserDto = append(roleUserDto, *dt)
	}
	return roleUserDto, nil
}

func (this *RoleUserModel) RelationUser(roleUserList []dto.RoleUserDto, userId string) error {
	db := orm.Dbobj.Begin()
	if db.Error != nil {
		logger.Error(db.Error)
		return db.Error
	}

	for _, val := range roleUserList {
		roleUser := new(entity.SysRoleUser)
		roleUser.CreateTime = time.Now()
		roleUser.CreateUserId = userId
		roleUser.Id = uuid.Random()
		roleUser.RoleId = val.RoleId
		roleUser.UserId = val.UserId
		err := db.Omit("update_user_id", "update_time").Create(&roleUser).Error
		if err != nil {
			logger.Error(err)
			db.Rollback()
			return err
		}
	}
	err := db.Commit().Error
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (this *RoleUserModel) UnRelationUser(userRoleList []dto.RoleUserDto) error {
	db := orm.Dbobj.Begin()
	if db.Error != nil {
		logger.Error(db.Error)
		return db.Error
	}
	for _, val := range userRoleList {
		err := db.Where("user_id = ?", val.UserId).Where("role_id = ?", val.RoleId).Delete(entity.SysRoleUser{}).Error
		if err != nil {
			logger.Error(err)
			db.Rollback()
			return err
		}
	}
	err := db.Commit().Error
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
