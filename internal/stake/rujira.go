package stake

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	RujiraContractAddr = "thor13g83nn5ef4qzqeafp0508dnvkvm0zqr3sj7eefcn5umu65gqluusrml5cr"
)

type RujiraStakeResolver struct {
	logger              *logrus.Logger
	thornodeBaseAddress string
	chainDecimal        int
	rujiraPrice         float64
	mu                  sync.RWMutex
}

func NewRujiraStakeResolver() *RujiraStakeResolver {
	return &RujiraStakeResolver{
		thornodeBaseAddress: "https://thornode.ninerealms.com",
		chainDecimal:        8,
		logger:              logrus.WithField("module", "stake_resolver").Logger,
	}
}

func (s *RujiraStakeResolver) GetRujiraAutoCompoundStake(address string) (float64, error) {
	url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", s.thornodeBaseAddress, address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error making GET request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return 0, fmt.Errorf("rate limited while fetching Rujira compound stake")
	}
	var result autoCompoundStakeResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}
	var balance int64
	for _, token := range result.Balances {
		if strings.EqualFold(token.Denom, "x/staking-x/ruji") {
			balance = token.Amount
			break
		}
	}
	return (float64(balance) * math.Pow10(-s.chainDecimal)) * s.GetRujiraPrice(), nil
}

func (s *RujiraStakeResolver) GetRujiraSimpleStake(address string) (float64, error) {
	param := `{ "account": { "addr": "` + address + `" } }`
	encodedParam := base64.StdEncoding.EncodeToString([]byte(param))
	url := fmt.Sprintf("%s/cosmwasm/wasm/v1/contract/%s/smart/%s", s.thornodeBaseAddress, RujiraContractAddr, encodedParam)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error making GET request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return 0, fmt.Errorf("rate limited while fetching Rujira simple stake")
	}
	var result simpleStakeResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}
	if result.Data.Addr == address {
		return (float64(result.Data.Bonded) * math.Pow10(-s.chainDecimal)) * s.GetRujiraPrice(), nil
	}
	return 0, nil
}

type simpleStakeResp struct {
	Data struct {
		Addr           string `json:"addr"`
		Bonded         int64  `json:"bonded,string"`
		PendingRevenue int64  `json:"pending_revenue,string"`
	} `json:"data"`
}

type autoCompoundStakeResp struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount int64  `json:"amount,string"`
	} `json:"balances"`
}

func (s *RujiraStakeResolver) SetRujiraPrice(price float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rujiraPrice = price

}
func (s *RujiraStakeResolver) GetRujiraPrice() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rujiraPrice
}
