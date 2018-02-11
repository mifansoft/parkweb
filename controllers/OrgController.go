package controllers

import (
	"errors"
	"io/ioutil"
	"mifanpark/entity"
	"mifanpark/models"
	"mifanpark/service/auth"
	"mifanpark/utilities/common"
	"mifanpark/utilities/groupcache"
	"mifanpark/utilities/hret"
	"mifanpark/utilities/i18n"
	"mifanpark/utilities/jwt"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/uuid"
	"mifanpark/utilities/validator"
	"net/http"
	"path/filepath"
	"time"

	"github.com/tealeg/xlsx"
)

type orgModel struct {
	model     *models.OrgModel
	importOrg chan int
}

var OrgCtl = &orgModel{
	model:     new(models.OrgModel),
	importOrg: make(chan int, 1),
}

func init() {
	groupcache.RegisterStaticFile("MiFanParkOrgPage", "./views/auth/org_page.tpl")
}

func (this *orgModel) OrgPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !auth.BasicAuth(r) {
		hret.Error(w, 403, i18n.NoAuth(r))
		return
	}
	rst, err := groupcache.GetStaticFile("MiFanParkOrgPage")
	if err != nil {
		hret.Error(w, 404, i18n.PageNotFound(r))
		return
	}

	hz, err := auth.ParseText(r, string(rst))
	if err != nil {
		hret.Error(w, 404, i18n.PageNotFound(r))
		return
	}
	hz.Execute(w, nil)
}

func (this *orgModel) GetOrgAll(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rst, err := this.model.GetOrgAll()
	if err != nil {
		logger.Error(err)
		hret.Error(w, 417, i18n.Get(r, "error_query_org"))
		return
	}
	hret.Json(w, rst)
}

func (this *orgModel) GetDetails(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	orgId := r.FormValue("orgId")
	if len(orgId) == 0 {
		logger.Error("机构ID参数为空")
		hret.Error(w, 421, "机构ID参数为空")
		return
	}
	orgDetailsDto, err := this.model.GetDetails(orgId)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 423, err.Error())
		return
	}
	hret.Json(w, orgDetailsDto)
}

func (this *orgModel) AddOrg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var org entity.SysOrg
	err := common.ParseForm(r, &org)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 423, err.Error())
		return
	}

	if !validator.IsAlnum(org.Code) {
		hret.Error(w, 421, i18n.Get(r, "error_org_code_format"), errors.New("error_org_code_format"))
		return
	}

	if validator.IsEmpty(org.Name) {
		hret.Error(w, 421, i18n.Get(r, "error_org_name_empty"), errors.New("error_org_name_empty"))
		return
	}

	if validator.IsEmpty(org.ParentId) {
		hret.Error(w, 421, i18n.Get(r, "error_org_parent_id_empty"), errors.New("error_org_parent_id_empty"))
		return
	}

	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}
	org.Id = uuid.Random()
	org.CreateTime = time.Now()
	org.CreateUserId = jclaim.LoginUser.Id
	err = this.model.AddOrg(org)
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *orgModel) DeleteOrg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	orgId := r.FormValue("orgId")
	if validator.IsEmpty(orgId) {
		hret.Error(w, 421, i18n.Get(r, "error_org_id_empty"), errors.New("error_org_id_empty"))
		return
	}
	err := this.model.DeleteOrg(orgId)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 418, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *orgModel) UpdateOrg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var sysOrg entity.SysOrg
	err := common.ParseForm(r, &sysOrg)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, err.Error())
		return
	}
	if validator.IsEmpty(sysOrg.Id) {
		hret.Error(w, 421, i18n.Get(r, "Id不能为空"), errors.New("Id不能为空"))
		return
	}
	if validator.IsEmpty(sysOrg.Code) {
		hret.Error(w, 421, i18n.Get(r, "组织机构代码不能为空"), errors.New("组织机构代码不能为空"))
		return
	}
	if validator.IsEmpty(sysOrg.Name) {
		hret.Error(w, 421, i18n.Get(r, "组织机构名称不能为空"), errors.New("组织机构名称不能为空"))
		return
	}
	if validator.IsEmpty(sysOrg.ParentId) {
		hret.Error(w, 421, i18n.Get(r, "组织机构父节点不能为空"), errors.New("组织机构父节点不能为空"))
		return
	}
	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}
	sysOrg.UpdateUserId = jclaim.LoginUser.Id
	sysOrg.UpdateTime = time.Now()
	err = this.model.UpdateOrg(sysOrg)
	if err != nil {
		hret.Error(w, 421, i18n.Get(r, err.Error()), err)
		return
	}
	hret.Success(w, i18n.Success(r))
}

