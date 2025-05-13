package servrewards

import (
	"github.com/serv-chain/serv/x/servrewards/keeper"
	"github.com/serv-chain/serv/x/servrewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// No specific actions needed at the beginning of the block for servrewards
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []sdk.ValidatorUpdate {
	// Check if this is the end of an epoch
	metrics := k.GetRewardMetrics(ctx)
	params := k.GetRewardParams(ctx)
	
	// If we've reached the end of an epoch, update rewards
	if uint64(ctx.BlockHeight())%params.EpochDuration == 0 {
		k.UpdateRewards(ctx)
		
		// Log epoch completion
		ctx.Logger().Info("SERV Rewards epoch completed", 
			"epoch", metrics.EpochNumber,
			"total_service_score", metrics.TotalServiceScore,
			"total_staked", metrics.TotalStaked)
	}
	
	return []sdk.ValidatorUpdate{}
}
