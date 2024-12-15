package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/balance"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/liquidity"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/utils"
)

// PointWorker is a worker that processes points
type PointWorker struct {
	logger                 *logrus.Logger
	storage                *Storage
	priceResolver          *PriceResolver
	balanceResolver        *balance.BalanceResolver
	lpResolver             *liquidity.LiquidityPositionResolver
	saverResolver          *liquidity.SaverPositionResolver
	startCoinID            int64
	wg                     *sync.WaitGroup
	stopChan               chan struct{}
	cfg                    *config.Config
	isJobInProgress        bool
	whitelistNFTCollection []models.NFTCollection
}

func NewPointWorker(cfg *config.Config, storage *Storage, priceResolver *PriceResolver, balanceResolver *balance.BalanceResolver) (*PointWorker, error) {

	if nil == storage {
		return nil, fmt.Errorf("storage is nil")
	}
	if nil == priceResolver {
		return nil, fmt.Errorf("priceResolver is nil")
	}

	return &PointWorker{
		logger:          logrus.WithField("module", "point_worker").Logger,
		storage:         storage,
		priceResolver:   priceResolver,
		balanceResolver: balanceResolver,
		lpResolver:      liquidity.NewLiquidtyPositionResolver(),
		saverResolver:   liquidity.NewSaverPositionResolver(),
		startCoinID:     cfg.Worker.StartID,
		stopChan:        make(chan struct{}),
		wg:              &sync.WaitGroup{},
		cfg:             cfg,
		whitelistNFTCollection: []models.NFTCollection{
			{
				Chain:             common.Ethereum,
				CollectionAddress: "0xa98b29a8f5a247802149c268ecf860b8308b7291",
				CollectionSlug:    "thorguards",
			},
		},
	}, nil
}

func (p *PointWorker) Run() error {
	p.wg.Add(1)
	go p.scheduler()
	return nil
}
func (p *PointWorker) scheduler() {
	p.logger.Info("start scheduler")
	defer p.logger.Info("scheduler stopped")
	defer p.wg.Done()
	for {
		select {
		case <-p.stopChan:
			return
		case <-time.After(time.Minute):
			p.ensureJobs()
		}
	}
}

func (p *PointWorker) ensureJobs() {
	lastJob, err := p.storage.GetLastJob()
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			p.logger.Errorf("failed to get last job: %v", err)
			return
		}
		// create a job for today
		lastJob = &models.Job{
			JobDate:    time.Now(),
			Multiplier: 1,
			IsSuccess:  false,
		}
		if err := p.storage.CreateJob(lastJob); err != nil {
			p.logger.Errorf("failed to create job: %v", err)
			return
		}
	}

	lastJob, err = p.storage.GetLastJob()
	if err != nil {
		p.logger.Errorf("failed to get last job: %v", err)
		return
	}

	if !lastJob.IsSuccess {
		p.startJob(lastJob)
		return
	}

	multiplier := lastJob.DaysSince()
	if multiplier < 1 {
		// last job has been finished , but not 24 hours yet
		return
	}
	newJob := &models.Job{
		JobDate:    time.Now(),
		Multiplier: multiplier,
		IsSuccess:  false,
	}
	if err := p.storage.CreateJob(newJob); err != nil {
		p.logger.Errorf("failed to create job: %v", err)
		return
	}
	p.startJob(newJob)
}

func (p *PointWorker) startJob(job *models.Job) {
	if p.isJobInProgress {
		return
	}
	if job.CurrentID == 0 {
		p.logger.Infof("start job %s", job.JobDate.Format("2006-01-02"))
	} else {
		p.logger.Infof("continue job %s from %d", job.JobDate.Format("2006-01-02"), job.CurrentID)
	}

	if job.CurrentVaultID == 0 {
		p.logger.Infof("start lp calculation job %s", job.JobDate.Format("2006-01-02"))
	} else {
		p.logger.Infof("continue lp calculation job %s from %d", job.JobDate.Format("2006-01-02"), job.CurrentVaultID)
	}

	if err := p.updateCoinPrice(); err != nil {
		p.logger.Errorf("failed to update coin prices: %e", err)
		return
	}

	p.wg.Add(1)
	workChan := make(chan models.CoinDBModel)
	// worker channel for lp calculation (key is vault id and value is vault addresses)
	positionWorkerChan := make(chan models.VaultAddress)
	go p.taskProvider(job, workChan, positionWorkerChan)

	// We have 2 type of concurrent workers, one for updating balance and one for updating position
	for i := 0; i < 2; i++ {
		p.wg.Add(1)
		idx := i
		go p.activePositionWorker(idx, positionWorkerChan, *job)
	}
	for i := 0; i < int(p.cfg.Worker.Concurrency); i++ {
		p.wg.Add(1)
		idx := i
		go p.taskWorker(idx, workChan, *job)
	}

}

