package types

// Config represents the configuration for the Babylon contract
type BabylonContractConfig struct {
	Network                       string `json:"network"`
	BabylonTag                    []byte `json:"babylon_tag"`
	BTCConfirmationDepth          uint32 `json:"btc_confirmation_depth"`
	CheckpointFinalizationTimeout uint32 `json:"checkpoint_finalization_timeout"`
	// NotifyCosmosZone indicates whether to send Cosmos zone messages notifying BTC-finalised headers.
	// NOTE: if set to true, then the Cosmos zone needs to integrate the corresponding message
	// handler as well
	NotifyCosmosZone bool `json:"notify_cosmos_zone"`
	// BTCStaking stores a BTC staking contract used for BTC re-staking if set
	BTCStaking string `json:"btc_staking,omitempty"`
	// BTCFinality stores a BTC finality contract used for BTC finality on the Consumer if set
	BTCFinality string `json:"btc_finality,omitempty"`
	// ConsumerName represents the name of the Consumer
	ConsumerName string `json:"consumer_name,omitempty"`
	// ConsumerDescription represents the description of the Consumer
	ConsumerDescription string `json:"consumer_description,omitempty"`
}

type BabylonContractInitMsg struct {
	// Network represents the Bitcoin network (mainnet, testnet, etc.)
	Network string `json:"network"`
	// BabylonTag is a string encoding four bytes used for identification / tagging of the Babylon zone.
	// NOTE: this is a hex string, not raw bytes
	BabylonTag string `json:"babylon_tag"`
	// BTCConfirmationDepth is the number of confirmations required for BTC headers
	BTCConfirmationDepth uint32 `json:"btc_confirmation_depth"`
	// CheckpointFinalizationTimeout is the timeout period for checkpoint finalization
	CheckpointFinalizationTimeout uint32 `json:"checkpoint_finalization_timeout"`
	// NotifyCosmosZone indicates whether to send Cosmos zone messages notifying BTC-finalised headers.
	// NOTE: If set to true, then the Cosmos zone needs to integrate the corresponding message handler
	// as well
	NotifyCosmosZone bool `json:"notify_cosmos_zone"`
	// BTCStakingCodeID is the code ID for the BTC staking contract, if set
	BTCStakingCodeID *uint64 `json:"btc_staking_code_id,omitempty"`
	// BTCStakingMsg is the instantiation message for the BTC staking contract.
	// This message is opaque to the Babylon contract, and depends on the specific staking contract
	// being instantiated
	BTCStakingMsg []byte `json:"btc_staking_msg,omitempty"`
	// BTCFinalityCodeID is the code ID for the BTC finality contract, if set
	BTCFinalityCodeID *uint64 `json:"btc_finality_code_id,omitempty"`
	// BTCFinalityMsg is the instantiation message for the BTC finality contract.
	// This message is opaque to the Babylon contract, and depends on the specific finality contract
	// being instantiated
	BTCFinalityMsg []byte `json:"btc_finality_msg,omitempty"`
	// Admin is the Wasm migration / upgrade admin of the BTC staking contract and the BTC finality contract
	Admin string `json:"admin,omitempty"`
	// ConsumerName represents the name of the Consumer
	ConsumerName string `json:"consumer_name,omitempty"`
	// ConsumerDescription represents the description of the Consumer
	ConsumerDescription string `json:"consumer_description,omitempty"`
}
