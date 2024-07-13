package address

import (
	"encoding/hex"
	"fmt"
	"strings"

	ss58 "github.com/ChainSafe/gossamer"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	tss "github.com/vultisig/mobile-tss-lib/tss"
	"golang.org/x/crypto/blake2b"
)

type ChainKeys struct {
	ChainName string
	PublicKey string
	Address   string
}

func ss58Encode(pubKey []byte, prefix byte) (string, error) {
	address := append([]byte{prefix}, pubKey...)

	hasher, err := blake2b.New(64, nil)
	if err != nil {
		return "", err
	}

	hasher.Write(address)
	checksum := hasher.Sum(nil)

	address = append(address, checksum[:2]...)

	return ss58.Encode(address)
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
	case "bitcoin", "bitcoin cash", "dash", "dogecoin", "litecoin":
		var net *chaincfg.Params
		var prefix string

		switch chainName {
		case "bitcoin":
			net = &chaincfg.MainNetParams
			prefix = "bc1"
		case "bitcoin cash":
			net = &chaincfg.MainNetParams
			prefix = "bitcoincash:"
		case "dash":
			net = &chaincfg.MainNetParams
			prefix = "X"
		case "dogecoin":
			net = &chaincfg.MainNetParams
			prefix = "D"
		case "litecoin":
			net = &chaincfg.MainNetParams
			prefix = "ltc1"
		}

		var address string
		witnessProgram := btcutil.Hash160(pubKeyBytes)
		conv, err := btcutil.NewAddressWitnessPubKeyHash(witnessProgram, net)
		if err != nil {
			return ChainKeys{}, err
		}
		address = conv.EncodeAddress()
		addressWithoutPrefix := strings.Split(address, "bc1")[1]
		keys.Address = prefix + addressWithoutPrefix
		fmt.Println(address, addressWithoutPrefix, prefix)
	case "ethereum", "arbitrum", "avalanche", "bsc", "base", "blast chain", "cronoschain", "optimism eth", "polygon", "zksync":
		pubKey, err := crypto.DecompressPubkey(pubKeyBytes)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = crypto.PubkeyToAddress(*pubKey).Hex()
	case "thorchain", "mayachain", "gaia", "kujira", "dydx":
		pubKeyHash := btcutil.Hash160(pubKeyBytes)
		prefix := map[string]string{
			"thorchain": "thor",
			"mayachain": "maya",
			"gaia":      "cosmos",
			"kujira":    "kujira",
			"dydx":      "dydx",
		}[chainName]
		addr, err := sdk.Bech32ifyAddressBytes(prefix, pubKeyHash)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = addr
	case "solana":
		keys.Address = base58.Encode(pubKeyBytes)
	case "polkadot":
		const polkadotPrefix byte = 0x00 // ss58 prefix for Polkadot mainnet
		keys.Address, err = ss58Encode(pubKeyBytes, polkadotPrefix)
		if err != nil {
			return ChainKeys{}, err
		}

	case "sui":
		// keys.Address = fmt.Sprintf("0x%x", pubKeyBytes)
	// case "dydx":
	// 	pubKey := ed25519.PublicKey(pubKeyBytes)
	// 	keys.Address = fmt.Sprintf("dydx%x", pubKey)
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
	{name: "arbitrum", derivePath: "m/44'/60'/0'/0/0"},
	{name: "avalanche", derivePath: "m/44'/60'/0'/0/0"},
	{name: "bsc", derivePath: "m/44'/60'/0'/0/0"},
	{name: "base", derivePath: "m/44'/60'/0'/0/0"},
	{name: "bitcoin cash", derivePath: "m/44'/145'/0'/0/0"},
	{name: "blast chain", derivePath: "m/44'/60'/0'/0/0"},
	{name: "cronoschain", derivePath: "m/44'/60'/0'/0/0"},
	{name: "dash", derivePath: "m/44'/5'/0'/0/0"},
	{name: "dogecoin", derivePath: "m/44'/3'/0'/0/0"},
	{name: "dydx", derivePath: "m/44'/118'/0'/0/0"},
	{name: "gaia", derivePath: "m/44'/118'/0'/0/0"},
	{name: "kujira", derivePath: "m/44'/118'/0'/0/0"},
	{name: "litecoin", derivePath: "m/84'/2'/0'/0/0"},
	{name: "optimism eth", derivePath: "m/44'/60'/0'/0/0"},
	{name: "polkadot", derivePath: "m/44'/354'/0'/0'/0'"},
	{name: "polygon", derivePath: "m/44'/60'/0'/0/0"},
	{name: "solana", derivePath: "m/44'/501'/0'/0'"},
	{name: "sui", derivePath: "m/44'/784'/0'/0'/0'"},
	{name: "zksync", derivePath: "m/44'/60'/0'/0/0"},
}
