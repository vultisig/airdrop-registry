package tasks

const (
	TypeVaultBalanceFetch        = "vault:balance:fetch"
	TypeBalanceFetch             = "balance:fetch"
	TypeBalanceFetchAll          = "balance:fetch_all"
	TypePointsCalculation        = "points:calculate"
	TypePriceFetch               = "price:fetch"
	TypePriceFetchAllActivePairs = "price:fetch_all_active_pairs"
)

// type VaultBalanceFetchPayload struct {
// 	ecdsa string
// 	EDDSA  string
// }

type BalanceFetchPayload struct {
	ecdsa   string
	EDDSA   string
	Chain   string
	Address string
}

type PointsCalculationPayload struct {
	ecdsa string
	EDDSA string
}

type PriceFetchPayload struct {
	Chain string
	Token string
}
