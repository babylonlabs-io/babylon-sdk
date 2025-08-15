package cosmwasm

import (
	"fmt"
)

// MustQueryBabylonContracts queries the Babylon module for all contract addresses and panics on error.
type BabylonContracts struct {
	BabylonContract        string
	BtcLightClientContract string
	BtcStakingContract     string
	BtcFinalityContract    string
}

// BSNContracts represents the hardcoded response structure for BSN contracts
type BSNContracts struct {
	BabylonContract        string `json:"babylon_contract"`
	BtcLightClientContract string `json:"btc_light_client_contract"`
	BtcStakingContract     string `json:"btc_staking_contract"`
	BtcFinalityContract    string `json:"btc_finality_contract"`
}

// QueryBSNContractsResponse represents the hardcoded gRPC response
type QueryBSNContractsResponse struct {
	BsnContracts *BSNContracts `json:"bsn_contracts"`
}

func (cc *CosmwasmConsumerController) MustQueryBabylonContracts() *BabylonContracts {
	contracts, err := cc.QueryBabylonContracts()
	if err != nil {
		panic(err)
	}
	return contracts
}

func (cc *CosmwasmConsumerController) QueryBabylonContracts() (*BabylonContracts, error) {
	// Hardcoded implementation for e2e tests
	// In a real e2e test environment, these contracts would be set up by the test infrastructure
	// For now, return an error indicating contracts need to be set up
	// TODO: This should be replaced with actual contract addresses once they're deployed in tests
	return nil, fmt.Errorf("Babylon contracts not yet implemented in hardcoded e2e test - contracts should be set up by test infrastructure")
}
