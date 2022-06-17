package async

import "github.com/astaxie/beego/logs"

type FilterFunc func (*NodeStatus) bool


func getFilterFuncs(filters map[string]string) []FilterFunc {
	var filterfuncs = make([]FilterFunc,0)
	for k,v := range filters {
		switch k {
		case "ghpbversion":
			filterfuncs = append(filterfuncs, getFilterGhpbVersion(v))
		case "boeversion":
			filterfuncs = append(filterfuncs, getFilterBoeVersion(v))
		case "coinbase":
			filterfuncs = append(filterfuncs, getFilterCoinbase(v))
		case "nodetype":
			filterfuncs = append(filterfuncs, getFilterNodeType(v))
		case "nodename":
			filterfuncs = append(filterfuncs, getFilterNodeName(v))
		case "mining":
			filterfuncs = append(filterfuncs, getFilterMining(v))
		case "status":
			filterfuncs = append(filterfuncs, getFilterByStatus(v))
		default:
			logs.Error("unknown filter type", "k", k)
		}
	}
	return filterfuncs
}

func getFilterGhpbVersion(version string) FilterFunc {
	return func (info *NodeStatus) bool {
		if info != nil && info.GhpbVersion == version {
			return true
		}
		return false
	}
}

func getFilterBoeVersion(version string) FilterFunc {
	return func (info *NodeStatus) bool {
		if info != nil && info.BoeVersion == version {
			return true
		}
		return false
	}
}
func getFilterCoinbase(Coinbase string) FilterFunc {
	return func (info *NodeStatus) bool {
		if info != nil && info.Coinbase == Coinbase {
			return true
		}
		return false
	}
}

func getFilterNodeType(nodetype string) FilterFunc {
	return func (info *NodeStatus) bool {
		if info != nil && info.NodeType == nodetype {
			return true
		}
		return false
	}
}

func getFilterNodeName(nodename string) FilterFunc {
	return func (info *NodeStatus) bool {
		if info != nil && info.NodeName == nodename {
			return true
		}
		return false
	}
}

func getFilterMining(mining string) FilterFunc {
	return func (info *NodeStatus) bool {
		if info != nil && info.Mining == mining {
			return true
		}
		return false
	}
}

func getFilterByStatus(status string) FilterFunc {
	return func (info *NodeStatus) bool {
		logs.Debug("check status ", "filter ", status, "input is", info.Status)
		if info != nil && info.Status == status {
			return true
		}
		return false
	}
}

