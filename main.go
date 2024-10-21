package main

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.base.org") 
	if err != nil {
		log.Fatalf("Failed to connect to the Base network: %v", err)
	}

	contractAddress := common.HexToAddress("0x76B4B28194170f9847Ae1566E44dCB4f2D97Ac24") 

	
	const abiJSON = //getABI

	contract, err := NewContract(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	value, err := contract.GetValue(nil)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Printf("Value from smart contract: %s\n", value.String())
}

type Contract struct {
	address common.Address
	abi     abi.ABI
	client  *ethclient.Client
}

func NewContract(address common.Address, client *ethclient.Client) (*Contract, error) {
	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}

	return &Contract{
		address: address,
		abi:     parsed,
		client:  client,
	}, nil
}

func (c *Contract) GetValue(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := c.abi.Call(opts, out, "totalUnderlying")
	return *ret0, err
}