func (p *PointWorker) Stop() {
	close(p.stopChan)
	p.wg.Wait()
}

func (p *PointWorker) taskProvider(job *models.Job, workChan chan models.CoinDBModel, positionWorkerChan chan models.VaultAddress) {
	defer p.wg.Done()
	p.isJobInProgress = true
	defer func() {
		p.isJobInProgress = false
	}()
	currentVaultId := job.CurrentVaultID
	currentID := uint64(job.CurrentID)
	// refresh bond providers
	if err := p.balanceResolver.GetTHORChainBondProviders(); err != nil {
		p.logger.Errorf("failed to get thorchain bond providers: %v", err)
	}
	if err := p.balanceResolver.GetTHORChainRuneProviders(); err != nil {
		p.logger.Errorf("failed to get thorchain rune providers: %v", err)
	}
	for {
		vaults, err := p.storage.GetVaultsWithPage(currentVaultId, 1000)
		if err != nil {
			p.logger.Errorf("failed to get vaults: %v", err)
			continue
		}
		if len(vaults) == 0 {
			p.logger.Info("no more vaults to process")
			break
		}
		for _, vault := range vaults {
			coins, err := p.storage.GetCoins(vault.ID)
			if err != nil {
				p.logger.Errorf("failed to get coins for vault: %v", err)
				continue
			}
			vaultAddress := models.NewVaultAddress(vault.ID)
			for _, coin := range coins {
				vaultAddress.SetAddress(coin.Chain, coin.Address)
			}
			if len(coins) > 0 {
				positionWorkerChan <- vaultAddress
			}
			currentVaultId = vault.ID
			job.CurrentVaultID = vault.ID
			if err := p.storage.UpdateJob(job); err != nil {
				p.logger.Errorf("failed to update job: %v", err)
			}
		}
	}
	for {
		coins, err := p.storage.GetCoinsWithPage(currentID, 1000)
		if err != nil {
			p.logger.Errorf("failed to get coins: %v", err)
			continue
		}

		for _, coin := range coins {
			currentID = uint64(coin.ID)
			job.CurrentID = int64(coin.ID)
			workChan <- coin
		}

		if len(coins) == 0 {
			p.logger.Info("no more coins to process, stopping task provider")
			// no more to process
			close(workChan)
			job.IsSuccess = true
		}

		if err := p.storage.UpdateJob(job); err != nil {
			p.logger.Errorf("failed to update job: %v", err)
		}

		if job.IsSuccess {
			if p.storage.UpdateVaultRanks() != nil {
				p.logger.Errorf("failed to update vault ranks: %v", err)
			}
			if p.storage.UpdateVaultBalance() != nil {
				p.logger.Errorf("failed to update vault balance: %v", err)
			}
			return
		}
	}
}
func (p *PointWorker) activePositionWorker(idx int, workerChan <-chan models.VaultAddress, job models.Job) {
	p.logger.Infof("active position worker %d started", idx)
	defer p.wg.Done()
	for {
		select {
		case <-p.stopChan:
			p.logger.Infof("active position worker %d stop signal received, stopping worker", idx)
			return
		case v, more := <-workerChan:
			if !more {
				return
			}
			if err := p.updatePosition(v, job.Multiplier); err != nil {
				p.logger.Errorf("failed to update position: %v", err)
			}
			if err := p.updateNFTBalance(v, job.Multiplier); err != nil {
				p.logger.Errorf("failed to update nft balance: %v", err)
			}
		}
	}
}
func (p *PointWorker) taskWorker(idx int, workerChan <-chan models.CoinDBModel, job models.Job) {
	p.logger.Infof("worker %d started", idx)
	defer p.wg.Done()
	for {
		select {
		case <-p.stopChan:
			p.logger.Infof("worker %d stop signal received, stopping worker", idx)
			return
		case t, more := <-workerChan:
			if !more {
				return
			}
			if err := p.updateBalance(t, job.Multiplier); err != nil {
				p.logger.Errorf("failed to update balance: %v", err)
			}
		}
	}
}

