package balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type SubscanResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	GeneratedAt int64  `json:"generated_at"`
	Data        struct {
		Account struct {
			Address        string `json:"address"`
			Balance        string `json:"balance"`
			Lock           string `json:"lock"`
			BalanceLock    string `json:"balance_lock"`
			IsEvmContract  bool   `json:"is_evm_contract"`
			AccountDisplay struct {
				Address string `json:"address"`
			} `json:"account_display"`
			SubstrateAccount   interface{} `json:"substrate_account"`
			EvmAccount         string      `json:"evm_account"`
			RegistrarInfo      interface{} `json:"registrar_info"`
			CountExtrinsic     int         `json:"count_extrinsic"`
			Reserved           string      `json:"reserved"`
			Bonded             string      `json:"bonded"`
			Unbonding          string      `json:"unbonding"`
			DemocracyLock      string      `json:"democracy_lock"`
			ConvictionLock     string      `json:"conviction_lock"`
			ElectionLock       string      `json:"election_lock"`
			StakingInfo        interface{} `json:"staking_info"`
			Nonce              int         `json:"nonce"`
			Role               string      `json:"role"`
			Stash              string      `json:"stash"`
			IsCouncilMember    bool        `json:"is_council_member"`
			IsTechcommMember   bool        `json:"is_techcomm_member"`
			IsRegistrar        bool        `json:"is_registrar"`
			IsFellowshipMember bool        `json:"is_fellowship_member"`
			IsModuleAccount    bool        `json:"is_module_account"`
			AssetsTag          interface{} `json:"assets_tag"`
			IsErc20            bool        `json:"is_erc20"`
			IsErc721           bool        `json:"is_erc721"`
			Vesting            interface{} `json:"vesting"`
			Proxy              struct{}    `json:"proxy"`
			Multisig           struct{}    `json:"multisig"`
			Delegate           interface{} `json:"delegate"`
		} `json:"account"`
	} `json:"data"`
}

func FetchPolkadotBalanceOfAddress(address string) (float64, error) {
	payload := fmt.Sprintf(`{"key":"%s"}`, address)
	response, err := http.Post(
		"https://polkadot.api.subscan.io/api/v2/scan/search",
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on Polkadot: %v", address, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var subscanResp SubscanResponse
	err = json.Unmarshal(body, &subscanResp)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	if subscanResp.Code != 0 {
		return 0, fmt.Errorf("error from subscan API: %s", subscanResp.Message)
	}

	balanceStr := subscanResp.Data.Account.Balance
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting balance to float: %v", err)
	}

	return balance, nil
}
