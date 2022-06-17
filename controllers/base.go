package controllers
import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

func (d *BaseController) ResponseInfo(code int, errMsg interface{}, result interface{}) {
	switch code {
	case 500:
		d.Data["json"] = map[string]interface{}{"error": "500", "err_msg": errMsg, "data": result}
	case 200:
		d.Data["json"] = map[string]interface{}{"error": "200", "err_msg": errMsg, "data": result}
	}
	d.ServeJSON()
}

