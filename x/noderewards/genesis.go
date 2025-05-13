package noderewards

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/noderewards/keeper"
	"github.com/serv-chain/serv/x/noderewards/types"
)

// InitGenesis initializes the noderewards module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set reward modifier parameters
	k.SetRewardModifier(ctx, genState.RewardModifier)
	
	// Set node performance metrics
	for _, performance := range genState.NodePerformances {
		k.SetNodePerformance(ctx, performance)
	}
}

// ExportGenesis returns the noderewards module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	rewardModifier := k.GetRewardModifier(ctx)
	
	// Get all node performances
	// In a real implementation, we would iterate over all node performances in the store
	// For now, we'll return an empty list
	nodePerformances := []types.NodePerformance{}
	
	return &types.GenesisState{
		RewardModifier:   rewardModifier,
		NodePerformances: nodePerformances,
	}
}
