package async

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var (
	lockedNodeInfo = sync.Map{} //make(map[string]LockedNodeInfo)
	mux = sync.RWMutex{}
)

type LockedNodeInfo struct {
	Name string `json:"name_eng"`
	Coinbase string `json:"coinbase"`
	VoteNumber int64 `json:"voteNumber"`
}

type requestNodeInfos struct {
	ErrCode int `json:"error_code"`
	Msg     string `json:"error_message"`
	Infos   []*LockedNodeInfo `json:"data"`
}


func getNodeInfoAndParse() {
	url := "https://vote.hpbnode.com/v1/votelist/"
	var info = &requestNodeInfos{}
	resp, err := http.Get(url)
	if err != nil {
		logs.Error("get node info from vote.hpbnode.com failed", "err", err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//Failed to read response.
		logs.Error("failed to read node info from response", "err", err)
		return
	}
	json.Unmarshal(data, info)
	if info.ErrCode != 0 {
		logs.Error("get node info from vote.hpb.node got error", "err", info.Msg)
		return
	}
	mux.Lock()
	defer mux.Unlock()
	lockedNodeInfo = sync.Map{} //make(map[string]LockedNodeInfo)
	for _, d := range info.Infos {
		logs.Debug("get node info", "info", d)
		lockedNodeInfo.Store(strings.ToLower(d.Coinbase), d)
	}
}

func getBoeNodeList() []*LockedNodeInfo {
	var boelist = make([]*LockedNodeInfo, 0)
	mux.RLock()
	defer mux.RUnlock()
	lockedNodeInfo.Range(func (k, v interface{}) bool {
		info := v.(*LockedNodeInfo)
		boelist = append(boelist, info)
		return true
	})
	return boelist
}
