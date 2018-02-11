<div class="row subsystem-header">
    <div class="pull-left">
        <span style="font-size: 14px;">角色资源信息</span>
    </div>
</div>
<div class="row subsystem-toolbar">
    <div class="pull-right">
        <span style="height:30px; line-height:30px; margin-top:10px;display:inline">角色名称 : {{.Name}}</span>
        <span id="h-role-resource-rel-role-id" class="hidden">{{.Id}}</span>
    </div>
</div>
<div class="subsystem-content" style="background-image: url('/static/images/hauth/pure_book.png');filter:'progid:DXImageTransform.Microsoft.AlphaImageLoader(sizingMethod='scale')';-moz-background-size:100% 100%;background-size:100% 100%;">
    <div id="h-domain-share-info-details" class="row">
        <div class="col-sm-6 col-md-6 col-lg-6" style="padding-left: 10%; padding-right: 2%;">
            <div class="col-ms-12 col-md-12 col-lg-12" style="margin-top: 3%">
                <div style="border-bottom: #598f56 solid 1px;height: 44px; line-height: 44px;">
                    <div class="pull-left">
                        <span><i class="icon-sitemap"> </i>已经被授权获取菜单资源:</span>
                    </div>
                    <div class="pull-right">
                        &nbsp;
                        {{if checkResIDAuth "2" "0105020520"}}
                        <button onclick="RoleResObj.revoke()" class="btn btn btn-danger btn-sm"><i class="icon-remove-circle"></i>&nbsp;撤销
                        </button>
                        {{end}}
                    </div>
                </div>
            </div>
            <div id="h-role-res-owner-resource" class="col-sm-12" style="overflow: auto">
            </div>
        </div>
        <div class="col-sm-6 col-md-6 col-lg-6" style="padding-left: 2%;padding-right: 10%;">
            <div class="col-ms-12 col-md-12 col-lg-12" style="margin-top: 3%;">
                <div style="border-bottom: #8f1121 solid 1px;height: 44px; line-height: 44px;">
                    <div class="pull-left">
                        <span><i class="icon-sitemap"> </i>尚未被授权菜单资源:</span>
                    </div>
                    <div class="pull-right">
                        &nbsp;
                        {{if checkResIDAuth "2" "0105020520"}}
                        <button onclick="RoleResObj.auth()" class="btn btn btn-info btn-sm"><i class="icon-plus-sign"></i>&nbsp;授权
                        </button>
                        {{end}}
                    </div>
                </div>
            </div>
            <div id="h-role-res-other-resource" class="col-sm-12" style="overflow: auto;">
            </div>
        </div>
    </div>
</div>
<script>
    var RoleResObj = {
        resource_self:function(){
            debugger;
            var roleId = $("#h-role-resource-rel-role-id").html()
            $.HAjaxRequest({
                url:"/auth/role/res/getroleres",
                type:"get",
                data:{typeId:0,roleId:roleId},
                success:function (data) {
                    var arr = new Array()
                    $(data).each(function (index, element) {
                        var ijs = {}
                        ijs.id=element.resId
                        ijs.text = element.name
                        ijs.parentId = element.parentId
                        arr.push(ijs)
                    });
                    $("#h-role-res-owner-resource").HtreeWithLine({
                        data:arr,
                        checkbox:true,
                    });
                },
            })
        },
        resource_other:function(){
            debugger;
            var roleId = $("#h-role-resource-rel-role-id").html()
            $.HAjaxRequest({
                url:"/auth/role/res/getroleres",
                type:"get",
                data:{typeId:1,roleId:roleId},
                success:function (data) {
                    var arr = new Array()
                    $(data).each(function (index, element) {
                        var ijs = {}
                        ijs.id=element.resId
                        ijs.text = element.name
                        ijs.parentId = element.parentId
                        arr.push(ijs)
                    });
                    $("#h-role-res-other-resource").HtreeWithLine({
                        data:arr,
                        checkbox:true,
                    })
                },
            })
        },
        revoke:function(){
            var roleId = $("#h-role-resource-rel-role-id").html();
            var resArr = new Array();
            $("#h-role-res-owner-resource").find("input[type='checkbox']").each(function (index, element) {
                if ($(element).is(":checked")) {
                    resArr.push(element.value);
                }
            });

            if (resArr.length == 0){
                $.Notify({
                    title:"温馨提示：",
                    message:"请在下列树形结构中选择需要撤销的资源",
                    type:"info",
                });
                return
            }
            $.Hconfirm({
                body:"点击确定删除角色拥有的资源",
                callback:function () {
                    $.HAjaxRequest({
                        url:"/auth/role/res/unauthorized",
                        type:"post",
                        dataType:"json",
                        data:{roleId:roleId,dataJson:JSON.stringify(resArr)},
                        success:function () {
                            $.Notify({
                                title:"操作成功",
                                message:"撤销资源权限成功",
                                type:"success",
                            });
                            RoleResObj.resource_self();
                            RoleResObj.resource_other();
                        },
                    })
                }
            })

        },
        auth:function(){
            var roleId = $("#h-role-resource-rel-role-id").html();
            var resArr = new Array();
            $("#h-role-res-other-resource").find("input[type='checkbox']").each(function (index, element) {
                if ($(element).is(":checked")) {
                    resArr.push(element.value);
                }
            });

            if (resArr.length == 0){
                $.Notify({
                    title:"温馨提示：",
                    message:"请在下列树形结构中选择需要撤销的资源",
                    type:"info",
                });
                return
            }
            $.Hconfirm({
                body:"点击确定给角色授予菜单资源",
                callback:function () {
                    $.HAjaxRequest({
                        url:"/auth/role/res/authorized",
                        type:"post",
                        data:{roleId:roleId,dataJson:JSON.stringify(resArr)},
                        success:function () {
                            $.Notify({
                                title:"操作成功",
                                message:"授权资源成功",
                                type:"success",
                            });
                            RoleResObj.resource_self()
                            RoleResObj.resource_other()
                        },
                    });
                }
            })
        },
    }

    $(document).ready(function(){
        debugger;
        var hwindow = document.documentElement.clientHeight;
        $("#h-domain-info-shareid").height(hwindow-130);
        $("#h-domain-share-info-details").height(hwindow-130);
        $("#h-role-getted-resource-info").height(hwindow-160);
        $("#h-role-ungetted-resource-info").height(hwindow-160);
        $("#h-role-res-owner-resource").height(hwindow-210);
        $("#h-role-res-other-resource").height(hwindow-210);
        RoleResObj.resource_self();
        RoleResObj.resource_other();
    });

</script>