package async

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var (
	peerInfos = sync.Map{} // map[peer.id]peerinfo
	peerMux = sync.RWMutex{}
)

func addPeerInfoToMap(info *PeerInfo) {
	peerInfos.Store(info.PeerId, info)
}



/*
{
      "id": "f432655fb5adcd10",
      "name": "",
      "version": "[1.0.9.2-stable/0]\u0026[N.A/0]",
      "coinbase": "0x8294CdfED3B7645154a02429a5F94948624fdC2b",
      "remote": "SynNode",
      "cap": "[hpb/100]",
      "network": {
        "local": "172.25.94.95:30308",
        "remote": "154.86.159.8:34172"
      },
      "start": "2022-06-17 12:16:48.889640035 +0800 CST m=+71.679437819",
      "beat": "1491",
      "mining": "",
      "hpb": {
        "handshakeTD": 1,
        "handshakeHD": "0x9c3704fab3915e36ef6f1e6353167c93ccf2486a4a2854dcaaa944066cf39966"
      }
    }
 */
type PeerInfo struct {
	PeerId string	`json:"id"`
	Name string		`json:"name"`
	Version string	`json:"version"`
	Coinbase string `json:"coinbase"`
	NodeType string `json:"remote"`
	StartTime string `json:"start"`
	Mining string   `json:"mining"`
}

type RequestPeerInfo struct {
	Data []*PeerInfo `json:"data"`
	ErrMsg string	`json:"err_msg"`
}

func getPeerInfo(bootnode string) (*RequestPeerInfo,error){
	url := fmt.Sprintf("http://%s:9000/nodeproxy/api/peerinfo", bootnode)
	var info = &RequestPeerInfo{}
	resp, err := http.Get(url)
	if err != nil {
		logs.Error("get peer info from bootnode failed", "err", err)
		return info, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//Failed to read response.
		logs.Error("failed to read peer info from response", "err", err)
		return info, err
	}
	json.Unmarshal(data, info)
	if len(info.ErrMsg) != 0 {
		logs.Error("get peer info from bootnode got error", "err", info.ErrMsg)
		return info, errors.New(info.ErrMsg)
	}
	peerMux.Lock()
	defer peerMux.Unlock()
	peerInfos = sync.Map{}
	for _, d := range info.Data {
		logs.Debug("get node peer", "info", d)
		peerInfos.Store(d.PeerId, d)
	}
	return info, nil
}

func getPeerInfoByCoinbase(coinbase string) []*PeerInfo {
	peerMux.RLock()
	defer peerMux.RUnlock()
	var info = make([]*PeerInfo,0)
	peerInfos.Range(func(k,v interface{}) bool {
		logs.Debug("find peer info by coinbase", "coinbase", coinbase, "key", k.(string))
		vinfo := v.(*PeerInfo)
		if strings.ToLower(coinbase) == strings.ToLower(vinfo.Coinbase) {
			info = append(info, vinfo)
		}
		return true
	})
	return info
}