package client

import (
	"sync"
	"time"

	wasmdparams "github.com/CosmWasm/wasmd/app/params"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"go.uber.org/zap"

	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmwasm-client/config"
	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmwasm-client/query"
	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmwasm-client/wasmclient"
)

type Client struct {
	mu sync.Mutex
	*query.QueryClient

	provider *wasmclient.CosmosProvider
	timeout  time.Duration
	logger   *zap.Logger
	cfg      *config.CosmwasmConfig
}

func New(cfg *config.CosmwasmConfig, chainName string, encodingCfg wasmdparams.EncodingConfig, logger *zap.Logger) (*Client, error) {
	var (
		zapLogger *zap.Logger
		err       error
	)

	// ensure cfg is valid
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// use the existing logger or create a new one if not given
	zapLogger = logger
	if zapLogger == nil {
		zapLogger, err = newRootLogger("console", true)
		if err != nil {
			return nil, err
		}
	}

	provider, err := cfg.ToCosmosProviderConfig().NewProvider(
		"", // TODO: set home path
		chainName,
	)
	if err != nil {
		return nil, err
	}

	cp := provider.(*wasmclient.CosmosProvider)
	cp.PCfg.KeyDirectory = cfg.KeyDirectory
	cp.Cdc = &encodingCfg

	// initialise Cosmos provider
	// NOTE: this will create a RPC client. The RPC client will be used for
	// submitting txs and making ad hoc queries. It won't create WebSocket
	// connection with wasmd node
	err = cp.Init()
	if err != nil {
		return nil, err
	}

	// create a queryClient so that the Client inherits all query functions
	// TODO: merge this RPC client with the one in `cp` after Cosmos side
	// finishes the migration to new RPC client
	// see https://github.com/strangelove-ventures/cometbft-client
	c, err := rpchttp.NewWithTimeout(cp.PCfg.RPCAddr, "/websocket", uint(cfg.Timeout.Seconds()))
	if err != nil {
		return nil, err
	}
	queryClient, err := query.NewWithClient(c, cfg.Timeout)
	if err != nil {
		return nil, err
	}

	return &Client{
		QueryClient: queryClient,
		provider:    cp,
		timeout:     cfg.Timeout,
		logger:      zapLogger,
		cfg:         cfg,
	}, nil
}

func (c *Client) GetConfig() *config.CosmwasmConfig {
	return c.cfg
}
