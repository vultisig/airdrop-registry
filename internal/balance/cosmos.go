package balance

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func (b *BalanceResolver) FetchThorchainBalanceOfAddress(address string) (float64, error) {
	if address == "" {
		return 0, fmt.Errorf("address cannot be empty")
	}
	url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", b.thornodeBaseAddress, address)
	runeBalance, err := b.fetchSpecificCosmosBalance(url, "rune", 8)
	if err != nil {
		return 0, fmt.Errorf("error fetching thorchain balance: %w", err)
	}
	// consider thorchain pooled rune
	pooledRune, ok := b.thorchainRuneProviders.Load(address)
	if ok {
		b.logger.Infof("address: %s, pooled rune: %v", address, pooledRune)
		if _, ok := pooledRune.(int64); ok {
			runeBalance += (float64)(pooledRune.(int64)) / math.Pow10(8)
		}
	}

	// consider thorchain bond
	bondValue, ok := b.thorchainBondProviders.Load(address)
	if !ok {
		return runeBalance, nil
	}
	b.logger.Infof("address: %s, bond: %s", address, bondValue)
	bond, err := strconv.ParseFloat(bondValue.(string), 64)
	if err != nil {
		b.logger.Errorf("failed to parse bond value: %v", err)
		return runeBalance, nil
	}
	return runeBalance + bond/math.Pow10(8), nil
}

func (b *BalanceResolver) FetchThorchainBalanceOfAddressAndAsset(address, asset string) (float64, error) {
	if strings.EqualFold(asset, "rujira") {
		asset = "x/ruji"
	}
	b.logger.Infof("Fetching Thorchain balance for address: %s and asset: %s", address, asset)
	if address == "" {
		return 0, fmt.Errorf("address cannot be empty")
	}
	url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", b.thornodeBaseAddress, address)
	balance, err := b.fetchSpecificCosmosBalance(url, asset, 8)
	if err != nil {
		return 0, fmt.Errorf("error fetching thorchain balance: %w", err)
	}
	return balance, nil
}

type THORNodeBondProvider struct {
	BondAddress string `json:"bond_address"`
	Bond        string `json:"bond"`
}
type THORNodeBondProviders struct {
	Providers []THORNodeBondProvider `json:"providers"`
}
type THORNode struct {
	BondProviders THORNodeBondProviders `json:"bond_providers"`
}

// GetTHORChainBondProviders fetches the bond providers from THORChain
func (b *BalanceResolver) GetTHORChainBondProviders() error {
	url := "https://thornode.ninerealms.com/thorchain/nodes"
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching bond providers from %s: %w", url, err)
	}

	defer b.closer(resp.Body)
	var nodes []THORNode
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}
	if len(nodes) == 0 {
		return nil
	}
	// clear all the existing bond providers
	b.thorchainBondProviders.Range(func(k, v interface{}) bool {
		b.thorchainBondProviders.Delete(k)
		return true
	})
	for _, node := range nodes {
		for _, item := range node.BondProviders.Providers {
			bond := 0.0
			existing, ok := b.thorchainBondProviders.Load(item.BondAddress)
			if ok {
				bond, err = strconv.ParseFloat(existing.(string), 64)
				if err != nil {
					b.logger.Errorf("failed to parse bond value: %v", err)
				}
			}
			newBond, err := strconv.ParseFloat(item.Bond, 64)
			if err != nil {
				b.logger.Errorf("failed to parse bond value: %v", err)
				continue
			}
			b.thorchainBondProviders.Store(item.BondAddress, strconv.FormatFloat(bond+newBond, 'f', -1, 64))
		}
	}

	b.thorchainBondProviders.Range(func(k, v interface{}) bool {
		b.logger.Infof("bond provider: %s, bond: %s", k, v)
		return true
	})
	return nil
}

type THORNodeRuneProviderResponse struct {
	RuneAddress string `json:"rune_address"`
	Value       int64  `json:"value,string"`
}

