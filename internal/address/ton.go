package address

import (
   
	"crypto/ed25519"
    "encoding/hex"
    "fmt"
    "github.com/xssnick/tonutils-go/ton/wallet"

)

func GetTonAddress(hexPublicKey string) (string, error) {
    pubKeyBytes, err := hex.DecodeString(hexPublicKey)
    if err != nil {
        return "", fmt.Errorf("invalid public key: %w", err)
    }
    fmt.Printf("Public Key Bytes: %x\n", pubKeyBytes)

    // Ensure the public key is the correct length
    if len(pubKeyBytes) != ed25519.PublicKeySize {
        fmt.Printf("Invalid public key length: expected %d bytes, got %d bytes\n", ed25519.PublicKeySize, len(pubKeyBytes))
        return "", fmt.Errorf("invalid public key length: expected %d bytes, got %d bytes", ed25519.PublicKeySize, len(pubKeyBytes))
    }

    // Convert to ed25519.PublicKey type
    pubKey := ed25519.PublicKey(pubKeyBytes)

    // Create a v4R2 wallet
    // The first argument is the wallet ID (use 0 if not needed)
    
	
	addressInstance, err := wallet.AddressFromPubKey(pubKey, wallet.V4R2, 698983191)
    if err != nil {
        fmt.Printf("Failed to create wallet: %v\n", err)
        return "", fmt.Errorf("failed to create wallet: %w", err)
    }

	addressInstance.SetBounce(false)

    return addressInstance.String(), nil
}
