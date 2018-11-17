package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"beetest2/models"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

//登陆注册控制器实现

// 注册页面
// http://192.168.6.108:8083/register
func (this *UserController)ShowReg(){

	//判断post get 请求方法内是否有此数据 如果有就返回给页面
	message:=this.GetString("message")
	if message!=""{
		this.Data["message"]=message
	}

	this.TplName="register.html"
}

//获取注册页面数据
func (this *UserController)HideReg() {


	//获取数据
	userName:=this.GetString("userName")
	pwd:=this.GetString("password")

	//校验数据   如出错则渲染该界面并提示错误
	if userName=="" || pwd==""{
		this.Data["message"]="数据写入有误"
		this.TplName="register.html"
		return
	}

	//插入数据库
	o:=orm.NewOrm()
	var user2 models.User2
	user2.UserName=userName
	user2.Pwd=pwd
	_,err:=o.Insert(&user2)
	if err!=nil{
		this.Data["message"]="数据有误或已经存在"
		this.TplName="register.html"
		return
	}
	//跳转页面
	this.Redirect("/login",302)

}


// 登陆页面
// http://192.168.6.108:8083/login
func (this *UserController)ShowLogin(){


	//获取用户上次登录时是否确定用户名
	Name:=this.Ctx.GetCookie("userName")

	//浏览器传值的时 中文无法传递   解密
	enc,_:=base64.StdEncoding.DecodeString(Name)

	if string(enc)!=""{
		this.Data["userName"]=string(enc)  //界面直接显示上次输入的数据
		this.Data["checkbox"]="checked"  //界面直接勾选记住密码   checked 为勾选
	}else{
		this.Data["userName"]=""
		this.Data["checkbox"]=""
	}

	//如其他页面出错该页面被调用并传递错误数据
	message:=this.GetString("message")
	if message!=""{
		this.Data["message"]=message
		this.TplName="login.html"
		return
	}else {
		this.TplName = "login.html"
	}
}

//获取登陆页面数据
func (this *UserController)HideLogin() {


	//获取数据
	userName:=this.GetString("userName")
	pwd:=this.GetString("password")

	//校验数据   如出错就跳转至注册页面 并将错误数据发送给跳转页面
	if userName=="" || pwd==""{
		message:="数据输入有误"
		this.Redirect("/register?message="+message,302)
		return
	}

	//插入数据库
	o:=orm.NewOrm()
	var user2 models.User2
	user2.UserName=userName
	user2.Pwd=pwd
	err:=o.Read(&user2,"userName")
	if err!=nil{
		message:="该用户未注册成功"
		this.Redirect("/register?message="+message,302)
		return
	}
	//信息匹配
	if user2.UserName!=userName || user2.Pwd!=pwd{
		message:="信息输入有误"
		this.Redirect("/register?message="+message,302)
		return
	}

	//记住用户名 this.Ctx.SetCookie()
	remember:=this.GetString("remember")
	//beego.Error(remember) //记住用户名 为  v：on
	if remember=="on"{
		beego.Error(remember)

		//浏览器传值的时 中文无法传递   加密
		Enc:=base64.StdEncoding.EncodeToString([]byte(userName))

		this.Ctx.SetCookie("userName",Enc,3600*1)
	}else{
		beego.Error(11)
		//浏览器传值的时 中文无法传递   加密
		userName=base64.StdEncoding.EncodeToString([]byte(userName))

		this.Ctx.SetCookie("userName",userName,-1)
	}

	//添加 session  并赋值
	this.SetSession("userName",userName)

	//跳转页面
	this.Redirect("/atc/acticle",302)

}

//新闻首界面 退出
func(this *UserController)ShowActicleOut(){

	//删除校验数据
	this.DelSession("userName")
	//退出跳转页面
	this.Redirect("/login",302)
}

