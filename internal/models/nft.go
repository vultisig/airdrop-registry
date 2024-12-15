package models

import "github.com/vultisig/airdrop-registry/internal/common"

type NFTCollection struct {
	Chain             common.Chain
	CollectionAddress string
	CollectionSlug    string
}
