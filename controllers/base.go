package controllers
import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

func (d *BaseController) ResponseInfo(code int, errMsg interface{}, result interface{}) {
	switch code {
	case 500:
		d.Data["json"] = map[string]interface{}{"code":code, "cnMsg": errMsg, "enMsg": errMsg, "data": result}
	case 200:
		d.Data["json"] = map[string]interface{}{"code":code, "cnMsg": "处理成功", "enMsg": "Success", "data": result}
	}
	d.ServeJSON()
}

