package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"beetest2/models"
	"math"
	"encoding/base64"
)

//控制器定义
type ActicleController struct {
	beego.Controller
}

//页面控制器实现  文章列表
func(this *ActicleController)ShowActicle(){

/*
	// 复杂版登录校验 获取 session 值
	userName:=this.GetSession("userName")
	if userName==nil{
		message:="未正确登录"
		this.Redirect("/login?message="+message,302)
		return
	}*/


	//如其他页面出错该页面被调用并传递错误数据
	message:=this.GetString("message")
	if message!=""{
		this.Data["message"]=message
		this.Layout="layout.html"
		this.TplName="index.html"
		return
	}



	//创建对象
	o:=orm.NewOrm()
	//新闻数据存储
	var acticle2 []models.Acticle2
	//类型数据存储
	var acticleType2 []models.ActicleType2

	//类型文件并展示
	o.QueryTable("ActicleType2").All(&acticleType2)
	this.Data["acticleType2"]=acticleType2


	//获取页面id /*避免用户url输入*/
	pageIndex,err:=this.GetInt("Id")
	if err!=nil{
		pageIndex=1
	}

	//分页显示
	//记录页数
	var Count int64
	//单页起始显示
	pageSize:=int64(2)
	//获取分页数据   当前数据为 2页 3，4 记录
	start:=(int64(pageIndex)-1)*pageSize


	//下拉框选择所看新闻类型
	//获取用户选择类型数据
	list:=this.GetString("select")

	if  list==""{

		//查询全部类型文章并显示
		o.QueryTable("Acticle2").Limit(pageSize,start).All(&acticle2)

		Count,_=o.QueryTable("Acticle2").Count()

	}else{

	   //查询选定类型文章并显示
		o.QueryTable("Acticle2").Limit(pageSize,start).RelatedSel("ActicleType2").Filter("ActicleType2__Acticletype",list).All(&acticle2)

		Count,_=o.QueryTable("Acticle2").RelatedSel("ActicleType2").Filter("ActicleType2__Acticletype",list).Count()
	}



	//分页数 向上取整
	mean:=math.Ceil(float64(Count)/2)


	//数据记录数
	this.Data["Count"]=Count
	//页数
	this.Data["mean"]=mean
	//当前页数
	this.Data["pageIndex"]=pageIndex
	//选定类型
	this.Data["list"]=list
	//查询到的数据传输给首页
	this.Data["acticle2"]=acticle2

	//获取登录用户名 作为登录后当前用户显示
	Name:=this.Ctx.GetCookie("userName")
	//浏览器传值的时 中文无法传递   解密
	enc,_:=base64.StdEncoding.DecodeString(Name)
	//用户名传输
	this.Data["userName"]=string(enc)


	//控制器代码
	this.Layout="layout.html"

	this.TplName="index.html"

}

//页面控制器实现   查看文章
func(this *ActicleController)ShowActicleRead(){

	//浏览器获取数据
	Id,err:=this.GetInt("Id")

	//校验数据
	if err!=nil{
		message:="请指定正确的文章"
		this.Redirect("/atc/acticle?message="+message,302)
		return
	}

	//处理数据
	o:=orm.NewOrm()
	var acticle2 models.Acticle2
	acticle2.Id=Id

	err=o.Read(&acticle2)
	if err!=nil{
		message:="该文章不存在"
		this.Redirect("/atc/acticle?message="+message,302)
		return
	}

	//  获取对应文章类型
	var acticleType2  models.ActicleType2
	//使用高级查询       查询数据库哪张表                where   数据库表的字段 = 已知的字段    数据传递给  &对应变量
	o.QueryTable("ActicleType2").Filter("Id",acticle2.ActicleType2.Id).All(&acticleType2)


	//----------查阅用户功能实现------------------
	//获取服务器存储的当前用户名
	userName:=this.GetSession("userName")

	//获取读取此文章的用户名
	//多对多操作对象      --------插入到 对象 的 那个字段    1
	m2m:=o.QueryM2M(&acticle2,"User2")

	//创建用户对象     --------插入什么数据   		2
	var use models.User2
	//判断数据    接口断言
	use.UserName=userName.(string)
	//查询该用户
	err=o.Read(&use,"UserName")
	if err!=nil{
		message:="当前用户不存在"
		this.Redirect("/atc/acticle?message="+message,302)
		return
	}

	beego.Error("添加阅读者测试1")

	//插入多对多关系
	//将当前用户名插入表阅读用户
	//m2m为我们查询表中阅读的用户对象
	// 利用对象去插入我们从浏览器获取的用户名  --------插入数据   3
	m2m.Add(use)


	beego.Error("添加阅读者测试2")

	//---------------查询表表内存储的用户对象 --------------


	/*
	//第一种 多对多查询  绑定  加载    直接在需要显示的字段去遍历数据
	o.LoadRelated(&acticle2,"User2")
	*/

	//第二种 多对多查询   正向插入，反向查询
	var  User2 []models.User2  //反向取出数据  select *form User2 (where 表id= Id)[目的确定查看的当前文章】
	o.QueryTable("User2").Filter("Acticle2__Acticle2__Id",Id).Distinct().All(&User2)



	//---------------查询多少用户查看该文章 --------------

	//LoadRelated 用于多对多查询是正向查询  绑定字段  加载   记录了每次查阅该文章的用户
	//但返回值为 记录了每次查阅该文章的用户 数量    适用于记录统计查阅量
	num,_:=o.LoadRelated(&acticle2,"User2")


	acticle2.ReadNum=int(num)
	//返回数据
	this.Data["acticle2"]=acticle2
	this.Data["acticleType2"]=acticleType2
	this.Data["User2"]=User2

	//页面组合
	this.Layout="layout.html"
	this.TplName="content.html"
}

