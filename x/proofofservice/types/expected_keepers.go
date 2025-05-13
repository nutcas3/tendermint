package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakingKeeper defines the expected staking keeper
type StakingKeeper interface {
	IsValidator(ctx sdk.Context, addr sdk.AccAddress) bool
	GetValidator(ctx sdk.Context, valAddr sdk.AccAddress) (validator StakingValidator, found bool)
	GetValidatorSigningInfo(ctx sdk.Context, consAddr sdk.ConsAddress) (SigningInfo, bool)
	SignedBlocksWindow(ctx sdk.Context) int64
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

// ProofOfServiceHooks event hooks for proof of service module
type ProofOfServiceHooks interface {
	AfterServiceProviderRegistered(ctx sdk.Context, provider string)
	AfterProofSubmitted(ctx sdk.Context, provider string, proofID string)
	AfterProofVerified(ctx sdk.Context, provider string, proofID string, score sdk.Int)
}
