package contract

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
)

// CustomMsg is a message sent from a smart contract to the Babylon module
type (
	CustomMsg struct {
		MintRewards *MintRewardsMsg `json:"mint_rewards,omitempty"`
	}
	// MintRewardsMsg mints the specified number of block rewards,
	// and sends them to the specified recipient (typically, the staking contract)
	MintRewardsMsg struct {
		Amount    wasmvmtypes.Coin `json:"amount"`
		Recipient string           `json:"recipient"`
	}
)
