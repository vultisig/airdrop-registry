package services

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/balance"
	"github.com/vultisig/airdrop-registry/internal/models"
)

// PointWorker is a worker that processes points
type PointWorker struct {
	logger          *logrus.Logger
	storage         *Storage
	priceResolver   *PriceResolver
	balanceResolver *balance.BalanceResolver
	startCoinID     int64
	workerChan      chan models.CoinDBModel
	wg              *sync.WaitGroup
	stopChan        chan struct{}
	cfg             *config.Config
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
		workerChan:      make(chan models.CoinDBModel, cfg.Worker.Concurrency),
		stopChan:        make(chan struct{}),
		wg:              &sync.WaitGroup{},
		cfg:             cfg,
	}, nil
}

func (p *PointWorker) Run() error {
	if err := p.updateCoinPrice(); err != nil {
		return fmt.Errorf("fail to update coin prices,err: %w", err)
	}
	p.wg.Add(1)
	go p.taskProvider()
	for i := 0; i < int(p.cfg.Worker.Concurrency); i++ {
		p.wg.Add(1)
		idx := i
		go p.taskWorker(idx)
	}
	return nil
}

func (p *PointWorker) Stop() {
	close(p.stopChan)
	p.wg.Wait()
}

func (p *PointWorker) taskProvider() {
	defer p.wg.Done()
	currentID := uint64(p.startCoinID)
	for {
		coins, err := p.storage.GetCoinsWithPage(currentID, 1000)
		if err != nil {
			p.logger.Errorf("failed to get coins: %v", err)
			continue
		}
		if len(coins) == 0 {
			p.logger.Info("no more coins to process, stopping task provider")
			close(p.workerChan)
			return
		}

		for _, coin := range coins {
			currentID = uint64(coin.ID)
			p.workerChan <- coin
		}
	}
}
func (p *PointWorker) taskWorker(idx int) {
	p.logger.Infof("worker %d started", idx)
	defer p.wg.Done()
	for {
		select {
		case <-p.stopChan:
			p.logger.Infof("worker %d stop signal received, stopping worker", idx)
			return
		case t, more := <-p.workerChan:
			if !more {
				return
			}
			if err := p.updateBalance(t); err != nil {
				p.logger.Errorf("failed to update balance: %v", err)
			}
		}
	}
}

func (p *PointWorker) updateBalance(coin models.CoinDBModel) error {
	coinBalance, err := p.balanceResolver.GetBalanceWithRetry(coin)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}
	if err := p.storage.UpdateCoinBalance(uint64(coin.ID), coinBalance); err != nil {
		return fmt.Errorf("failed to update coin balance: %w", err)
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
	for _, coinIden := range coinIdentities {
		key := fmt.Sprintf("%s-%s", coinIden.Chain, coinIden.Ticker)
		price, ok := coinPrices[key]
		if !ok {
			continue
		}
		if err := p.storage.UpdateCoinPrice(coinIden.Chain, coinIden.Ticker, price); err != nil {
			p.logger.Errorf("failed to update coin price: %s-%s, err: %v", coinIden.Chain, coinIden.Ticker, err)
			// log the error and move on
			continue
		}
	}
	defer p.logger.Info("finish updating coin prices")
	return nil
}
