package balance

import (
	"github.com/sirupsen/logrus"
)

const vultisigApiProxy = "https://api.vultisig.com"

// BalanceResolver is to fetch address balances
type BalanceResolver struct {
	logger *logrus.Logger
}

func NewBalanceResolver() (*BalanceResolver, error) {
	return &BalanceResolver{
		logger: logrus.WithField("module", "balance_resolver").Logger,
	}, nil
}
