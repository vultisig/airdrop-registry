package liquidity

import (
	"context"
	_ "embed"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

//go:embed weweabi.json
var weweAbi string

//go:embed wewevaultabi.json
var weweVaultAbi string

//go:embed wewepriceabi.json
var wewePriceAbi string

var (
	contractAddress = common.HexToAddress("0x76B4B28194170f9847Ae1566E44dCB4f2D97Ac24")
	vaultAddress    = common.HexToAddress("0x3Fd7957D9F98D46c755685B67dFD8505468A7Cb6")
	weweUsdcAddress = common.HexToAddress("0x6f71796114b9cdaef29a801bc5cdbcb561990eeb")
)

type EVMClient interface {
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

type weweLpResolver struct {
	weweCache  *cache.Cache
	logger     *logrus.Logger
	ethAddress string
}

func NewWeWeLpResolver() *weweLpResolver {
	return &weweLpResolver{
		weweCache:  cache.New(time.Hour, time.Hour),
		logger:     logrus.WithField("module", "wewe_lp_resolver").Logger,
		ethAddress: "https://mainnet.base.org",
	}
}

type UnderlyingOutput struct {
	Amount0   *big.Int `json:"amount0"`
	Amount1   *big.Int `json:"amount1"`
	Fee0      *big.Int `json:"fee0"`
	Fee1      *big.Int `json:"fee1"`
	LeftOver0 *big.Int `json:"leftOver0"`
	LeftOver1 *big.Int `json:"leftOver1"`
}

func (w *weweLpResolver) GetLiquidityPosition(address string) (float64, error) {
	userAddress := common.HexToAddress(address)
	client, err := ethclient.Dial(w.ethAddress)
	if err != nil {
		w.logger.Errorf("Failed to connect to the Ethereum client: %v", err)
		return 0, err
	}
	defer client.Close()

	wewetotalPoolBalance, usdcTotalPoolBalance, err := w.fetchPoolInfo(client)
	if err != nil {
		w.logger.Errorf("Failed to fetch total pool balance: %v", err)
		return 0, err
	}

	totalSupply, err := w.fetchTotalSupply(client)
	if err != nil {
		w.logger.Errorf("Failed to fetch total supply and balance: %v", err)
		return 0, err
	}
	share, err := w.fetchPoolShare(client, userAddress)
	if err != nil {
		w.logger.Errorf("Failed to fetch balance: %v", err)
		return 0, err
	}

	userWEWEAmount := new(big.Int).Mul(share, wewetotalPoolBalance)
	userWEWEAmount.Div(userWEWEAmount, totalSupply)
	userUSDCAmount := new(big.Int).Mul(share, usdcTotalPoolBalance)
	userUSDCAmount.Div(userUSDCAmount, totalSupply)

	wewePrice, err := w.fetchWEWEPrice(client)
	if err != nil {
		w.logger.Errorf("Failed to fetch WEWE price: %v", err)
		return 0, err
	}

	wewePriceFloat := new(big.Float).SetFloat64(wewePrice)
	userWEWEAmountFloat := new(big.Float).SetInt(userWEWEAmount)
	userWEWEUSDAmountFloat := new(big.Float).Mul(userWEWEAmountFloat, wewePriceFloat)
	userWEWEUSDAmount, _ := userWEWEUSDAmountFloat.Int(nil)

	scaledUSDC := new(big.Int).Mul(userUSDCAmount, big.NewInt(1e12)) //wewe has 18 decimals and usdc has 6decimals
	totalUserUsdBalance := new(big.Int).Add(userWEWEUSDAmount, scaledUSDC)
	totalUserUsdBalance = totalUserUsdBalance.Div(totalUserUsdBalance, big.NewInt(1e12))

	value, _ := totalUserUsdBalance.Float64()
	value = value / math.Pow10(6)
	return value, err

}

type Slot0 struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}

func (w *weweLpResolver) fetchWEWEPrice(client *ethclient.Client) (float64, error) {
	if price, ok := w.weweCache.Get("wewePrice"); ok {
		return price.(float64), nil
	}

	parsed, err := abi.JSON(strings.NewReader(wewePriceAbi))
	if err != nil {
		return 0, fmt.Errorf("failed to parse ABI: %v", err)
	}

	data, err := parsed.Pack("slot0")
	if err != nil {
		return 0, fmt.Errorf("failed to pack data: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &weweUsdcAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to call contract: %v", err)
	}

	var slot0 Slot0
	err = parsed.UnpackIntoInterface(&slot0, "slot0", result)
	if err != nil {
		return 0, fmt.Errorf("failed to unpack result: %v", err)
	}

	tickFloat, _ := new(big.Float).SetInt(slot0.Tick).Float64()
	price := math.Pow(1.0001, tickFloat) * math.Pow(10, 12)
	w.weweCache.Set("wewePrice", price, cache.DefaultExpiration)
	return price, nil
}

func (w *weweLpResolver) fetchPoolInfo(client EVMClient) (*big.Int, *big.Int, error) {
	if wewetotalPoolBalance, ok := w.weweCache.Get("weweTotalPoolBalance"); ok {
		if usdcTotalPoolBalance, ok := w.weweCache.Get("usdcTotalPoolBalance"); ok {
			return wewetotalPoolBalance.(*big.Int), usdcTotalPoolBalance.(*big.Int), nil
		}
	}

	parsed, err := abi.JSON(strings.NewReader(weweVaultAbi))
	if err != nil {
		w.logger.Errorf("Failed to parse ABI: %v", err)
		return nil, nil, err
	}

	data, err := parsed.Pack("totalUnderlyingWithFeesAndLeftOver", vaultAddress)
	if err != nil {
		w.logger.Errorf("Failed to pack data: %v", err)
		return nil, nil, err
	}

	msgAux := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msgAux, nil)
	if err != nil {
		w.logger.Errorf("Failed to call contract: %v", err)
		return nil, nil, err
	}

	var underlying struct {
		Underlying UnderlyingOutput
	}
	err = parsed.UnpackIntoInterface(&underlying, "totalUnderlyingWithFeesAndLeftOver", result)
	if err != nil {
		w.logger.Errorf("Failed to unpack result: %v", err)
		return nil, nil, err
	}
	amount0 := new(big.Int).Add(underlying.Underlying.Amount0, underlying.Underlying.Fee0)
	amount0.Add(amount0, underlying.Underlying.LeftOver0)

	amount1 := new(big.Int).Add(underlying.Underlying.Amount1, underlying.Underlying.Fee1)
	amount1.Add(amount1, underlying.Underlying.LeftOver1)
	w.weweCache.Set("weweTotalPoolBalance", amount0, cache.DefaultExpiration)
	w.weweCache.Set("usdcTotalPoolBalance", amount1, cache.DefaultExpiration)

	return amount0, amount1, nil
}

