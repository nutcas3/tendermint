package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakingKeeper defines the expected staking keeper
type StakingKeeper interface {
	GetAllValidators(ctx sdk.Context) []StakingValidator
	GetValidator(ctx sdk.Context, addr sdk.AccAddress) (validator StakingValidator, found bool)
	GetValidatorSigningInfo(ctx sdk.Context, consAddr sdk.ConsAddress) (SigningInfo, bool)
	SignedBlocksWindow(ctx sdk.Context) int64
}

// DistrKeeper defines the expected distribution keeper
type DistrKeeper interface {
	AllocateTokensToValidator(ctx sdk.Context, val sdk.ValAddress, tokens sdk.DecCoins)
}

// ProofOfServiceKeeper defines the expected proof of service keeper
type ProofOfServiceKeeper interface {
	GetServiceScore(ctx sdk.Context, addr string) sdk.Int
	GetTotalServiceScore(ctx sdk.Context) sdk.Int
}

// StakingValidator defines the expected validator interface
type StakingValidator interface {
	GetOperator() sdk.ValAddress
	GetConsAddr() sdk.ConsAddress
}

// SigningInfo defines the expected signing info interface
type SigningInfo interface {
	GetMissedBlocksCounter() int64
}

// NodeRewardsHooks event hooks for node rewards module
type NodeRewardsHooks interface {
	AfterNodePerformanceUpdated(ctx sdk.Context, validatorAddr string)
}
