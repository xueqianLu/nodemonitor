package async

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xueqianLu/nodemonitor/ethrpc"
	"math"
	"math/big"
	"sort"
	"sync"
)

var (
	roundinfos  = sync.Map{} // round and roundinfo.
	startnumber = int64(0)
	client      = new(ethclient.Client)
	hpbrpc      = new(ethrpc.EthRpc)
)

func init() {
	var err error
	client, err = ethclient.Dial(beego.AppConfig.String("electedrpc"))
	if err != nil {
		panic("connect electedrpc failed")
	}
	hpbrpc = ethrpc.NewRPC(beego.AppConfig.String("electedrpc"))
	if hpbrpc == nil {
		panic("connect electedrpc failed")
	}
}

func addNodeLose(blocknumber int64, node common.Address) {

	roundNumber := int64(math.Floor(float64(blocknumber/200))) * 200
	round, exist := roundinfos.Load(roundNumber)
	if !exist {
		infos := NewSafeMap()
		infos.AddLose(node, int(blocknumber))
		roundinfos.Store(roundNumber, infos)
	} else {
		infos := round.(*SafeMap)
		infos.AddLose(node, int(blocknumber))
	}
}

func getMinerInfo() {
	ctx := context.Background()
	if startnumber == 0 {
		// init
		startnumber, _ = beego.AppConfig.Int64("electedstart")
	}
	blocknumber, err := client.BlockNumber(ctx)
	if err != nil {
		beego.Error("get block number failed", "err", err)
		return
	}
	for int64(blocknumber) >= startnumber {
		number := big.NewInt(startnumber)
		header, err := client.HeaderByNumber(ctx, number)
		if err != nil {
			beego.Error("get header by number failed", "err", err)
			return
		}
		elected, err := hpbrpc.ElectedMiner(number)
		if err != nil {
			beego.Error("get header by number failed", "err", err)
			return
		}
		var hexaddr string
		json.Unmarshal(elected.(json.RawMessage), &hexaddr)
		beego.Info("get elected", "miner is ", hexaddr, "block is ", startnumber)
		electedaddr := common.HexToAddress(hexaddr)
		if bytes.Compare(electedaddr.Bytes(), common.Address{}.Bytes()) == 0 {
			beego.Error("not get elected addr, wait next time")
			return
		}
		if bytes.Compare(header.Coinbase.Bytes(), electedaddr.Bytes()) != 0 && header.Difficulty.Int64() == 1 {
			addNodeLose(startnumber, electedaddr)
		}
		startnumber++
	}
}

type RoundsInfo struct {
	RoundNumber int64 `json:"round"`
	RoundInfo
}
type RoundsInfos []RoundsInfo

func (p RoundsInfos) Len() int { return len(p) }
func (p RoundsInfos) Less(i, j int) bool {
	return p[i].RoundNumber > p[j].RoundNumber
}
func (p RoundsInfos) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func getRoundInfo(blocknumber int64, round int64) RoundsInfos {
	if blocknumber > startnumber || blocknumber == 0 {
		blocknumber = startnumber
	}
	if round > 50 || round <= 0 {
		round = 10
	}
	result := make(RoundsInfos, 0)

	nodelist := getBoeNodeList()
	nodename := make(map[common.Address]string)
	for _, node := range nodelist {
		nodename[common.HexToAddress(node.Coinbase)] = node.Name
	}

	roundNumber := int64(math.Floor(float64(blocknumber/200))) * 200
	for i := int64(0); i < round; i++ {
		roundNumber -= 200 * i
		infos, exist := roundinfos.Load(roundNumber)
		if !exist {
			break
		}
		roundinfo := RoundsInfo{
			RoundNumber: roundNumber,
			RoundInfo:   make(RoundInfo, 0),
		}
		infos.(*SafeMap).Range(func(addr common.Address, lose int, blocks []int) {
			name := nodename[addr]
			nodeinfo := NodeInfo{
				Address:   addr,
				Name:      name,
				LoseBlock: lose,
				LostBlock: make([]int, len(blocks)),
			}
			copy(nodeinfo.LostBlock, blocks)
			roundinfo.RoundInfo = append(roundinfo.RoundInfo, nodeinfo)
		})
		sort.Sort(roundinfo.RoundInfo)
		result = append(result, roundinfo)
	}
	sort.Sort(result)
	return result
}
