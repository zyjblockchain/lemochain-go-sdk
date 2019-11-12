package transaction

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/LemoFoundationLtd/lemochain-core/chain/params"
	"github.com/LemoFoundationLtd/lemochain-core/chain/types"
	"github.com/LemoFoundationLtd/lemochain-core/common"
	"github.com/LemoFoundationLtd/lemochain-core/common/crypto"
	"lemo_kms/common/log"
	"math/big"
	"time"
)

const (
	chainID = 100
)

func hexTxECDSA(priv string) (*ecdsa.PrivateKey, error) {
	if len(priv) > 1 {
		if priv[:2] == "0x" || priv[:2] == "0X" {
			priv = priv[2:]
		}
	}
	return crypto.HexToECDSA(priv)
}

func OrdinaryTx(hexPriv, from, to, amount, message string) (*types.Transaction, error) {
	private, err := hexTxECDSA(hexPriv)
	if err != nil {
		fmt.Println("SendOrdinaryTx HexToECDSA err: ", err)
		return nil, err
	}
	sender, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}
	recipient, err := common.StringToAddress(to)
	if err != nil {
		return nil, err
	}
	amountNum := common.Lemo2Mo(amount)
	gasLimit := uint64(50000)
	gasPrice := new(big.Int).SetUint64(1e9)
	rawTx := types.NewTransaction(sender, recipient, amountNum, gasLimit, gasPrice, nil, params.OrdinaryTx, chainID, uint64(time.Now().Unix()+20*60), "", message)
	// 交易签名
	signer := types.MakeSigner()
	return signer.SignTx(rawTx, private)
}

// 创建candidate交易
func CreateCandidateTx(hexPriv, from, amount, message string, candidateInfo map[string]string) (*types.Transaction, error) {
	private, err := hexTxECDSA(hexPriv)
	if err != nil {
		return nil, err
	}
	sender, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}
	amountNum := common.Lemo2Mo(amount)
	gasLimit := uint64(1000000)
	gasPrice := new(big.Int).SetUint64(1e9)
	data, err := json.Marshal(candidateInfo)
	if err != nil {
		return nil, err
	}
	rawTx := types.NoReceiverTransaction(sender, amountNum, gasLimit, gasPrice, data, params.RegisterTx, chainID, uint64(time.Now().Unix()+20*60), "", message)
	signer := types.MakeSigner()
	return signer.SignTx(rawTx, private)
}

// 投票交易
func VoteCandidateTx(hexPriv, from, to, message string) (*types.Transaction, error) {
	private, err := hexTxECDSA(hexPriv)
	if err != nil {
		return nil, err
	}
	sender, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}
	recipient, err := common.StringToAddress(to)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(50000)
	gasPrice := new(big.Int).SetUint64(1e9)
	rawTx := types.NewTransaction(sender, recipient, nil, gasLimit, gasPrice, nil, params.VoteTx, chainID, uint64(time.Now().Unix()+20*60), "", message)
	signer := types.MakeSigner()
	return signer.SignTx(rawTx, private)
}

// 创建箱子交易
func BoxTx(hexPriv, from string, subTxs types.Transactions) (*types.Transaction, error) {
	private, err := hexTxECDSA(hexPriv)
	if err != nil {
		return nil, err
	}
	sender, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(50000)
	gasPrice := new(big.Int).SetUint64(1e9)
	data, _ := types.MarshalBoxData(subTxs)
	rawTx := types.NoReceiverTransaction(sender, nil, gasLimit, gasPrice, data, params.BoxTx, chainID, uint64(time.Now().Unix()+15*60), "", "boxTx")
	signer := types.MakeSigner()
	return signer.SignTx(rawTx, private)
}

