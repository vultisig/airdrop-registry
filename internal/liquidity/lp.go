package liquidity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type LiquidityPositionResolver struct {
	logger            *logrus.Logger
	thorwalletBaseURL string
	arbitrumRPCURL    string
	wewelpResolver    *weweLpResolver
	tgtPrice          float64
	mu                sync.RWMutex
}

func NewLiquidtyPositionResolver() *LiquidityPositionResolver {
	return &LiquidityPositionResolver{
		logger:            logrus.WithField("module", "liquidity_position_resolver").Logger,
		thorwalletBaseURL: "https://api-v2-prod.thorwallet.org",
		arbitrumRPCURL:    "https://arbitrum-one-rpc.publicnode.com",
		wewelpResolver:    NewWeWeLpResolver(),
	}
}

type poolPositionResponse struct {
	RuneOrCacaoAddedUsd float64 `json:"runeOrCacaoAddedUsd,string"`
	AssetAddedUsd       float64 `json:"assetAddedUsd,string"`
}

// fetch Thorchain and Maya LP position from Thorwallet api
func (l *LiquidityPositionResolver) GetLiquidityPosition(address string) (float64, error) {
	if address == "" {
		return 0, fmt.Errorf("address cannot be empty")
	}
	url := fmt.Sprintf("%s/pools/positions?addresses=%s", l.thorwalletBaseURL, address)
	resp, err := http.Get(url)
	if err != nil {
		l.logger.Errorf("error fetching liquidity position from %s: %e", url, err)
		return 0, fmt.Errorf("error fetching liquidity position from %s: %e", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		l.logger.Errorf("error fetching liquidity position from %s: %s", url, resp.Status)
		return 0, fmt.Errorf("error fetching liquidity position from %s: %s", url, resp.Status)
	}
	var positions map[string][]poolPositionResponse
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading liquidity position response: %e", err)
	}
	l.logger.Infof("response(%s) : %s", url, string(buf))
	if err := json.Unmarshal(buf, &positions); err != nil {
		return 0, fmt.Errorf("error decoding liquidity position response: %e", err)
	}
	var totalLiquidity float64
	if positions == nil {
		l.logger.Errorf("no liquidity position found for address %s", address)
		return 0, nil
	}
	for _, v := range positions {
		for _, p := range v {
			totalLiquidity += p.RuneOrCacaoAddedUsd + p.AssetAddedUsd
		}
	}
	return totalLiquidity, nil
}

type tgtLPPositionResponse struct {
	StakeAmount float64 `json:"stakedAmount,string"`
	Reward      float64 `json:"reward,string"`
}

type RpcParams struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type RpcResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

func (l *LiquidityPositionResolver) closer(closer io.Closer) {
	if err := closer.Close(); err != nil {
		l.logger.Error(err)
	}
}

var ErrRateLimited = errors.New("rate limited")

func (l *LiquidityPositionResolver) GetTGTStakePosition(address string) (float64, error) {
	if len(address) < 2 || !strings.HasPrefix(address, "0x") {
		return 0, fmt.Errorf("invalid address: must start with 0x, got: %s", address)
	}
	// Remove the '0x' prefix
	address = address[2:]
	// Create parameters array
	from := fmt.Sprintf("0X%040s", "0")
	data := fmt.Sprintf("0xf2801fe7%024s%s%024s%s", "", address, "", address)
	TGTStakeContract := "0x6745c897ab1f4fda9f7700e8be6ea2ee03672759"
	params := []interface{}{
		map[string]interface{}{
			"from": from,
			"data": data,
			"to":   TGTStakeContract,
		},
	}
	// Create RPC request
	rpcReq := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_call",
		Params:  params,
		Id:      1,
	}
	// Convert RPC request to JSON
	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return 0, fmt.Errorf("error marshalling RPC request: %w", err)
	}
	resp, err := http.Post(l.arbitrumRPCURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, fmt.Errorf("error fetching from %s: %w", l.arbitrumRPCURL, err)
	}
	defer l.closer(resp.Body)
	if resp.StatusCode == http.StatusTooManyRequests {
		// rate limited, need to backoff and then retry
		return 0, ErrRateLimited
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching user info of address %s on Arbitrum: %s", address, resp.Status)
	}

	var rpcResp RpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}
	// Remove the '0x' prefix
	result := rpcResp.Result[2:]
	if len(result) < 128 {
		return 0, fmt.Errorf("unexpected RPC result length: %d", len(result))
	}
	// Split the hex string into two 32-byte parts (64 hex characters each)
	stakedAmount := result[:64]
	reward := result[64:]

	// Decode both parts into big integers (uint256)
	TGTStakedInt := new(big.Int)
	TGTRewardInt := new(big.Int) // Decode into big integers (uint256)
	TGTStakedInt.SetString(stakedAmount, 16)
	TGTStakedFloat, _ := new(big.Float).SetInt(TGTStakedInt).Float64()
	TGTRewardInt.SetString(reward, 16)
	TGTRewardFloat, _ := new(big.Float).SetInt(TGTRewardInt).Float64()

	return (TGTStakedFloat / 1e18) + (TGTRewardFloat / 1e18), nil
}

func (l *LiquidityPositionResolver) SetTGTPrice(price float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tgtPrice = price
}

func (l *LiquidityPositionResolver) GetTGTPrice() float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.tgtPrice
}

func (l *LiquidityPositionResolver) GetWeWeLPPosition(address string) (float64, error) {
	return l.wewelpResolver.GetLiquidityPosition(address)
}
