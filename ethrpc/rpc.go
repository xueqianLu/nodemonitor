package ethrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	ErrNotConnect = errors.New("RPC not connected")
)

type EthRpc struct {
	url       string
	rpcclient *rpc.Client
}

func NewRPC(url string) *EthRpc {
	ethRpc := new(EthRpc)
	c, err := rpc.Dial(url)
	if err != nil {
		beego.Error("eth rpc init failed", "dial url error", err)
	} else {
		ethRpc.url = url
		ethRpc.rpcclient = c
	}
	return ethRpc
}

func (ethRpc *EthRpc) RandomByNumber(hex string) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_getRandom", hex)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (ethRpc *EthRpc) BlockByNumber(number *big.Int, tx bool) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_getBlockByNumber", hexutil.EncodeBig(number), tx)
	if err != nil {
		return nil, err
	}
	//beego.Info("got block by number", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) BlockNumber() (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_blockNumber")
	if err != nil {
		return nil, err
	}
	//beego.Info("got blocknumber", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) BlockTransactionCountByNumber(number *big.Int) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_getBlockTransactionCountByNumber", hexutil.EncodeBig(number))
	if err != nil {
		return nil, err
	}
	//beego.Info("eth_getBlockTransactionCountByNumber", "result ", string(raw))
	return raw, nil
}
func (ethRpc *EthRpc) TransactionByHash(hash common.Hash) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, err
	}
	//beego.Info("eth_getTransactionByHash", "result ", string(raw))
	return raw, nil
}
func (ethRpc *EthRpc) TransactionByBlockNumberAndIndex(number *big.Int, index uint) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_getTransactionByBlockNumberAndIndex", hexutil.EncodeBig(number), hexutil.Uint64(index))
	if err != nil {
		return nil, err
	}
	//beego.Info("eth_getTransactionByBlockNumberAndIndex", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) TransactionCount(addr common.Address, number *big.Int) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_getTransactionCount", addr, toBlockNumArg(number))
	if err != nil {
		return nil, err
	}
	//beego.Info("eth_getTransactionCount", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) SendRawTransaction(data []byte) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_sendRawTransaction", hexutil.Encode(data))
	if err != nil {
		return err.Error(), errors.New("rpc request failed")
	}
	//beego.Info("eth_sendRawTransaction", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) TransactionReceipt(txhash common.Hash) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_getTransactionReceipt", txhash)
	if err != nil {
		return err.Error(), errors.New("request failed")
	}
	//beego.Info("eth_getTransactionReceipt", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) CallContract(msg ethereum.CallMsg, blockNumber *big.Int) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return err.Error(), errors.New("request failed")
	}
	//beego.Info("eth_call", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) GetCode(account common.Address, blockNumber *big.Int) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_getCode", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	//beego.Info("eth_getCode", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) StorageAt(account common.Address, key common.Hash, blockNumber *big.Int) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_getStorageAt", account, key, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	//beego.Info("eth_getStorageAt", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) NetPeers() (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "admin_peers")
	if err != nil {
		return nil, err
	}
	//beego.Info("eth_gasPrice", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) HpbNodes() (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "prometheus_getHpbNodes", "latest")
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (ethRpc *EthRpc) HpbCadNodes() (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "prometheus_getCandidateNodes", "latest")
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (ethRpc *EthRpc) GasPrice() (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_gasPrice")
	if err != nil {
		return nil, err
	}
	//beego.Info("eth_gasPrice", "result ", string(raw))
	return raw, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func (ethRpc *EthRpc) EstimateGas(msg ethereum.CallMsg) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "eth_estimateGas", toCallArg(msg))
	if err != nil {
		return err.Error(), errors.New("request failed")
	}
	//beego.Info("eth_estimateGas", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) ElectedMiner(blockNumber *big.Int) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var err = ethRpc.rpcclient.Call(&raw, "prometheus_getElectedMiner", toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	beego.Info("prometheus_getElectedMiner", "result ", string(raw))
	return raw, nil
}

func (ethRpc *EthRpc) GetLogs(fq ethereum.FilterQuery) (interface{}, error) {
	if ethRpc.rpcclient == nil {
		return nil, ErrNotConnect
	}

	var raw json.RawMessage
	var param, err = toFilterArg(fq)
	if err != nil {
		return nil, err
	}
	err = ethRpc.rpcclient.Call(&raw, "eth_getLogs", param)
	if err != nil {
		return err.Error(), errors.New("request failed")
	}
	//beego.Info("eth_getLogs", "result ", string(raw))
	return raw, nil
}

func toCallArg(msg ethereum.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = toBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = toBlockNumArg(q.ToBlock)
	}
	return arg, nil
}
