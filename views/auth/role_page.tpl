<div class="row subsystem-header">
    <div class="pull-left">
        <span style="font-size: 14px;">角色管理</span>
    </div>
</div>
<div class="row subsystem-toolbar">
    <div class="pull-left" style="height: 44px; line-height: 44px; width: 260px;">
    </div>
    <div class="pull-right">
        &nbsp;
        {{if checkResIDAuth "2" "0105020200"}}
        <button onclick="RoleObj.add()" class="btn btn btn-info btn-sm">
            <i class="icon-plus"> 新增</i>
        </button>
        {{end}}
        {{if checkResIDAuth "2" "0105020300"}}
        <button onclick="RoleObj.edit()" class="btn btn btn-info btn-sm">
            <i class="icon-edit"> 编辑</i>
        </button>
        {{end}}
        {{if checkResIDAuth "2" "0105020400"}}
        <button onclick="RoleObj.delete()" class="btn btn btn-danger btn-sm">
            <i class="icon-trash"> 删除</i>
        </button>
        {{end}}
    </div>
</div>
<div id="h-role-info" class="subsystem-content">
    <table id="h-role-info-table-details"
           data-toggle="table"
           data-striped="true"
           data-side-pagination="client"
           data-click-to-select="true"
           data-pagination="true"
           data-page-size="30"
           data-page-list="[20,30, 50, 100, 200]"
           data-search="false">
        <thead>
        <tr>
            <th data-field="state" data-checkbox="true"></th>
            <th data-field="id" class="hidden">id</th>
            <th data-align="center" data-field="name">角色名称</th>
            <th data-align="center" data-field="status">状态</th>
            <th data-align="center" data-field="createUser">创建人</th>
            <th data-align="center" data-field="createTime">创建时间</th>
            <th data-align="center" data-field="updateUser">修改人</th>
            <th data-align="center" data-field="updateTime">修改时间</th>
            <th data-field="state-handle" data-align="center" data-formatter="RoleObj.formatter">资源操作</th>
        </tr>
        </thead>
    </table>
</div>