func (p *PointWorker) updatePosition(vaultAddress models.VaultAddress, multiplier int64) error {
	newlp, err := p.fetchPosition(vaultAddress)
	if err != nil {
		p.logger.Errorf("failed to fetch position for vault id %d , using old position: %v", vaultAddress.GetVaultID(), err)
		oldLp, err := p.storage.GetLPValue(vaultAddress.GetVaultID())
		if err != nil {
			return fmt.Errorf("failed to get vault: %w", err)
		}
		newlp = oldLp
	} else {
		p.logger.Infof("new lp value for vault %d is %d", vaultAddress.GetVaultID(), newlp)
		if err := p.storage.UpdateLPValue(vaultAddress.GetVaultID(), newlp); err != nil {
			p.logger.Errorf("failed to update lp value: %v", err)
		}
	}
	newPoints := int64(newlp * multiplier)
	if newlp == 0 {
		return nil
	}
	if err := p.storage.IncreaseVaultTotalPoints(vaultAddress.GetVaultID(), newPoints); err != nil {
		return fmt.Errorf("failed to increase vault total points: %w", err)
	}
	return nil
}

func (p *PointWorker) updateNFTBalance(vaultAddress models.VaultAddress, multiplier int64) error {
	var nftValue int64
	nftValue, err := p.fetchNFTValue(vaultAddress)
	if err != nil {
		p.logger.Errorf("failed to fetch nft value for vault id %d , using old nft value: %v", vaultAddress.GetVaultID(), err)
		nftValue, err = p.storage.GetNFTValue(vaultAddress.GetVaultID())
		if err != nil {
			return fmt.Errorf("failed to get vault: %w", err)
		}
	} else {
		p.logger.Infof("new nft value for vault %d is %d", vaultAddress.GetVaultID(), nftValue)
		if err := p.storage.UpdateNFTValue(vaultAddress.GetVaultID(), nftValue); err != nil {
			p.logger.Errorf("failed to update nft value: %v", err)
		}
	}
	newPoints := int64(nftValue * multiplier)
	if newPoints == 0 {
		return nil
	}
	if err := p.storage.IncreaseVaultTotalPoints(vaultAddress.GetVaultID(), newPoints); err != nil {
		return fmt.Errorf("failed to increase vault total points: %w", err)
	}
	return nil
}

