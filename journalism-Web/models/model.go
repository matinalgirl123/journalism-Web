package models

import(
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

//存储用户登陆数据   N个用户对N个表
type User2 struct {
	Id int `orm:"pk;auto"`
	UserName string `orm:"unique"`
	Pwd string
	Acticle2 []*Acticle2 `orm:"rel(m2m)"`

}


//存储添加的新闻信息    N个类型对1个表
type Acticle2 struct {
	Id int `orm:"pk;auto"`
	TitleName string `orm:"size(100)"`
	Time time.Time `orm:"auto_now"`
	ReadNum int `orm:"dafault(0) null"`
	Content string  `orm:"size(500)"`
	Aimg string `orm:"size(100)"`
	User2 []*User2 `orm:"reverse(many)"`
	ActicleType2 *ActicleType2 `orm:"null;rel(fk);on_delete(do_nothing)"`

}

//新闻类型
type ActicleType2 struct {
	Id int `orm:"pk;auto"`
	Acticletype string `orm:"unique;size(100)"`
	Acticle2 []*Acticle2 `orm:"reverse(many)"`

}

func init(){
	//连接服务器
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/mytest01?charset=utf8")

	//注册表
	orm.RegisterModel(new(User2),new(Acticle2),new(ActicleType2))

	//生成表
	orm.RunSyncdb("default",false,true)

	}