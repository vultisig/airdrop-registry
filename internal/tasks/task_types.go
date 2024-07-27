package tasks

const (
	TypeBalanceFetch            = "balance:fetch"
	TypeBalanceFetchParent      = "balance:fetch:parent"
	TypePointsCalculation       = "points:calculate"
	TypePointsCalculationParent = "points:calculate:parent"
	TypePriceFetch              = "price:fetch"
	TypePriceFetchParent        = "price:fetch:parent"
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
