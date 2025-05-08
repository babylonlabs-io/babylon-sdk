package types

const (
	// ModuleName defines the module name.
	ModuleName = "babylon"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "memory:babylon"

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

var (
	// ParamsKey is the prefix for the module parameters
	ParamsKey = []byte{0x1}

	// BTCStakingGaugeKey is the prefix for the BTC staking gauge
	BTCStakingGaugeKey = []byte{0x2}
)
