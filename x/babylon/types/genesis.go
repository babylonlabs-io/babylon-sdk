package types

// ValidateGenesis does basic validation on genesis state
func ValidateGenesis(gs *GenesisState) error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}
	if gs.BsnContracts != nil && gs.BsnContracts.IsSet() {
		if err := gs.BsnContracts.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

// NewGenesisState constructor
func NewGenesisState(params Params, bsnContracts *BSNContracts) *GenesisState {
	return &GenesisState{
		Params:       params,
		BsnContracts: bsnContracts,
	}
}

// DefaultGenesisState default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), &BSNContracts{})
}
