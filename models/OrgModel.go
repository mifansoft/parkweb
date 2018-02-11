package models

import (
	"errors"
	"mifanpark/dto"
	"mifanpark/entity"
	"mifanpark/utilities/common"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
	"mifanpark/utilities/validator"
)

type OrgModel struct {
}

func (this *OrgModel) GetOrgAll() ([]entity.SysOrg, error) {
	var rst []entity.SysOrg
	if err := orm.Dbobj.Find(&rst).Error; err != nil {
		logger.Error(err)
		return nil, err
	}
	return rst, nil
}

func (this *OrgModel) GetDetails(orgId string) (dto.OrgDetailsDto, error) {
	var currentOrg entity.SysOrg
	err := orm.Dbobj.Where("id = ?", orgId).First(&currentOrg).Error
	if err != nil {
		logger.Error(err)
		return dto.OrgDetailsDto{}, err
	}
	orgDetailsDTO := new(dto.OrgDetailsDto)
	orgDetailsDTO.Code = currentOrg.Code
	orgDetailsDTO.Name = currentOrg.Name
	if !currentOrg.CreateTime.IsZero() {
		orgDetailsDTO.CreateTime = common.DateToString(currentOrg.CreateTime)
	}
	if !currentOrg.UpdateTime.IsZero() {
		orgDetailsDTO.UpdateTime = common.DateToString(currentOrg.UpdateTime)
	}
	if !validator.IsEmpty(currentOrg.CreateUserId) {
		var createUserName string
		if err = orm.Dbobj.Table("sys_user").Select("name").Where("id = ?", currentOrg.CreateUserId).Row().Scan(&createUserName); err != nil {
			logger.Error(err)
			return dto.OrgDetailsDto{}, err
		}
		orgDetailsDTO.CreateUser = createUserName
	}
	if !validator.IsEmpty(currentOrg.UpdateUserId) {
		var updateUserName string
		if err = orm.Dbobj.Table("sys_user").Select("name").Where("id = ?", currentOrg.UpdateUserId).Row().Scan(&updateUserName); err != nil {
			logger.Error(err)
			return dto.OrgDetailsDto{}, err
		}
		orgDetailsDTO.UpdateUser = updateUserName
	}
	if !validator.IsEmpty(currentOrg.ParentId) && currentOrg.ParentId != "system_root" {
		var orgName string
		if err = orm.Dbobj.Table("sys_org").Select("name").Where("id = ?", currentOrg.ParentId).Row().Scan(&orgName); err != nil {
			logger.Error(err)
			return dto.OrgDetailsDto{}, err
		}
		orgDetailsDTO.ParentName = orgName
	} else {
		orgDetailsDTO.ParentName = currentOrg.ParentId
	}
	return *orgDetailsDTO, err
}

func (this *OrgModel) AddOrg(org entity.SysOrg) error {
	err := orm.Dbobj.Omit("update_user_id", "update_time").Create(&org).Error
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (this *OrgModel) DeleteOrg(orgId string) error {
	db := orm.Dbobj.Begin()
	if db.Error != nil {
		logger.Error(db.Error)
		return db.Error
	}
	sublist, err := this.GetSubOrgList(orgId)
	if err != nil {
		return err
	}
	sublist = append(sublist, entity.SysOrg{Id: orgId})
	for _, val := range sublist {
		err = db.Delete(val).Error
		if err != nil {
			logger.Error(err)
			db.Rollback()
			return err
		}
	}
	if err = db.Commit().Error; err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (this *OrgModel) GetSubOrgList(orgId string) ([]entity.SysOrg, error) {
	if validator.IsEmpty(orgId) {
		logger.Error("error_org_id_empty")
		return nil, errors.New("error_org_id_empty")
	}
	var rst []entity.SysOrg
	all, err := this.GetOrgAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	for _, val := range all {
		if val.ParentId == orgId {
			rst = append(rst, val)
			break
		}
	}
	this.dfs(all, orgId, &rst)
	return rst, nil
}

func (this *OrgModel) dfs(node []entity.SysOrg, orgId string, rst *[]entity.SysOrg) {
	for _, val := range node {
		if val.ParentId == orgId {
			*rst = append(*rst, val)
			if val.ParentId == val.Id {
				logger.Error("当前机构与上级机构编码一致，逻辑错误，退出递归")
				return
			}
			this.dfs(node, val.Id, rst)
		}
	}
}

func (this *OrgModel) UpdateOrg(sysOrg entity.SysOrg) error {
	db := orm.Dbobj.Begin()
	if db.Error != nil {
		logger.Error(db.Error)
		return db.Error
	}

	var org entity.SysOrg
	if err := db.Select("id", "parent_id").Where("id = ?", sysOrg.Id).First(&org).Error; err != nil {
		logger.Error(err)
		return err
	}
	rst, err := this.GetOrgChildById(org.Id)
	if err != nil {
		return err
	}
	for _, val := range rst {
		if err := db.Model(&val).Update("parent_id", org.ParentId).Error; err != nil {
			logger.Error(err)
			db.Rollback()
			return err
		}
	}
	sqlStr := map[string]interface{}{
		"name":           sysOrg.Name,
		"parent_id":      sysOrg.ParentId,
		"update_user_id": sysOrg.UpdateUserId,
		"update_time":    sysOrg.UpdateTime,
	}
	if err := orm.Dbobj.Model(&sysOrg).Updates(sqlStr).Error; err != nil {
		logger.Error(err)
		db.Rollback()
		return err
	}
	if err := db.Commit().Error; err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (this *OrgModel) GetOrgChildById(orgId string) ([]entity.SysOrg, error) {
	if validator.IsEmpty(orgId) {
		logger.Error(orgId)
		return nil, errors.New("error_org_id_is_empty")
	}
	var orgList []entity.SysOrg
	if err := orm.Dbobj.Where("parent_id = ?", orgId).Find(&orgList).Error; err != nil {
		logger.Error(err)
		return nil, err
	}
	return orgList, nil
}

func (this *OrgModel) ImportOrg(sysOrgList []entity.SysOrg) error {
	db := orm.Dbobj.Begin()
	if db.Error != nil {
		logger.Error(db.Error)
		return db.Error
	}
	for _, val := range sysOrgList {
		if err := db.Omit("update_user_id", "update_time").Create(&val).Error; err != nil {
			logger.Error(err)
			return err
		}
	}
	return db.Error
}
