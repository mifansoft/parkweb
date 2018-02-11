package models

import (
	"mifanpark/dto"
	"mifanpark/entity"
	"mifanpark/utilities/common"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
	"mifanpark/utilities/validator"
)

type RoleModel struct {
}

func (this *RoleModel) GetRoleAll() ([]dto.RoleDto, error) {
	var roleDto []dto.RoleDto
	var roleList []entity.SysRole
	err := orm.Dbobj.Find(&roleList).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for _, val := range roleList {
		dt := new(dto.RoleDto)
		if !validator.IsEmpty(val.StatusId) {
			if err := orm.Dbobj.Table("sys_role_status").Select("name").Where("id = ?", val.StatusId).Row().Scan(&dt.Status); err != nil {
				logger.Error(err)
				return nil, err
			}
		}
		if !validator.IsEmpty(val.CreateUserId) {
			if err := orm.Dbobj.Table("sys_user").Select("name").Where("id = ?", val.CreateUserId).Row().Scan(&dt.CreateUser); err != nil {
				logger.Error(err)
				return nil, err
			}
		}
		if !validator.IsEmpty(val.UpdateUserId) {
			if err := orm.Dbobj.Table("sys_user").Select("name").Where("id = ?", val.UpdateUserId).Row().Scan(&dt.UpdateUser); err != nil {
				logger.Error(err)
				return nil, err
			}
		}
		if !val.CreateTime.IsZero() {
			dt.CreateTime = common.DateToString(val.CreateTime)
		}
		if !val.UpdateTime.IsZero() {
			dt.UpdateTime = common.DateToString(val.UpdateTime)
		}
		dt.Id = val.Id
		dt.Name = val.Name
		roleDto = append(roleDto, *dt)
	}
	return roleDto, nil
}

func (this *RoleModel) AddRole(role entity.SysRole) error {
	err := orm.Dbobj.Omit("update_user_id", "update_time").Create(&role).Error
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (this *RoleModel) UpdateRole(role entity.SysRole) error {
	err := orm.Dbobj.Model(&role).Where("id = ?", role.Id).Updates(&role).Error
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (this *RoleModel) DeleteRole(roleList []dto.RoleDto) error {
	db := orm.Dbobj.Begin()
	if db.Error != nil {
		logger.Error(db.Error)
		return db.Error
	}
	for _, val := range roleList {
		if err := db.Where("id = ?", val.Id).Delete(entity.SysRole{}).Error; err != nil {
			logger.Error(err)
			db.Rollback()
			return err
		}
	}
	if err := db.Commit().Error; err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (this *RoleModel) GetRoleById(roleId string) (entity.SysRole, error) {
	var rst entity.SysRole
	err := orm.Dbobj.Where("id = ?", roleId).First(&rst).Error
	if err != nil {
		logger.Error(err)
	}
	return rst, err
}

func (this *RoleModel) GetRoleOther(userId string) ([]entity.SysRole, error) {
	var role []entity.SysRole
	err := orm.Dbobj.Where("not exists(select 1 from sys_role_user where sys_role_user.user_id = ? and sys_role.id = sys_role_user.role_id)", userId).
		Find(&role).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return role, nil
}