<script>
    var RoleObj={
        getUserRelationPage:function (roleId, name) {
            var name = name + "关联用户信息";
            Hutils.openTab({
                url:"/auth/role/user/relationpage?roleId=" + roleId,
                id:"resourcedetails999988899",
                title:name,
                error:function (m) {
                    $.Notify({
                        title:"温馨提示：",
                        message:"权限不足",
                        type:"danger",
                    })
                }
            })
        },
        getRoleResPage:function(id,name){
            var name = name + "资源信息";
            Hutils.openTab({
                url:"/auth/role/details?id=" + id,
                id:"resourcedetails999988899",
                title:name,
                error:function (m) {
                    $.Notify({
                        title:"温馨提示：",
                        message:"权限不足",
                        type:"danger",
                    })
                }
            })
        },
        formatter:function(value,rows,index){
            var html = "-";
            {{if checkResIDAuth "2" "0105020500"}}
                html = '<span class="h-td-btn btn-primary btn-xs" onclick="RoleObj.getRoleResPage(\''+rows.id+'\',\''+ rows.name+'\')">功能权限</span>';
            {{end}}
            {{if checkResIDAuth "2" "0105020600"}}
                html += '&nbsp;&nbsp;<span class="h-td-btn btn-success btn-xs" onclick="RoleObj.getUserRelationPage(\''+rows.id+'\',\''+ rows.name+'\')">关联用户</span>';
            {{end}}
            return html;
        },
        add:function () {
            $.Hmodal({
                header:"新增角色",
                body:$("#role_input_form").html(),
                width:"420px",
                preprocess:function () {
                    $("#h-role-add-status").Hselect({
                        height:"30px",
                        value:"0",
                    });
                },
                callback:function (hmode) {
                    $.HAjaxRequest({
                        url:"/auth/role/addrole",
                        type:"post",
                        data:$("#h-role-add-info").serialize(),
                        success:function () {
                            $.Notify({
                                title:"操作成功",
                                message:"新增角色信息成功",
                                type:"success",
                            });
                            $(hmode).remove();
                            $("#h-role-info-table-details").bootstrapTable('refresh');
                        },
                    })
                }
            })
        },
        edit:function () {
            var rows = $("#h-role-info-table-details").bootstrapTable('getSelections');
            if (rows.length==0){
                $.Notify({
                    title:"温馨提示",
                    message:"您没有选择需要编辑的角色信息",
                    type:"info",
                });
            }else if (rows.length==1){
                var id = rows[0].id;
                var name = rows[0].name;
                var status = rows[0].status;

                $.Hmodal({
                    header:"编辑角色信息",
                    body:$("#role_modify_form").html(),
                    width:"420px",
                    preprocess:function () {
                        $("#h-role-modify-role-name").val(name);
                        $("#h-role-modify-role-status-cd").Hselect({
                            height:"30px",
                            value:status,
                        });
                    },
                    callback:function (hmode) {
                        var reqData = $("#h-role-modify-info").serialize() + "&id=" + id
                        $.HAjaxRequest({
                            url:"/auth/role/updaterole",
                            type:"put",
                            data:reqData,
                            success:function () {
                                $.Notify({
                                    title:"操作成功",
                                    message:"编辑角色信息成功",
                                    type:"success",
                                });
                                $(hmode).remove();
                                $("#h-role-info-table-details").bootstrapTable('refresh')
                            },
                        })
                    },
                })

            }else {
                $.Notify({
                    title:"温馨提示",
                    message:"您选择了多行角色信息，不知道确定要编辑哪一行",
                    type:"info",
                });
            }
        },
        delete:function () {
           var rows = $("#h-role-info-table-details").bootstrapTable('getSelections');
           if (rows.length==0){
               $.Notify({
                   title:"温馨提示",
                   message:"您没有选择需要删除的角色信息",
                   type:"info",
               });
               return;
           }else{
               $.Hconfirm({
                   callback:function () {
                       $.HAjaxRequest({
                           url:"/auth/role/deleterole",
                           type:"post",
                           data:{dataJson:JSON.stringify(rows)},
                           success:function () {
                               $.Notify({
                                   title:"操作成功",
                                   message:"删除角色信息成功",
                                   type:"success",
                               });
                               $("#h-role-info-table-details").bootstrapTable('refresh')
                           },
                       })
                   },
                   body:"确认要删除选中的角色吗"
               })
           }
        },
    };

    $(document).ready(function(obj){
        $("#h-role-info-table-details").bootstrapTable({
            url:"/auth/role/getall",
        });
    });
</script>

<script type="text/html" id="role_input_form">
    <form class="row" id="h-role-add-info">

        <div class="col-sm-12" style="margin-top: 12px;">
            <label class="h-label" style="width: 100%;">角色名称：</label>
            <input placeholder="1..30位汉字，字母，数字组成" type="text" class="form-control" name="name" style="width: 100%;height: 30px;line-height: 30px;">
        </div>
        <div class="col-sm-12" style="margin-top: 12px;">
            <label class="h-label" style="width: 100%;">状　态：</label>
            <select id="h-role-add-status" name="status"  class="form-control" style="width: 100%;height: 30px;line-height: 30px;">
                <option value="0">正常</option>
                <option value="1">失效</option>
            </select>
        </div>
    </form>
</script>

<script type="text/html" id="role_modify_form">
    <form class="row" id="h-role-modify-info">
        <div class="col-sm-12" style="margin-top: 12px;">
            <label class="h-label" style="width: 100%;">角色名称：</label>
            <input id="h-role-modify-role-name" placeholder="角色名称" type="text" class="form-control" name="name" style="width: 100%;height: 30px;line-height: 30px;">

        </div>
        <div class="col-sm-12" style="margin-top: 12px;">
            <label class="h-label" style="width: 100%;">状　态：</label>
            <select id="h-role-modify-role-status-cd" name="status"  class="form-control" style="width: 100%;height: 30px;line-height: 30px;">
                <option value="0">正常</option>
                <option value="1">失效</option>
            </select>
        </div>
    </form>
</script>