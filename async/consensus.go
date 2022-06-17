package async

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xueqianLu/nodemonitor/ethrpc"
	"strings"
	"sync"
)

var (
	nodestype = sync.Map{}
	ntmux = sync.RWMutex{}
)

type RequestHPNodesInfo struct {
	Coinbase []string
}

type RequestCadNodesInfo struct {
	Coinbase []string `json:"cadaddresses"`
	Number int64 	`json:"number"`
}


func getHpNodes() {
	var hpinfo = []string{}
	var cadinfo = RequestCadNodesInfo{}
	{
		hpnodes, err := ethrpc.HpbNodes()
		if err != nil {
			logs.Error("got hpbnodes from rpc failed", "err", err)
			return
		}

		err = json.Unmarshal(hpnodes.(json.RawMessage), &hpinfo)
		if err != nil {
			logs.Error("unmarshal hpbnodes info failed", "err", err)
			return
		}

		for _, coinbase := range hpinfo {
			logs.Debug("got coinbase is hpnode", coinbase)
		}
	}
	{
		hpnodes, err := ethrpc.HpbCadNodes()
		if err != nil {
			logs.Error("got HpbCadNodes from rpc failed", "err", err)
			return
		}

		err = json.Unmarshal(hpnodes.(json.RawMessage), &cadinfo)
		if err != nil {
			logs.Error("unmarshal HpbCadNodes info failed", "err", err)
			return
		}

		for _, coinbase := range cadinfo.Coinbase {
			logs.Debug("got coinbase is cadnode", coinbase)
		}
	}
	ntmux.Lock()
	defer ntmux.Unlock()
	nodestype = sync.Map{}
	for _, hp := range hpinfo {
		nodestype.Store(strings.ToLower(hp), "hpnode")
	}
	for _, cad := range cadinfo.Coinbase {
		nodestype.Store(strings.ToLower(cad), "prenode")
	}
}

func getTypeByCoinbase(coinbase string) string {
	ntmux.RLock()
	defer ntmux.RUnlock()
	v, exist := nodestype.Load(strings.ToLower(coinbase))
	if !exist {
		return "syncnode"
	}
	return v.(string)
}