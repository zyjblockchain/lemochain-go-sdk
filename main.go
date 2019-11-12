package main

import (
	"fmt"
	"github.com/LemoFoundationLtd/lemochain-core/chain/types"
	"github.com/lemochain-go-sdk/kmd"
	"github.com/lemochain-go-sdk/transaction"
	"math/rand"
	"strconv"
	"time"
)

const (
	addr01     = "Lemo83GN72GYH2NZ8BA729Z9TCT7KQ5FC3CR6DJG"
	addrPriv01 = "0xc21b6b2fbf230f665b936194d14da67187732bf9d28768aef1a3cbb26608f8aa"

	addr02     = "Lemo845TCPPH2H6SFJ9Q3ASTWFYZSWFY4N2AYA93"
	addrPriv02 = "0x518e38c58721a50e0bf39b59022e6c47a6de874e952675668cbfae266a8ae579"

	addr03     = "Lemo83FTJZS4CTHWWGG355BK2Z63G8R2JP3DA2DA"
	addrPriv03 = "0x6fdd60a57a1a6521a2ccbce5650da418245cf80a64f2e6cdde9ace9c36c3a992"
)

func main1() {
	// 创建client
	chainUrl := "https://distribution.lemolabs.com"
	c := kmd.NewClient(chainUrl)
	// // 1. 发送普通交易
	// result,err := c.SendOrdinaryTx(addrPriv01, addr01, addr02, "90", "test lemochain go-sdk")
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// 	return
	// }
	// fmt.Println("普通交易结果：", result)
	// // 2. 发送创建候选节点交易
	// info := make(map[string]string)
	// info[types.IsCandidateNode] = "true"
	// info[types.CandidateKeyNodeID] = "0xceab900b63027e3fb248ad04b046f38efff1f667ba5464f7bee64d7b302378dbc2cf8f11c106aa21a662e172d1efdbe6d95f7399f91e4d4449681e2e969673bc"
	// info[types.CandidateKeyHost] = "www.baidu.com"
	// info[types.CandidateKeyPort] = "7001"
	// info[types.CandidateKeyIncomeAddress] = addr02
	// result,err = c.SendCreateCandidateTx(addrPriv02, addr02,"5000000", "go-sdk create candidate test", info)
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// 	return
	// }
	// fmt.Println("注册候选节点交易结果：", result)

	// // 3. 投票交易
	// result,err := c.SendVoteCandidateTx(addrPriv02,addr02,addr02,"go-sdk vote tx test")
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// 	return
	// }
	// fmt.Println("投票交易结果：", result)

	// 4. 箱子交易
	rand.Seed(time.Now().UnixNano())
	subTxs := make(types.Transactions, 0)
	// 创建10个普通交易的转账
	for i := 0; i < 100; i++ {
		message := strconv.Itoa(rand.Int()) + "subTx"
		tx, err := transaction.OrdinaryTx(addrPriv01, addr01, addr02, "1", message)
		if err != nil {
			continue
		}
		subTxs = append(subTxs, tx)
	}
	// // 创建一个注册候选节点交易
	// info := make(map[string]string)
	// info[types.IsCandidateNode] = "true"
	// info[types.CandidateKeyNodeID] = "0xceab900b63027e3fb248ad04b046f38efff1f667ba5464f7bee64d7b302378dbc2cf8f11c106aa21a662e172d1efdbe6d95f7399f91e4d4449681e2e969673bc"
	// info[types.CandidateKeyHost] = "www.google.com"
	// info[types.CandidateKeyPort] = "7002"
	// info[types.CandidateKeyIncomeAddress] = addr03
	// candidateTx,err := transaction.CreateCandidateTx(addrPriv03,addr03,"5000000","注册一个测试候选节点", info)
	// if err != nil {
	// 	fmt.Println("1err: ", err)
	// 	return
	// }
	// subTxs = append(subTxs,candidateTx)
	result, err := c.SendBoxTx(addrPriv02, addr02, subTxs)
	if err != nil {
		fmt.Println("2err: ", err)
		return
	}
	fmt.Println("box交易结果：", result)
}
