package utils

import "math"

func CalculateShare(averageBalance, totalUSD float64) float64 {
	if totalUSD <= 0 || averageBalance <= 0 {
		return 0
	}

	share := (averageBalance / totalUSD) * 100

	// Ensure the share is between 0 and 100
	if share < 0 {
		return 0
	} else if share > 100 {
		return 100
	}

	// Avoid NaN by checking for infinity
	if math.IsInf(share, 0) || math.IsNaN(share) {
		return 0
	}

	return share
}
