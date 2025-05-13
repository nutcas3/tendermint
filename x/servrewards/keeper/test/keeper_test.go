package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/serv-chain/serv/x/servrewards/keeper"
	"github.com/serv-chain/serv/x/servrewards/types"
)

// Setup initializes a test keeper with mock dependencies
func Setup(t *testing.T) (*keeper.Keeper, sdk.Context, *MockBankKeeper, *MockStakingKeeper, *MockPosKeeper) {
	// Initialize keepers
	bankKeeper := NewMockBankKeeper()
	stakingKeeper := NewMockStakingKeeper()
	posKeeper := NewMockPosKeeper()

	// Initialize codec
	encodingConfig := MakeTestEncodingConfig()
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := sdk.NewKVStoreKey(types.MemStoreKey)

	// Initialize params keeper and subspace
	paramsKeeper := initParamsKeeper(
		encodingConfig.Marshaler,
		encodingConfig.Amino,
		storeKey,
		memStoreKey,
	)
	subspace := paramsKeeper.Subspace(types.ModuleName)

	// Create test keeper
	k := keeper.NewKeeper(
		encodingConfig.Marshaler,
		storeKey,
		subspace,
		bankKeeper,
		stakingKeeper,
		posKeeper,
	)

	// Create test context
	ctx := sdk.NewContext(
		initKVStore(t, storeKey),
		tmproto.Header{Height: 1, Time: time.Now().UTC()},
		false,
		nil,
	)

	// Initialize params
	subspace.SetParamSet(ctx, &types.Params{})

	return k, ctx, bankKeeper, stakingKeeper, posKeeper
}

// TestGetRewardMetrics tests the GetRewardMetrics function
func TestGetRewardMetrics(t *testing.T) {
	k, ctx, _, _, _ := Setup(t)

	// Test default metrics
	metrics := k.GetRewardMetrics(ctx)
	require.Equal(t, types.DefaultRewardMetrics(), metrics)

	// Set custom metrics
	customMetrics := types.RewardMetrics{
		TotalServiceScore: sdk.NewInt(1000),
		TotalStaked:       sdk.NewInt(5000),
		EpochNumber:       5,
	}
	k.SetRewardMetrics(ctx, customMetrics)

	// Test getting custom metrics
	metrics = k.GetRewardMetrics(ctx)
	require.Equal(t, customMetrics, metrics)
}

// TestGetRewardParams tests the GetRewardParams function
func TestGetRewardParams(t *testing.T) {
	k, ctx, _, _, _ := Setup(t)

	// Test default params
	params := k.GetRewardParams(ctx)
	require.Equal(t, types.DefaultRewardParams(), params)

	// Set custom params
	customParams := types.RewardParams{
		ServiceScoreWeight: sdk.NewDecWithPrec(7, 1), // 0.7
		StakingWeight:      sdk.NewDecWithPrec(3, 1), // 0.3
		RewardPerEpoch:     sdk.NewInt(2000000),
		EpochDuration:      200,
	}
	k.SetRewardParams(ctx, customParams)

	// Test getting custom params
	params = k.GetRewardParams(ctx)
	require.Equal(t, customParams, params)
}

// TestGetAccumulatedRewards tests the GetAccumulatedRewards function
func TestGetAccumulatedRewards(t *testing.T) {
	k, ctx, _, _, _ := Setup(t)

	addr := "cosmos1abcdef"

	// Test default rewards (zero)
	rewards := k.GetAccumulatedRewards(ctx, addr)
	require.Equal(t, addr, rewards.Address)
	require.True(t, rewards.Rewards.IsZero())
	require.Equal(t, uint64(0), rewards.LastClaim)

	// Set custom rewards
	customRewards := types.AccumulatedRewards{
		Address:   addr,
		Rewards:   sdk.NewInt(1000),
		LastClaim: 5,
	}
	k.SetAccumulatedRewards(ctx, customRewards)

	// Test getting custom rewards
	rewards = k.GetAccumulatedRewards(ctx, addr)
	require.Equal(t, customRewards, rewards)
}

