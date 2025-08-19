package types

const (
	EventTypeUnbond                     = "instant_unbond"
	EventTypeDelegate                   = "instant_delegate"
	EventTypeFeeCollectorError          = "fee_collector_error"
	EventTypeContractCommunicationError = "contract_communication_error"
)

const (
	AttributeKeyContractAddr = "virtual_staking_contract"
	AttributeKeyValidator    = "validator"
	AttributeKeyDelegator    = "delegator"
	AttributeKeyError        = "error"
	AttributeKeyHeight       = "height"
	AttributeKeyPhase        = "phase"
)
