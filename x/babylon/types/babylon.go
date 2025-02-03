package types

import "encoding/json"

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

// NewInitMsg creates the init message for the Babylon contract
func NewInitMsg(
	network string,
	babylonTag string,
	btcConfirmationDepth uint32,
	checkpointFinalizationTimeout uint32,
	notifyCosmosZone bool,
	btcStakingCodeId uint64,
	btcStakingInitMsgBytes []byte,
	btcFinalityCodeId uint64,
	btcFinalityInitMsgBytes []byte,
	consumerName string,
	consumerDescription string,
	admin string,
) ([]byte, error) {
	initMsg := map[string]interface{}{
		"network":                         network,
		"babylon_tag":                     babylonTag,
		"btc_confirmation_depth":          btcConfirmationDepth,
		"checkpoint_finalization_timeout": checkpointFinalizationTimeout,
		"notify_cosmos_zone":              notifyCosmosZone,
		"btc_staking_code_id":             btcStakingCodeId,
		"btc_staking_msg":                 btcStakingInitMsgBytes,
		"consumer_name":                   consumerName,
		"consumer_description":            consumerDescription,
		"btc_finality_code_id":            btcFinalityCodeId,
		"btc_finality_msg":                btcFinalityInitMsgBytes,
	}
	if len(admin) > 0 {
		initMsg["admin"] = admin
	}
	initMsgBytes, err := json.Marshal(initMsg)
	if err != nil {
		return nil, err
	}
	return initMsgBytes, nil
}
