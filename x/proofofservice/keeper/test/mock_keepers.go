package test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/proofofservice/types"
)

// MockStakingKeeper is a mock for the staking keeper
// used in testing the Proof-of-Service module

type MockStakingKeeper struct {
	Validators map[string]bool
}

// NewMockStakingKeeper creates a new instance of MockStakingKeeper
func NewMockStakingKeeper() *MockStakingKeeper {
	return &MockStakingKeeper{
		Validators: make(map[string]bool),
	}
}

// SetValidator sets a validator in the mock keeper
func (k *MockStakingKeeper) SetValidator(addr sdk.AccAddress, isValidator bool) {
	k.Validators[addr.String()] = isValidator
}

// IsValidator checks if an address is a validator
func (k *MockStakingKeeper) IsValidator(addr sdk.AccAddress) bool {
	return k.Validators[addr.String()]
}

// MockKeeper is a mock for the Proof-of-Service keeper
// used in testing

type MockKeeper struct {
	storeKey sdk.StoreKey
	cdc      codec.Codec
}

// NewMockKeeper creates a new instance of MockKeeper
func NewMockKeeper(storeKey sdk.StoreKey, cdc codec.Codec) *MockKeeper {
	return &MockKeeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// GetServiceProvider retrieves a service provider from the store
func (k *MockKeeper) GetServiceProvider(ctx sdk.Context, address string) (types.ServiceProvider, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetServiceProviderKey(address))
	if bz == nil {
		return types.ServiceProvider{}, false
	}
	var provider types.ServiceProvider
	k.cdc.MustUnmarshal(bz, &provider)
	return provider, true
}

// SetServiceProvider sets a service provider in the store
func (k *MockKeeper) SetServiceProvider(ctx sdk.Context, provider types.ServiceProvider) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&provider)
	store.Set(types.GetServiceProviderKey(provider.Address), bz)
}

// GetProof retrieves a proof from the store
func (k *MockKeeper) GetProof(ctx sdk.Context, provider, proofID string) (types.Proof, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetProofKey(provider, proofID))
	if bz == nil {
		return types.Proof{}, false
	}
	var proof types.Proof
	k.cdc.MustUnmarshal(bz, &proof)
	return proof, true
}

// SetProof sets a proof in the store
func (k *MockKeeper) SetProof(ctx sdk.Context, proof types.Proof) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&proof)
	store.Set(types.GetProofKey(proof.Provider, proof.ProofID), bz)
}
