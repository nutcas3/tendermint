package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/serv-chain/serv/x/servrewards/types"
)

// Keeper of the servrewards store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	bankKeeper       types.BankKeeper
	stakingKeeper    types.StakingKeeper
	posKeeper        types.ProofOfServiceKeeper
	hooks            types.ServRewardsHooks
}

// NewKeeper creates a new servrewards Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
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
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		posKeeper:     posKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks sets the servrewards hooks
func (k *Keeper) SetHooks(h types.ServRewardsHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set servrewards hooks twice")
	}
	k.hooks = h
	return k
}

// GetRewardMetrics returns the current reward metrics
func (k Keeper) GetRewardMetrics(ctx sdk.Context) types.RewardMetrics {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardMetricsKey)
	if bz == nil {
		return types.DefaultRewardMetrics()
	}

	var metrics types.RewardMetrics
	k.cdc.MustUnmarshal(bz, &metrics)
	return metrics
}

// SetRewardMetrics sets the current reward metrics
func (k Keeper) SetRewardMetrics(ctx sdk.Context, metrics types.RewardMetrics) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&metrics)
	store.Set(types.RewardMetricsKey, bz)
}

// GetRewardParams returns the current reward parameters
func (k Keeper) GetRewardParams(ctx sdk.Context) types.RewardParams {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardParamsKey)
	if bz == nil {
		return types.DefaultRewardParams()
	}

	var params types.RewardParams
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetRewardParams sets the current reward parameters
func (k Keeper) SetRewardParams(ctx sdk.Context, params types.RewardParams) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.RewardParamsKey, bz)
}

// GetAccumulatedRewards returns the accumulated rewards for an address
func (k Keeper) GetAccumulatedRewards(ctx sdk.Context, addr string) types.AccumulatedRewards {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAccumulatedRewardsKey(addr)
	bz := store.Get(key)
	if bz == nil {
		return types.AccumulatedRewards{
			Address:   addr,
			Rewards:   sdk.ZeroInt(),
			LastClaim: 0,
		}
	}

	var rewards types.AccumulatedRewards
	k.cdc.MustUnmarshal(bz, &rewards)
	return rewards
}

// SetAccumulatedRewards sets the accumulated rewards for an address
func (k Keeper) SetAccumulatedRewards(ctx sdk.Context, rewards types.AccumulatedRewards) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAccumulatedRewardsKey(rewards.Address)
	bz := k.cdc.MustMarshal(&rewards)
	store.Set(key, bz)
}

// CalculateRewards calculates rewards for an address based on service score and staking amount
func (k Keeper) CalculateRewards(ctx sdk.Context, addr string) sdk.Int {
	params := k.GetRewardParams(ctx)
	metrics := k.GetRewardMetrics(ctx)
	
	// Get service score from proof of service module
	serviceScore := k.posKeeper.GetServiceScore(ctx, addr)
	
	// Get staking amount from staking module
	stakingAmount := k.stakingKeeper.GetDelegatorStake(ctx, sdk.MustAccAddressFromBech32(addr))
	
	// Calculate rewards based on service score and staking amount
	serviceReward := sdk.NewDecFromInt(params.RewardPerEpoch).
		Mul(params.ServiceScoreWeight).
		Mul(sdk.NewDecFromInt(serviceScore)).
		Quo(sdk.NewDecFromInt(metrics.TotalServiceScore))
	
	stakingReward := sdk.NewDecFromInt(params.RewardPerEpoch).
		Mul(params.StakingWeight).
		Mul(sdk.NewDecFromInt(stakingAmount)).
		Quo(sdk.NewDecFromInt(metrics.TotalStaked))
	
	totalReward := serviceReward.Add(stakingReward).TruncateInt()
	
	return totalReward
}

// ClaimRewards claims accumulated rewards for an address
func (k Keeper) ClaimRewards(ctx sdk.Context, addr string) (sdk.Int, error) {
	rewards := k.GetAccumulatedRewards(ctx, addr)
	metrics := k.GetRewardMetrics(ctx)
	
	// Check if rewards have already been claimed for this epoch
	if rewards.LastClaim == metrics.EpochNumber {
		return sdk.ZeroInt(), fmt.Errorf("rewards already claimed for this epoch")
	}
	
	// Mint coins to the address
	coins := sdk.NewCoins(sdk.NewCoin("serv", rewards.Rewards))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return sdk.ZeroInt(), err
	}
	
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(addr), coins); err != nil {
		return sdk.ZeroInt(), err
	}
	
	// Update accumulated rewards
	claimedAmount := rewards.Rewards
	rewards.Rewards = sdk.ZeroInt()
	rewards.LastClaim = metrics.EpochNumber
	k.SetAccumulatedRewards(ctx, rewards)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewardClaimed,
			sdk.NewAttribute(types.AttributeKeyAddress, addr),
			sdk.NewAttribute(types.AttributeKeyAmount, claimedAmount.String()),
			sdk.NewAttribute(types.AttributeKeyEpoch, fmt.Sprintf("%d", metrics.EpochNumber)),
		),
	)
	
	return claimedAmount, nil
}

// UpdateRewards updates accumulated rewards for all addresses at the end of an epoch
func (k Keeper) UpdateRewards(ctx sdk.Context) {
	metrics := k.GetRewardMetrics(ctx)
	metrics.EpochNumber++
	
	// Update total service score and total staked
	metrics.TotalServiceScore = k.posKeeper.GetTotalServiceScore(ctx)
	metrics.TotalStaked = k.stakingKeeper.GetTotalBondedTokens(ctx)
	
	// Update metrics
	k.SetRewardMetrics(ctx, metrics)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEpochCompleted,
			sdk.NewAttribute(types.AttributeKeyEpoch, fmt.Sprintf("%d", metrics.EpochNumber)),
			sdk.NewAttribute(types.AttributeKeyTotalServiceScore, metrics.TotalServiceScore.String()),
			sdk.NewAttribute(types.AttributeKeyTotalStaked, metrics.TotalStaked.String()),
		),
	)
	
	// Call hooks if set
	if k.hooks != nil {
		k.hooks.AfterEpochCompleted(ctx, metrics.EpochNumber)
	}
}
