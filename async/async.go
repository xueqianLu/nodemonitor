package async

import (
	"github.com/ethereum/go-ethereum/common"
	"sort"
	"strings"
	"time"
)

func init() {
	go syncNodeInfo()
	go syncPeerInfo()
	go syncConsensusNode()
	go syncMinerInfo()
}

func syncNodeInfo() {
	getNodeInfoAndParse()

	timer := time.NewTicker(time.Minute * 5)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			getNodeInfoAndParse()
		}
	}
}

func syncPeerInfo() {
	bootnodes := []string{"47.75.213.166", "47.94.20.30", "47.254.133.46", "47.88.60.227", "47.75.213.166"}

	for _, node := range bootnodes {
		getPeerInfo(node)
	}

	timer := time.NewTicker(time.Minute * 2)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			for _, node := range bootnodes {
				go getPeerInfo(node)
			}
		}
	}
}

func syncConsensusNode() {
	getHpNodes()
	timer := time.NewTicker(time.Minute * 5)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			go getHpNodes()
		}
	}
}

type NodeStatus struct {
	PeerId      string `json:"peerid"`
	NodeName    string `json:"name"`
	Coinbase    string `json:"coinbase"`
	NodeType    string `json:"nodetype"`
	Status      string `json:"status"`
	GhpbVersion string `json:"ghpbversion"`
	BoeVersion  string `json:"boeversion"`
	Mining      string `json:"mining"`
	VoteNumber  int64  `json:"vote"`
}

func parseVersion(peerversion string) (string, string) {
	sep := "&"
	versions := strings.Split(peerversion, sep)
	if len(versions) >= 2 {
		return versions[0], versions[1]
	} else if len(versions) == 1 {
		return versions[0], ""
	}
	return "", ""
}

func filterall(allinfo map[string]*NodeStatus, filterfunc []FilterFunc) map[string]*NodeStatus {
	var filtered = make(map[string]*NodeStatus)
	for k, info := range allinfo {
		passed := true
		for _, f := range filterfunc {
			if !f(info) {
				passed = false
			}
		}
		if passed {
			filtered[k] = info
		}
	}
	return filtered
}

type sortNodeStatus []*NodeStatus

func (v sortNodeStatus) Len() int      { return len(v) }
func (v sortNodeStatus) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v sortNodeStatus) Less(i, j int) bool {
	return v[i].VoteNumber > v[j].VoteNumber
}

type NodeInfos struct {
	HpNumber    int           `json:"hpcount"`
	HpOfflines  int           `json:"hpoffline"`
	HpNotMining int           `json:"hpnotmining"`
	Infos       []*NodeStatus `json:"nodeinfos"`
}

func GetAllNodeStatus(filter map[string]string) interface{} {
	var result = new(NodeInfos)

	filterfuncs := getFilterFuncs(filter)

	var allinfo = make(map[string]*NodeStatus)
	nodelist := getBoeNodeList()
	for _, node := range nodelist {
		status := &NodeStatus{
			NodeName:   node.Name,
			Coinbase:   node.Coinbase,
			VoteNumber: node.VoteNumber,
		}
		status.NodeType = getTypeByCoinbase(node.Coinbase)
		peerinfo := getPeerInfoByCoinbase(node.Coinbase)
		if len(peerinfo) > 0 {
			info := peerinfo[0]
			status.Status = "online"
			status.PeerId = info.PeerId
			status.Mining = info.Mining
			status.GhpbVersion, status.BoeVersion = parseVersion(info.Version)
		} else {
			status.Status = "offline"
			status.Mining = "false"
		}
		allinfo[node.Coinbase] = status
	}
	for _, node := range allinfo {
		if node.NodeType == "hpnode" {
			result.HpNumber += 1
			if node.Status == "offline" {
				result.HpOfflines += 1
			} else if node.Mining == "false" {
				result.HpNotMining += 1
			}
		}
	}
	//log.Println("result is ", result)
	filtered := filterall(allinfo, filterfuncs)
	var res = make([]*NodeStatus, 0, len(filtered))
	for _, v := range filtered {
		res = append(res, v)
	}
	sort.Sort(sortNodeStatus(res))

	result.Infos = res

	return result
}

type NodeInfo struct {
	Address   common.Address `json:"nodeaddress"`
	Name      string         `json:"nodename"`
	LoseBlock int            `json:"loseblock"`
	LostBlock []int          `json:"lostblocks"`
}
type RoundInfo []NodeInfo

func (p RoundInfo) Len() int { return len(p) }
func (p RoundInfo) Less(i, j int) bool {
	return p[i].LoseBlock > p[j].LoseBlock
}
func (p RoundInfo) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func syncMinerInfo() {
	getMinerInfo()
	timer := time.NewTicker(time.Minute * 5)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			go getMinerInfo()
		}
	}
}

func GetBlockLoseInfo(number int64, round int) RoundsInfos {
	return getRoundInfo(number, int64(round))
}
