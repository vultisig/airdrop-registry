package address

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	tss "github.com/vultisig/mobile-tss-lib/tss"
)

type ChainKeys struct {
	ChainName string
	PublicKey string
	Address   string
}

func GenerateChainKeys(chainName, hexPubKey, hexChainCode, path string, isEdDSA bool) (ChainKeys, error) {
	derivedKey, err := tss.GetDerivedPubKey(hexPubKey, hexChainCode, path, isEdDSA)
	if err != nil {
		return ChainKeys{}, err
	}

	keys := ChainKeys{
		ChainName: chainName,
		PublicKey: derivedKey,
	}

	pubKeyBytes, err := hex.DecodeString(derivedKey)
	if err != nil {
		return ChainKeys{}, err
	}

	switch chainName {
	case "bitcoin":
		net := &chaincfg.MainNetParams
		addressPubKey, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKeyBytes), net)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = addressPubKey.EncodeAddress()
	case "ethereum":
		pubKey, err := crypto.DecompressPubkey(pubKeyBytes)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = crypto.PubkeyToAddress(*pubKey).Hex()
	case "thorchain":
		pubKeyHash := btcutil.Hash160(pubKeyBytes)
		thorAddr, err := sdk.Bech32ifyAddressBytes("thor", pubKeyHash)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = thorAddr

	case "mayachain":
		pubKeyHash := btcutil.Hash160(pubKeyBytes)
		mayaAddr, err := sdk.Bech32ifyAddressBytes("maya", pubKeyHash)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = mayaAddr
	default:
		return ChainKeys{}, fmt.Errorf("unsupported chain: %s", chainName)
	}

	return keys, nil
}

func GenerateSupportedChainAddresses(hexPubKey, hexChainCode string) (map[string]string, error) {
	addresses := make(map[string]string)

	for _, chain := range supportedChains {
		keys, err := GenerateChainKeys(chain.name, hexPubKey, hexChainCode, chain.derivePath, false)
		if err != nil {
			return nil, fmt.Errorf("error generating address for %s: %w", chain.name, err)
		}
		addresses[chain.name] = keys.Address
	}

	return addresses, nil
}

var supportedChains = []struct {
	name       string
	derivePath string
}{
	{name: "bitcoin", derivePath: "m/84'/0'/0'/0/0"},
	{name: "ethereum", derivePath: "m/44'/60'/0'/0/0"},
	{name: "thorchain", derivePath: "m/44'/931'/0'/0/0"},
	{name: "mayachain", derivePath: "m/44'/931'/0'/0/0"},
}
