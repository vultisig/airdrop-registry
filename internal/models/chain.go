package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Chain int

const (
	THORChain Chain = iota
	Solana
	Ethereum
	Avalanche
	BscChain
	Bitcoin
	BitcoinCash
	Litecoin
	Dogecoin
	GaiaChain
	Kujira
	Dash
	MayaChain
	Arbitrum
	Base
	Optimism
	Polygon
	Blast
	CronosChain
	Sui
	Polkadot
	Zksync
	Dydx
)

var chainToString = map[Chain]string{
	THORChain:   "THORChain",
	Solana:      "Solana",
	Ethereum:    "Ethereum",
	Avalanche:   "Avalanche",
	BscChain:    "BSC",
	Bitcoin:     "Bitcoin",
	BitcoinCash: "BitcoinCash",
	Litecoin:    "Litecoin",
	Dogecoin:    "Dogecoin",
	GaiaChain:   "Cosmos",
	Kujira:      "Kujira",
	Dash:        "Dash",
	MayaChain:   "MayaChain",
	Arbitrum:    "Arbitrum",
	Base:        "Base",
	Optimism:    "Optimism",
	Polygon:     "Polygon",
	Blast:       "Blast",
	CronosChain: "CronosChain",
	Sui:         "Sui",
	Polkadot:    "Polkadot",
	Zksync:      "Zksync",
	Dydx:        "Dydx",
}

func (c Chain) String() string {
	if str, ok := chainToString[c]; ok {
		return str
	}
	return "UNKNOWN"
}
func (c Chain) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Chain) UnmarshalJSON(data []byte) error {
	var chainStr string
	if err := json.Unmarshal(data, &chainStr); err != nil {
		return err
	}
	for key, value := range chainToString {
		if value == chainStr {
			*c = key
			return nil
		}
	}
	return nil
}
func (c Chain) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *Chain) Scan(value interface{}) error {
	if value == nil {
		*c = 0
		return nil
	}

	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Chain enum: %v", value)
	}
	for key, value := range chainToString {
		if value == string(str) {
			*c = key
			return nil
		}
	}
	return nil
}
