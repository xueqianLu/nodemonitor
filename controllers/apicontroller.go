package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xueqianLu/nodemonitor/async"
)

type Controller struct {
	BaseController
}

func (e *Controller) NodeInfo() {
	filter := make(map[string]string)

	if e.Ctx.Input != nil {

		logs.Info("get request body", "value", string(e.Ctx.Input.RequestBody))
		err := json.Unmarshal(e.Ctx.Input.RequestBody, &filter)
		if err != nil {
			logs.Error("parse request param failed", "err", err)
		}
	}
	var result = async.GetAllNodeStatus(filter)
	e.ResponseInfo(200, nil, result)
}

