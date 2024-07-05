package address

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	bip32 "github.com/tyler-smith/go-bip32"
)

type ChainKeys struct {
	ChainName string
	PublicKey string
	Address   string
}

func parseDerivationPath(path string) ([]uint32, error) {
	var result []uint32
	parts := strings.Split(path, "/")
	if parts[0] != "m" {
		return nil, fmt.Errorf("invalid derivation path")
	}
	for _, part := range parts[1:] {
		var hardened uint32
		if strings.HasSuffix(part, "'") {
			return nil, fmt.Errorf("cannot derive hardened key from public key")
		}
		index, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid derivation path segment: %v", err)
		}
		result = append(result, uint32(index)+hardened)
	}
	return result, nil
}

func deriveKey(masterKey *bip32.Key, derivePath string) (*bip32.Key, error) {
	key := masterKey
	segments, err := parseDerivationPath(derivePath)
	if err != nil {
		return nil, err
	}
	for _, segment := range segments {
		key, err = key.NewChildKey(segment)
		if err != nil {
			return nil, err
		}
	}
	return key, nil
}

func GenerateChainKeys(chainName, derivePath string, masterKey *bip32.Key) (ChainKeys, error) {
	derivedKey, err := deriveKey(masterKey, derivePath)
	if err != nil {
		return ChainKeys{}, err
	}

	nonHardenedPubKey := derivedKey.PublicKey().Key
	publicKeyHex := hex.EncodeToString(nonHardenedPubKey)

	keys := ChainKeys{
		ChainName: chainName,
		PublicKey: publicKeyHex,
	}

	switch chainName {
	case "bitcoin":
		net := &chaincfg.MainNetParams
		pubKeyBytes := nonHardenedPubKey
		addressPubKey, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKeyBytes), net)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = addressPubKey.EncodeAddress()
	case "ethereum":
		pubKeyBytes, err := hex.DecodeString(publicKeyHex)
		if err != nil {
			return ChainKeys{}, err
		}
		pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
		if err != nil {
			return ChainKeys{}, err
		}
		keys.Address = crypto.PubkeyToAddress(*pubKey).Hex()
	case "thorchain":
		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount("thor", "thorpub")
		compressedPubkey := secp256k1.PubKey{Key: nonHardenedPubKey}
		addr := sdk.AccAddress(compressedPubkey.Address().Bytes())
		keys.Address = addr.String()
	case "mayachain":
		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount("maya", "mayapub")
		compressedPubkey := secp256k1.PubKey{Key: nonHardenedPubKey}
		addr := sdk.AccAddress(compressedPubkey.Address().Bytes())
		keys.Address = addr.String()
	default:
		return ChainKeys{}, fmt.Errorf("unsupported chain: %s", chainName)
	}

	return keys, nil
}

var supportedChains = []struct {
	name       string
	derivePath string
}{
	{name: "bitcoin", derivePath: "m/0/0"},
	{name: "ethereum", derivePath: "m/0/0"},
	{name: "thorchain", derivePath: "m/0/0"},
	{name: "mayachain", derivePath: "m/0/0"},
}
