package kmd

import (
	"github.com/LemoFoundationLtd/lemochain-core/chain/types"
	"github.com/LemoFoundationLtd/lemochain-core/common"
	"github.com/lemochain-go-sdk/transaction"
	"lemo_kms/common/log"
	"math/big"
)

type Client struct {
	chainUrl string
}

func NewClient(url string) *Client {
	return &Client{chainUrl: url}
}

// SendOrdinaryTx 发送普通交易，返回交易hash
func (c *Client) SendOrdinaryTx(hexPriv, from, to, amount, message string) (interface{}, error) {
	signedTx, err := transaction.OrdinaryTx(hexPriv, from, to, amount, message)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

// SendCreateCandidateTx
func (c *Client) SendCreateCandidateTx(hexPriv, from, amount, message string, candidateInfo map[string]string) (interface{}, error) {
	signedTx, err := transaction.CreateCandidateTx(hexPriv, from, amount, message, candidateInfo)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

// SendVoteCandidateTx
func (c *Client) SendVoteCandidateTx(hexPriv, from, to, message string) (interface{}, error) {
	signedTx, err := transaction.VoteCandidateTx(hexPriv, from, to, message)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

func (c *Client) SendBoxTx(hexPriv, from string, subTxs types.Transactions) (interface{}, error) {
	signedTx, err := transaction.BoxTx(hexPriv, from, subTxs)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

// SendCreateAssetTx
func (c *Client) SendCreateAssetTx(hexPriv, from, message string, category, Decimal uint32, IsDivisible, IsReplenishable bool, profile map[string]string) (interface{}, error) {
	assetInfo := &types.Asset{
		Category:        category,
		IsDivisible:     IsDivisible,
		Decimal:         Decimal,
		IsReplenishable: IsReplenishable,
		Profile:         profile,
	}
	signedTx, err := transaction.CreateAssetTx(hexPriv, from, message, assetInfo)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

// SendIssueAssetTx
func (c *Client) SendIssueAssetTx(hexPriv, from, to, message string, assetCode, metaData, amount string) (interface{}, error) {
	num, b := new(big.Int).SetString(amount, 10)
	if !b {
		log.Errorf("new(big.Int).SetString(amount, 10) failed")
		return nil, nil
	}
	issueInfo := &types.IssueAsset{
		AssetCode: common.HexToHash(assetCode),
		MetaData:  metaData,
		Amount:    num,
	}
	signedTx, err := transaction.IssueAssetTx(hexPriv, from, to, message, issueInfo)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

// SendReplenishAssetTx
func (c *Client) SendReplenishAssetTx(hexPriv, from, to, message string, assetCode, assetId, amount string) (interface{}, error) {
	num, b := new(big.Int).SetString(amount, 10)
	if !b {
		log.Errorf("new(big.Int).SetString(amount, 10) failed")
		return nil, nil
	}
	replenishInfo := &types.ReplenishAsset{
		AssetCode: common.HexToHash(assetCode),
		AssetId:   common.HexToHash(assetId),
		Amount:    num,
	}
	signedTx, err := transaction.ReplenishAssetTx(hexPriv, from, to, message, replenishInfo)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

// SendModifyAssetProfileTx
func (c *Client) SendModifyAssetProfileTx(hexPriv, from, message string, assetCode string, updateProfile map[string]string) (interface{}, error) {
	modifyAssetInfo := &types.ModifyAssetInfo{
		AssetCode:     common.HexToHash(assetCode),
		UpdateProfile: updateProfile,
	}
	signedTx, err := transaction.ModifyAssetProfileTx(hexPriv, from, message, modifyAssetInfo)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

// SendTransferAssetTx
func (c *Client) SendTransferAssetTx(hexPriv, from, to, message string, assetId, amount string) (interface{}, error) {
	num, b := new(big.Int).SetString(amount, 10)
	if !b {
		log.Errorf("new(big.Int).SetString(amount, 10) failed")
		return nil, nil
	}
	tradingAsset := &types.TradingAsset{
		AssetId: common.HexToHash(assetId),
		Value:   num,
	}
	signedTx, err := transaction.TransferAssetTx(hexPriv, from, to, message, tradingAsset)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

// SendModifyMultiTx
/*
signers: [{
                    address: 'Lemo83HKZK68JQZDRGS5PWT2ZBSKR5CRADCSJB9B',
                    weight: 50,
                }, {
                    address: 'Lemo83W3DBN8QASNAR2D5386QSNGC8DAN8TSRK53',
                    weight: 30,
                }, {
                    address: 'Lemo83P3377KJHWD3DGQ4T82QDB8QBG69456BQZ8',
                    weight: 50,
                }],
*/
type SignerAddr struct {
	Address string
	Weight  uint8
}

func (c *Client) SendModifyMultiTx(hexPriv, from, to, message string, signers []SignerAddr) (interface{}, error) {
	ss := make(types.Signers, 0)
	for _, signer := range signers {
		addr, _ := common.StringToAddress(signer.Address)
		s := types.SignAccount{
			Address: addr,
			Weight:  signer.Weight,
		}
		ss = append(ss, s)
	}
	signedTx, err := transaction.ModifyMultiTx(hexPriv, from, to, message, ss)
	if err != nil {
		return nil, err
	}
	// 发送交易到glemo
	return SendTx(signedTx, c.chainUrl)
}

// GetBalance
func (c *Client) GetBalance(address string) (string, error) {
	return getBalance(address, c.chainUrl)
}

func (c *Client) GetAccount(address string) (map[string]interface{}, error) {
	return getAccount(address, c.chainUrl)
}

func (c *Client) IsCandidateAcc(address string) (bool, error) {
	mapAcc, err := getAccount(address, c.chainUrl)
	if err != nil {
		return false, err
	}
	if val, ok := mapAcc["candidate"]; ok {
		if newVal, ok := val.(map[string]interface{}); ok {
			for key, v := range newVal["profile"].(map[string]interface{}) {
				if key == "isCandidate" {
					if v.(string) == "true" {
						return true, nil
					} else {
						return false, nil
					}

				}
			}
		}
		return false, nil
	} else {
		return false, nil
	}
}
