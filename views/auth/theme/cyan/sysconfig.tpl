<div id="wrapper" class="navbar-static-side mifan-theme-background" style="margin: 0px;width: 100%;">
    <div class="H-logo-area" style="margin: 0px; padding: 0 0 0 30px;">
        系统配置管理
    </div>
    <div class="col-sm-12 col-md-12 col-lg-12" id="H-main-menu" style="margin-bottom: 60px; overflow: auto;">
        <div id="h-system-service" class="col-sm-12 col-md-6 col-lg-4">
        </div>
        <div id="h-mas-service" class="col-sm-12 col-md-6 col-lg-4">
        </div>
        <div id="h-other-service"  class="col-sm-12 col-md-6 col-lg-4">
        </div>
    </div>
</div>

<div id="page-wrapper" class="gray-bg col-sm-12 col-md-12 col-lg-12"
     style="margin:0px;padding: 0px;display: none;">
    <div id="h-main-content"
         style="padding: 0px; margin: 0px;position: relative; overflow: auto;">
    </div>
</div>

<script type="text/javascript">
    NProgress.start();
    /*
     * 调整页面宽度和高度
     * */
    $(document).ready(function(){
        var hwindow = document.documentElement.clientHeight;
        $("#wrapper").height(hwindow);
        $("#page-wrapper").height(hwindow-36);
        $("#H-main-menu").height(hwindow-96);
        $("#h-main-content").height(hwindow-36);
    });

    $(document).ready(function(){
        Hutils.initMenu(1,'0100000000',"资源管理","用户与权限","系统审计");
        $("#page-wrapper").show();
        NProgress.done();
    });

    window.onresize = function(){
        var hh = document.documentElement.clientHeight;
        $("#wrapper").height(hh);
        $("#page-wrapper").height(hh-36);
        $("#H-main-menu").height(hh-96);
        $("#h-main-content").height(hh-36);
    }
</script>