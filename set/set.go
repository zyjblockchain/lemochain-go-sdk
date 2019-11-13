package set

import (
	"bytes"
	"github.com/LemoFoundationLtd/lemochain-core/common"
	"github.com/LemoFoundationLtd/lemochain-core/common/crypto"
	"github.com/LemoFoundationLtd/lemochain-core/common/log"
	"github.com/lemochain-go-sdk/kmd"
	"github.com/lemochain-go-sdk/utils"
	"math/big"
	"sync"
	"time"
)

type Acc struct {
	Address common.Address
	Private string
	Balance *big.Int
}

// 存储具有余额的账户
type AccSet struct {
	Acc        map[common.Address]*Acc // key为address
	MinBalance common.AddressSlice     // 存储最小balance
	MaxBalance common.AddressSlice     // 存储最大的balance
	Addresses  common.AddressSlice     // address的一个队列
	sync.Mutex
}

func NewAccSet() *AccSet {
	accSet := &AccSet{
		Acc:        make(map[common.Address]*Acc),
		MinBalance: make(common.AddressSlice, 0),
		MaxBalance: make(common.AddressSlice, 0),
		Addresses:  make(common.AddressSlice, 0),
		Mutex:      sync.Mutex{},
	}
	return accSet
}

// Push
func (a *AccSet) Push(address common.Address, private string, balance *big.Int) {
	a.Lock()
	defer a.Unlock()
	if balance.Cmp(big.NewInt(0)) == 0 {
		return
	}
	// 存在则不能重复被push
	if _, ok := a.Acc[address]; ok {
		return
	} else {
		// 设置acc
		acc := &Acc{
			Address: address,
			Private: private,
			Balance: balance,
		}
		a.Acc[address] = acc
		// 进入栈
		a.Addresses = append(a.Addresses, address)
		// 设置 minBalance
		if len(a.MinBalance) == 0 {
			a.MinBalance = append(a.MinBalance, address)
		} else {
			minB := a.Acc[a.MinBalance[len(a.MinBalance)-1]].Balance
			// 判断balance大小
			if balance.Cmp(minB) == -1 {
				a.MinBalance = append(a.MinBalance, address)
				return
			}
		}
		// 设置maxBalance
		if len(a.MaxBalance) == 0 {
			a.MaxBalance = append(a.MaxBalance, address)
		} else {
			maxB := a.Acc[a.MaxBalance[len(a.MaxBalance)-1]].Balance
			// 判断大小
			if balance.Cmp(maxB) == 1 {
				a.MaxBalance = append(a.MaxBalance, address)
			}
		}
	}
}

// Pop 出栈
func (a *AccSet) Pop() *Acc {
	if len(a.Addresses) == 0 {
		return nil
	}
	addr := a.Addresses[len(a.Addresses)-1]
	a.Addresses = a.Addresses[:len(a.Addresses)-1]
	acc, _ := a.Acc[addr]
	// 判断是否为当前的minBalance
	minAddr := a.MinBalance[len(a.MinBalance)-1]
	if bytes.Compare(minAddr.Bytes(), addr.Bytes()) == 0 {
		// 表示pop的是当前的最小balance
		a.MinBalance = a.MinBalance[:len(a.MinBalance)-1]
	}
	maxAddr := a.MaxBalance[len(a.MaxBalance)-1]
	if bytes.Compare(maxAddr.Bytes(), addr.Bytes()) == 0 {
		a.MaxBalance = a.MaxBalance[:len(a.MaxBalance)-1]
	}
	delete(a.Acc, addr)
	return acc
}
func (a *AccSet) Size() int {
	return len(a.Addresses)
}

func (a *AccSet) BatchPush(num int, amount string) {
	client := kmd.NewClient(utils.ChainUrl)
	accKeys := make([]*crypto.AccountKey, 0)
	for i := 0; i < num; i++ {
		accKey := utils.GenerateAddress()
		accKeys = append(accKeys, accKey)
		client.SendOrdinaryTx(utils.GodPrivate, utils.GodAddr, accKey.Address.String(), amount, "batch set Address")
	}
	time.Sleep(time.Second * 3)
	for _, accKey := range accKeys {
		// 定时查询
		temp := 0
		for {
			balance, err := client.GetBalance(accKey.Address.String())
			if err != nil {
				log.Errorf("getBalance error: %v", err)
				break
			}
			if balance != "0" {
				// 存入
				bb, _ := new(big.Int).SetString(balance, 10)
				a.Push(accKey.Address, accKey.Private, bb)
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
