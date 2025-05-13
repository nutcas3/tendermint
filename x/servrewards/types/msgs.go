package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgClaimReward       = "claim_reward"
	TypeMsgUpdateRewardParams = "update_reward_params"
)

var _ sdk.Msg = &MsgClaimReward{}
var _ sdk.Msg = &MsgUpdateRewardParams{}

// MsgClaimReward defines a message for claiming accumulated rewards
type MsgClaimReward struct {
	Claimer string `json:"claimer"`
}

// NewMsgClaimReward creates a new MsgClaimReward instance
func NewMsgClaimReward(claimer string) *MsgClaimReward {
	return &MsgClaimReward{
		Claimer: claimer,
	}
}

// Route implements sdk.Msg
func (msg MsgClaimReward) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgClaimReward) Type() string {
	return TypeMsgClaimReward
}

// ValidateBasic implements sdk.Msg
func (msg MsgClaimReward) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Claimer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid claimer address: %s", err)
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgClaimReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements sdk.Msg
func (msg MsgClaimReward) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Claimer)
	return []sdk.AccAddress{addr}
}

// MsgUpdateRewardParams defines a message for updating reward parameters (governance)
type MsgUpdateRewardParams struct {
	Authority          string  `json:"authority"`
	ServiceScoreWeight sdk.Dec `json:"service_score_weight"`
	StakingWeight      sdk.Dec `json:"staking_weight"`
	RewardPerEpoch     sdk.Int `json:"reward_per_epoch"`
	EpochDuration      uint64  `json:"epoch_duration"`
}

// NewMsgUpdateRewardParams creates a new MsgUpdateRewardParams instance
func NewMsgUpdateRewardParams(
	authority string,
	serviceScoreWeight sdk.Dec,
	stakingWeight sdk.Dec,
	rewardPerEpoch sdk.Int,
	epochDuration uint64,
) *MsgUpdateRewardParams {
	return &MsgUpdateRewardParams{
		Authority:          authority,
		ServiceScoreWeight: serviceScoreWeight,
		StakingWeight:      stakingWeight,
		RewardPerEpoch:     rewardPerEpoch,
		EpochDuration:      epochDuration,
	}
}

// Route implements sdk.Msg
func (msg MsgUpdateRewardParams) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgUpdateRewardParams) Type() string {
	return TypeMsgUpdateRewardParams
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateRewardParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %s", err)
	}

	if msg.ServiceScoreWeight.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service score weight cannot be negative")
	}

	if msg.StakingWeight.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "staking weight cannot be negative")
	}

	if msg.RewardPerEpoch.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "reward per epoch cannot be negative")
	}

	if msg.EpochDuration == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "epoch duration must be positive")
	}

	// Ensure weights sum to 1
	sumWeights := msg.ServiceScoreWeight.Add(msg.StakingWeight)
	if !sumWeights.Equal(sdk.NewDec(1)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service score weight and staking weight must sum to 1")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgUpdateRewardParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateRewardParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}
