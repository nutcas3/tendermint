package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/proofofservice/keeper"
	"github.com/serv-chain/serv/x/proofofservice/types"
)

// Setup initializes a test keeper with mock dependencies
func Setup(t *testing.T) (*keeper.Keeper, sdk.Context, *MockStakingKeeper) {
	// Initialize keepers
	stakingKeeper := NewMockStakingKeeper()

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
		stakingKeeper,
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

	return k, ctx, stakingKeeper
}

// TestGetServiceParams tests the GetServiceParams function
func TestGetServiceParams(t *testing.T) {
	k, ctx, _ := Setup(t)

	// Test default params
	params := k.GetServiceParams(ctx)
	require.Equal(t, types.DefaultServiceParams(), params)

	// Set custom params
	customParams := types.ServiceParams{
		MinVerifications:    5,
		ScoreDecayRate:      sdk.NewDecWithPrec(2, 1), // 0.2
		ProofValidityPeriod: 200,
		MaxProofsPerEpoch:   10,
	}
	k.SetServiceParams(ctx, customParams)

	// Test getting custom params
	params = k.GetServiceParams(ctx)
	require.Equal(t, customParams, params)
}

// TestRegisterServiceProvider tests the RegisterServiceProvider function
func TestRegisterServiceProvider(t *testing.T) {
	k, ctx, _ := Setup(t)

	provider := "cosmos1abcdef"
	serviceType := "storage"
	metadata := "{\"capacity\":\"1TB\",\"region\":\"us-east\"}"

	// Register service provider
	err := k.RegisterServiceProvider(ctx, provider, serviceType, metadata)
	require.NoError(t, err)

	// Check that provider was registered
	serviceProvider, found := k.GetServiceProvider(ctx, provider)
	require.True(t, found)
	require.Equal(t, provider, serviceProvider.Address)
	require.Equal(t, serviceType, serviceProvider.ServiceType)
	require.Equal(t, metadata, serviceProvider.Metadata)
	require.True(t, serviceProvider.Active)

	// Check that service score was initialized
	score := k.GetServiceScore(ctx, provider)
	require.True(t, score.IsZero())

	// Try to register the same provider again
	err = k.RegisterServiceProvider(ctx, provider, serviceType, metadata)
	require.Error(t, err)
	require.Contains(t, err.Error(), "service provider already registered")
}

// TestSubmitProof tests the SubmitProof function
func TestSubmitProof(t *testing.T) {
	k, ctx, _ := Setup(t)

	provider := "cosmos1abcdef"
	serviceType := "storage"
	metadata := "{\"capacity\":\"1TB\",\"region\":\"us-east\"}"
	proofID := "proof-123"
	evidence := "hash-of-evidence-data"

	// First register the provider
	err := k.RegisterServiceProvider(ctx, provider, serviceType, metadata)
	require.NoError(t, err)

	// Submit proof
	err = k.SubmitProof(ctx, provider, serviceType, proofID, evidence)
	require.NoError(t, err)

	// Check that proof was submitted
	proof, found := k.GetProof(ctx, provider, proofID)
	require.True(t, found)
	require.Equal(t, proofID, proof.ProofID)
	require.Equal(t, provider, proof.Provider)
	require.Equal(t, serviceType, proof.ServiceType)
	require.Equal(t, evidence, proof.Evidence)
	require.False(t, proof.Verified)
	require.Empty(t, proof.VerifiedBy)
	require.True(t, proof.Score.IsZero())

	// Try to submit the same proof again
	err = k.SubmitProof(ctx, provider, serviceType, proofID, evidence)
	require.Error(t, err)
	require.Contains(t, err.Error(), "proof already submitted")

	// Try to submit proof for unregistered provider
	err = k.SubmitProof(ctx, "cosmos1xyz", serviceType, "proof-456", evidence)
	require.Error(t, err)
	require.Contains(t, err.Error(), "service provider not registered")
}