// 创建资产
func CreateAssetTx(hexPriv, from, message string, assetInfo *types.Asset) (*types.Transaction, error) {
	private, err := hexTxECDSA(hexPriv)
	if err != nil {
		return nil, err
	}
	sender, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(100000)
	gasPrice := new(big.Int).SetUint64(1e9)
	data, err := json.Marshal(assetInfo)
	if err != nil {
		log.Errorf("json marshal error: %v", err)
		return nil, err
	}
	rawTx := types.NoReceiverTransaction(sender, nil, gasLimit, gasPrice, data, params.CreateAssetTx, chainID, uint64(time.Now().Unix()+20*60), "", message)
	signer := types.MakeSigner()
	return signer.SignTx(rawTx, private)
}

// 发行资产
func IssueAssetTx(hexPriv, from, to, message string, issueInfo *types.IssueAsset) (*types.Transaction, error) {
	private, err := hexTxECDSA(hexPriv)
	if err != nil {
		return nil, err
	}
	sender, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}
	recipient, err := common.StringToAddress(to)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(100000)
	gasPrice := new(big.Int).SetUint64(1e9)
	data, err := json.Marshal(issueInfo)
	if err != nil {
		log.Errorf("json marshal error: %v", err)
		return nil, err
	}
	rawTx := types.NewTransaction(sender, recipient, nil, gasLimit, gasPrice, data, params.IssueAssetTx, chainID, uint64(time.Now().Unix()+20*60), "", message)
	signer := types.MakeSigner()
	return signer.SignTx(rawTx, private)
}

// 增发资产
func ReplenishAssetTx(hexPriv, from, to, message string, replenishInfo *types.ReplenishAsset) (*types.Transaction, error) {
	private, err := hexTxECDSA(hexPriv)
	if err != nil {
		return nil, err
	}
	sender, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}
	recipient, err := common.StringToAddress(to)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(100000)
	gasPrice := new(big.Int).SetUint64(1e9)
	data, err := json.Marshal(replenishInfo)
	if err != nil {
		log.Errorf("json marshal error: %v", err)
		return nil, err
	}
	rawTx := types.NewTransaction(sender, recipient, nil, gasLimit, gasPrice, data, params.ReplenishAssetTx, chainID, uint64(time.Now().Unix()+20*60), "", message)
	signer := types.MakeSigner()
	return signer.SignTx(rawTx, private)
}

// 修改资产注册信息
func ModifyAssetProfileTx(hexPriv, from, message string, modifyAssetInfo *types.ModifyAssetInfo) (*types.Transaction, error) {
	private, err := hexTxECDSA(hexPriv)
	if err != nil {
		return nil, err
	}
	sender, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(100000)
	gasPrice := new(big.Int).SetUint64(1e9)
	data, err := json.Marshal(modifyAssetInfo)
	if err != nil {
		log.Errorf("json marshal error: %v", err)
		return nil, err
	}
	rawTx := types.NoReceiverTransaction(sender, nil, gasLimit, gasPrice, data, params.ModifyAssetTx, chainID, uint64(time.Now().Unix()+20*60), "", message)
	signer := types.MakeSigner()
	return signer.SignTx(rawTx, private)
}

// 交易资产交易
func TransferAssetTx(hexPriv, from, to, message string, tradingAsset *types.TradingAsset) (*types.Transaction, error) {
	private, err := hexTxECDSA(hexPriv)
	if err != nil {
		return nil, err
	}
	sender, err := common.StringToAddress(from)
	if err != nil {
		return nil, err
	}
	recipient, err := common.StringToAddress(to)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(100000)
	gasPrice := new(big.Int).SetUint64(1e9)
	data, err := json.Marshal(tradingAsset)
	if err != nil {
		log.Errorf("json marshal error: %v", err)
		return nil, err
	}
	rawTx := types.NewTransaction(sender, recipient, nil, gasLimit, gasPrice, data, params.TransferAssetTx, chainID, uint64(time.Now().Unix()+20*60), "", message)
	signer := types.MakeSigner()
	return signer.SignTx(rawTx, private)
}
