package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// StakingKeeper defines the expected staking keeper
type StakingKeeper interface {
	GetDelegatorStake(ctx sdk.Context, delegator sdk.AccAddress) sdk.Int
	GetTotalBondedTokens(ctx sdk.Context) sdk.Int
}

// ProofOfServiceKeeper defines the expected proof of service keeper
type ProofOfServiceKeeper interface {
	GetServiceScore(ctx sdk.Context, addr string) sdk.Int
	GetTotalServiceScore(ctx sdk.Context) sdk.Int
}

// ServRewardsHooks event hooks for servrewards module
type ServRewardsHooks interface {
	AfterEpochCompleted(ctx sdk.Context, epochNumber uint64)
}
