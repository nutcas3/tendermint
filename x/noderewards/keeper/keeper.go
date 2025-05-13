package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/serv-chain/serv/x/noderewards/types"
)

// Keeper of the noderewards store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	stakingKeeper    types.StakingKeeper
	distrKeeper      types.DistrKeeper
	posKeeper        types.ProofOfServiceKeeper
	hooks            types.NodeRewardsHooks
}

// NewKeeper creates a new noderewards Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	stakingKeeper types.StakingKeeper,
	distrKeeper types.DistrKeeper,
	posKeeper types.ProofOfServiceKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramstore:    ps,
		stakingKeeper: stakingKeeper,
		distrKeeper:   distrKeeper,
		posKeeper:     posKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks sets the noderewards hooks
func (k *Keeper) SetHooks(h types.NodeRewardsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set noderewards hooks twice")
	}
	k.hooks = h
	return k
}

// GetRewardModifier returns the current reward modifier parameters
func (k Keeper) GetRewardModifier(ctx sdk.Context) types.RewardModifier {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardModifierKey)
	if bz == nil {
		return types.DefaultRewardModifier()
	}

	var modifier types.RewardModifier
	k.cdc.MustUnmarshal(bz, &modifier)
	return modifier
}

// SetRewardModifier sets the current reward modifier parameters
func (k Keeper) SetRewardModifier(ctx sdk.Context, modifier types.RewardModifier) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&modifier)
	store.Set(types.RewardModifierKey, bz)
}

// GetNodePerformance returns the performance metrics for a validator node
func (k Keeper) GetNodePerformance(ctx sdk.Context, validatorAddr string) types.NodePerformance {
	store := ctx.KVStore(k.storeKey)
	key := types.GetNodePerformanceKey(validatorAddr)
	
	bz := store.Get(key)
	if bz == nil {
		return types.DefaultNodePerformance(validatorAddr)
	}
	
	var performance types.NodePerformance
	k.cdc.MustUnmarshal(bz, &performance)
	return performance
}

// SetNodePerformance sets the performance metrics for a validator node
func (k Keeper) SetNodePerformance(ctx sdk.Context, performance types.NodePerformance) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetNodePerformanceKey(performance.ValidatorAddr)
	bz := k.cdc.MustMarshal(&performance)
	store.Set(key, bz)
}

// UpdateNodePerformance updates the performance metrics for a validator node
func (k Keeper) UpdateNodePerformance(ctx sdk.Context, validatorAddr string) {
	// Get current performance
	performance := k.GetNodePerformance(ctx, validatorAddr)
	
	// Update service score from proof of service module
	performance.ServiceScore = k.posKeeper.GetServiceScore(ctx, validatorAddr)
	
	// Update uptime from slashing module (via staking keeper)
	validator, found := k.stakingKeeper.GetValidator(ctx, sdk.MustAccAddressFromBech32(validatorAddr))
	if found {
		// Calculate uptime based on signing info
		signInfo, found := k.stakingKeeper.GetValidatorSigningInfo(ctx, validator.GetConsAddr())
		if found {
			missedBlocks := signInfo.MissedBlocksCounter
			windowSize := k.stakingKeeper.SignedBlocksWindow(ctx)
			uptime := sdk.OneDec().Sub(sdk.NewDec(missedBlocks).QuoInt64(windowSize))
			performance.UptimePercent = uptime
		}
	}
	
	// Update response time (this would typically come from monitoring data)
	// For this implementation, we'll use a placeholder value
	performance.ResponseTime = sdk.NewInt(100) // 100ms
	
	// Update last update height
	performance.LastUpdateHeight = ctx.BlockHeight()
	
	// Save updated performance
	k.SetNodePerformance(ctx, performance)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNodePerformanceUpdated,
			sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr),
			sdk.NewAttribute(types.AttributeKeyServiceScore, performance.ServiceScore.String()),
			sdk.NewAttribute(types.AttributeKeyUptime, performance.UptimePercent.String()),
			sdk.NewAttribute(types.AttributeKeyResponseTime, performance.ResponseTime.String()),
		),
	)
	
	// Call hooks if set
	if k.hooks != nil {
		k.hooks.AfterNodePerformanceUpdated(ctx, validatorAddr)
	}
}

// CalculateRewardModifier calculates the reward modifier for a validator
func (k Keeper) CalculateRewardModifier(ctx sdk.Context, validatorAddr string) sdk.Dec {
	performance := k.GetNodePerformance(ctx, validatorAddr)
	modifier := k.GetRewardModifier(ctx)
	
	// Get total service score for normalization
	totalServiceScore := k.posKeeper.GetTotalServiceScore(ctx)
	if totalServiceScore.IsZero() {
		totalServiceScore = sdk.OneInt() // Avoid division by zero
	}
	
	// Normalize service score (0-1)
	normalizedServiceScore := sdk.NewDecFromInt(performance.ServiceScore).QuoInt(totalServiceScore)
	
	// Response time score (lower is better, 1000ms is considered worst case)
	responseTimeScore := sdk.OneDec().Sub(sdk.NewDecFromInt(performance.ResponseTime).QuoInt64(1000))
	if responseTimeScore.IsNegative() {
		responseTimeScore = sdk.ZeroDec()
	}
	
	// Calculate weighted score
	weightedScore := normalizedServiceScore.Mul(modifier.ServiceScoreWeight).
		Add(performance.UptimePercent.Mul(modifier.UptimeWeight)).
		Add(responseTimeScore.Mul(modifier.ResponseTimeWeight))
	
	// Scale to min-max range
	// 0 score = min modifier, 1 score = max modifier
	scaledModifier := modifier.MinModifier.Add(
		weightedScore.Mul(modifier.MaxModifier.Sub(modifier.MinModifier)),
	)
	
	// Ensure modifier is within bounds
	if scaledModifier.LT(modifier.MinModifier) {
		return modifier.MinModifier
	}
	if scaledModifier.GT(modifier.MaxModifier) {
		return modifier.MaxModifier
	}
	
	return scaledModifier
}

// ModifyValidatorReward modifies the reward for a validator based on performance
func (k Keeper) ModifyValidatorReward(ctx sdk.Context, validatorAddr string, baseReward sdk.Dec) sdk.Dec {
	modifier := k.CalculateRewardModifier(ctx, validatorAddr)
	modifiedReward := baseReward.Mul(modifier)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewardModified,
			sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr),
			sdk.NewAttribute(types.AttributeKeyBaseReward, baseReward.String()),
			sdk.NewAttribute(types.AttributeKeyModifier, modifier.String()),
			sdk.NewAttribute(types.AttributeKeyModifiedReward, modifiedReward.String()),
		),
	)
	
	return modifiedReward
}

// UpdateAllNodePerformances updates performance metrics for all validators
func (k Keeper) UpdateAllNodePerformances(ctx sdk.Context) {
	validators := k.stakingKeeper.GetAllValidators(ctx)
	for _, validator := range validators {
		k.UpdateNodePerformance(ctx, validator.GetOperator().String())
	}
}
