<div class="row subsystem-header">
    <div class="pull-left">
        <span style="font-size: 14px;">组织架构管理</span>
    </div>
</div>
<div class="row subsystem-toolbar">
    <div class="pull-left" style="height: 44px; line-height: 44px; width: 260px;">
    </div>
    <div class="pull-right">
        &nbsp;
        {{if checkResIDAuth "2" "0103020200"}}
        <button onclick="OrgObj.add()" class="btn btn-info btn-sm" title="新增机构信息">
            <i class="icon-plus"> 新增</i>
        </button>
        {{end}}
        {{if checkResIDAuth "2" "0103020600"}}
        <button onclick="OrgObj.upload()" class="btn btn-info btn-sm" title="导入机构信息">
        <span class="icon-edit"> 导入</span>
        </button>
        {{end}}
        {{if checkResIDAuth "2" "0103020500"}}
        <button onclick="OrgObj.download()" class="btn btn-info btn-sm" title="导出机构信息">
        <span class="icon-trash"> 导出</span>
        </button>
        {{end}}
        {{if checkResIDAuth "2" "0103020300"}}
        <button onclick="OrgObj.edit()" class="btn btn-info btn-sm" title="编辑机构信息">
            <span class="icon-edit"> 编辑</span>
        </button>
        {{end}}
        {{if checkResIDAuth "2" "0103020400"}}
        <button onclick="OrgObj.delete()" class="btn btn-danger btn-sm" title="删除机构信息">
            <span class="icon-trash"> 删除</span>
        </button>
        {{end}}
    </div>
</div>
<div class="subsystem-content">
    <div class="row">
        <div class="col-sm-12 col-md-12 col-lg-7">
            <div id="h-org-tree-info" class="thumbnail">
                <div class="col-ms-12 col-md-12 col-lg-12">
                    <div style="border-bottom: #598f56 solid 1px;height: 44px; line-height: 44px;">
                        <div class="pull-left">
                            <span><i class="icon-sitemap"> </i>组织架构树</span>
                        </div>
                        <div class="pull-right">
                        <span>
                            <i class=" icon-search" style="margin-top: 15px;"></i>&nbsp;
                    </span>
                        </div>
                    </div>
                </div>
                <div id="h-org-tree-info-list" class="col-sm-12 col-md-12 col-lg-12"
                     style="padding:15px 5px;overflow: auto">
                </div>
            </div>
        </div>
        <div class="col-sm-12 col-md-12 col-lg-5" style="padding-left: 0px;">
            <div id="h-org-table-info" class="thumbnail" style="padding-left:15px; padding-right: 15px;">
                <div style="border-bottom: #006c8f solid 2px;height: 44px; line-height: 44px;">
                    <div class="pull-left" style="height: 44px; line-height: 44px;">
                        <span class="icon-info-sign" style="height: 44px;line-height: 44px;">&nbsp;机构详细信息</span>
                    </div>
                </div>
                <table class="table table-bordered table-condensed" style="margin-top: 6px;">
                    <tr style="height: 36px; line-height: 36px;">
                        <td style="width: 90px;background-color: #fafafa;text-align: right;padding-right: 15px;vertical-align: middle;">机构编码
                        </td>
                        <td id="h-org-row-org-unit-id" style="font-weight: 600;vertical-align: middle;padding-left: 15px;">-
                        </td>
                    </tr>
                    <tr style="height: 36px; line-height: 36px;">
                        <td style="background-color: #fafafa;text-align: right;padding-right: 15px;vertical-align: middle;">机构名称
                        </td>
                        <td id="h-org-row-org-unit-desc" style="font-weight: 600;vertical-align: middle;padding-left: 15px;">-
                        </td>
                    </tr>
                    <tr style="height: 36px; line-height: 36px;">
                        <td style="background-color: #fafafa;text-align: right;padding-right: 15px;vertical-align: middle;">上级机构
                        </td>
                        <td id="h-org-row-up-org-id" style="font-weight: 600;vertical-align: middle;padding-left: 15px;">-
                        </td>
                    </tr>
                    <tr style="height: 36px; line-height: 36px;">
                        <td style="background-color: #fafafa;text-align: right;padding-right: 15px;vertical-align: middle;">创建日期
                        </td>
                        <td id="h-org-row-create-time" style="font-weight: 600;vertical-align: middle;padding-left: 15px;">-
                        </td>
                    </tr>
                    <tr style="height: 36px; line-height: 36px;">
                        <td style="background-color: #fafafa;text-align: right;padding-right: 15px;vertical-align: middle;">创建人
                        </td>
                        <td id="h-org-row-create-user" style="font-weight: 600;vertical-align: middle;padding-left: 15px;">-
                        </td>
                    </tr>
                    <tr style="height: 36px; line-height: 36px;">
                        <td style="background-color: #fafafa;text-align: right;padding-right: 15px;vertical-align: middle;">修改日期
                        </td>
                        <td id="h-org-row-modify-time" style="font-weight: 600;vertical-align: middle;padding-left: 15px;">-
                        </td>
                    </tr>
                    <tr style="height: 36px; line-height: 36px;">
                        <td style="background-color: #fafafa;text-align: right;padding-right: 15px;vertical-align: middle;">修改人
                        </td>
                        <td id="h-org-row-modify-user" style="font-weight: 600;vertical-align: middle;padding-left: 15px;">-
                        </td>
                    </tr>
                </table>
            </div>
        </div>
    </div>
