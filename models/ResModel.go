package models

import (
	"mifanpark/entity"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
)

type ResModel struct {
	Mtheme ThemeResModel
}

func (this *ResModel) GetChildrenRes(resId string) ([]entity.SysRes, error) {
	rst, err := this.GetResourceAll()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	var ret []entity.SysRes
	this.getChildrenResourceInfo(rst, resId, &ret)
	return ret, nil
}

func (this *ResModel) GetResourceAll() ([]entity.SysRes, error) {
	var rst []entity.SysRes
	err := orm.Dbobj.Table("sys_res").Joins("left join sys_res_attr on sys_res.res_attr_id = sys_res_attr.id").
		Joins("left join sys_res_type on sys_res.res_type_id = sys_res_type.id").Find(&rst).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return rst, err
}

func (this *ResModel) getChildrenResourceInfo(allRes []entity.SysRes, resId string, rst *[]entity.SysRes) {
	for _, val := range allRes {
		if val.ParentId == resId {
			*rst = append(*rst, val)
			if val.Id == val.ParentId {
				logger.Error("层级关系错误,不允许上级菜单域当前菜单编码一致,当前菜单编码:", val.Id, "上级菜单编码:", val.ParentId)
				return
			}
			this.getChildrenResourceInfo(allRes, val.Id, rst)
		}
	}
}
