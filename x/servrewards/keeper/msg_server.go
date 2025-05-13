package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/serv-chain/serv/x/servrewards/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServer returns an implementation of the MsgServer interface
// for the servrewards module.
func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// ClaimReward implements the MsgServer.ClaimReward method.
func (m msgServer) ClaimReward(goCtx context.Context, msg *types.MsgClaimReward) (*types.MsgClaimRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message sender
	_, err := sdk.AccAddressFromBech32(msg.Claimer)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid claimer address: %s", err)
	}

	// Claim rewards
	amount, err := m.Keeper.ClaimRewards(ctx, msg.Claimer)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Claimer),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return &types.MsgClaimRewardResponse{
		Amount: amount,
	}, nil
}

// UpdateRewardParams implements the MsgServer.UpdateRewardParams method.
func (m msgServer) UpdateRewardParams(goCtx context.Context, msg *types.MsgUpdateRewardParams) (*types.MsgUpdateRewardParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Ensure the message sender has the authority to update parameters
	// In a real implementation, this would check if the sender is the governance module account
	// or has another appropriate permission
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %s", err)
	}

	// Check if the authority is the governance module account
	govModuleAddr := govtypes.NewModuleAddress(govtypes.ModuleName)
	if !authority.Equals(govModuleAddr) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only the governance module can update reward parameters")
	}

	// Validate parameters
	if msg.ServiceScoreWeight.IsNegative() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service score weight cannot be negative")
	}

	if msg.StakingWeight.IsNegative() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking weight cannot be negative")
	}

	if msg.RewardPerEpoch.IsNegative() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "reward per epoch cannot be negative")
	}

	if msg.EpochDuration == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "epoch duration must be positive")
	}

	// Ensure weights sum to 1
	sumWeights := msg.ServiceScoreWeight.Add(msg.StakingWeight)
	if !sumWeights.Equal(sdk.NewDec(1)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("service score weight and staking weight must sum to 1, got: %s", sumWeights))
	}

	// Update parameters
	params := types.RewardParams{
		ServiceScoreWeight: msg.ServiceScoreWeight,
		StakingWeight:      msg.StakingWeight,
		RewardPerEpoch:     msg.RewardPerEpoch,
		EpochDuration:      msg.EpochDuration,
	}
	m.Keeper.SetRewardParams(ctx, params)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Authority),
			sdk.NewAttribute(types.AttributeKeyServiceScoreWeight, msg.ServiceScoreWeight.String()),
			sdk.NewAttribute(types.AttributeKeyStakingWeight, msg.StakingWeight.String()),
			sdk.NewAttribute(types.AttributeKeyRewardPerEpoch, msg.RewardPerEpoch.String()),
			sdk.NewAttribute(types.AttributeKeyEpochDuration, fmt.Sprintf("%d", msg.EpochDuration)),
		),
	})

	return &types.MsgUpdateRewardParamsResponse{}, nil
}
