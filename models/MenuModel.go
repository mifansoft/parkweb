package models

import (
	"encoding/json"
	"mifanpark/entity"
	"mifanpark/utilities/common"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
)

type MenuModel struct {
	roleUserModel      RoleUserModel
	userThemeModel     UserThemeModel
	themeResourceModel ThemeResModel
	roleResModel       RoleResModel
	resourceModel      ResModel
}

const redirect = `
<script type="text/javascript">
    $.Hconfirm({
		cancelBtn:false,
        header:"获取页面失败",
        body:"获取页面失败，请与管理员联系",
        callback:function () {
            window.location.href="/"
        }
    })
</script>
`

func (this *MenuModel) GetHomePageMenus(resId, typeId, userId string) ([]byte, error) {
	//获取用户主题信息
	theme, err := this.userThemeModel.GetUserTheme(userId)
	if err != nil {
		return nil, err
	}

	//获取这个主题下的所有资源信息
	themeRes, err := this.themeResourceModel.GetThemeRes(theme.ThemeId)
	if err != nil {
		return nil, err
	}

	var mres = make(map[string]entity.SysRes)
	var rst []entity.HomeMenu
	//如果是超级管理员，获取所有菜单信息
	if common.IsAdmin(userId) {
		resData, err := this.resourceModel.GetChildrenRes(resId)
		if err != nil {
			return nil, err
		}
		for _, val := range resData {
			if val.ResTypeId == typeId {
				mres[val.Id] = val
			}
		}
	} else {
		//获取这个用户的角色信息
		roles, err := this.roleUserModel.GetRoleByUserId(userId)
		if err != nil {
			return nil, err
		}

		var roleList []string
		for _, val := range roles {
			roleList = append(roleList, val.Id)
		}

		roleRes, err := this.roleResModel.GetRoleResWithUser(roleList, resId, typeId)
		if err != nil {
			return nil, err
		}
		for _, val := range roleRes {
			mres[val.Id] = val
		}
	}
	for _, tres := range themeRes {
		if val, ok := mres[tres.ResId]; ok {
			var menu entity.HomeMenu
			menu.ResId = tres.Id
			menu.ParentId = val.ParentId
			menu.ResName = val.Name
			menu.GroupId = tres.GroupId
			menu.ResBgColor = tres.BgColor
			menu.ResClass = tres.ResClass
			menu.ResImg = tres.ResImg
			menu.ResUrl = tres.ResUrl
			menu.OpenTypeId = tres.OpenTypeId
			menu.NewIFrame = tres.NewIFrame
			rst = append(rst, menu)
		}
	}
	return json.Marshal(rst)
}

func (this *MenuModel) GetResUrl(userId, themeResId string) (string, error) {
	var url string
	err := orm.Dbobj.Table("sys_user_theme").Select("distinct sys_theme_res.res_url").
		Joins("inner join sys_theme_res on sys_user_theme.theme_id = sys_theme_res.theme_id").
		Joins("inner join sys_res on sys_theme_res.res_id = sys_res.id").
		Where("sys_user_theme.user_id = ?", userId).
		Where("sys_theme_res.id = ?", themeResId).
		Where("sys_res.res_type_id = ?", "0").Row().Scan(&url)

	if err != nil {
		logger.Error(err)
		return redirect, err
	}
	return url, nil
}
