package types

const (
	// ModuleName defines the module name
	ModuleName = "servrewards"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_servrewards"
)

var (
	// RewardMetricsKey is the key to store reward metrics
	RewardMetricsKey = []byte{0x01}

	// AccumulatedRewardsPrefix is the prefix for storing accumulated rewards
	AccumulatedRewardsPrefix = []byte{0x02}

	// RewardParamsKey is the key to store reward parameters
	RewardParamsKey = []byte{0x03}
)

// GetAccumulatedRewardsKey returns the key for storing accumulated rewards for an address
func GetAccumulatedRewardsKey(addr string) []byte {
	return append(AccumulatedRewardsPrefix, []byte(addr)...)
}
