package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const abiJSON = `[{"inputs":[{"internalType":"contract IUniswapV3Factory","name":"factory_","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"user","type":"address"},{"indexed":false,"internalType":"uint256","name":"burnAmount0","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"burnAmount1","type":"uint256"}],"name":"LPBurned","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint24[]","name":"feeTiers","type":"uint24[]"}],"name":"LogAddPools","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address[]","name":"routers","type":"address[]"}],"name":"LogBlacklistRouters","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"receiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"burnAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount0Out","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1Out","type":"uint256"}],"name":"LogBurn","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"fee0","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"fee1","type":"uint256"}],"name":"LogCollectedFees","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"receiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"mintAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount0In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1In","type":"uint256"}],"name":"LogMint","type":"event"},{"anonymous":false,"inputs":[{"components":[{"components":[{"internalType":"uint128","name":"liquidity","type":"uint128"},{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"}],"internalType":"struct PositionLiquidity[]","name":"burns","type":"tuple[]"},{"components":[{"internalType":"uint128","name":"liquidity","type":"uint128"},{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"}],"internalType":"struct PositionLiquidity[]","name":"mints","type":"tuple[]"},{"components":[{"internalType":"bytes","name":"payload","type":"bytes"},{"internalType":"address","name":"router","type":"address"},{"internalType":"uint256","name":"amountIn","type":"uint256"},{"internalType":"uint256","name":"expectedMinReturn","type":"uint256"},{"internalType":"bool","name":"zeroForOne","type":"bool"}],"internalType":"struct SwapPayload","name":"swap","type":"tuple"},{"internalType":"uint256","name":"minBurn0","type":"uint256"},{"internalType":"uint256","name":"minBurn1","type":"uint256"},{"internalType":"uint256","name":"minDeposit0","type":"uint256"},{"internalType":"uint256","name":"minDeposit1","type":"uint256"}],"indexed":false,"internalType":"struct Rebalance","name":"rebalanceParams","type":"tuple"},{"indexed":false,"internalType":"uint256","name":"swapDelta0","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"swapDelta1","type":"uint256"}],"name":"LogRebalance","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address[]","name":"pools","type":"address[]"}],"name":"LogRemovePools","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"minter","type":"address"}],"name":"LogRestrictedMint","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"init0","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"init1","type":"uint256"}],"name":"LogSetInits","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"newManager","type":"address"}],"name":"LogSetManager","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint16","name":"managerFeeBPS","type":"uint16"}],"name":"LogSetManagerFeeBPS","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address[]","name":"routers","type":"address[]"}],"name":"LogWhitelistRouters","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"amount0","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1","type":"uint256"}],"name":"LogWithdrawManagerBalance","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[{"internalType":"uint24[]","name":"feeTiers_","type":"uint24[]"}],"name":"addPools","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address[]","name":"routers_","type":"address[]"}],"name":"blacklistRouters","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"burnAmount_","type":"uint256"},{"internalType":"address","name":"receiver_","type":"address"}],"name":"burn","outputs":[{"internalType":"uint256","name":"amount0","type":"uint256"},{"internalType":"uint256","name":"amount1","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"collectFees","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"factory","outputs":[{"internalType":"contract IUniswapV3Factory","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeManager","outputs":[{"internalType":"contract IFeeManager","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getPools","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getRanges","outputs":[{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getRouters","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"init0","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"init1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"name_","type":"string"},{"internalType":"string","name":"symbol_","type":"string"},{"components":[{"internalType":"uint24[]","name":"feeTiers","type":"uint24[]"},{"internalType":"address","name":"token0","type":"address"},{"internalType":"address","name":"token1","type":"address"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"init0","type":"uint256"},{"internalType":"uint256","name":"init1","type":"uint256"},{"internalType":"address","name":"manager","type":"address"},{"internalType":"address[]","name":"routers","type":"address[]"}],"internalType":"struct InitializePayload","name":"params_","type":"tuple"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"manager","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"managerBalance0","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"managerBalance1","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"managerFeeBPS","outputs":[{"internalType":"uint16","name":"","type":"uint16"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"mintAmount_","type":"uint256"},{"internalType":"address","name":"receiver_","type":"address"}],"name":"mint","outputs":[{"internalType":"uint256","name":"amount0","type":"uint256"},{"internalType":"uint256","name":"amount1","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"components":[{"internalType":"uint128","name":"liquidity","type":"uint128"},{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"}],"internalType":"struct PositionLiquidity[]","name":"burns","type":"tuple[]"},{"components":[{"internalType":"uint128","name":"liquidity","type":"uint128"},{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"}],"internalType":"struct PositionLiquidity[]","name":"mints","type":"tuple[]"},{"components":[{"internalType":"bytes","name":"payload","type":"bytes"},{"internalType":"address","name":"router","type":"address"},{"internalType":"uint256","name":"amountIn","type":"uint256"},{"internalType":"uint256","name":"expectedMinReturn","type":"uint256"},{"internalType":"bool","name":"zeroForOne","type":"bool"}],"internalType":"struct SwapPayload","name":"swap","type":"tuple"},{"internalType":"uint256","name":"minBurn0","type":"uint256"},{"internalType":"uint256","name":"minBurn1","type":"uint256"},{"internalType":"uint256","name":"minDeposit0","type":"uint256"},{"internalType":"uint256","name":"minDeposit1","type":"uint256"}],"internalType":"struct Rebalance","name":"rebalanceParams_","type":"tuple"}],"name":"rebalance","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"pools_","type":"address[]"}],"name":"removePools","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"restrictedMint","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"feeManager_","type":"address"}],"name":"setFeeManager","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"init0_","type":"uint256"},{"internalType":"uint256","name":"init1_","type":"uint256"}],"name":"setInits","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"manager_","type":"address"}],"name":"setManager","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint16","name":"managerFeeBPS_","type":"uint16"}],"name":"setManagerFeeBPS","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"minter_","type":"address"}],"name":"setRestrictedMint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"token0","outputs":[{"internalType":"contract IERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"token1","outputs":[{"internalType":"contract IERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"amount0Owed_","type":"uint256"},{"internalType":"uint256","name":"amount1Owed_","type":"uint256"},{"internalType":"bytes","name":"","type":"bytes"}],"name":"uniswapV3MintCallback","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"routers_","type":"address[]"}],"name":"whitelistRouters","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"withdrawManagerBalance","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
const abiJSONAux = `[{"inputs":[{"internalType":"contract IUniswapV3Factory","name":"factory_","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"factory","outputs":[{"internalType":"contract IUniswapV3Factory","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range[]","name":"ranges_","type":"tuple[]"},{"internalType":"address","name":"token0_","type":"address"},{"internalType":"address","name":"token1_","type":"address"},{"internalType":"address","name":"vaultV2_","type":"address"}],"name":"token0AndToken1ByRange","outputs":[{"components":[{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Amount[]","name":"amount0s","type":"tuple[]"},{"components":[{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Amount[]","name":"amount1s","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range[]","name":"ranges_","type":"tuple[]"},{"internalType":"address","name":"token0_","type":"address"},{"internalType":"address","name":"token1_","type":"address"},{"internalType":"address","name":"vaultV2_","type":"address"}],"name":"token0AndToken1PlusFeesByRange","outputs":[{"components":[{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Amount[]","name":"amount0s","type":"tuple[]"},{"components":[{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Amount[]","name":"amount1s","type":"tuple[]"},{"components":[{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Amount[]","name":"fee0s","type":"tuple[]"},{"components":[{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Amount[]","name":"fee1s","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"contract IArrakisV2","name":"vault_","type":"address"}],"name":"totalLiquidity","outputs":[{"components":[{"internalType":"uint128","name":"liquidity","type":"uint128"},{"components":[{"internalType":"int24","name":"lowerTick","type":"int24"},{"internalType":"int24","name":"upperTick","type":"int24"},{"internalType":"uint24","name":"feeTier","type":"uint24"}],"internalType":"struct Range","name":"range","type":"tuple"}],"internalType":"struct PositionLiquidity[]","name":"liquidities","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"contract IArrakisV2","name":"vault_","type":"address"}],"name":"totalUnderlying","outputs":[{"internalType":"uint256","name":"amount0","type":"uint256"},{"internalType":"uint256","name":"amount1","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"contract IArrakisV2","name":"vault_","type":"address"},{"internalType":"uint160","name":"sqrtPriceX96_","type":"uint160"}],"name":"totalUnderlyingAtPrice","outputs":[{"internalType":"uint256","name":"amount0","type":"uint256"},{"internalType":"uint256","name":"amount1","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"contract IArrakisV2","name":"vault_","type":"address"}],"name":"totalUnderlyingWithFees","outputs":[{"internalType":"uint256","name":"amount0","type":"uint256"},{"internalType":"uint256","name":"amount1","type":"uint256"},{"internalType":"uint256","name":"fee0","type":"uint256"},{"internalType":"uint256","name":"fee1","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"contract IArrakisV2","name":"vault_","type":"address"}],"name":"totalUnderlyingWithFeesAndLeftOver","outputs":[{"components":[{"internalType":"uint256","name":"amount0","type":"uint256"},{"internalType":"uint256","name":"amount1","type":"uint256"},{"internalType":"uint256","name":"fee0","type":"uint256"},{"internalType":"uint256","name":"fee1","type":"uint256"},{"internalType":"uint256","name":"leftOver0","type":"uint256"},{"internalType":"uint256","name":"leftOver1","type":"uint256"}],"internalType":"struct UnderlyingOutput","name":"underlying","type":"tuple"}],"stateMutability":"view","type":"function"}]`
const COMMON_POOL_CONTRACT_ABI = `[
	{
		"inputs": [],
		"name": "slot0",
		"outputs": [
			{ "internalType": "uint160", "name": "sqrtPriceX96", "type": "uint160" },
			{ "internalType": "int24", "name": "tick", "type": "int24" },
			{ "internalType": "uint16", "name": "observationIndex", "type": "uint16" },
			{ "internalType": "uint16", "name": "observationCardinality", "type": "uint16" },
			{ "internalType": "uint16", "name": "observationCardinalityNext", "type": "uint16" },
			{ "internalType": "uint8", "name": "feeProtocol", "type": "uint8" },
			{ "internalType": "bool", "name": "unlocked", "type": "bool" }
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "liquidity",
		"outputs": [{ "internalType": "uint128", "name": "", "type": "uint128" }],
		"stateMutability": "view",
		"type": "function"
	}
]`

var (
	contractAddress     = common.HexToAddress("0x76B4B28194170f9847Ae1566E44dCB4f2D97Ac24")
	vaultAddress        = common.HexToAddress("0x3Fd7957D9F98D46c755685B67dFD8505468A7Cb6")
	ArrakisTokenAddress = common.HexToAddress("0x3Fd7957D9F98D46c755685B67dFD8505468A7Cb6")
	testUser            = common.HexToAddress("0x6159DfAEbea7522a493ad1d402B3B5aaFB8e1E37")
	weweUsdcAddress     = common.HexToAddress("0x6f71796114b9cdaef29a801bc5cdbcb561990eeb")
)

type UnderlyingOutput struct {
	Amount0   *big.Int `json:"amount0"`
	Amount1   *big.Int `json:"amount1"`
	Fee0      *big.Int `json:"fee0"`
	Fee1      *big.Int `json:"fee1"`
	LeftOver0 *big.Int `json:"leftOver0"`
	LeftOver1 *big.Int `json:"leftOver1"`
}

func main() {
	client, err := ethclient.Dial("https://mainnet.base.org")
	if err != nil {
		log.Fatalf("Failed to connect to the Base network: %v", err)
	}

	wewetotalPoolBalance, usdcTotalPoolBalance, err := fetchAuxInfo(client)
	if err != nil {
		log.Fatalf("Failed to fetch aux info: %v", err)
	}

	fmt.Printf("WEWE total pool balance: %s\n", wewetotalPoolBalance.String())
	fmt.Printf("USDC total pool balance: %s\n", usdcTotalPoolBalance.String())

	totalSupply, balance, err := fetchTotalSupplyAndBalance(client, testUser)
	if err != nil {
		log.Fatalf("Failed to fetch total supply and balance: %v", err)
	}

	userWEWEAmount := new(big.Int).Mul(balance, wewetotalPoolBalance)
	userWEWEAmount.Div(userWEWEAmount, totalSupply)
	userUSDCAmount := new(big.Int).Mul(balance, usdcTotalPoolBalance)
	userUSDCAmount.Div(userUSDCAmount, totalSupply)

	fmt.Printf("User WEWE amount: %s\n", userWEWEAmount.String()) //wewe has 18 decimals
	fmt.Printf("User USDC amount: %s\n", userUSDCAmount.String()) //usdc has 6 decimals

	wewePrice, err := FetchWEWEPrice(client)
	if err != nil {
		log.Fatalf("Failed to fetch WEWE price: %v", err)
	}

	fmt.Printf("WEWE price: %f\n", wewePrice)

	// Convert WEWE price to big.Float
	wewePriceFloat := new(big.Float).SetFloat64(wewePrice)

	// Convert userWEWEAmount to big.Float
	userWEWEAmountFloat := new(big.Float).SetInt(userWEWEAmount)

	// Multiply to get USD amount
	userWEWEUSDAmountFloat := new(big.Float).Mul(userWEWEAmountFloat, wewePriceFloat)

	// Convert result back to big.Int (rounding down)
	userWEWEUSDAmount, _ := userWEWEUSDAmountFloat.Int(nil)

	fmt.Printf("User WEWE USD amount: %s\n", userWEWEUSDAmount.String())
	scaledUSDC := new(big.Int).Mul(userUSDCAmount, big.NewInt(1e12))
	totalUserUsdBalance := new(big.Int).Add(userWEWEUSDAmount, scaledUSDC)

	fmt.Printf("Total user USD balance: %s\n", totalUserUsdBalance.String())
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

func FetchWEWEPrice(client *ethclient.Client) (float64, error) {
	parsed, err := abi.JSON(strings.NewReader(COMMON_POOL_CONTRACT_ABI))
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
	return price, nil
}

func fetchAuxInfo(client *ethclient.Client) (*big.Int, *big.Int, error) {
	parsedAux, err := abi.JSON(strings.NewReader(abiJSONAux))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	dataAux, err := parsedAux.Pack("totalUnderlyingWithFeesAndLeftOver", vaultAddress)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}

	msgAux := ethereum.CallMsg{
		To:   &contractAddress,
		Data: dataAux,
	}

	resultAux, err := client.CallContract(context.Background(), msgAux, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	var underlying struct {
		Underlying UnderlyingOutput
	}
	err = parsedAux.UnpackIntoInterface(&underlying, "totalUnderlyingWithFeesAndLeftOver", resultAux)
	if err != nil {
		log.Fatalf("Failed to unpack result: %v", err)
	}
	amount0 := new(big.Int).Add(underlying.Underlying.Amount0, underlying.Underlying.Fee0)
	amount0.Add(amount0, underlying.Underlying.LeftOver0)

	amount1 := new(big.Int).Add(underlying.Underlying.Amount1, underlying.Underlying.Fee1)
	amount1.Add(amount1, underlying.Underlying.LeftOver1)

	return amount0, amount1, nil
}

func fetchTotalSupplyAndBalance(client *ethclient.Client, userAddress common.Address) (*big.Int, *big.Int, error) {
	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}
	data, err := parsed.Pack("totalSupply")
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &ArrakisTokenAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}
	var totalSupply *big.Int
	err = parsed.UnpackIntoInterface(&totalSupply, "totalSupply", result)
	if err != nil {
		log.Fatalf("Failed to unpack result: %v", err)
	}
	fmt.Printf("Total Supply: %s\n", totalSupply.String())

	dataBalance, err := parsed.Pack("balanceOf", userAddress)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}

	msgBalance := ethereum.CallMsg{
		To:   &ArrakisTokenAddress,
		Data: dataBalance,
	}

	resultBalance, err := client.CallContract(context.Background(), msgBalance, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}
	var balance *big.Int
	err = parsed.UnpackIntoInterface(&balance, "balanceOf", resultBalance)
	if err != nil {
		log.Fatalf("Failed to unpack result: %v", err)
	}
	fmt.Printf("Balance: %s\n", balance.String())

	return totalSupply, balance, nil
}
