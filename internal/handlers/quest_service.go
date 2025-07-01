package handlers

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/services"
)

// store eth address of all vaults in memory (for 10K vaults, it will take around 800KB of memory)
type QuestService struct {
	sync.RWMutex
	logger          *logrus.Logger
	storage         *services.Storage
	ethAddressStore map[string]uint
}

func NewQuestService(storage *services.Storage) (QuestService, error) {
	questService := QuestService{
		logger:          logrus.WithField("module", "quest_service").Logger,
		storage:         storage,
		ethAddressStore: make(map[string]uint),
	}
	if err := questService.initialize(); err != nil {
		return QuestService{}, err
	}
	return questService, nil
}

// initialize the quest service
// step1: Try to fetch all eth address from coin store
// step2: Try to generate missing eth address from vaults
func (q *QuestService) initialize() error {
	q.Lock()
	defer q.Unlock()
	q.ethAddressStore = make(map[string]uint)

	//step1
	var coinPageId uint64
	for {
		coins, err := q.storage.GetCoinsWithPage(coinPageId, 10000)
		if err != nil {
			return err
		}
		if coinPageId > 0 {
			break
		}
		if len(coins) == 0 {
			break
		}
		coinPageId = uint64(coins[len(coins)-1].ID + 1)
		for _, coin := range coins {
			for _, evmChain := range common.EVMChains {
				if coin.Chain == evmChain {
					q.ethAddressStore[coin.Address] = coin.VaultID
					break
				}
			}
		}
	}

	//step2
	var vaultPageId uint
	for {
		vaults, err := q.storage.GetVaultsWithPage(vaultPageId, 1000)
		if err != nil {
			return err
		}
		if len(vaults) == 0 {
			break
		}
		if vaultPageId > 0 {
			break
		}
		vaultPageId = vaults[len(vaults)-1].ID + 1
		for _, vault := range vaults {
			exits := false
			for _, v := range q.ethAddressStore {
				if v == vault.ID {
					exits = true
					break
				}
			}
			if !exits {
				//generate eth address from vault
				ethAddress, err := vault.GetAddress(common.Ethereum)
				if err != nil {
					q.logger.Errorf("failed to get eth address for vault %d: %v", vault.ID, err)
					continue
				}
				q.ethAddressStore[ethAddress] = vault.ID
			}
		}
	}
	return nil
}

func (q *QuestService) Exists(ethAddress string) bool {
	q.RLock()
	defer q.RUnlock()
	_, exists := q.ethAddressStore[ethAddress]
	return exists
}

func (q *QuestService) Add(vault models.Vault) {
	q.Lock()
	defer q.Unlock()
	ethAddress, err := vault.GetAddress(common.Ethereum)
	if err == nil {
		q.ethAddressStore[ethAddress] = vault.ID
	}
}

func (q *QuestService) Remove(vaultId uint) {
	q.Lock()
	defer q.Unlock()
	for ethAddress, id := range q.ethAddressStore {
		if id == vaultId {
			delete(q.ethAddressStore, ethAddress)
		}
	}
}
