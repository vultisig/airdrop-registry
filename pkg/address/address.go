package address

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	bchchaincfg "github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchutil"
	ltcchaincfg "github.com/ltcsuite/ltcd/chaincfg"
	ltcutil "github.com/ltcsuite/ltcd/ltcutil"
	"github.com/vultisig/airdrop-registry/pkg/utils"
	tss "github.com/vultisig/mobile-tss-lib/tss"
)

type ChainKeys struct {
	ChainName string
	PublicKey string
	Address   string
}

func GenerateChainKeys(chainName, hexPubKeyECDSA, hexPubKeyEdDSA, hexChainCode, path string) (ChainKeys, error) {
	var pubKeyBytes []byte
	var err error

	keys := ChainKeys{
		ChainName: chainName,
	}

	switch chainName {
	case "solana", "sui", "polkadot":
		if hexPubKeyEdDSA == "" {
			return ChainKeys{}, fmt.Errorf("EdDSA public key required for %s", chainName)
		}
		pubKeyBytes, err = hex.DecodeString(hexPubKeyEdDSA)
		if err != nil {
			return ChainKeys{}, fmt.Errorf("invalid EdDSA public key: %w", err)
		}
		keys.PublicKey = hexPubKeyEdDSA
	default:
		if hexPubKeyECDSA == "" || hexChainCode == "" {
			return ChainKeys{}, fmt.Errorf("ECDSA public key and chain code required for %s", chainName)
		}
		derivedKey, err := tss.GetDerivedPubKey(hexPubKeyECDSA, hexChainCode, path, false)
		if err != nil {
			return ChainKeys{}, err
		}
		pubKeyBytes, err = hex.DecodeString(derivedKey)
		if err != nil {
			return ChainKeys{}, fmt.Errorf("invalid derived ECDSA public key: %w", err)
		}
		keys.PublicKey = derivedKey
	}

	switch chainName {
	case "bitcoin", "dash", "dogecoin":
		var net *chaincfg.Params
		var prefix byte

		switch chainName {
		case "bitcoin":
			net = &chaincfg.MainNetParams
		case "dash":
			net = &chaincfg.MainNetParams
			prefix = 76
		case "dogecoin":
			net = &chaincfg.MainNetParams
			prefix = 30
		}

		if chainName == "bitcoin" {
			witnessProgram := btcutil.Hash160(pubKeyBytes)
			conv, err := btcutil.NewAddressWitnessPubKeyHash(witnessProgram, net)
			if err != nil {
				return ChainKeys{}, err
			}
			keys.Address = conv.EncodeAddress()
		} else {
			pubKey, err := btcec.ParsePubKey(pubKeyBytes)
			if err != nil {
				return ChainKeys{}, err
			}
			pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
			address := make([]byte, 25)
			address[0] = prefix
			copy(address[1:], pubKeyHash)
			checksum := chainhash.DoubleHashB(address[:21])
			copy(address[21:], checksum[:4])
			keys.Address = base58.Encode(address)
		}

	case "bitcoincash":
		pubKey, err := btcec.ParsePubKey(pubKeyBytes)
		if err != nil {
			return ChainKeys{}, err
		}
		pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
		addr, err := bchutil.NewAddressPubKeyHash(pubKeyHash, &bchchaincfg.MainNetParams)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = addr.EncodeAddress()
		keys.Address = fmt.Sprintf("bitcoincash:%s", keys.Address)

	case "litecoin":
		net := &ltcchaincfg.MainNetParams
		witnessProgram := btcutil.Hash160(pubKeyBytes)
		conv, err := ltcutil.NewAddressWitnessPubKeyHash(witnessProgram, net)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = conv.EncodeAddress()

	case "ethereum", "arbitrum", "avalanche", "bsc", "base", "blast", "cronos", "optimism", "polygon", "zksync":
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
		keys.Address = utils.SS58Encode(pubKeyBytes, 0)

	case "sui":
		keys.Address = "0x" + hex.EncodeToString(pubKeyBytes)

	default:
		return ChainKeys{}, fmt.Errorf("unsupported chain: %s", chainName)
	}

	return keys, nil
}

func GenerateSupportedChainAddresses(hexPubKeyECDSA, hexPubKeyEdDSA, hexChainCode string) (map[string]string, error) {
	addresses := make(map[string]string)

	for _, chain := range supportedChains {
		var keys ChainKeys
		var err error

		switch chain.name {
		case "solana", "sui", "polkadot":
			keys, err = GenerateChainKeys(chain.name, "", hexPubKeyEdDSA, "", chain.derivePath)
		default:
			keys, err = GenerateChainKeys(chain.name, hexPubKeyECDSA, "", hexChainCode, chain.derivePath)
		}

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
	{name: "bitcoincash", derivePath: "m/44'/145'/0'/0/0"},
	{name: "blast", derivePath: "m/44'/60'/0'/0/0"},
	{name: "cronos", derivePath: "m/44'/60'/0'/0/0"},
	{name: "dash", derivePath: "m/44'/5'/0'/0/0"},
	{name: "dogecoin", derivePath: "m/44'/3'/0'/0/0"},
	{name: "dydx", derivePath: "m/44'/118'/0'/0/0"},
	{name: "gaia", derivePath: "m/44'/118'/0'/0/0"},
	{name: "kujira", derivePath: "m/44'/118'/0'/0/0"},
	{name: "litecoin", derivePath: "m/84'/2'/0'/0/0"},
	{name: "optimism", derivePath: "m/44'/60'/0'/0/0"},
	{name: "polygon", derivePath: "m/44'/60'/0'/0/0"},
	{name: "zksync", derivePath: "m/44'/60'/0'/0/0"},
	// EDDSA
	{name: "solana", derivePath: ""},
	{name: "sui", derivePath: ""},
	{name: "polkadot", derivePath: ""},
}
