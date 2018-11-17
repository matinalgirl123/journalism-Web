package routers

import (
	"beetest2/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//设置路由过滤器 防止用户未登录直接访问页面
	beego.InsertFilter("/atc/*",beego.BeforeExec,FuncFilter)

	//路由
    beego.Router("/", &controllers.MainController{})
	//注册界面
    beego.Router("/register", &controllers.UserController{},"get:ShowReg;post:HideReg")
	//登陆界面
    beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:HideLogin")
	//新闻首界面
	beego.Router("/atc/acticle", &controllers.ActicleController{},"get:ShowActicle")
	//新闻首界面 查看
	beego.Router("/atc/acticleread", &controllers.ActicleController{},"get:ShowActicleRead")
	//新闻首界面 编辑
	beego.Router("/atc/acticleupdate", &controllers.ActicleController{},"get:ShowUpdate;post:HideUpdate")
	//添加文章面
	beego.Router("/atc/add", &controllers.ActicleController{},"get:ShowAdd;post:HideAdd")
	//新闻首界面 删除
	beego.Router("/atc/acticledel", &controllers.ActicleController{},"get:ShowActicleDel")
	//分类界面 删除
	beego.Router("/atc/acticleremove", &controllers.ActicleController{},"get:ShowActicleRemove")
	//添加分类界面
	beego.Router("/atc/addType", &controllers.ActicleController{},"get:ShowAddType;post:HideAddType")
	//新闻首界面 退出
	beego.Router("/atc/acticlerout", &controllers.UserController{},"get:ShowActicleOut")


	}

	//路由控制器的函数获取 session 内存储的数据    登录校验
	var FuncFilter=func(cxt *context.Context){

		userName:=cxt.Input.Session("userName")
		if userName==nil{
			message:="未正确登录"
			cxt.Redirect(302,"/login?message="+message,)
			return
			}
	}


/*

//插入数据库
o:=orm.NewOrm()
var user2 models.User2
user2.UserName=userName
user2.Pwd=pwd
err:=o.Read(&user2,"userName")
if err!=nil{
message:="数据输入有误"
this.Redirect("/register?message="+message,302)
return
}
*/
