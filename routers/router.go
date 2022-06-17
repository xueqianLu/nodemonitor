package routers

import (
	"github.com/astaxie/beego"
	"github.com/xueqianLu/nodemonitor/controllers"
)

func init() {
	ns := beego.NewNamespace("/monitor",
		beego.NSNamespace("api",
			//用户信息
			beego.NSRouter("/nodeinfo", &controllers.Controller{}, "post:NodeInfo"),
		),
	)
	beego.AddNamespace(ns)
}
