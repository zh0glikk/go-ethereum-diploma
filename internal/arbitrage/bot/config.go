package bot

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/internal/arbitrage/bot/configger"
	"github.com/ethereum/go-ethereum/internal/arbitrage/bot/runners/blocklistener"
	"github.com/ethereum/go-ethereum/internal/arbitrage/bot/runners/mempool"
	"github.com/ethereum/go-ethereum/internal/arbitrage/bot/runners/simulator"
	"github.com/ethereum/go-ethereum/internal/arbitrage/bot/runners/tracer"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"math/big"
	"strconv"
	"sync"
)

type ServiceConfig struct {
	b     ethapi.Backend
	chain *core.BlockChain

	ctx context.Context
	wg  *sync.WaitGroup

	memPools      []mempool.MemPool
	tracer        *tracer.Tracer
	simulator     *simulator.Simulator
	blockListener *blocklistener.BlockListener

	maxDepth          int
	maxBruteTime      int64
	splitParam        *big.Int
	initialSplitParam *big.Int

	contractAddress common.Address

	infuraEthClient *ethclient.Client
}

func NewServiceConfig(eth *eth.Ethereum, ctx context.Context, wg *sync.WaitGroup) (*ServiceConfig, error) {
	cfg := &ServiceConfig{
		b:     eth.APIBackend,
		chain: eth.BlockChain(),
		ctx:   ctx,
		wg:    wg,
	}
	var err error

	if configger.Get("ARB_INFURA_RPC") == "" {
		return nil, errors.New("ARB_INFURA_RPC is required")
	} else {
		cfg.infuraEthClient, err = ethclient.Dial(configger.Get("ARB_INFURA_RPC"))
		if err != nil {
			return nil, errors.New("failed to get ARB_INFURA_RPC")
		}
	}

	var memPools []mempool.MemPool
	if configger.Get("ARB_INTERNAL_MEMPOOL_ENABLE") == "true" {
		memPools = append(memPools, mempool.NewInternalMemPool(cfg.b))
	}
	if configger.Get("ARB_INFURA_MEMPOOL_URL") != "" {
		memPools = append(memPools, mempool.NewInfuraMemPool(configger.Get("ARB_INFURA_MEMPOOL_URL")))
	}
	cfg.memPools = memPools
	cfg.tracer = tracer.NewTracer(cfg.b)
	cfg.blockListener = blocklistener.NewBlockListener(cfg.chain)

	// simulation params
	if configger.Get("ARB_CONTRACT_ADDRESS") == "" {
		return nil, errors.New("sando contract address not provided")
	}
	cfg.contractAddress = common.HexToAddress(configger.Get("ARB_CONTRACT_ADDRESS"))
	cfg.maxDepth, err = strconv.Atoi(configger.Get("ARB_SIMULATION_MAX_DEPTH"))
	if err != nil {
		return nil, errors.New("ARB_SIMULATION_MAX_DEPTH provided")
	}
	cfg.maxBruteTime, err = strconv.ParseInt(configger.Get("ARB_SIMULATION_MAX_BRUTE_TIME"), 10, 64)
	if err != nil {
		return nil, errors.New("ARB_SIMULATION_MAX_BRUTE_TIME not provided")
	}

	splitParam, err := strconv.ParseInt(configger.Get("ARB_SIMULATION_SPLIT_PARAM"), 10, 64)
	if err != nil {
		return nil, errors.New("ARB_SIMULATION_SPLIT_PARAM not provided")
	}
	cfg.splitParam = big.NewInt(splitParam)

	initialSplitParam, err := strconv.ParseInt(configger.Get("ARB_SIMULATION_INITIAL_SPLIT_PARAM"), 10, 64)
	if err != nil {
		return nil, errors.New("ARB_SIMULATION_INITIAL_SPLIT_PARAM not provided")
	}
	cfg.initialSplitParam = big.NewInt(initialSplitParam)

	cfg.simulator = simulator.NewSimulator(
		cfg.b,
		cfg.contractAddress,
		cfg.maxDepth,
		cfg.maxBruteTime,
		cfg.splitParam,
		cfg.initialSplitParam)

	return cfg, nil
}
