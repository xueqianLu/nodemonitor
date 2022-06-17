package async

import (
	"strings"
	"time"
)

func init() {
	go syncNodeInfo()
	go syncPeerInfo()
	go syncConsensusNode()
}

func syncNodeInfo() {
	getNodeInfoAndParse()

	timer := time.NewTicker(time.Minute * 30)
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

	for _,node := range bootnodes {
		getPeerInfo(node)
	}

	timer := time.NewTicker(time.Minute * 10)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			for _,node := range bootnodes {
				go getPeerInfo(node)
			}
		}
	}
}


func syncConsensusNode() {
	getHpNodes()
	timer := time.NewTicker(time.Minute * 30)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			go getHpNodes()
		}
	}
}

type NodeStatus struct {
	PeerId string `json:"peerid"`
	NodeName string `json:"name"`
	Coinbase string `json:"coinbase"`
	NodeType string `json:"nodetype"`
	Status string `json:"status"`
	GhpbVersion string `json:"ghpbversion"`
	BoeVersion string `json:"boeversion"`
	Mining string `json:"mining"`
	VoteNumber int64 `json:"vote"`

}

func parseVersion(peerversion string) (string,string) {
	sep := "&"
	versions := strings.Split(peerversion, sep)
	if len(versions) >= 2 {
		return versions[0], versions[1]
	} else if len(versions) == 1 {
		return versions[0], ""
	}
	return "", ""
}

func filterall(allinfo map[string]*NodeStatus, filterfunc []FilterFunc ) map[string]*NodeStatus {
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


func GetAllNodeStatus(filter map[string]string) map[string]*NodeStatus {
	filterfuncs := getFilterFuncs(filter)

	var allinfo = make(map[string]*NodeStatus)
	nodelist := getBoeNodeList()
	for _, node := range nodelist {
		status := &NodeStatus{
			NodeName: node.Name,
			Coinbase: node.Coinbase,
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
	filtered := filterall(allinfo, filterfuncs)
	return filtered

}