func (p *PointWorker) fetchPosition(vaultAddress models.VaultAddress) (int64, error) {
	backoffRetry := utils.NewBackoffRetry(5)
	address := strings.Join(vaultAddress.GetAllAddress(), ",")
	p.logger.Infof("start to update position for vault: %d,  address: %s ", vaultAddress.GetVaultID(), address)
	tcmayalp, err := backoffRetry.RetryWithBackoff(p.lpResolver.GetLiquidityPosition, address)
	if err != nil {
		return 0, fmt.Errorf("failed to get tc/maya liquidity position for vault:%d : %w", vaultAddress.GetVaultID(), err)
	}
	p.logger.Infof("tc/maya liquidity position for vault %d is %f", vaultAddress.GetVaultID(), tcmayalp)
	tgtPrice, err := p.priceResolver.GetCoinGeckoPrice("thorwallet", "usd")
	if err != nil {
		return 0, fmt.Errorf("failed to get tgt price: %w", err)
	}
	p.lpResolver.SetTGTPrice(tgtPrice)
	tgtlp, err := backoffRetry.RetryWithBackoff(p.lpResolver.GetTGTStakePosition, vaultAddress.GetEVMAddress())
	if err != nil {
		return 0, fmt.Errorf("failed to get tgt stake position for vault:%d : %w", vaultAddress.GetVaultID(), err)
	}
	p.logger.Infof("tgt stake position for vault %d is %f", vaultAddress.GetVaultID(), tgtlp)
	wewelp, err := backoffRetry.RetryWithBackoff(p.lpResolver.GetWeWeLPPosition, vaultAddress.GetEVMAddress())
	if err != nil {
		return 0, fmt.Errorf("failed to get wewel liquidity position for vault:%d : %w", vaultAddress.GetVaultID(), err)
	}
	p.logger.Infof("wewel liquidity position for vault %d is %f", vaultAddress.GetVaultID(), wewelp)
	saver, err := backoffRetry.RetryWithBackoff(p.saverResolver.GetSaverPosition, address)
	if err != nil {
		return 0, fmt.Errorf("failed to get saver position for vault:%d : %w", vaultAddress.GetVaultID(), err)
	}
	p.logger.Infof("saver position for vault %d is %f", vaultAddress.GetVaultID(), saver)
	newLP := tcmayalp + tgtlp + wewelp + saver
	return int64(newLP), nil
}
func (p *PointWorker) fetchNFTValue(vault models.VaultAddress) (int64, error) {
	sum := float64(0)
	for _, nft := range p.whitelistNFTCollection {
		address := vault.GetAddress(nft.Chain)
		if address != "" {
			balance, err := p.balanceResolver.GetBalanceWithRetry(models.CoinDBModel{CoinBase: models.CoinBase{
				Chain:           nft.Chain,
				Address:         address,
				ContractAddress: nft.CollectionAddress,
				Decimals:        0,
				IsNative:        false,
			}})
			if err != nil {
				return 0, fmt.Errorf("failed to get balance for address:%s : %v", address, err)
			}
			price, err := p.priceResolver.GetOpenSeaCollectionMinPrice(nft.CollectionSlug)
			if err != nil {
				return 0, fmt.Errorf("failed to get price for collection:%s : %v", nft.CollectionSlug, err)
			}
			sum += balance * price
		}
	}
	return int64(sum), nil
}
func (p *PointWorker) updateBalance(coin models.CoinDBModel, multiplier int64) error {
	p.logger.Infof("start to update balance for chain: %s, ticker: %s, address: %s ", coin.Chain, coin.Ticker, coin.Address)
	coinBalance, err := p.balanceResolver.GetBalanceWithRetry(coin)
	if err != nil {
		p.logger.Errorf("failed to get balance for address:%s : %v", coin.Address, err)
		prevBalance, errP := strconv.ParseFloat(coin.Balance, 64)
		if errP != nil {
			return fmt.Errorf("failed to parse previous balance: %w", errP)
		}
		// server failed to get the latest balance , assume his previous balance is correct and use it to accumulate points
		coinBalance = prevBalance
	} else {
		if err := p.storage.UpdateCoinBalance(uint64(coin.ID), coinBalance); err != nil {
			return fmt.Errorf("failed to update coin balance: %w", err)
		}
	}
	if coin.PriceUSD == "" {
		coin.PriceUSD = "0"
	}
	// increase vault's point
	price, err := strconv.ParseFloat(coin.PriceUSD, 64)
	if err != nil {
		return fmt.Errorf("failed to parse coin price: %w", err)
	}
	newPoints := int64(coinBalance * price * float64(multiplier))
	if newPoints == 0 {
		return nil
	}
	if err := p.storage.IncreaseVaultTotalPoints(coin.VaultID, newPoints); err != nil {
		return fmt.Errorf("failed to increase vault total points: %w", err)
	}
	return nil
}

func (p *PointWorker) updateCoinPrice() error {
	p.logger.Info("start to update coin prices")
	coinIdentities, err := p.storage.GetUniqueCoins()
	if err != nil {
		return fmt.Errorf("failed to get unique coins: %w", err)
	}
	p.logger.Infof("got %d unique coins", len(coinIdentities))
	if len(coinIdentities) == 0 {
		return nil
	}
	coinPrices, err := p.priceResolver.GetAllTokenPrices(coinIdentities)
	if err != nil {
		return fmt.Errorf("failed to get all token prices: %w", err)
	}
	p.logger.Infof("%+v", coinPrices)
	for id, coinIden := range coinPrices {
		if err := p.storage.UpdateCoinPriceByCMCID(id, coinIden); err != nil {
			p.logger.Errorf("failed to update coin price: %d, err: %v", id, err)
			// log the error and move on
			continue
		}
	}
	cacaoPrice, err := p.priceResolver.GetMidgardCacaoPrices()
	if err != nil {
		p.logger.Errorf("failed to get CACAO price: %v", err)
	} else {
		if err := p.storage.UpdateCoinPrice(common.MayaChain, "CACAO", cacaoPrice); err != nil {
			p.logger.Errorf("failed to update CACAO price: %v", err)
		}
	}

	vthorPrice, err := p.priceResolver.GetLiFiPrice("eth", "0x815C23eCA83261b6Ec689b60Cc4a58b54BC24D8D")
	if err != nil {
		p.logger.Errorf("failed to get VTHOR price: %v", err)
	} else {
		if err := p.storage.UpdateCoinPrice(common.Ethereum, "vTHOR", vthorPrice); err != nil {
			p.logger.Errorf("failed to update VTHOR price: %v", err)
		}
	}
	mayaPrice := float64(40)
	if err := p.storage.UpdateCoinPrice(common.MayaChain, "MAYA", mayaPrice); err != nil {
		p.logger.Errorf("failed to update VTHOR price: %v", err)
	}

	defer p.logger.Info("finish updating coin prices")
	return nil
}