//删除文章
func(this *ActicleController)ShowActicleDel(){
	//浏览器获取数据
	Id,err:=this.GetInt("Id")

	//校验数据
	if err!=nil{
		message:="请指定删除正确的文章"
		this.Redirect("/atc/acticle?message="+message,302)
		return
	}

	//处理数据
	o:=orm.NewOrm()
	var acticle2 models.Acticle2
	acticle2.Id=Id

	err=o.Read(&acticle2)
	if err!=nil{
		message:="该文章不存在"
		this.Redirect("/atc/acticle?message="+message,302)
		return
	}

		//删除文章
	_,err=o.Delete(&acticle2)
	if err!=nil{
		message:="该文章不存在"
		this.Redirect("/atc/acticle?message="+message,302)
		return
	}

			//页面组合
		this.Redirect("/atc/acticle",302)
}

//编辑文章
func(this *ActicleController)ShowUpdate(){

	//由主界面获取Id传递过来  用于展示更新数据界面
	Id,_:=this.GetInt("Id")
	this.Data["acticle2Id"]=Id

	//插入数据库
	o:=orm.NewOrm()
	var acticle2 models.Acticle2
	acticle2.Id=Id

	//数据校验
	err:=o.Read(&acticle2)
	if err!=nil{
		message:="更新页面有误"
		this.Redirect("/register?message="+message,302)
		return
	}

	this.Data["acticle2"]=acticle2
	this.Layout="layout.html"
	this.TplName="update.html"


}

//编辑文章 更新
func(this *ActicleController)HideUpdate() {

	//获取数据
	articleId,_:=this.GetInt("Id")
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	//图片处理
	fileName:=FileGain(this,"uploadname")
	//校验数据
	if articleId==0 ||  articleName=="" ||   content=="" ||   fileName=="" {

		beego.Error(articleId)
		beego.Error(articleName)
		beego.Error(content)
		beego.Error(fileName)

		message := "该更新文章不存在"
		this.Redirect("/atc/acticle?message="+message, 302)
		return
	}

	o:=orm.NewOrm()
	var acticle2 models.Acticle2
	acticle2.Id = articleId

	err:=o.Read(&acticle2)
	if err != nil {
		message := "该更新文章不存在"
		this.Redirect("/atc/acticle?message="+message, 302)
		return
	}

	//处理数据
	acticle2.TitleName= articleName
	acticle2.Aimg=fileName
	acticle2.Content=content
	_,err=o.Update(&acticle2)
	if err != nil {
		message := "文章更新失败"
		beego.Error(err)
		this.Redirect("/atc/acticle?message="+message, 302)
		return
	}

	//返回数据
	this.Redirect("/atc/acticle", 302)

}

// 添加文章实现
func(this *ActicleController)ShowAdd(){

	err:=this.GetString("message")
	if err!=""{
		this.Data["message"]=err
		this.Layout="layout.html"
		this.TplName="add.html"
		return
	}

	//查询文章所有类型并进行展示
	o:=orm.NewOrm()
	var acticleType2  []models.ActicleType2
	o.QueryTable("ActicleType2").All(&acticleType2)

	//查询到的列表数据传输给首页
	this.Data["acticleType2"]=acticleType2

	//控制器代码
	this.Layout="layout.html"

	this.TplName="add.html"
	beego.Error("测试11")
}

