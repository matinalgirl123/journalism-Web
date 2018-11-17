package main

import (
	_ "beetest2/routers"
	"github.com/astaxie/beego"
	_"beetest2/models"
)

func main() {
	beego.AddFuncMap("PagingTop",PagingTop)
	beego.AddFuncMap("PagingDown",PagingDown)
	beego.Run()
}


func PagingTop( Num int)int {

	if Num<=1{
		return 1
	}
	return  Num-1
}

func PagingDown( Num int,mean float64)int {

	if Num+1>int(mean){
		return int(mean)
	}
	return  Num+1
}