func (w *weweLpResolver) fetchTotalSupply(client EVMClient) (*big.Int, error) {
	if totalSupply, ok := w.weweCache.Get("totalSupply"); ok {
		return totalSupply.(*big.Int), nil
	}
	parsed, err := abi.JSON(strings.NewReader(weweAbi))
	if err != nil {
		w.logger.Errorf("Failed to parse ABI: %v", err)
		return nil, err
	}
	data, err := parsed.Pack("totalSupply")
	if err != nil {
		w.logger.Errorf("Failed to pack data: %v", err)
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &vaultAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		w.logger.Error("Failed to call contract: %v", err)
		return nil, err
	}
	var totalSupply *big.Int
	err = parsed.UnpackIntoInterface(&totalSupply, "totalSupply", result)
	if err != nil {
		w.logger.Errorf("Failed to unpack result: %v", err)
		return nil, err
	}
	w.weweCache.Set("totalSupply", totalSupply, cache.DefaultExpiration)
	return totalSupply, nil
}

func (w *weweLpResolver) fetchPoolShare(client EVMClient, userAddress common.Address) (*big.Int, error) {
	parsed, err := abi.JSON(strings.NewReader(weweAbi))
	if err != nil {
		w.logger.Errorf("Failed to parse ABI: %v", err)
		return nil, err
	}
	dataBalance, err := parsed.Pack("balanceOf", userAddress)
	if err != nil {
		w.logger.Errorf("Failed to pack data: %v", err)
		return nil, err
	}

	msgBalance := ethereum.CallMsg{
		To:   &vaultAddress,
		Data: dataBalance,
	}

	resultBalance, err := client.CallContract(context.Background(), msgBalance, nil)
	if err != nil {
		w.logger.Errorf("Failed to call contract: %v", err)
		return nil, err
	}
	var balance *big.Int
	err = parsed.UnpackIntoInterface(&balance, "balanceOf", resultBalance)
	if err != nil {
		w.logger.Errorf("Failed to unpack result: %v", err)
		return nil, err
	}
	return balance, nil
}
