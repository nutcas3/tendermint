package types

const (
	// ModuleName defines the module name
	ModuleName = "proofofservice"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_proofofservice"
)

var (
	// ServiceProviderPrefix is the prefix for storing service providers
	ServiceProviderPrefix = []byte{0x01}

	// ServiceProofPrefix is the prefix for storing service proofs
	ServiceProofPrefix = []byte{0x02}

	// ServiceScorePrefix is the prefix for storing service scores
	ServiceScorePrefix = []byte{0x03}

	// TotalServiceScoreKey is the key for storing total service score
	TotalServiceScoreKey = []byte{0x04}

	// ServiceParamsKey is the key for storing service parameters
	ServiceParamsKey = []byte{0x05}
)

// GetServiceProviderKey returns the key for storing a service provider
func GetServiceProviderKey(addr string) []byte {
	return append(ServiceProviderPrefix, []byte(addr)...)
}

// GetServiceProofKey returns the key for storing a service proof
func GetServiceProofKey(addr string, proofID string) []byte {
	addrKey := append(ServiceProofPrefix, []byte(addr)...)
	return append(addrKey, []byte(proofID)...)
}

// GetServiceScoreKey returns the key for storing a service score
func GetServiceScoreKey(addr string) []byte {
	return append(ServiceScorePrefix, []byte(addr)...)
}
