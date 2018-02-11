package models

import (
	"mifanpark/dto"
	"mifanpark/entity"
	"mifanpark/utilities/common"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm"
)

type LoginModel struct {
}

func (this *LoginModel) Login(rdto dto.AuthDto) (entity.SysUser, dto.AuthDto) {
	user, err := this.GetUserByAccount(rdto.Account)
	rdto.Success = false
	if err != nil {
		rdto.Code = 402
		rdto.Msg = "用户不存在"
		return user, rdto
	}

	if user.StatusId != "0" {
		rdto.Code = 406
		rdto.Msg = "用户状态被锁定，请联系管理员解锁"
		return user, rdto
	}

	if user.ErrorCount > 6 {
		err = this.ForbidUser(user)
		if err != nil {
			logger.Error("冻结用户失败,用户Id:", user.Id, err)
			rdto.Code = 405
			rdto.Msg = "冻结用户失败"
			return user, rdto
		}
		rdto.Code = 403
		rdto.Msg = "用户已锁定，请联系管理员解锁"
		return user, rdto
	}

	if user.Account == rdto.Account && user.Password == rdto.Password {
		err = this.UpdateErrorCount(0, user)
		if err != nil {
			logger.Error("修改用户错误次数失败,用户Id:", user.Id, err)
			rdto.Code = 405
			rdto.Msg = "修改用户错误次数失败"
			return user, rdto
		}
		err = this.UpdateLoginInfo(rdto.LoginIp, user)
		if err != nil {
			logger.Error("修改用户登录信息失败,用户Id:", user.Id, err)
			rdto.Code = 405
			rdto.Msg = "修改用户登录信息失败"
			return user, rdto
		}
		rdto.Success = true
		rdto.Code = 200
		return user, rdto
	} else {
		err = this.UpdateErrorCount(user.ErrorCount+1, user)
		if err != nil {
			logger.Error("修改用户错误次数失败,用户Id:", user.Id, err)
			rdto.Code = 405
			rdto.Msg = "修改用户错误次数失败"
			return user, rdto
		}
		rdto.Code = 405
		rdto.Msg = "用户密码错误"
		return user, rdto
	}
}

func (this *LoginModel) GetUserByAccount(account string) (entity.SysUser, error) {
	var user entity.SysUser
	err := orm.Dbobj.Where("account = ?", account).First(&user).Error
	if err != nil {
		logger.Error(err)
	}
	return user, err
}

func (this *LoginModel) ForbidUser(user entity.SysUser) error {
	err := orm.Dbobj.Model(&user).Where("id", user.Id).Update("user_status_id", "1").Error
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (this *LoginModel) UpdateErrorCount(count int, user entity.SysUser) error {
	err := orm.Dbobj.Model(&user).Where("id = ?", user.Id).Update("error_count", count).Error
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (this *LoginModel) UpdateLoginInfo(loginIp string, user entity.SysUser) error {
	sqlStr := map[string]interface{}{
		"login_ip":        loginIp,
		"last_login_time": common.CurTime(),
	}
	err := orm.Dbobj.Model(&user).Where("id = ?", user.Id).Updates(sqlStr).Error
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (this *LoginModel) GetDefaultTheme(userId string) string {
	var url string
	err := orm.Dbobj.Table("sys_index_page").Select("sys_index_page").
		Joins("inner join sys_user_theme on sys_index_page.theme_id = sys_user_theme.theme_id").
		Where("sys_user_theme.user_id = ?", userId).Row().Scan(&url)
	if err != nil {
		logger.Error(err)
		url = "./views/auth/theme/default/index.tpl"
	}
	return url
}
