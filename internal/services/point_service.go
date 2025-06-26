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
	"github.com/vultisig/airdrop-registry/internal/volume"
)

const MinBalanceForValidReferral = 50 // 50 USDT
// PointWorker is a worker that processes points
type PointWorker struct {
	logger                 *logrus.Logger
	storage                *Storage
	priceResolver          *PriceResolver
	balanceResolver        *balance.BalanceResolver
	lpResolver             *liquidity.LiquidityPositionResolver
	saverResolver          *liquidity.SaverPositionResolver
	referralResolver       *ReferralResolverService
	volumeResolver         *volume.VolumeResolver
	startCoinID            int64
	wg                     *sync.WaitGroup
	stopChan               chan struct{}
	cfg                    *config.Config
	isJobInProgress        bool
	isVolumeFetched        bool // flag to indicate if volume fetched successfully
	whitelistNFTCollection []models.NFTCollection
}

func NewPointWorker(cfg *config.Config, storage *Storage, priceResolver *PriceResolver, balanceResolver *balance.BalanceResolver, volumeResolver *volume.VolumeResolver, referralResolver *ReferralResolverService) (*PointWorker, error) {

	if nil == storage {
		return nil, fmt.Errorf("storage is nil")
	}
	if nil == priceResolver {
		return nil, fmt.Errorf("priceResolver is nil")
	}

	return &PointWorker{
		logger:           logrus.WithField("module", "point_worker").Logger,
		storage:          storage,
		priceResolver:    priceResolver,
		balanceResolver:  balanceResolver,
		lpResolver:       liquidity.NewLiquidtyPositionResolver(),
		referralResolver: referralResolver,
		saverResolver:    liquidity.NewSaverPositionResolver(),
		volumeResolver:   volumeResolver,
		startCoinID:      cfg.Worker.StartID,
		stopChan:         make(chan struct{}),
		wg:               &sync.WaitGroup{},
		cfg:              cfg,
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
	//default value for lastVolumeFetch is first of June 2025
	lastVolumeFetch := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC).Unix()
	lastVolumeJob, err := p.storage.GetLastVolumeFetch()
	if err == nil {
		lastVolumeFetch = models.GetDate(lastVolumeJob.JobDate)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		p.logger.Errorf("failed to get last volume fetch: %e", err)
		return
	}

	// TODO: make sure logic for from/to is correct
	if err := p.volumeResolver.LoadVolume(lastVolumeFetch, models.GetDate(job.JobDate)); err != nil {
		p.logger.Errorf("failed to load volume: %e", err)
	} else {
		p.logger.Infof("volume fetch completed successfully (from %d to %d)", lastVolumeFetch, models.GetDate(job.JobDate))
		p.isVolumeFetched = true
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
		for i, vault := range vaults {
			if vault.CurrentSeasonID < p.cfg.GetCurrentSeason().ID {
				p.logger.Infof("vault %d is not in current season, commiting old season points", vault.ID)
				if err := p.storage.CommitSeasonPoints(vault, p.cfg.GetCurrentSeason().ID); err != nil {
					p.logger.Errorf("failed to commit season points for vault %d: %v", vault.ID, err)
					continue
				}
			}
			// Fetch referral count
			vaults[i].ReferralCount, err = p.getValidReferralCount(vault.ECDSA, vault.EDDSA)
			if err != nil {
				p.logger.Errorf("failed to get referral count for vault %d: %v", vault.ID, err)
				continue
			}

			err = p.storage.UpdateReferralCount(&vaults[i])
			if err != nil {
				p.logger.Errorf("failed to update referral count for vault %d: %v", vault.ID, err)
				continue
			}

			coins, err := p.storage.GetCoins(vault.ID)
			if err != nil {
				p.logger.Errorf("failed to get coins for vault: %v", err)
				continue
			}
			var totalVolume float64
			address := make(map[string]interface{})
			//generate vault address for all chains
			for _, chain := range common.GetAllChains() {
				//generate address for the given chains
				addr, err := vault.GetAddress(chain)
				if err != nil {
					p.logger.Errorf("failed to get address for vault %d on chain %s: %v", vault.ID, chain, err)
					continue
				}
				found := false
				for _, coin := range coins {
					if coin.Address == addr {
						found = true
					}
				}
				if !found {
					// if address not found in coins, add it
					coins = append(coins, models.CoinDBModel{
						CoinBase: models.CoinBase{
							Chain:    chain,
							Address:  addr,
							IsNative: true,
						},
						VaultID: vault.ID,
					})
				}
			}

			// fetch volume for each coin
			for _, coin := range coins {
				if _, ok := address[coin.Address]; ok {
					continue // skip if address already exists
				}
				coinVolume := p.volumeResolver.GetVolume(coin.Address)
				if coinVolume > 0 {
					totalVolume += coinVolume
				}
				address[coin.Address] = nil
			}
			err = p.storage.UpdateVolume(vault.ID, totalVolume)
			if err != nil {
				p.logger.Errorf("failed to update volume for vault %d: %v", vault.ID, err)
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
			break
		}
	}
	if job.IsSuccess {
		if err := p.storage.UpdateVaultBalance(); err != nil {
			p.logger.Errorf("failed to update vault balance: %v", err)
		}
		if p.cfg.GetCurrentSeason().ID > 0 {
			p.logger.Infof("update vaults total point based on new formula for season %d", p.cfg.GetCurrentSeason().ID)
			if err := p.storage.UpdateVaultTotalPoints(); err != nil {
				p.logger.Errorf("failed to update vault total points: %v", err)
			}
			if err := p.updateVaultsMilestone(); err != nil {
				p.logger.Errorf("failed to update vaults milestones: %v", err)
			}
		}
		if err := p.storage.UpdateVaultRanks(); err != nil {
			p.logger.Errorf("failed to update vault ranks: %v", err)
		}
	}
	if p.isVolumeFetched {
		err := p.storage.UpdateIsVolumeFetched(job)
		if err != nil {
			//TODO: handler error properly
			p.logger.Errorf("failed to update is_volume_fetched: %v", err)
		} else {
			p.logger.Infof("volume fetched successfully, updated job %d", job.ID)
		}
	}
}
func (p *PointWorker) updateVaultsMilestone() error {
	startId := uint(0)
	for {
		vaults, err := p.storage.GetVaultsWithPage(startId, 1000)
		if err != nil {
			p.logger.Errorf("failed to get vaults: %v", err)
			return fmt.Errorf("failed to get vaults: %w", err)
		}
		if len(vaults) == 0 {
			break
		}
		for _, vault := range vaults {
			for i := 0; i < len(p.cfg.GetCurrentSeason().Milestones); i++ {
				// if total points is greater than or equal to milestone minimum
				if vault.TotalPoints >= float64(p.cfg.GetCurrentSeason().Milestones[i].Minimum) {
					// if this milestone is locked
					if vault.NextMilestoneID <= i {
						// unlock milestone: update vault total points and next milestone id
						p.storage.UpdateVaultMilestone(vault.ID, i+1, float64(p.cfg.GetCurrentSeason().Milestones[i].Prize))
					}
				}
			}
			startId = vault.ID
		}
	}
	p.logger.Info("all vaults processed for milestones")
	return nil
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
	newPoints := float64(newlp * multiplier)
	if newlp == 0 {
		return nil
	}
	if err := p.storage.IncreaseVaultTotalValue(vaultAddress.GetVaultID(), newPoints); err != nil {
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
	newPoints := float64(nftValue * multiplier)
	if newPoints == 0 {
		return nil
	}
	if err := p.storage.IncreaseVaultTotalValue(vaultAddress.GetVaultID(), newPoints); err != nil {
		return fmt.Errorf("failed to increase vault total points: %w", err)
	}
	return nil
}

func (p *PointWorker) fetchPosition(vaultAddress models.VaultAddress) (int64, error) {
	backoffRetry := utils.NewBackoffRetry(5)
	address := strings.Join(vaultAddress.GetAllAddress(), ",")
	p.logger.Infof("start to update position for vault: %d,  address: %s ", vaultAddress.GetVaultID(), address)

	tcyPrice, err := p.priceResolver.GetMidgardPrices("THOR.TCY")
	if err != nil {
		return 0, fmt.Errorf("failed to get tcy price: %w", err)
	}
	p.lpResolver.SetTCYPrice(tcyPrice)

	tcmayalp, err := backoffRetry.RetryWithBackoff(p.lpResolver.GetLiquidityPosition, address)
	if err != nil {
		return 0, fmt.Errorf("failed to get tc/maya liquidity position for vault:%d : %w", vaultAddress.GetVaultID(), err)
	}
	p.logger.Infof("tc/maya liquidity position for vault %d is %f", vaultAddress.GetVaultID(), tcmayalp)

	saver, err := backoffRetry.RetryWithBackoff(p.saverResolver.GetSaverPosition, address)
	if err != nil {
		return 0, fmt.Errorf("failed to get saver position for vault:%d : %w", vaultAddress.GetVaultID(), err)
	}
	p.logger.Infof("saver position for vault %d is %f", vaultAddress.GetVaultID(), saver)

	tcyStake, err := backoffRetry.RetryWithBackoff(p.lpResolver.GetTCYStakePosition, vaultAddress.GetAddress(common.THORChain))
	if err != nil {
		return 0, fmt.Errorf("failed to get tcy stake position for vault:%d : %w", vaultAddress.GetVaultID(), err)
	}
	p.logger.Infof("tcy stake position for vault %d is %f", vaultAddress.GetVaultID(), tcyStake)

	newLP := tcmayalp + saver + tcyStake
	return int64(newLP), nil
}
func (p *PointWorker) fetchNFTValue(vault models.VaultAddress) (int64, error) {
	sum := float64(0)
	for _, nft := range p.whitelistNFTCollection {
		address := vault.GetAddress(nft.Chain)
		if address != "" {
			token := models.CoinDBModel{CoinBase: models.CoinBase{
				Chain:           nft.Chain,
				Address:         address,
				ContractAddress: nft.CollectionAddress,
				Decimals:        0,
				IsNative:        false,
			}}
			balance, err := p.balanceResolver.GetBalanceWithRetry(token)
			if err != nil {
				return 0, fmt.Errorf("failed to get balance for address:%s : %v", address, err)
			}
			price, err := p.priceResolver.GetOpenSeaCollectionMinPrice(nft.CollectionSlug)
			if err != nil {
				return 0, fmt.Errorf("failed to get price for collection:%s : %v", nft.CollectionSlug, err)
			}
			seasonMultiplier := p.getSeasonMultiplierForNFT(token)
			sum += balance * float64(seasonMultiplier) * price
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
	seasonMultiplier := p.getSeasonMultiplierForCoin(coin)
	newPoints := float64(coinBalance * price * float64(multiplier) * float64(seasonMultiplier))
	if newPoints == 0 {
		return nil
	}
	if err := p.storage.IncreaseVaultTotalValue(coin.VaultID, newPoints); err != nil {
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
	kweenPrice, err := p.priceResolver.GetCoinGeckoPrice("kween", "usd")
	if err != nil {
		p.logger.Errorf("failed to get KWEEN price: %v", err)
	} else {
		if err := p.storage.UpdateCoinPrice(common.Solana, "KWEEN", kweenPrice); err != nil {
			p.logger.Errorf("failed to update KWEEN price: %v", err)
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
	tcyPrice, err := p.priceResolver.GetMidgardPrices("THOR.TCY")
	if err != nil {
		p.logger.Errorf("failed to get TCY price: %v", err)
	} else {
		if err := p.storage.UpdateCoinPrice(common.THORChain, "THOR.TCY", tcyPrice); err != nil {
			p.logger.Errorf("failed to update TCY price: %v", err)
		}
	}

	defer p.logger.Info("finish updating coin prices")
	return nil
}

func (p *PointWorker) getValidReferralCount(ecdsaKey string, eddsaKey string) (int64, error) {
	referrals, err := p.referralResolver.GetReferrals(ecdsaKey, eddsaKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get referrals: %w", err)
	}

	var cnt int64
	for _, r := range referrals {
		// User can not refer himself
		if r.WalletPublicKeyEcdsa == ecdsaKey && r.WalletPublicKeyEddsa == eddsaKey {
			continue
		}
		if r.WalletPublicKeyEcdsa == "" || r.WalletPublicKeyEddsa == "" {
			continue
		}
		v, err := p.storage.GetVault(r.WalletPublicKeyEcdsa, r.WalletPublicKeyEddsa)
		if err != nil || v == nil {
			p.logger.Warnf("Referral vault not found for ECDSA: %s, EDDSA: %s",
				r.WalletPublicKeyEcdsa, r.WalletPublicKeyEddsa)
			continue
		}
		if v.Balance+v.LPValue+v.NFTValue >= MinBalanceForValidReferral {
			cnt++
		}
	}

	return cnt, nil
}

func (p *PointWorker) getSeasonMultiplierForCoin(coin models.CoinDBModel) float64 {
	for _, token := range p.cfg.GetCurrentSeason().Tokens {
		if token.Chain == coin.Chain.String() && token.Name == coin.Ticker && coin.ContractAddress == token.ContractAddress {
			return token.Multiplier
		}
	}
	return 1
}

func (p *PointWorker) getSeasonMultiplierForNFT(coin models.CoinDBModel) float64 {
	for _, collection := range p.cfg.GetCurrentSeason().NFTs {
		if collection.Chain == coin.Chain.String() && collection.ContractAddress == coin.ContractAddress {
			return collection.Multiplier
		}
	}
	return 1
}
