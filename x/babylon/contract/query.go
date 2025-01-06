package contract

// CustomQuery is a query request from a smart contract to the Babylon module
type CustomQuery struct {
	Params *ParamsQuery `json:"params,omitempty"`
}

// ParamsQuery requests the current module parameters
type ParamsQuery struct{}

// ParamsResponse contains the current module parameters
type ParamsResponse struct {
	BabylonContractCodeId      uint64 `json:"babylon_contract_code_id,omitempty"`
	BtcStakingContractCodeId   uint64 `json:"btc_staking_contract_code_id,omitempty"`
	BtcFinalityContractCodeId  uint64 `json:"btc_finality_contract_code_id,omitempty"`
	BabylonContractAddress     string `json:"babylon_contract_address,omitempty"`
	BtcStakingContractAddress  string `json:"btc_staking_contract_address,omitempty"`
	BtcFinalityContractAddress string `json:"btc_finality_contract_address,omitempty"`
	MaxGasBeginBlocker         uint32 `json:"max_gas_begin_blocker,omitempty"`
}
