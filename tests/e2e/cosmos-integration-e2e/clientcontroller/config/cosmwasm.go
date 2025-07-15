package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmwasm-client/config"
)

type CosmwasmConfig struct {
	Key            string        `long:"key" description:"name of the key to sign transactions with"`
	ChainID        string        `long:"chain-id" description:"chain id of the chain to connect to"`
	RPCAddr        string        `long:"rpc-address" description:"address of the rpc server to connect to"`
	GRPCAddr       string        `long:"grpc-address" description:"address of the grpc server to connect to"`
	AccountPrefix  string        `long:"acc-prefix" description:"account prefix to use for addresses"`
	KeyringBackend string        `long:"keyring-type" description:"type of keyring to use"`
	GasAdjustment  float64       `long:"gas-adjustment" description:"adjustment factor when using gas estimation"`
	GasPrices      string        `long:"gas-prices" description:"comma separated minimum gas prices to accept for transactions"`
	MinGasAmount   uint64        `long:"min-gas-amount" description:"minimum gas amount for transactions"`
	MaxGasAmount   uint64        `long:"max-gas-amount" description:"maximum gas amount for transactions"`
	KeyDirectory   string        `long:"key-dir" description:"directory to store keys in"`
	Debug          bool          `long:"debug" description:"flag to print debug output"`
	Timeout        time.Duration `long:"timeout" description:"client timeout when doing queries"`
	BlockTimeout   time.Duration `long:"block-timeout" description:"block timeout when waiting for block events"`
	OutputFormat   string        `long:"output-format" description:"default output when printint responses"`
	SignModeStr    string        `long:"sign-mode" description:"sign mode to use"`
}

func (cfg *CosmwasmConfig) Validate() error {
	if _, err := url.Parse(cfg.RPCAddr); err != nil {
		return fmt.Errorf("rpc-address is not correctly formatted: %w", err)
	}

	if _, err := url.Parse(cfg.GRPCAddr); err != nil {
		return fmt.Errorf("grpc-address is not correctly formatted: %w", err)
	}

	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if cfg.BlockTimeout < 0 {
		return fmt.Errorf("block-timeout can't be negative")
	}

	return nil
}

func DefaultCosmwasmConfig() *CosmwasmConfig {
	return &CosmwasmConfig{
		Key:            "validator",
		ChainID:        "wasmd-test",
		RPCAddr:        "http://localhost:26677",
		GRPCAddr:       "http://localhost:9092",
		AccountPrefix:  "wasm",
		KeyringBackend: "test",
		GasAdjustment:  1.3,
		GasPrices:      "1ustake",
		Debug:          true,
		Timeout:        20 * time.Second,
		BlockTimeout:   1 * time.Minute,
		OutputFormat:   "direct",
		SignModeStr:    "",
	}
}

func (cfg *CosmwasmConfig) ToQueryClientConfig() *config.CosmwasmConfig {
	return &config.CosmwasmConfig{
		Key:              cfg.Key,
		ChainID:          cfg.ChainID,
		RPCAddr:          cfg.RPCAddr,
		GRPCAddr:         cfg.GRPCAddr,
		AccountPrefix:    cfg.AccountPrefix,
		KeyringBackend:   cfg.KeyringBackend,
		GasAdjustment:    cfg.GasAdjustment,
		GasPrices:        cfg.GasPrices,
		MinGasAmount:     cfg.MinGasAmount,
		MaxGasAmount:     cfg.MaxGasAmount,
		KeyDirectory:     cfg.KeyDirectory,
		Debug:            cfg.Debug,
		Timeout:          cfg.Timeout,
		BlockTimeout:     cfg.BlockTimeout,
		OutputFormat:     cfg.OutputFormat,
		SignModeStr:      cfg.SignModeStr,
		SubmitterAddress: "",
	}
}