func (this *orgModel) ExportOrg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !auth.BasicAuth(r) {
		hret.Error(w, 403, i18n.NoAuth(r))
		return
	}
	w.Header().Set("Content-Type", "application/vnd.ms-excel")

	rst, err := this.model.GetOrgAll()
	if err != nil {
		hret.Error(w, 417, i18n.Get(r, "error_query_org"))
		return
	}

	var sheet *xlsx.Sheet
	file, err := xlsx.OpenFile(filepath.Join("./views", "uploadTemplate", "OrgExportTemplate.xlsx"))
	if err != nil {
		file = xlsx.NewFile()
		sheet, err = file.AddSheet("机构信息")
		if err != nil {
			logger.Error(err)
			hret.Error(w, 421, i18n.Get(r, "error_org_sheet"))
			return
		}

		{
			row := sheet.AddRow()
			cell1 := row.AddCell()
			cell1.Value = "机构编码"
			cell2 := row.AddCell()
			cell2.Value = "机构名称"
			cell3 := row.AddCell()
			cell3.Value = "上级编码"
			cell4 := row.AddCell()
			cell4.Value = "所属域"
			cell5 := row.AddCell()
			cell5.Value = "创建日期"
			cell6 := row.AddCell()
			cell6.Value = "创建人"
			cell7 := row.AddCell()
			cell7.Value = "修改日期"
			cell8 := row.AddCell()
			cell8.Value = "修改人"
		}
	} else {
		sheet = file.Sheet["机构信息"]
		if sheet == nil {
			hret.Error(w, 421, i18n.Get(r, "error_org_sheet"))
			return
		}
	}

	for _, val := range rst {
		row := sheet.AddRow()
		cell1 := row.AddCell()
		cell1.Value = val.Code
		cell1.SetStyle(sheet.Rows[1].Cells[0].GetStyle())
		cell2 := row.AddCell()
		cell2.Value = val.Name
		cell2.SetStyle(sheet.Rows[1].Cells[1].GetStyle())
		cell3 := row.AddCell()
		cell3.Value, _ = common.GetKey(val.ParentId, 2)
		cell3.SetStyle(sheet.Rows[1].Cells[2].GetStyle())
		cell4 := row.AddCell()
		if !val.CreateTime.IsZero() {
			cell4.Value = common.DateToString(val.CreateTime)
		} else {
			cell4.Value = ""
		}
		cell4.SetStyle(sheet.Rows[1].Cells[4].GetStyle())
		cell5 := row.AddCell()
		cell5.Value = val.CreateUserId
		cell5.SetStyle(sheet.Rows[1].Cells[5].GetStyle())
		cell6 := row.AddCell()
		if !val.UpdateTime.IsZero() {
			cell6.Value = common.DateToString(val.UpdateTime)
		} else {
			cell6.Value = ""
		}
		cell6.SetStyle(sheet.Rows[1].Cells[6].GetStyle())
		cell7 := row.AddCell()
		cell7.Value = val.UpdateUserId
		cell7.SetStyle(sheet.Rows[1].Cells[7].GetStyle())
	}

	if len(sheet.Rows) >= 3 {
		sheet.Rows = append(sheet.Rows[0:1], sheet.Rows[2:]...)
	}
	file.Write(w)
}

func (this *orgModel) ImportOrg(w http.ResponseWriter, r *http.Request) {
	if len(this.importOrg) != 0 {
		hret.Success(w, i18n.Get(r, "error_org_upload_wait"))
		return
	}

	jclaim, err := jwt.ParseHttp(r)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 403, i18n.Disconnect(r))
		return
	}

	this.importOrg <- 1
	defer func() {
		<-this.importOrg
	}()

	r.ParseForm()
	fd, _, err := r.FormFile("file")
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, i18n.Get(r, "error_org_read_import_file"))
		return
	}

	result, err := ioutil.ReadAll(fd)
	if err != nil {
		logger.Error(err)
		hret.Error(w, 421, i18n.Get(r, "error_org_read_import_file"))
		return
	}

	file, err := xlsx.OpenBinary(result)
	sheet, ok := file.Sheet["机构信息"]
	if !ok {
		logger.Error("没有找到'机构信息'这个sheet页")
		hret.Error(w, 421, i18n.Get(r, "error_org_sheet"))
		return
	}

	var sysOrgList []entity.SysOrg
	for index, val := range sheet.Rows {
		if index > 0 {
			var sysOrg entity.SysOrg
			sysOrg.Id = uuid.Random()
			sysOrg.Code = val.Cells[0].Value
			sysOrg.Name = val.Cells[1].Value
			sysOrg.ParentId = val.Cells[2].Value
			sysOrg.CreateUserId = jclaim.LoginUser.Id
			sysOrg.CreateTime = time.Now()
			if sysOrg.Id == sysOrg.ParentId {
				hret.Error(w, 421, i18n.Get(r, "as_of_date_import_org_equal_id"))
				return
			}
			if !validator.IsAlnum(sysOrg.Code) {
				hret.Error(w, 421, i18n.Get(r, "机构编码必须由1-30位字母,数字组成"), errors.New("机构编码必须由1-30位字母,数字组成"))
				return
			}
			if validator.IsEmpty(sysOrg.Name) {
				hret.Error(w, 421, i18n.Get(r, "error_org_id_desc_empty"), errors.New("error_org_id_desc_empty"))
				return
			}
			if validator.IsEmpty(sysOrg.ParentId) {
				hret.Error(w, 421, i18n.Get(r, "error_org_up_id_empty"), errors.New("error_org_up_id_empty"))
				return
			}
			sysOrgList = append(sysOrgList, sysOrg)
		}
	}

	err = this.model.ImportOrg(sysOrgList)
}
