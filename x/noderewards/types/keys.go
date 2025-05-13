package types

const (
	// ModuleName defines the module name
	ModuleName = "noderewards"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_noderewards"
)

var (
	// NodePerformancePrefix is the prefix for storing node performance metrics
	NodePerformancePrefix = []byte{0x01}

	// RewardModifierKey is the key for storing reward modifier parameters
	RewardModifierKey = []byte{0x02}
)

// GetNodePerformanceKey returns the key for storing node performance metrics
func GetNodePerformanceKey(validatorAddr string) []byte {
	return append(NodePerformancePrefix, []byte(validatorAddr)...)
}
