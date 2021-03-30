package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"time"

	cc "github.com/QuarkChain/goqkcclient/client"
)

var (
	qkcClient   = cc.NewClient("http://128.199.215.162:38391")
	pri1        = "0x82283e556e0d8cae13dc13a691c6a1cdc67ccd68216a14ae8e79ddd909d08a74"
	addr3       = "0x6D29f28D3F1FCeB5788D445eD1C267e291D6F1F4"
	wei         = new(big.Int).Mul(new(big.Int).SetUint64(1000000000), new(big.Int).SetUint64(1000000000))
	nonceManger = newQkcNoncePool()
)

func main() {
	Transfer()
	SubmitContract()
	CallContract()
}

func Transfer() {
	prvkey, _ := crypto.ToECDSA(common.FromHex(pri1))
	from := crypto.PubkeyToAddress(prvkey.PublicKey)
	tx, err := qkcClient.CreateTransaction(nonceManger.getNonce(from), &cc.QkcAddress{Recipient: from, FullShardKey: 0}, &cc.QkcAddress{Recipient: common.HexToAddress(addr3), FullShardKey: 0}, new(big.Int).Mul(new(big.Int).SetUint64(5), wei), uint64(3000000), new(big.Int).SetUint64(1000000000), nil)

	if err != nil {
		fmt.Println(err.Error())
	}
	tx, err = cc.SignTx(tx, prvkey)
	if err != nil {
		fmt.Println(err.Error())
	}

	txid, err := qkcClient.SendTransaction(tx)
	if err != nil {
		fmt.Println("SendTransaction error: ", err.Error())
	}
	time.Sleep(30 * time.Second)
	_txid, _ := cc.ByteToTransactionId(txid)

	rs, err := qkcClient.GetTransactionReceipt(_txid)
	fmt.Println("Transfer err", err)
	fmt.Println("Transfer result", rs)
}

func SubmitContract() {
	prvkey, _ := crypto.ToECDSA(common.FromHex(pri1))
	from := crypto.PubkeyToAddress(prvkey.PublicKey)
	code := common.Hex2Bytes("608060405234801561001057600080fd5b5060405161047238038061047283398101604052805160008054600160a060020a0319163317905501805161004c906001906020840190610053565b50506100ee565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061009457805160ff19168380011785556100c1565b828001600101855582156100c1579182015b828111156100c15782518255916020019190600101906100a6565b506100cd9291506100d1565b5090565b6100eb91905b808211156100cd57600081556001016100d7565b90565b610375806100fd6000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166342cbb15c811461005b578063a413686214610082578063cfae3217146100dd575b600080fd5b34801561006757600080fd5b50610070610167565b60408051918252519081900360200190f35b34801561008e57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100db94369492936024939284019190819084018382808284375094975061016c9650505050505050565b005b3480156100e957600080fd5b506100f261021c565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561012c578181015183820152602001610114565b50505050905090810190601f1680156101595780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b435b90565b805161017f9060019060208401906102b1565b507f37782675f6a443a36b946fb3a2029cdb9e3bcd39a26f2f32cc695f6c3ff5af3f816040518080602001828103825283818151815260200191508051906020019080838360005b838110156101df5781810151838201526020016101c7565b50505050905090810190601f16801561020c5780820380516001836020036101000a031916815260200191505b509250505060405180910390a150565b60018054604080516020601f600260001961010087891615020190951694909404938401819004810282018101909252828152606093909290918301828280156102a75780601f1061027c576101008083540402835291602001916102a7565b820191906000526020600020905b81548152906001019060200180831161028a57829003601f168201915b5050505050905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106102f257805160ff191683800117855561031f565b8280016001018555821561031f579182015b8281111561031f578251825591602001919060010190610304565b5061032b92915061032f565b5090565b61016991905b8082111561032b57600081556001016103355600a165627a7a723058203fadf97d16b56cc4a72867ed2881d4507508c000280680efe5dc2db4f591ad830029")
	tx, err := qkcClient.CreateTransaction(nonceManger.getNonce(from), &cc.QkcAddress{Recipient: from, FullShardKey: 0}, nil, new(big.Int), uint64(3000000), new(big.Int).SetUint64(1000000000), code)
	if err != nil {
		fmt.Println(err.Error())
	}
	tx, err = cc.SignTx(tx, prvkey)
	if err != nil {
		fmt.Println(err.Error())
	}

	txid, err := qkcClient.SendTransaction(tx)
	if err != nil {
		fmt.Println("SendTransaction error: ", err.Error())
	}
	time.Sleep(30 * time.Second)
	_txid, _ := cc.ByteToTransactionId(txid)

	rs, err := qkcClient.GetTransactionReceipt(_txid)
	fmt.Println("SubmitContract err", err)
	fmt.Println("SubmitContract result", rs)
}
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	rs := bytesBuffer.Bytes()
	noneZeroIndex := 0
	for index := 0; index < len(rs); index++ {
		if rs[index] == 0 {
			noneZeroIndex = index
		}
	}
	return rs[noneZeroIndex+1:]
}
func makeGreetMsg(str string) []byte {
	rs := make([]byte, 0)
	rs = append(rs, common.Hex2Bytes("a41368620000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000")...)
	rs = append(rs, IntToBytes(len(str))...)
	rs = append(rs, []byte(str)...)

	ll := len(rs)
	for index := ll; index < 100; index++ {
		rs = append(rs, byte(0))
	}
	return rs

}
func CallContract() {
	prvkey, _ := crypto.ToECDSA(common.FromHex(pri1))
	from := crypto.PubkeyToAddress(prvkey.PublicKey)
	contractAddress := common.HexToAddress("0xc3c2b0d7f4A1239EAA77B49Cf1F17502eb882086")
	tx, err := qkcClient.CreateTransaction(nonceManger.getNonce(from), &cc.QkcAddress{Recipient: from, FullShardKey: 0}, &cc.QkcAddress{Recipient: contractAddress, FullShardKey: 0}, new(big.Int), uint64(3000000), new(big.Int).SetUint64(1000000000), makeGreetMsg("qkc666231321321321"))
	if err != nil {
		fmt.Println(err.Error())
	}
	tx, err = cc.SignTx(tx, prvkey)
	if err != nil {
		fmt.Println(err.Error())
	}

	txid, err := qkcClient.SendTransaction(tx)
	if err != nil {
		fmt.Println("SendTransaction error: ", err.Error())
	}
	time.Sleep(30 * time.Second)
	_txid, _ := cc.ByteToTransactionId(txid)

	rs, err := qkcClient.GetTransactionReceipt(_txid)
	fmt.Println("CallContract err", err)
	fmt.Println("CallContract receipt", rs)
}

type qkcNoncePool struct {
	curr map[common.Address]uint64
}

func newQkcNoncePool() *qkcNoncePool {
	return &qkcNoncePool{curr: make(map[common.Address]uint64, 0)}
}

func (q *qkcNoncePool) getNonce(from common.Address) uint64 {
	defer func() {
		q.curr[from]++
	}()
	if data, ok := q.curr[from]; ok {
		return data
	}

	nonce, err := qkcClient.GetNonce(&cc.QkcAddress{
		Recipient:    from,
		FullShardKey: 0,
	})
	if err != nil {
		panic(err)
	}
	q.curr[from] = nonce
	return nonce
}
