package proofofservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/proofofservice/keeper"
	"github.com/serv-chain/serv/x/proofofservice/types"
)

// InitGenesis initializes the proofofservice module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set service parameters
	k.SetServiceParams(ctx, genState.ServiceParams)
	
	// Register service providers
	for _, provider := range genState.ServiceProviders {
		// In a real implementation, we would use the keeper's RegisterServiceProvider method
		// For now, we'll just store the provider directly
		store := ctx.KVStore(k.GetStoreKey())
		key := types.GetServiceProviderKey(provider.Address)
		bz := k.GetCodec().MustMarshal(&provider)
		store.Set(key, bz)
		
		// Initialize service score
		serviceScore := types.ServiceScore{
			Provider:    provider.Address,
			Score:       sdk.ZeroInt(),
			LastUpdated: 0,
		}
		
		scoreKey := types.GetServiceScoreKey(provider.Address)
		scoreBz := k.GetCodec().MustMarshal(&serviceScore)
		store.Set(scoreKey, scoreBz)
	}
	
	// Store service proofs
	for _, proof := range genState.ServiceProofs {
		proofKey := types.GetServiceProofKey(proof.Provider, proof.ProofID)
		bz := k.GetCodec().MustMarshal(&proof)
		ctx.KVStore(k.GetStoreKey()).Set(proofKey, bz)
	}
	
	// Set total service score
	totalScoreBz := k.GetCodec().MustMarshal(&genState.TotalServiceScore)
	ctx.KVStore(k.GetStoreKey()).Set(types.TotalServiceScoreKey, totalScoreBz)
}

// ExportGenesis returns the proofofservice module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	serviceParams := k.GetServiceParams(ctx)
	
	// Get all service providers
	// In a real implementation, we would iterate over all providers in the store
	// For now, we'll return an empty list
	serviceProviders := []types.ServiceProvider{}
	
	// Get all service proofs
	// In a real implementation, we would iterate over all proofs in the store
	// For now, we'll return an empty list
	serviceProofs := []types.ServiceProof{}
	
	// Get total service score
	totalServiceScore := k.GetTotalServiceScore(ctx)
	
	return &types.GenesisState{
		ServiceParams:     serviceParams,
		ServiceProviders:  serviceProviders,
		ServiceProofs:     serviceProofs,
		TotalServiceScore: totalServiceScore,
	}
}
