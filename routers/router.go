package routers

import (
	"github.com/astaxie/beego"
	"github.com/xueqianLu/nodemonitor/controllers"
)

func init() {
	ns := beego.NewNamespace("/nodemonitor",
		beego.NSNamespace("api",
			//用户信息
			beego.NSRouter("/nodeinfo", &controllers.Controller{}, "*:NodeInfo"),
			beego.NSRouter("/hpnodeinfo", &controllers.Controller{}, "get:HpNodeInfo"),
		),
	)
	beego.AddNamespace(ns)
}
