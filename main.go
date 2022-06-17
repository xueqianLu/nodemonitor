package main

import (
	"github.com/astaxie/beego"
	_ "github.com/xueqianLu/nodemonitor/ethrpc"
	_ "github.com/xueqianLu/nodemonitor/routers"
	_ "github.com/xueqianLu/nodemonitor/async"
)

func main() {
	beego.Run()
}