</div>
<script type="text/javascript">

    var OrgObj = {
        $table:$('#h-org-info-table-details'),
        /*
        * 新增机构信息,只能在自己拥有写入权限的域中新增机构信息
        * */
        add:function(){
            $.Hmodal({
                body:$("#org_input_form").html(),
                width:"420px",
                header:"新增机构",
                callback:function(hmode){
                    $.HAjaxRequest({
                        url:"/auth/org/addorg",
                        data:$("#h-org-add-info").serialize(),
                        type:"post",
                        success:function (data) {
                            $.Notify({
                                message:"插入机构信息成功",
                                type:"success",
                            });
                            $(hmode).remove();
                            var domain_id = $("#h-org-domain-list").val()
                            OrgObj.tree(domain_id)
                        },
                    })
                },
                preprocess:function() {
                    $.getJSON("/auth/org/getall",function(data){
                        var arr = new Array();
                        $(data).each(function(index,element){
                            var ijs = {};
                            ijs.id=element.id;
                            ijs.text=element.name;
                            ijs.parentId=element.parentId;
                            arr.push(ijs)
                        });

                        var ijs = {};
                        ijs.id="system_root";
                        ijs.text="机构树根节点";
                        ijs.parentId="######mf#####";
                        arr.push(ijs)

                        $("#h-org-up-id").Hselect({
                            data:arr,
                            height:"30px",
                            value:$("#h-org-tree-info-list").attr("data-selected"),
                        })
                    })
                }
            })
        },
        /*
        * 编辑处理函数,
        * 当右侧table中没有机构被选中时,默认会编辑左侧被选中的机构
        * 如果左侧也没有机构被选中,则提示没有任何被选中的机构
        * 这个函数只在edit函数中被调用
        * */
        handle_edit:function(row){
            $.Hmodal({
                body:$("#org_modify_form").html(),
                header:"修改机构信息",
                width:"420px",
                preprocess:function () {
                    /*
                     * 初始化下拉框中机构信息
                     * */
                    $.getJSON("/auth/org/getall",function(data){
                        var arr = new Array()
                        $(data).each(function(index,element){
                            if(row.id == "0") return false;
                            if(row.id == element.id) return true;
                            var ijs = {};
                            ijs.id=element.id;
                            ijs.text=element.name;
                            ijs.parentId=element.parentId;
                            arr.push(ijs);
                        });
                        var ijs = {};
                        ijs.id="system_root";
                        ijs.text="机构树根节点";
                        ijs.parentId="######mf#####";
                        arr.push(ijs);

                        $("#h-modify-org-up-id").Hselect({
                            data:arr,
                            value:row.parentId,
                            height:"30px",
                        });
                    });

                    /*
                     * 在编辑框中，填上目前的机构信息。
                     * */
                    var code = row.code;
                    var name = row.name;

                    $("#h-modify-org-id").val(code);
                    $("#h-modify-org-name").val(name);

                },
                callback:function(hmode){
                    var dt = $("#h-org-modify-info").serialize();
                    dt = dt + "&id=" + $("#h-org-tree-info-list").attr("data-selected");
                    $.HAjaxRequest({
                        url:"/auth/org/updateorg",
                        type:"put",
                        data:dt,
                        success:function (data) {
                            $.Notify({
                                title:"温馨提示：",
                                message:"修改机构信息成功",
                                type:"success",
                            });
                            $(hmode).remove();
                            OrgObj.tree();
                        },
                    })
                }
            })
        },
        /*
        * 机构编辑按钮,当点击页面上的就编辑按钮时,
        * 会调用此函数
        * */
        edit:function(){
            var id = $("#h-org-row-org-unit-id").html();
            if (id.trim() != "-"){
                var row = {};
                row.id = $("#h-org-tree-info-list").attr("data-selected");
                row.code = $("#h-org-row-org-unit-id").html();
                row.name = $("#h-org-row-org-unit-desc").html();
                row.parentId = $("#h-org-row-up-org-id").html();
                row.createTime = $("#h-org-row-create-time").html();
                row.createUserId = $("#h-org-row-create-user").html();
                row.updateTime = $("#h-org-row-modify-time").html();
                row.updateUserId = $("#h-org-row-modify-user").html();
                OrgObj.handle_edit(row);
            } else {
                $.Notify({
                    message:"请在机构树中选择需要编辑的机构",
                    type:"warning",
                })
            }
        },
        delete:function(){
            var selected_id = $("#h-org-tree-info-list").attr("data-selected");
            if (selected_id == undefined || selected_id == null || selected_id.length==0) {
                $.Notify({
                    message: "请在列表中选择一个需要编辑的机构",
                    type: "warning",
                });
                return
            }

            $.Hconfirm({
                callback:function () {
                    $.HAjaxRequest({
                        url:"/auth/org/deleteorg",
                        type:"post",
                        data:{orgId:selected_id},
                        success:function () {
                            $.Notify({
                                title:"操作成功",
                                message:"删除机构信息成功",
                                type:"success",
                            });
                            OrgObj.tree()
                        },
                    })
                },
                body:"点击确认删除选中的机构信息"
            })
        },
        download:function(){
            var domain_id = $("#h-org-domain-list").val()
            var x=new XMLHttpRequest();
            x.open("GET", "/auth/org/exportorg?domain_id="+domain_id, true);
            x.responseType = 'blob';
            x.onload=function(e){
                download(x.response, "机构信息.xlsx", "application/vnd.ms-excel");
            };
            x.send();
        },
        upload:function(param){
            $.Hupload({
                url:"/auth/org/importorg",
                header:"导入机构信息",
                callback:function () {
                    var domain_id = $("#h-org-domain-list").val();
                    OrgObj.tree(domain_id)
                },
            })
        },
        tree:function(){
          $.getJSON("/auth/org/getall",function(data){
              if (data.length==0){
                  $.Notify({
                      title:"温馨提示",
                      message:"您选择的域中没有机构信息",
                      type:"info",
                  });
                  OrgObj.$table.bootstrapTable('load',[])
                  $("#h-org-tree-info-list").Htree({
                      data:[],
                  })
              } else {
                  var arr = new Array()
                  $(data).each(function(index,element){
                      var ijs = {}
                      ijs.id = element.id
                      ijs.text = element.name
                      ijs.parentId = element.parentId
                      arr.push(ijs)
                  });
                  $("#h-org-tree-info-list").HtreeWithLine({
                      data:arr,
                      onChange:function(obj){
                          var id = $(obj).attr("data-id")
                          $.getJSON("/auth/org/details",{orgId:id},function(data){
                                $("#h-org-row-org-unit-id").html(data.code);
                                $("#h-org-row-org-unit-desc").html(data.name);
                                $("#h-org-row-up-org-id").html(data.parentName);
                                $("#h-org-row-create-time").html(data.createTime);
                                $("#h-org-row-create-user").html(data.createUser);
                                $("#h-org-row-modify-time").html(data.updateTime);
                                $("#h-org-row-modify-user").html(data.updateUser);
                          });
                      }
                  });
              }
          })
        },
    };

    $(document).ready(function(){
        var hwindow = document.documentElement.clientHeight;
        $("#h-org-tree-info").height(hwindow-156);
        $("#h-org-table-info").height(hwindow-156);
        $("#h-org-tree-info-list").height(hwindow-228);
        $('#h-org-info-table-details').bootstrapTable({
            height:hwindow-156,
        });
        OrgObj.tree();
    });