// TestCalculateRewards tests the CalculateRewards function
func TestCalculateRewards(t *testing.T) {
	k, ctx, _, stakingKeeper, posKeeper := Setup(t)

	addr := "cosmos1abcdef"
	addrAcc, _ := sdk.AccAddressFromBech32(addr)

	// Set up test data
	params := types.RewardParams{
		ServiceScoreWeight: sdk.NewDecWithPrec(6, 1), // 0.6
		StakingWeight:      sdk.NewDecWithPrec(4, 1), // 0.4
		RewardPerEpoch:     sdk.NewInt(1000),
		EpochDuration:      100,
	}
	k.SetRewardParams(ctx, params)

	metrics := types.RewardMetrics{
		TotalServiceScore: sdk.NewInt(1000),
		TotalStaked:       sdk.NewInt(10000),
		EpochNumber:       1,
	}
	k.SetRewardMetrics(ctx, metrics)

	// Mock service score and staking amount
	posKeeper.SetServiceScore(addr, sdk.NewInt(100))
	stakingKeeper.SetDelegatorStake(addrAcc, sdk.NewInt(1000))

	// Calculate rewards
	reward := k.CalculateRewards(ctx, addr)

	// Expected calculation:
	// serviceReward = 1000 * 0.6 * 100 / 1000 = 60
	// stakingReward = 1000 * 0.4 * 1000 / 10000 = 40
	// totalReward = 60 + 40 = 100
	require.Equal(t, sdk.NewInt(100), reward)
}

// TestClaimRewards tests the ClaimRewards function
func TestClaimRewards(t *testing.T) {
	k, ctx, bankKeeper, _, _ := Setup(t)

	addr := "cosmos1abcdef"
	addrAcc, _ := sdk.AccAddressFromBech32(addr)

	// Set up test data
	metrics := types.RewardMetrics{
		TotalServiceScore: sdk.NewInt(1000),
		TotalStaked:       sdk.NewInt(10000),
		EpochNumber:       5,
	}
	k.SetRewardMetrics(ctx, metrics)

	rewards := types.AccumulatedRewards{
		Address:   addr,
		Rewards:   sdk.NewInt(1000),
		LastClaim: 0,
	}
	k.SetAccumulatedRewards(ctx, rewards)

	// Claim rewards
	claimed, err := k.ClaimRewards(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(1000), claimed)

	// Check that rewards were minted and sent
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("serv", sdk.NewInt(1000))), bankKeeper.MintedCoins)
	require.Equal(t, addrAcc, bankKeeper.SentCoinsToAddr)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin("serv", sdk.NewInt(1000))), bankKeeper.SentCoins)

	// Check that accumulated rewards were updated
	updatedRewards := k.GetAccumulatedRewards(ctx, addr)
	require.True(t, updatedRewards.Rewards.IsZero())
	require.Equal(t, metrics.EpochNumber, updatedRewards.LastClaim)

	// Try to claim again in the same epoch
	_, err = k.ClaimRewards(ctx, addr)
	require.Error(t, err)
	require.Contains(t, err.Error(), "rewards already claimed for this epoch")
}

// TestUpdateRewards tests the UpdateRewards function
func TestUpdateRewards(t *testing.T) {
	k, ctx, _, _, posKeeper := Setup(t)

	// Set up test data
	metrics := types.RewardMetrics{
		TotalServiceScore: sdk.NewInt(1000),
		TotalStaked:       sdk.NewInt(10000),
		EpochNumber:       5,
	}
	k.SetRewardMetrics(ctx, metrics)

	// Mock total service score and total staked
	posKeeper.SetTotalServiceScore(sdk.NewInt(2000))

	// Update rewards
	k.UpdateRewards(ctx)

	// Check that metrics were updated
	updatedMetrics := k.GetRewardMetrics(ctx)
	require.Equal(t, uint64(6), updatedMetrics.EpochNumber)
	require.Equal(t, sdk.NewInt(2000), updatedMetrics.TotalServiceScore)
}
