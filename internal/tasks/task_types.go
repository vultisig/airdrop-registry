package tasks

const (
	TypeVaultBalanceFetch = "vault:balance:fetch"
	TypeBalanceFetch      = "balance:fetch"
	TypePointsCalculation = "points:calculate"
	TypePriceFetch        = "price:fetch"
)

// type VaultBalanceFetchPayload struct {
// 	ECCDSA string
// 	EDDSA  string
// }

type BalanceFetchPayload struct {
	ECCDSA  string
	EDDSA   string
	Chain   string
	Address string
}

type PointsCalculationPayload struct {
	ECCDSA string
	EDDSA  string
}

type PriceFetchPayload struct {
	Chain string
	Token string
}
