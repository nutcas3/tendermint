package servrewards

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/servrewards/keeper"
	"github.com/serv-chain/serv/x/servrewards/types"
)

// InitGenesis initializes the servrewards module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set reward metrics
	k.SetRewardMetrics(ctx, genState.RewardMetrics)
	
	// Set reward parameters
	k.SetRewardParams(ctx, genState.RewardParams)
	
	// Set accumulated rewards
	for _, reward := range genState.AccumulatedRewards {
		k.SetAccumulatedRewards(ctx, reward)
	}
}

// ExportGenesis returns the servrewards module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	rewardMetrics := k.GetRewardMetrics(ctx)
	rewardParams := k.GetRewardParams(ctx)
	
	// Get all accumulated rewards
	// In a real implementation, we would iterate over all accumulated rewards in the store
	// For now, we'll return an empty list
	accumulatedRewards := []types.AccumulatedRewards{}
	
	return &types.GenesisState{
		RewardMetrics:      rewardMetrics,
		RewardParams:       rewardParams,
		AccumulatedRewards: accumulatedRewards,
	}
}
