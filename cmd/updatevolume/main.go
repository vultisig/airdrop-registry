package main

import (
	"fmt"

	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

func main() {
	ecdsaKey := "03b89475a593ca3f9855dd976043991eb9cae344dbc5117388fc5ba8e1cd7376b8"
	eddsaKey := "e2cb7f4cb662f36e19d7bc07606c09e4c372fa6a3d95bce7910254fb2e06cab7"
	hexChainCode := "1d56947b5d04d32e13332c0ec308123ca82bac779f6304116820d9742fe1d4e3"
	vault := &models.Vault{
		ECDSA:        ecdsaKey,
		EDDSA:        eddsaKey,
		HexChainCode: hexChainCode,
	}
	addresses := make(map[common.Chain]string)
	for _, chain := range common.GetAllChains() {
		// generate address
		addr, err := vault.GetAddress(chain)
		if err != nil {
			panic(err)
		}
		addresses[chain] = addr
		fmt.Printf("Chain: %s, Address: %s\n", chain, addr)
	}
}
