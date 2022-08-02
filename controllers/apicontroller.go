package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xueqianLu/nodemonitor/async"
	"strconv"
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

func (e *Controller) HpNodeInfo() {
	filter := make(map[string]string)
	filter["nodetype"] = "hpnode"
	var result = async.GetAllNodeStatus(filter)
	e.ResponseInfo(200, nil, result)
}

func (e *Controller) BlockLoseInfo() {
	block, _ := strconv.Atoi(e.Ctx.Input.Query("number"))
	round, _ := strconv.Atoi(e.Ctx.Input.Query("round"))
	var result = async.GetBlockLoseInfo(int64(block), round)
	e.ResponseInfo(200, nil, result)
}
