package services

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/balance"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

// PointWorker is a worker that processes points
type PointWorker struct {
	logger          *logrus.Logger
	storage         *Storage
	priceResolver   *PriceResolver
	balanceResolver *balance.BalanceResolver
	startCoinID     int64
	wg              *sync.WaitGroup
	stopChan        chan struct{}
	cfg             *config.Config
	isJobInProgress bool
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
		startCoinID:     cfg.Worker.StartID,
		stopChan:        make(chan struct{}),
		wg:              &sync.WaitGroup{},
		cfg:             cfg,
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

	if err := p.updateCoinPrice(); err != nil {
		p.logger.Errorf("failed to update coin prices: %e", err)
		return
	}

	p.wg.Add(1)
	workChan := make(chan models.CoinDBModel)
	go p.taskProvider(job, workChan)
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

func (p *PointWorker) taskProvider(job *models.Job, workChan chan models.CoinDBModel) {
	defer p.wg.Done()
	p.isJobInProgress = true
	defer func() {
		p.isJobInProgress = false
	}()
	currentID := uint64(job.CurrentID)
	// refresh bond providers
	if err := p.balanceResolver.GetTHORChainBondProviders(); err != nil {
		p.logger.Errorf("failed to get thorchain bond providers: %v", err)
	}
	if err := p.balanceResolver.GetTHORChainRuneProviders(); err != nil {
		p.logger.Errorf("failed to get thorchain rune providers: %v", err)
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

func (p *PointWorker) updateBalance(coin models.CoinDBModel, multiplier int64) error {
	p.logger.Infof("start to update balance for chain: %s, ticker: %s, address: %s ", coin.Chain, coin.Ticker, coin.Address)
	coinBalance, err := p.balanceResolver.GetBalanceWithRetry(coin)
	if err != nil {
		return fmt.Errorf("failed to get balance for address:%s : %w", coin.Address, err)
	}
	if err := p.storage.UpdateCoinBalance(uint64(coin.ID), coinBalance); err != nil {
		return fmt.Errorf("failed to update coin balance: %w", err)
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
	cacaoPrice, err := p.priceResolver.GetCoinGeckoPrice("cacao", "usd")
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
