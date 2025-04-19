package tokens

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

//go:embed predefined_tokens.json
var predefinedTokens string

type PredefinedTokenDiscoveryService struct {
	predefinedTokens []models.CoinBase
}


func NewPredefinedTokenService() AutoDiscoveryService {
	var tokens []models.CoinBase
	if err := json.Unmarshal([]byte(predefinedTokens), &tokens); err != nil {
		panic(fmt.Errorf("failed to unmarshal predefined tokens: %w", err))
	}
	return &PredefinedTokenDiscoveryService{
		predefinedTokens: tokens ,
	}
}

func (p *PredefinedTokenDiscoveryService) Search(coin models.CoinBase) (models.CoinBase, error) {
	for _, token := range p.predefinedTokens {
		if token.Chain == coin.Chain && token.ContractAddress == coin.ContractAddress {
			return token, nil
		}
	}
	return models.CoinBase{}, fmt.Errorf("token not found: chain=%s, address=%s", coin.Chain, coin.ContractAddress)
}

func (p *PredefinedTokenDiscoveryService) Discover(address string, chain common.Chain) ([]models.CoinBase, error) {
	return nil, fmt.Errorf("Discover method not implemented for PredefinedTokenDiscoveryService")
}