// TestVerifyProof tests the VerifyProof function
func TestVerifyProof(t *testing.T) {
	k, ctx, stakingKeeper := Setup(t)

	provider := "cosmos1abcdef"
	serviceType := "storage"
	metadata := "{\"capacity\":\"1TB\",\"region\":\"us-east\"}"
	proofID := "proof-123"
	evidence := "hash-of-evidence-data"
	validator := "cosmos1validator"
	validatorAddr, _ := sdk.AccAddressFromBech32(validator)

	// Set up validator
	stakingKeeper.SetValidator(validatorAddr, true)

	// Register provider and submit proof
	err := k.RegisterServiceProvider(ctx, provider, serviceType, metadata)
	require.NoError(t, err)
	err = k.SubmitProof(ctx, provider, serviceType, proofID, evidence)
	require.NoError(t, err)

	// Verify proof
	err = k.VerifyProof(ctx, validator, provider, proofID, true, 80)
	require.NoError(t, err)

	// Check that proof was updated
	proof, found := k.GetProof(ctx, provider, proofID)
	require.True(t, found)
	require.Contains(t, proof.VerifiedBy, validator)
	
	// Proof shouldn't be verified yet (not enough verifications)
	require.False(t, proof.Verified)
	require.True(t, proof.Score.IsZero())

	// Add more verifications to reach minimum
	params := k.GetServiceParams(ctx)
	for i := 1; i < int(params.MinVerifications); i++ {
		validatorAddr := fmt.Sprintf("cosmos1validator%d", i)
		stakingKeeper.SetValidator(sdk.MustAccAddressFromBech32(validatorAddr), true)
		err = k.VerifyProof(ctx, validatorAddr, provider, proofID, true, 80)
		require.NoError(t, err)
	}

	// Check that proof is now verified
	proof, found = k.GetProof(ctx, provider, proofID)
	require.True(t, found)
	require.True(t, proof.Verified)
	require.Equal(t, sdk.NewIntFromUint64(80), proof.Score)

	// Check that service score was updated
	score := k.GetServiceScore(ctx, provider)
	require.Equal(t, sdk.NewIntFromUint64(80), score)

	// Try to verify with the same validator again
	err = k.VerifyProof(ctx, validator, provider, proofID, true, 90)
	require.Error(t, err)
	require.Contains(t, err.Error(), "validator has already verified this proof")

	// Try to verify with non-validator
	err = k.VerifyProof(ctx, "cosmos1nonvalidator", provider, proofID, true, 90)
	require.Error(t, err)
	require.Contains(t, err.Error(), "address is not a validator")

	// Try to verify non-existent proof
	err = k.VerifyProof(ctx, validator, provider, "non-existent", true, 90)
	require.Error(t, err)
	require.Contains(t, err.Error(), "proof not found")
}

// TestDecayServiceScores tests the DecayServiceScores function
func TestDecayServiceScores(t *testing.T) {
	k, ctx, _ := Setup(t)

	// Register providers and set scores
	providers := []string{"cosmos1a", "cosmos1b", "cosmos1c"}
	scores := []int64{100, 200, 300}
	
	for i, provider := range providers {
		err := k.RegisterServiceProvider(ctx, provider, "storage", "")
		require.NoError(t, err)
		
		// Manually set service score
		store := ctx.KVStore(k.GetStoreKey())
		scoreKey := types.GetServiceScoreKey(provider)
		serviceScore := types.ServiceScore{
			Provider:    provider,
			Score:       sdk.NewInt(scores[i]),
			LastUpdated: 0,
		}
		bz := k.GetCodec().MustMarshal(&serviceScore)
		store.Set(scoreKey, bz)
	}
	
	// Set total service score
	totalScore := sdk.NewInt(600) // 100 + 200 + 300
	totalScoreBz := k.GetCodec().MustMarshal(&totalScore)
	ctx.KVStore(k.GetStoreKey()).Set(types.TotalServiceScoreKey, totalScoreBz)
	
	// Set decay rate to 10%
	params := k.GetServiceParams(ctx)
	params.ScoreDecayRate = sdk.NewDecWithPrec(1, 1) // 0.1
	k.SetServiceParams(ctx, params)
	
	// Decay scores
	k.DecayServiceScores(ctx)
	
	// Check decayed scores
	for i, provider := range providers {
		score := k.GetServiceScore(ctx, provider)
		expectedScore := sdk.NewInt(int64(float64(scores[i]) * 0.9))
		require.Equal(t, expectedScore, score)
	}
	
	// Check total score was recalculated
	totalScore = k.GetTotalServiceScore(ctx)
	expectedTotal := sdk.NewInt(540) // (100 + 200 + 300) * 0.9
	require.Equal(t, expectedTotal, totalScore)
}
