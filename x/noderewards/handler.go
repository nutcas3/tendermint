package noderewards

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/serv-chain/serv/x/noderewards/keeper"
	"github.com/serv-chain/serv/x/noderewards/types"
)

// NewHandler returns a handler for "noderewards" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	// Node rewards module doesn't have any messages to handle directly
	// It primarily works through hooks and BeginBlocker/EndBlocker

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}
