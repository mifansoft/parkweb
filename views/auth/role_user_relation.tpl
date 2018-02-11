<div class="row subsystem-header">
    <div class="pull-left">
        <span style="font-size: 14px;">角色用户关联关系</span>
    </div>
    <div class="pull-right">  
        <span style="height: 30px; line-height: 30px; margin-top: 7px;display: inline"
              class="pull-left">角色名称 : {{.Name}}</span>
        <span id="h-role-resource-rel-role-id" class="hidden">{{.Id}}</span>
    </div>
</div>
<div class="row subsystem-toolbar">
    <div class="pull-right">
        &nbsp;
        {{if checkResIDAuth "2" "0105020610"}}
        <button onclick="UserRoleObj.add()" class="btn btn-info btn-sm">
            <i class="icon-plus"> 添加关联用户</i>
        </button>
        {{end}}
        {{if checkResIDAuth "2" "0105020620"}}
        <button onclick="UserRoleObj.delete()" class="btn btn-danger btn-sm">
            <i class="icon-trash"> 批量取消关联用户</i>
        </button>
        {{end}}
    </div>
</div>
<div class="subsystem-content">
    <table id="h-role-user-relation-table-details"
           class="table"
           data-toggle="table"
           data-unique-id="userId"
           data-side-pagination="client"
           data-click-to-select="true"
           data-pagination="true"
           data-striped="true"
           data-page-size="30"
           data-show-refresh="false"
           data-page-list="[20,30, 50, 100, 200]"
           data-search="false">
        <thead>
        <tr>
            <th data-field="state" data-checkbox="true"></th>
            <th data-field="account" data-sortable="true">账号</th>
            <th data-field="name" data-sortable="true">用户名称</th>
            <th data-field="orgCode" data-sortable="true">机构编码</th>
            <th data-field="orgName" data-sortable="true">机构描述</th>
            <th data-align="center" data-field="authUser" data-sortable="true">授权人</th>
            <th data-align="center" data-field="authTime" data-sortable="true">授权时间</th>
            <th data-align="center" data-formatter="UserRoleObj.formatter">操作</th>
        </tr>
        </thead>
    </table>
</div>
<script>
    $(document).ready(function(){
        $("#h-role-user-relation-table-details").bootstrapTable({
            method: "get",
            url:"/auth/role/user/getrelation",
            queryParams:function (params) {
                params.roleId = $("#h-role-resource-rel-role-id").html();
                return params;
            }
        })
    });
    
    var UserRoleObj = {
        formatter:function (value,row,index) {
            var html = "-";
            {{if checkResIDAuth "2" "0105020620"}}
                html = '<span class="h-td-btn btn-danger btn-xs" onclick="UserRoleObj.deleteRow(\''+ row.userId+'\')">取消关联</span>';
            {{end}}
            return html;
        },
        add:function () {
            $.Hmodal({
                header:"添加关联用户",
                body:$("#h-role-user-add-src").html(),
                preprocess:function () {
                    $("#h-role-user-relation-add").bootstrapTable({
                        queryParams:function (param) {
                            param.roleId = $("#h-role-resource-rel-role-id").html();
                            return param;
                        }
                    });
                },
                callback:function (hmode) {
                    var rows = $("#h-role-user-relation-add").bootstrapTable('getSelections');
                    if (rows.length == 0){
                        $.Notify({
                            message:"请选择需要关联的用户",
                            type:"warning",
                        });
                        return;
                    }

                    var arr = new Array();
                    $(rows).each(function (index, element) {
                        var e = {};
                        e.userId = element.userId;
                        e.roleId = $("#h-role-resource-rel-role-id").html();
                        arr.push(e);
                    });

                    $.HAjaxRequest({
                        url: "/auth/role/user/relationuser",
                        type: "post",
                        data: {dataJson:JSON.stringify(arr)},
                        success: function () {
                            $.Notify({
                                message: "角色关联新用户成功",
                                type: "success",
                            });
                            $(hmode).remove();
                            $("#h-role-user-relation-table-details").bootstrapTable('refresh');
                        },
                    })
                }
            })
        },
        delete:function () {
            var rows = $("#h-role-user-relation-table-details").bootstrapTable("getSelections");
            if (rows.length == 0){
                $.Notify({
                    message:"请在下表中选择需要取消关联的用户",
                    type:"warning",
                });
                return;
            }

            var arr = new Array();
            $(rows).each(function (index, element) {
                var e = {};
                e.userId = element.userId;
                e.roleId = $("#h-role-resource-rel-role-id").html();
                arr.push(e);
            });

            $.Hconfirm({
                body:"点击确定取消用户与角色关联关系",
                callback:function () {
                    $.HAjaxRequest({
                        url:"/auth/role/user/unrelationuser",
                        type:"post",
                        data:{dataJson:JSON.stringify(arr)},
                        success:function () {
                            $("#h-role-user-relation-table-details").bootstrapTable('refresh');
                        },
                    })
                }
            })
        },
        deleteRow:function (userId) {
            var arr = new Array();
            var e = {};
            e.userId = userId;
            e.roleId = $("#h-role-resource-rel-role-id").html();
            arr.push(e);
            $.HAjaxRequest({
                url:"/auth/role/user/unrelationuser",
                type:"post",
                data:{dataJson:JSON.stringify(arr)},
                success:function () {
                    $("#h-role-user-relation-table-details").bootstrapTable('refresh');
                },
            })
        }
    }
</script>

<script id="h-role-user-add-src" type="text/html">
    <table id="h-role-user-relation-add"
           class="table"
           data-toggle="table"
           data-unique-id="userId"
           data-side-pagination="client"
           data-url="/auth/role/user/getrelationother"
           data-click-to-select="true"
           data-pagination="false"
           data-striped="true"
           data-show-refresh="false"
           data-page-list="[20, 50, 100, 200]"
           data-search="false">
        <thead>
        <tr>
            <th data-field="state" data-checkbox="true"></th>
            <th data-field="userId" class="hidden">id</th>
            <th data-field="account" data-sortable="true">账号</th>
            <th data-field="name" data-sortable="true">用户名称</th>
            <th data-field="orgCode" data-sortable="true">机构编码</th>
            <th data-field="orgName" data-sortable="true">机构描述</th>
        </tr>
        </thead>
    </table>
</script>