func (b *BalanceResolver) GetTHORChainRuneProviders() error {
	url := fmt.Sprintf("%s/thorchain/rune_providers", b.thornodeBaseAddress)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching bond providers from %s: %w", url, err)
	}

	defer b.closer(resp.Body)
	var runeProviders []THORNodeRuneProviderResponse
	if err := json.NewDecoder(resp.Body).Decode(&runeProviders); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}
	//clear all the existing rune providers
	b.thorchainRuneProviders.Range(func(k, v interface{}) bool {
		b.thorchainRuneProviders.Delete(k)
		return true
	})
	if len(runeProviders) == 0 {
		return nil
	}
	for _, provider := range runeProviders {
		//discard rune provider with 0 value
		if provider.Value > 0 {
			b.thorchainRuneProviders.Store(provider.RuneAddress, provider.Value)
		}
	}

	b.thorchainRuneProviders.Range(func(k, v interface{}) bool {
		b.logger.Infof("rune provider: %s, value: %d", k, v)
		return true
	})
	return nil
}

func (b *BalanceResolver) GetLP(address string) (float64, error) {
	return 0.0, nil
}

func (b *BalanceResolver) FetchMayachainCacoBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://mayanode.mayachain.info/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "cacao", 10)
}
func (b *BalanceResolver) FetchMayachainMayaBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://mayanode.mayachain.info/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "maya", 4)
}

func (b *BalanceResolver) FetchCosmosBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://cosmos-rest.publicnode.com/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "uatom", 6)
}

func (b *BalanceResolver) FetchKujiraBalanceOfAddress(address string, denom string, decimals int) (float64, error) {
	url := fmt.Sprintf("%s/%s", b.kujiraBalanceBaseAddress, address)
	return b.fetchSpecificCosmosBalance(url, denom, decimals)
}

func (b *BalanceResolver) FetchOsmosisBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://osmosis-rest.publicnode.com/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "uosmo", 6)
}

func (b *BalanceResolver) FetchDydxBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://dydx-rest.publicnode.com/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "adydx", 18)
}

func (b *BalanceResolver) FetchTerraBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://terra-lcd.publicnode.com/cosmos/bank/v1beta1/spendable_balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "uluna", 6)
}

func (b *BalanceResolver) FetchTerraClassicBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://terra-classic-lcd.publicnode.com/cosmos/bank/v1beta1/spendable_balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "uluna", 6)
}

func (b *BalanceResolver) FetchNobleBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://noble-api.polkachu.com/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "uusdc", 6)
}

func (b *BalanceResolver) FetchAkashBalanceOfAddress(address string) (float64, error) {
	url := fmt.Sprintf("https://akash-rest.publicnode.com/cosmos/bank/v1beta1/balances/%s", address)
	return b.fetchSpecificCosmosBalance(url, "uakt", 6)
}

type CosmosData struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
}

func (b *BalanceResolver) fetchSpecificCosmosBalance(url, denom string, decimals int) (float64, error) {
	if url == "" {
		return 0, fmt.Errorf("url cannot be empty")
	}
	if denom == "" {
		return 0, fmt.Errorf("denom cannot be empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error fetching balance from %s: %w", url, err)
	}
	defer b.closer(resp.Body)
	if resp.StatusCode == http.StatusTooManyRequests {
		// rate limited, need to backoff and then retry
		return 0, ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching balance from %s: %s", url, resp.Status)
	}
	var result CosmosData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var balance float64
	if len(result.Balances) == 0 {
		return 0, nil
	}
	for _, b := range result.Balances {
		if strings.EqualFold(b.Denom, denom) {
			balance, err = strconv.ParseFloat(b.Amount, 64)
			if err != nil {
				return 0, fmt.Errorf("error converting balance to float: %v", err)
			}
			break
		}
	}

	return balance / math.Pow10(decimals), nil
}
