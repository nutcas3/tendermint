package noderewards

import (
	"github.com/serv-chain/serv/x/noderewards/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Update node performance metrics at the beginning of each block
	// This ensures that rewards distributed in this block use the latest metrics
	k.UpdateAllNodePerformances(ctx)
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []sdk.ValidatorUpdate {
	// No specific actions needed at the end of the block for noderewards
	// The actual reward modification happens through hooks into the distribution module
	return []sdk.ValidatorUpdate{}
}
