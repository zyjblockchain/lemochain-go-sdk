package main

import (
	"github.com/LemoFoundationLtd/lemochain-core/chain/types"
	"github.com/LemoFoundationLtd/lemochain-core/common/crypto"
	"github.com/lemochain-go-sdk/kmd"
	"github.com/lemochain-go-sdk/set"
	"github.com/lemochain-go-sdk/transaction"
	"github.com/lemochain-go-sdk/utils"
	"log"
	"math/big"
	"time"
)

// func main1() {
// 	// 创建client
// 	chainUrl := "https://distribution.lemolabs.com"
// 	c := kmd.NewClient(chainUrl)
// 	// // 1. 发送普通交易
// 	// result,err := c.SendOrdinaryTx(addrPriv01, addr01, addr02, "90", "test lemochain go-sdk")
// 	// if err != nil {
// 	// 	fmt.Println("err: ", err)
// 	// 	return
// 	// }
// 	// fmt.Println("普通交易结果：", result)
// 	// // 2. 发送创建候选节点交易
// 	// info := make(map[string]string)
// 	// info[types.IsCandidateNode] = "true"
// 	// info[types.CandidateKeyNodeID] = "0xceab900b63027e3fb248ad04b046f38efff1f667ba5464f7bee64d7b302378dbc2cf8f11c106aa21a662e172d1efdbe6d95f7399f91e4d4449681e2e969673bc"
// 	// info[types.CandidateKeyHost] = "www.baidu.com"
// 	// info[types.CandidateKeyPort] = "7001"
// 	// info[types.CandidateKeyIncomeAddress] = addr02
// 	// result,err = c.SendCreateCandidateTx(addrPriv02, addr02,"5000000", "go-sdk create candidate test", info)
// 	// if err != nil {
// 	// 	fmt.Println("err: ", err)
// 	// 	return
// 	// }
// 	// fmt.Println("注册候选节点交易结果：", result)
//
// 	// // 3. 投票交易
// 	// result,err := c.SendVoteCandidateTx(addrPriv02,addr02,addr02,"go-sdk vote tx test")
// 	// if err != nil {
// 	// 	fmt.Println("err: ", err)
// 	// 	return
// 	// }
// 	// fmt.Println("投票交易结果：", result)
//
// 	// 4. 箱子交易
// 	rand.Seed(time.Now().UnixNano())
// 	subTxs := make(types.Transactions, 0)
// 	// 创建10个普通交易的转账
// 	for i := 0; i < 100; i++ {
// 		message := strconv.Itoa(rand.Int()) + "subTx"
// 		tx, err := transaction.OrdinaryTx(addrPriv01, addr01, addr02, "1", message)
// 		if err != nil {
// 			continue
// 		}
// 		subTxs = append(subTxs, tx)
// 	}
// 	// // 创建一个注册候选节点交易
// 	// info := make(map[string]string)
// 	// info[types.IsCandidateNode] = "true"
// 	// info[types.CandidateKeyNodeID] = "0xceab900b63027e3fb248ad04b046f38efff1f667ba5464f7bee64d7b302378dbc2cf8f11c106aa21a662e172d1efdbe6d95f7399f91e4d4449681e2e969673bc"
// 	// info[types.CandidateKeyHost] = "www.google.com"
// 	// info[types.CandidateKeyPort] = "7002"
// 	// info[types.CandidateKeyIncomeAddress] = addr03
// 	// candidateTx,err := transaction.CreateCandidateTx(addrPriv03,addr03,"5000000","注册一个测试候选节点", info)
// 	// if err != nil {
// 	// 	fmt.Println("1err: ", err)
// 	// 	return
// 	// }
// 	// subTxs = append(subTxs,candidateTx)
// 	result, err := c.SendBoxTx(addrPriv02, addr02, subTxs)
// 	if err != nil {
// 		fmt.Println("2err: ", err)
// 		return
// 	}
// 	fmt.Println("box交易结果：", result)
// }

func main() {
	// 1. 发送箱子交易
	randBoxTx(5)
	// 2. 发送注册候选节点交易
	randRegisterCandidateTx(5)
	// 3. 发送
}

// 随机创建资产交易

