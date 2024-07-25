package tasks

const (
	TypeVaultBalanceFetch        = "vault:balance:fetch"
	TypeBalanceFetch             = "balance:fetch"
	TypeBalanceFetchAll          = "balance:fetch_all"
	TypePointsCalculation        = "points:calculate"
	TypePointsCalculationAll     = "points:calculate:all"
	TypePriceFetch               = "price:fetch"
	TypePriceFetchAllActivePairs = "price:fetch_all_active_pairs"
)

type BalanceFetchPayload struct {
	ECDSA   string
	EDDSA   string
	Chain   string
	Address string
}

type PointsCalculationPayload struct {
	ECDSA string
	EDDSA string
}

type PriceFetchPayload struct {
	Chain string
	Token string
}
