package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/serv-chain/serv/x/proofofservice/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServer returns an implementation of the MsgServer interface
// for the proofofservice module.
func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// RegisterService implements the MsgServer.RegisterService method.
func (m msgServer) RegisterService(goCtx context.Context, msg *types.MsgRegisterService) (*types.MsgRegisterServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message sender
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address: %s", err)
	}

	// Register service provider
	err = m.Keeper.RegisterServiceProvider(ctx, msg.Provider, msg.ServiceType, msg.Metadata)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyServiceType, msg.ServiceType),
		),
	})

	return &types.MsgRegisterServiceResponse{}, nil
}

// SubmitProof implements the MsgServer.SubmitProof method.
func (m msgServer) SubmitProof(goCtx context.Context, msg *types.MsgSubmitProof) (*types.MsgSubmitProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message sender
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address: %s", err)
	}

	// Submit proof
	err = m.Keeper.SubmitProof(ctx, msg.Provider, msg.ServiceType, msg.ProofId, msg.Evidence)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyServiceType, msg.ServiceType),
			sdk.NewAttribute(types.AttributeKeyProofID, msg.ProofId),
		),
	})

	return &types.MsgSubmitProofResponse{}, nil
}

// VerifyProof implements the MsgServer.VerifyProof method.
func (m msgServer) VerifyProof(goCtx context.Context, msg *types.MsgVerifyProof) (*types.MsgVerifyProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message sender
	_, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address: %s", err)
	}

	// Validate the provider address
	_, err = sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address: %s", err)
	}

	// Verify proof
	err = m.Keeper.VerifyProof(ctx, msg.Validator, msg.Provider, msg.ProofId, msg.IsVerified, msg.Score)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyProvider, msg.Provider),
			sdk.NewAttribute(types.AttributeKeyProofID, msg.ProofId),
			sdk.NewAttribute(types.AttributeKeyVerified, sdk.FormatBool(msg.IsVerified)),
			sdk.NewAttribute(types.AttributeKeyScore, sdk.FormatUint(msg.Score)),
		),
	})

	return &types.MsgVerifyProofResponse{}, nil
}