// 随机发送箱子交易
func randBoxTx(accountNum int) {
	accSet := set.NewAccSet()
	// 1. 创建地址，每个地址分配lemo
	accSet.BatchPush(accountNum, "1")
	log.Println("初始化地址数量：", accSet.Size())
	client := kmd.NewClient(utils.ChainUrl)
	fromSet := make([]*set.Acc, 0)
	for i := 0; i < accSet.Size()/2; i++ {
		fromSet = append(fromSet, accSet.Pop())
	}
	log.Println("fromSet: ", len(fromSet))
	toSet := make([]*set.Acc, 0)
	for i := 0; i < accSet.Size(); i++ {
		toSet = append(toSet, accSet.Pop())
	}
	log.Println("toSet: ", len(toSet))
	subTxs := make(types.Transactions, 0)
	for i, v := range fromSet {
		if i == len(toSet) {
			break
		}
		tx, _ := transaction.OrdinaryTx(v.Private, v.Address.String(), toSet[i].Address.String(), moToLemo(new(big.Int).Div(v.Balance, big.NewInt(2))), "box tx zyj")
		subTxs = append(subTxs, tx)
	}
	log.Println(len(subTxs))
	result, err := client.SendBoxTx(utils.GodPrivate, utils.GodAddr, subTxs)
	if err != nil {
		log.Println("error: ", err)
	}
	log.Println("result: ", result)
}

// 首先转账到注册者，然后注册成功之后又注销，注销之后把钱转回到16亿账户
func randRegisterCandidateTx(num int) {
	client := kmd.NewClient(utils.ChainUrl)
	accKeys := make([]*crypto.AccountKey, 0)
	for i := 0; i < num; i++ {
		accKey := utils.GenerateAddress()
		_, err := client.SendOrdinaryTx(utils.GodPrivate, utils.GodAddr, accKey.Address.String(), "5000001", "500万交易")
		if err == nil {
			accKeys = append(accKeys, accKey)
		}
	}
	time.Sleep(time.Second * 3)
	candidateAcc := make([]*crypto.AccountKey, 0) // 存储注册候选节点的账户
	for _, accKey := range accKeys {
		temp := 0
		for {
			balance, err := client.GetBalance(accKey.Address.String())
			if err != nil {
				log.Println("getBalance error: %v", err)
				break
			}
			if balance != "0" {
				// 开始注册候选节点
				info := make(map[string]string)
				info[types.CandidateKeyIsCandidate] = "true"
				info[types.CandidateKeyNodeID] = accKey.Public
				info[types.CandidateKeyHost] = "www.google.com"
				info[types.CandidateKeyPort] = "7002"
				result, err := client.SendCreateCandidateTx(accKey.Private, accKey.Address.String(), "5000000", "randRegisterCandidateTx zyj 测试", info)
				if err != nil {
					log.Println("error: ", err)
					continue
				}
				log.Println("注册候选节点交易hash： ", result)
				candidateAcc = append(candidateAcc, accKey)
				break
			} else {
				if temp == 3 {
					break
				}
				temp++
			}
			time.Sleep(time.Second * 3)
		}
	}

	time.Sleep(time.Second * 10)
	// 注销候选节点并把余额转回到16亿账户地址
	refundAcc := make([]*crypto.AccountKey, 0)
	for _, accKey := range candidateAcc {
		isCandidate, err := client.IsCandidateAcc(accKey.Address.String())
		log.Println("isCandidate: ", isCandidate)
		if err != nil {
			log.Println("error: ", err)
			continue
		}
		if isCandidate {
			// 进行注销候选节点
			info := make(map[string]string)
			info[types.CandidateKeyIsCandidate] = "false"
			info[types.CandidateKeyNodeID] = accKey.Public
			info[types.CandidateKeyHost] = "www.google.com"
			info[types.CandidateKeyPort] = "7002"
			result, err := client.SendCreateCandidateTx(accKey.Private, accKey.Address.String(), "0", "randRegisterCandidateTx zyj 测试", info)
			if err != nil {
				log.Println("error: ", err)
				continue
			}
			refundAcc = append(refundAcc, accKey)
			log.Println("注销候选节点交易hash： ", result)
		}
	}

	time.Sleep(time.Second * 5)
	// 把退回的押金转回给16亿账户
	for _, accKey := range refundAcc {
		temp := 0
		for {
			balance, err := client.GetBalance(accKey.Address.String())
			if err != nil {
				log.Println("getBalance error: %v", err)
				break
			}
			if balance != "0" || len(balance) > 22 {
				// 转给16亿账户 // todo 这里优化为代付交易gas交易进行转账
				result, err := client.SendOrdinaryTx(accKey.Private, accKey.Address.String(), utils.GodAddr, "5000000", "kkkkk")
				if err != nil {
					log.Println("error: ", err)
					continue
				}
				log.Println("转账交易hash： ", result)
				break
			} else {
				if temp == 3 {
					break
				}
				temp++
			}
			time.Sleep(time.Second * 3)
		}
	}
}

func moToLemo(num *big.Int) string {
	return num.String()[:len(num.String())-18]
}