//添加文章
func(this *ActicleController)HideAdd(){

	beego.Error("添加文章测试1")
	//获取页面数据
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	select1:=this.GetString("select")

	//调用函数 处理图片数据 =返回图片存储路径
	filename:=FileGain(this,"uploadname")

	beego.Error("添加文章测试2")
	//插入数据库
	//获取对象
	o:=orm.NewOrm()
	var acticleb2 models.Acticle2

	beego.Error("filename",filename)

	//判断数据为空
	if  articleName=="" || content=="" || filename=="" ||select1=="" {
		message:="数据为空"
		this.Redirect("/atc/add?message="+message,302)
		return

	}
	beego.Error("添加文章测试3")

	//对象赋值
	acticleb2.TitleName=articleName
	acticleb2.Content=content
	acticleb2.Aimg=filename

	//添加文章类型插入
	var acticleType2 models.ActicleType2
	acticleType2.Acticletype=select1
	err:=o.Read(&acticleType2,"Acticletype")
	if err!=nil{
		errmsg:="数据类型查询失败"
		beego.Error(err)
		this.Redirect("/atc/add?errmsg="+errmsg,302)
		beego.Error("添加文章测试失败跳转")
		return
	}

	//结构体指针赋值 &赋值
	acticleb2.ActicleType2=&acticleType2

	beego.Error("添加文章测试4")

	//数据插入
	_,err=o.Insert(&acticleb2)
	if err!=nil{
		errmsg:="数据插入失败"
		beego.Error(err)
		this.Redirect("/atc/add?errmsg="+errmsg,302)
		beego.Error("添加文章测试失败跳转")
		return
	}
	beego.Error("添加文章测试5")

	this.Redirect("/atc/acticle",302)

}

//添加分类实现   添加分类
func(this *ActicleController)ShowAddType(){

	o:=orm.NewOrm()

	var acticleType2 []models.ActicleType2

	o.QueryTable("ActicleType2").All(&acticleType2)

   this.Data["acticleType2"]=acticleType2
	//控制器代码
	this.Layout="layout.html"

	this.TplName="addType.html"
}

//添加分类实现   添加分类 处理
func(this *ActicleController)HideAddType(){

	beego.Error("添加分类1")

	typeName:=this.GetString("typeName")
	if typeName==""{
		this.Data["message"]="数据为空"
		this.Layout="layout.html"
		this.TplName="addType.html"
		return
	}
	beego.Error(typeName)
	beego.Error("添加分类2")
	//插入数据库
	o:=orm.NewOrm()

	var acticleType2 models.ActicleType2

	acticleType2.Acticletype=typeName

	beego.Error("添加分类3")


	_,err:=o.Insert(&acticleType2)
	if err!=nil{
		this.Data["message"]="数据插入失败"
		this.Layout="layout.html"
		this.TplName="addType.html"

		return
	}
	beego.Error(acticleType2.Acticletype)

	beego.Error("添加分类4")

	this.Data["typeName"]=acticleType2
	//控制器代码
	this.Redirect("/atc/addType",302)
}

//添加分类实现   删除分类
func(this *ActicleController)ShowActicleRemove(){

	id,err:=this.GetInt("Id")
	if err!=nil{
		this.Data["message"]="请插入正确的删除对象"
		this.Layout="layout.html"
		this.TplName="addType.html"
		return
	}

	o:=orm.NewOrm()

	var acticleType2 models.ActicleType2

	acticleType2.Id=id
	err=o.Read(&acticleType2)
	if err!=nil{
		this.Data["message"]="请插入正确的删除对象"
		this.Layout="layout.html"
		this.TplName="addType.html"
		return
	}
	_,err=o.Delete(&acticleType2)
	if err!=nil{
		this.Data["message"]="删除失败"
		this.Layout="layout.html"
		this.TplName="addType.html"
		return
	}


	this.Redirect("/atc/addType",302)

}

//图片数据处理存储
func FileGain(this *ActicleController,uploadname string)string{

	beego.Error("图片数据测试1")
	//获取网页数据
	//文件流,文件头[size,格式,名],错误
	file,head,err:=this.GetFile(uploadname)

	//校验数据
	if err!=nil{
		this.Data["message"]="获取数据方式err"
		return ""
	}

	defer file.Close()

	beego.Error("图片数据测试2")
	//格式
	filesuffix:=path.Ext(head.Filename)

	if  filesuffix!=".jpg" && filesuffix!=".png" {

		this.Data["message"]="获取数据格式err"
		return ""
	}

	beego.Error("图片数据测试3")
	//大小  KB
	if head.Size>100000{
		beego.Error("图片数据测试size")
		this.Data["message"]="数据大小err"
		return ""
	}

	//文件名防止重名
	//文件创建 [预防图片名相似 造成的覆盖]
	//以当前时间转为字符串 + 后缀名 =文件名
	fileNema:=time.Now().Format("2006-01-02-15-04-05")+filesuffix


	//处理数据
	//存储函数    网页上获取的数据    存储本地路径
	this.SaveToFile(uploadname,"./static/img/"+fileNema)

	//数据库存储路径
	filepath:= "/static/img/"+fileNema

	beego.Error("图片数据ok")
	beego.Error(filepath)
     return  filepath
}
