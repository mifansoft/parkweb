package models

import (
	"mifanpark/dto"
	"mifanpark/entity"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
	"mifanpark/utilities/uuid"
	"sync"
	"time"
)

type RoleResModel struct {
	resource ResModel
	lock     sync.RWMutex
}

func (this *RoleResModel) CheckUrlAuth(userId string, url, method string) bool {
	var resId []string
	err := orm.Dbobj.Select("res_id").Where("user_id", userId).Where("res_url", url).Where("method", method).Scan(&resId).Error
	if err != nil {
		logger.Error(err)
		return false
	}
	if len(resId) > 0 {
		return true
	}
	logger.Error("insufficient privileges", "user id is :", userId, "api is :", url, "request method is:", method)
	return false
}

func (this *RoleResModel) GetRoleResWithUser(roleIds []string, resId ...string) ([]entity.SysRes, error) {
	var rst []entity.SysRes
	var roleRes map[string]string = make(map[string]string)
	for _, val := range roleIds {
		tmp, err := this.getUserRoleResByRoleId(val)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		for _, p := range tmp {
			roleRes[p.ResId] = ""
		}
	}

	var rstRes []entity.SysRes
	if len(resId) == 1 {
		var err error
		rstRes, err = this.resource.GetChildrenRes(resId[0])
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	} else if len(resId) == 2 {
		tmp, err := this.resource.GetChildrenRes(resId[0])
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		for _, val := range tmp {
			if val.ResTypeId == resId[1] {
				rstRes = append(rstRes, val)
			}
		}
	} else {
		var err error
		rstRes, err = this.resource.GetResourceAll()
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	for _, res := range rstRes {
		if _, ok := roleRes[res.Id]; ok {
			rst = append(rst, res)
		}
	}
	return rst, nil
}

func (this *RoleResModel) getUserRoleResByRoleId(roleId string) ([]entity.SysRoleRes, error) {
	var rst []entity.SysRoleRes
	err := orm.Dbobj.Where("role_id = ?", roleId).Find(&rst).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return rst, err
}

func (this *RoleResModel) CheckResIDAuth(userId, resId string) bool {
	cnt := 0
	err := orm.Dbobj.Table("sys_role_user").Joins("inner join sys_role_res on sys_role_user.role_id = sys_role_res.role_id").
		Where("sys_role_user.user_id", userId).Where("sys_role_res.res_id", resId).Count(&cnt).Error
	if err != nil {
		logger.Error(err)
		return false
	}
	if cnt > 0 {
		return true
	}
	return false
}

func (this *RoleResModel) GetRoleResByRoleId(roleId string) ([]dto.RoleResDto, error) {
	var roleResList []entity.SysRoleRes
	var roleResDto []dto.RoleResDto
	var resList []entity.SysRes

	err := orm.Dbobj.Select("role_id,res_id").Where("role_id = ?", roleId).Find(&roleResList).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = orm.Dbobj.Find(&resList).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	for _, roleRes := range roleResList {
		for _, res := range resList {
			if roleRes.ResId == res.Id {
				dt := new(dto.RoleResDto)
				dt.ResId = res.Id
				dt.Name = res.Name
				dt.ParentId = res.ParentId
				roleResDto = append(roleResDto, *dt)
				break
			}
		}
	}
	return roleResDto, nil
}

func (this *RoleResModel) UnGetRoleResByRoleId(roleId string) ([]dto.RoleResDto, error) {
	var roleResList []entity.SysRoleRes
	var resList []entity.SysRes
	var roleResDto []dto.RoleResDto
	err := orm.Dbobj.Select("role_id,res_id").Where("role_id = ?", roleId).Find(&roleResList).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = orm.Dbobj.Find(&resList).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var diff = make(map[string]dto.RoleResDto)
	for _, val := range resList {
		dt := new(dto.RoleResDto)
		dt.ResId = val.Id
		dt.ParentId = val.ParentId
		dt.Name = val.Name
		diff[dt.ResId] = *dt
	}

	for _, val := range roleResList {
		delete(diff, val.ResId)
	}

	tmp := this.searchParent(diff, resList)
	if len(tmp) != 0 {
		for _, val := range tmp {
			diff[val.ResId] = val
		}
		tmp = this.searchParent(diff, resList)
	}
	for _, val := range diff {
		roleResDto = append(roleResDto, val)
	}
	return roleResDto, nil
}

func (this *RoleResModel) searchParent(diff map[string]dto.RoleResDto, all []entity.SysRes) []dto.RoleResDto {
	var ret []dto.RoleResDto
	for _, val := range diff {
		if _, ok := diff[val.ParentId]; !ok {
			for _, res := range all {
				if res.Id == val.ResId {
					dt := new(dto.RoleResDto)
					dt.ResId = res.Id
					dt.Name = res.Name
					dt.ParentId = res.ParentId
					ret = append(ret, *dt)
				}
			}
		}
	}
	return ret
}

func (this *RoleResModel) Authorized(roleId, userId string, resList []string) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	all, err := this.resource.GetResourceAll()
	if err != nil {
		return err
	}

	getted, err := this.GetRoleResByRoleId(roleId)
	if err != nil {
		return err
	}

	var mp = make(map[string]bool)
	for _, val := range getted {
		mp[val.ResId] = true
	}

	db := orm.Dbobj.Begin()
	if db.Error != nil {
		logger.Error(db.Error)
		return db.Error
	}

	for _, val := range resList {
		var rst []dto.RoleResDto
		this.parent(all, val, &rst)
		if len(rst) == 0 {
			if _, yes := mp[val]; !yes {
				mp[val] = true
				var roleRes = new(entity.SysRoleRes)
				roleRes.Id = uuid.Random()
				roleRes.RoleId = roleId
				roleRes.ResId = val
				roleRes.CreateTime = time.Now()
				roleRes.CreateUserId = userId
				if err := db.Omit("update_user_id", "update_time").Create(&roleRes).Error; err != nil {
					logger.Error(err)
					db.Rollback()
					return err
				}
			}
		} else {
			for _, dt := range rst {
				if _, yes := mp[dt.ResId]; !yes {
					mp[dt.ResId] = true
					var roleRes = new(entity.SysRoleRes)
					roleRes.Id = uuid.Random()
					roleRes.RoleId = roleId
					roleRes.ResId = dt.ResId
					roleRes.CreateTime = time.Now()
					roleRes.CreateUserId = userId
					if err := db.Omit("update_user_id", "update_time").Create(&roleRes).Error; err != nil {
						logger.Error(err)
						db.Rollback()
						return err
					}
				}
			}
		}
	}
	if err := db.Commit().Error; err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (this *RoleResModel) parent(all []entity.SysRes, resId string, ret *[]dto.RoleResDto) {
	for _, val := range all {
		if val.Id == resId {
			dt := new(dto.RoleResDto)
			dt.ResId = val.Id
			dt.Name = val.Name
			dt.ParentId = val.ParentId
			*ret = append(*ret, *dt)
			if val.ParentId != val.Id {
				this.parent(all, val.ParentId, ret)
			}
		}
	}
}

func (this *RoleResModel) UnAuthorized(roleId string, resList []string) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	var mp = make(map[string]bool)

	for _, val := range resList {
		mp[val] = true
	}

	all, err := this.GetRoleResByRoleId(roleId)
	if err != nil {
		return err
	}

	db := orm.Dbobj.Begin()
	if db.Error != nil {
		logger.Error(db.Error)
		return db.Error
	}

	for _, val := range resList {
		var rst []dto.RoleResDto
		this.dfs(all, val, &rst)
		if len(rst) == 0 {
			if err := db.Where("role_id = ?", roleId).Where("res_id = ?", val).Delete(entity.SysRoleRes{}).Error; err != nil {
				logger.Error(err)
				db.Rollback()
				return err
			}
		} else {
			var flag = true
			for _, dt := range rst {
				if _, yes := mp[dt.ResId]; !yes {
					flag = false
					break
				}
			}
			if flag {
				if err := db.Where("role_id = ?", roleId).Where("res_id = ?", val).Delete(entity.SysRoleRes{}).Error; err != nil {
					logger.Error(err)
					db.Rollback()
					return err
				}
			}
		}
	}
	if err := db.Commit().Error; err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (this *RoleResModel) dfs(rst []dto.RoleResDto, resId string, ret *[]dto.RoleResDto) {
	for _, val := range rst {
		if resId == val.ParentId {
			*ret = append(*ret, val)
			if val.ResId != val.ParentId {
				this.dfs(rst, val.ResId, ret)
			}
		}
	}
}
