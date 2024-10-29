package contract

// CustomMsg is a message sent from a smart contract to the Babylon module
type (
	CustomMsg struct {
		MintRewards *MintRewardsMsg `json:"mint_rewards,omitempty"`
	}
	// MintRewardsMsg mints block rewards to the staking contract
	MintRewardsMsg struct{}
)