</script>

<script type="text/html" id="org_input_form">
    <form class="row"  id="h-org-add-info">
        <div class="col-sm-12">
            <label class="h-label" style="width:100%;">组织机构代码：</label>
            <input placeholder="请输入1-30位数字，字母（必填）" name="code" type="text" class="form-control" style="width: 100%;height: 30px;line-height: 30px;">
        </div>
        <div class="col-sm-12" style="margin-top: 23px">
            <label class="h-label" style="width: 100%;">组织机构名称：</label>
            <input placeholder="请输入1-60位汉字，字母，数字（必填）" type="name" class="form-control" name="name" style="width: 100%;height: 30px;line-height: 30px;">
        </div>
        <div class="col-sm-12" style="margin-top: 23px">
            <label class="h-label" style="width: 100%;">上级组织机构代码：</label>
            <select id="h-org-up-id" name="parentId" style="width: 100%;height: 30px;line-height: 30px;">
            </select>
        </div>
    </form>
</script>

<script type="text/html" id="org_modify_form">
    <form class="row" id="h-org-modify-info">
        <div class="col-sm-12" style="margin-top: 12px">
            <label class="h-label" style="width:100%;">组织部门代码：</label>
            <input id="h-modify-org-id" readonly="readonly" placeholder="1-30位字母、数字组成" name="code" type="text" class="form-control" style="width: 100%;height: 30px;line-height: 30px;">
        </div>
        <div class="col-sm-12" style="margin-top: 23px;">
            <label class="h-label" style="width: 100%;">组织部门名称：</label>
            <input id="h-modify-org-name" placeholder="组织部门描述信息" type="text" class="form-control" name="name" style="width: 100%;height: 30px;line-height: 30px;">
        </div>
        <div class="col-sm-12" style="margin-top: 23px;">
            <label class="h-label" style="width: 100%;">上级组织部门代码：</label>
            <select id="h-modify-org-up-id" name="parentId" style="width: 100%;height: 30px;line-height: 30px;">
            </select>
        </div>
    </form>
</script>