package models

import (
	"errors"
	"mifanpark/dto"
	"mifanpark/entity"
	"mifanpark/utilities/common"
	"mifanpark/utilities/crypto/aes"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
	"mifanpark/utilities/validator"
	"time"
)

type UserModel struct {
	morg OrgModel
}

func (this *UserModel) GetOrgIdByAccount(account string) (string, error) {
	var orgId string
	db := orm.Dbobj.Select("org_id").Where("account", account).Scan(&orgId)
	if db.Error != nil {
		logger.Error(db.Error)
	}
	return orgId, db.Error
}

func (this *UserModel) GetUserAll() ([]entity.SysUser, error) {
	var rst []entity.SysUser
	db := orm.Dbobj.Find(&rst)
	if db.Error != nil {
		logger.Error(db.Error)
	}
	return rst, db.Error
}

func (this *UserModel) GetUserAllDto() ([]dto.UserDto, error) {
	var ret []dto.UserDto
	rst, err := this.GetUserAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for _, val := range rst {
		var statusName string
		if err := orm.Dbobj.Table("sys_user_status").Select("name").Where("id = ?", val.StatusId).Row().Scan(&statusName); err != nil {
			logger.Error(err)
			return nil, err
		}
		var orgName string
		if err := orm.Dbobj.Table("sys_org").Select("name").Where("id = ?", val.OrgId).Row().Scan(&orgName); err != nil {
			logger.Error(err)
			return nil, err
		}
		var updateUserName string
		if !validator.IsEmpty(val.UpdateUserId) {
			if err := orm.Dbobj.Table("sys_user").Select("name").Where("id = ?", val.UpdateUserId).Row().Scan(&updateUserName); err != nil {
				logger.Error(err)
				return nil, err
			}
		}
		tmp := new(dto.UserDto)
		tmp.UserId = val.Id
		tmp.Account = val.Account
		tmp.Name = val.Name
		tmp.Status = statusName
		tmp.OrgName = orgName
		tmp.Phone = val.Phone
		tmp.Email = val.Email
		if !val.CreateTime.IsZero() {
			tmp.CreateTime = common.DateToString(val.CreateTime)
		}
		tmp.UpdateUser = orgName
		if !val.UpdateTime.IsZero() {
			tmp.UpdateTime = common.DateToString(val.UpdateTime)
		}
		ret = append(ret, *tmp)
	}
	return ret, nil
}

func (this *UserModel) AddUser(user entity.SysUser) error {
	err := orm.Dbobj.Omit("update_user_id", "update_time", "last_login_time").Create(&user).Error
	if err != nil {
		logger.Error(err)
	}
	return nil
}

func (this *UserModel) UpdateUser(user entity.SysUser) error {
	err := orm.Dbobj.Model(&user).Omit("last_login_time").Where("account = ?", user.Account).Updates(&user).Error
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (this *UserModel) DeleteUser(userList []dto.UserDto) error {
	db := orm.Dbobj.Begin()
	if db.Error != nil {
		logger.Error(db.Error)
		return db.Error
	}
	for _, val := range userList {
		if err := db.Where("account = ?", val.Account).Delete(entity.SysUser{}).Error; err != nil {
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

func (this *UserModel) UpdatePassword(account, password, updateUserId string) error {
	encryptPwd, err := aes.Encrypt(password)
	if err != nil {
		return err
	}
	sqlStr := map[string]interface{}{
		"password":       encryptPwd,
		"update_user_id": updateUserId,
		"update_time":    time.Now(),
	}
	err = orm.Dbobj.Model(entity.SysUser{}).Where("account = ?", account).Updates(sqlStr).Error
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (this *UserModel) GetUserByAccount(account string) (entity.SysUser, error) {
	var user entity.SysUser
	if validator.IsEmpty(account) {
		return user, errors.New("account_is_empty")
	}
	err := orm.Dbobj.Where("account = ?", account).First(&user).Error
	if err != nil {
		logger.Error(err)
	}
	return user, err
}

func (this *UserModel) UpdateStatus(account, status, updateUserId string) error {
	sqlStr := map[string]interface{}{
		"status_id":      status,
		"update_user_id": updateUserId,
		"update_time":    time.Now(),
	}
	err := orm.Dbobj.Model(entity.SysUser{}).Where("account = ?", account).Updates(sqlStr).Error
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (this *UserModel) Search(orgId, statusId string) ([]dto.UserDto, error) {
	var userDto []dto.UserDto
	var userList []entity.SysUser
	if !validator.IsEmpty(orgId) && !validator.IsEmpty(statusId) {
		err := orm.Dbobj.Where("org_id = ?", orgId).Where("status_id = ?", statusId).Find(&userList).Error
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	} else if !validator.IsEmpty(orgId) {
		err := orm.Dbobj.Where("org_id = ?", orgId).Find(&userList).Error
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	} else if !validator.IsEmpty(statusId) {
		err := orm.Dbobj.Where("status_id = ?", statusId).Find(&userList).Error
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}
	for _, val := range userList {
		var statusName string
		if err := orm.Dbobj.Table("sys_user_status").Select("name").Where("id = ?", val.StatusId).Row().Scan(&statusName); err != nil {
			logger.Error(err)
			return nil, err
		}
		var orgName string
		if err := orm.Dbobj.Table("sys_org").Select("name").Where("id = ?", val.OrgId).Row().Scan(&orgName); err != nil {
			logger.Error(err)
			return nil, err
		}
		var updateUserName string
		if !validator.IsEmpty(val.UpdateUserId) {
			if err := orm.Dbobj.Table("sys_user").Select("name").Where("id = ?", val.UpdateUserId).Row().Scan(&updateUserName); err != nil {
				logger.Error(err)
				return nil, err
			}
		}
		tmp := new(dto.UserDto)
		tmp.UserId = val.Id
		tmp.Account = val.Account
		tmp.Name = val.Name
		tmp.Status = statusName
		tmp.OrgName = orgName
		tmp.Phone = val.Phone
		tmp.Email = val.Email
		if !val.CreateTime.IsZero() {
			tmp.CreateTime = common.DateToString(val.CreateTime)
		}
		tmp.UpdateUser = orgName
		if !val.UpdateTime.IsZero() {
			tmp.UpdateTime = common.DateToString(val.UpdateTime)
		}
		userDto = append(userDto, *tmp)
	}
	return userDto, nil
}
