package tasks

const (
	TypeBalanceFetch      = "balance:fetch"
	TypePointsCalculation = "points:calculate"
)

type BalanceFetchPayload struct {
	VaultID string
}

type PointsCalculationPayload struct {
	VaultID string
}
