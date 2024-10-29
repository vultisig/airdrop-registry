package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Chain int

const (
	Undefined Chain = iota
	THORChain
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
var EVMChains = []Chain{
	Ethereum,
	Avalanche,
	BscChain,
	Base,
	Arbitrum,
	Optimism,
	Polygon,
	Blast,
	CronosChain,
	Zksync,
}

var chainDerivePath = map[Chain]string{
	Bitcoin:     "m/84'/0'/0'/0/0",
	Ethereum:    "m/44'/60'/0'/0/0",
	THORChain:   "m/44'/931'/0'/0/0",
	MayaChain:   "m/44'/931'/0'/0/0",
	Arbitrum:    "m/44'/60'/0'/0/0",
	Avalanche:   "m/44'/60'/0'/0/0",
	BscChain:    "m/44'/60'/0'/0/0",
	Base:        "m/44'/60'/0'/0/0",
	BitcoinCash: "m/44'/145'/0'/0/0",
	Blast:       "m/44'/60'/0'/0/0",
	CronosChain: "m/44'/60'/0'/0/0",
	Dash:        "m/44'/5'/0'/0/0",
	Dogecoin:    "m/44'/3'/0'/0/0",
	Dydx:        "m/44'/118'/0'/0/0",
	GaiaChain:   "m/44'/118'/0'/0/0",
	Kujira:      "m/44'/118'/0'/0/0",
	Litecoin:    "m/84'/2'/0'/0/0",
	Optimism:    "m/44'/60'/0'/0/0",
	Polygon:     "m/44'/60'/0'/0/0",
	Zksync:      "m/44'/60'/0'/0/0",
	Solana:      "",
	Sui:         "",
	Polkadot:    "",
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

func (c *Chain) GetDerivePath() string {
	if str, ok := chainDerivePath[*c]; ok {
		return str
	}
	return ""
}

func (c *Chain) IsEdDSA() bool {
	if *c == Solana || *c == Sui || *c == Polkadot {
		return true
	}
	return false
}
