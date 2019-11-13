package utils

import "github.com/LemoFoundationLtd/lemochain-core/common/crypto"

func GenerateAddress() *crypto.AccountKey {
	account, _ := crypto.GenerateAddress()
	return account
}
