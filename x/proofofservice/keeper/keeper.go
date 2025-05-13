package keeper

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/serv-chain/serv/x/proofofservice/types"
)

// Keeper of the proofofservice store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	stakingKeeper types.StakingKeeper
	hooks         types.ProofOfServiceHooks
}

// NewKeeper creates a new proofofservice Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	stakingKeeper types.StakingKeeper,
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
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks sets the proofofservice hooks
func (k *Keeper) SetHooks(h types.ProofOfServiceHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set proofofservice hooks twice")
	}
	k.hooks = h
	return k
}

// GetServiceParams returns the current service parameters
func (k Keeper) GetServiceParams(ctx sdk.Context) types.ServiceParams {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ServiceParamsKey)
	if bz == nil {
		return types.DefaultServiceParams()
	}

	var params types.ServiceParams
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetServiceParams sets the current service parameters
func (k Keeper) SetServiceParams(ctx sdk.Context, params types.ServiceParams) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ServiceParamsKey, bz)
}

// RegisterServiceProvider registers a new service provider
func (k Keeper) RegisterServiceProvider(ctx sdk.Context, provider string, serviceType string, metadata string) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetServiceProviderKey(provider)
	
	// Check if provider is already registered
	if store.Has(key) {
		return fmt.Errorf("service provider already registered")
	}
	
	// Create and store service provider
	serviceProvider := types.ServiceProvider{
		Address:     provider,
		ServiceType: serviceType,
		Metadata:    metadata,
		RegisteredAt: ctx.BlockTime(),
		Active:      true,
	}
	
	bz := k.cdc.MustMarshal(&serviceProvider)
	store.Set(key, bz)
	
	// Initialize service score
	serviceScore := types.ServiceScore{
		Provider:    provider,
		Score:       sdk.ZeroInt(),
		LastUpdated: 0,
	}
	
	scoreKey := types.GetServiceScoreKey(provider)
	scoreBz := k.cdc.MustMarshal(&serviceScore)
	store.Set(scoreKey, scoreBz)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeServiceProviderRegistered,
			sdk.NewAttribute(types.AttributeKeyProvider, provider),
			sdk.NewAttribute(types.AttributeKeyServiceType, serviceType),
		),
	)
	
	// Call hooks if set
	if k.hooks != nil {
		k.hooks.AfterServiceProviderRegistered(ctx, provider)
	}
	
	return nil
}

// GetServiceProvider returns a service provider by address
func (k Keeper) GetServiceProvider(ctx sdk.Context, provider string) (types.ServiceProvider, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetServiceProviderKey(provider)
	
	bz := store.Get(key)
	if bz == nil {
		return types.ServiceProvider{}, false
	}
	
	var serviceProvider types.ServiceProvider
	k.cdc.MustUnmarshal(bz, &serviceProvider)
	return serviceProvider, true
}

// SubmitProof submits a new proof of service
func (k Keeper) SubmitProof(ctx sdk.Context, provider string, serviceType string, proofID string, evidence string) error {
	store := ctx.KVStore(k.storeKey)
	
	// Check if provider is registered
	providerKey := types.GetServiceProviderKey(provider)
	if !store.Has(providerKey) {
		return fmt.Errorf("service provider not registered")
	}
	
	// Check if proof already exists
	proofKey := types.GetServiceProofKey(provider, proofID)
	if store.Has(proofKey) {
		return fmt.Errorf("proof already submitted")
	}
	
	// Create and store service proof
	serviceProof := types.ServiceProof{
		ProofID:     proofID,
		Provider:    provider,
		ServiceType: serviceType,
		Evidence:    evidence,
		Timestamp:   ctx.BlockTime(),
		Verified:    false,
		VerifiedBy:  []string{},
		Score:       sdk.ZeroInt(),
	}
	
	bz := k.cdc.MustMarshal(&serviceProof)
	store.Set(proofKey, bz)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProofSubmitted,
			sdk.NewAttribute(types.AttributeKeyProvider, provider),
			sdk.NewAttribute(types.AttributeKeyProofID, proofID),
			sdk.NewAttribute(types.AttributeKeyServiceType, serviceType),
		),
	)
	
	// Call hooks if set
	if k.hooks != nil {
		k.hooks.AfterProofSubmitted(ctx, provider, proofID)
	}
	
	return nil
}

// GetProof returns a proof by provider and proofID
func (k Keeper) GetProof(ctx sdk.Context, provider string, proofID string) (types.ServiceProof, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetServiceProofKey(provider, proofID)
	
	bz := store.Get(key)
	if bz == nil {
		return types.ServiceProof{}, false
	}
	
	var proof types.ServiceProof
	k.cdc.MustUnmarshal(bz, &proof)
	return proof, true
}

// VerifyProof verifies a proof of service
func (k Keeper) VerifyProof(ctx sdk.Context, validator string, provider string, proofID string, isVerified bool, score uint64) error {
	store := ctx.KVStore(k.storeKey)
	
	// Check if validator is a valid validator
	if !k.stakingKeeper.IsValidator(ctx, sdk.MustAccAddressFromBech32(validator)) {
		return fmt.Errorf("address is not a validator")
	}
	
	// Get the proof
	proofKey := types.GetServiceProofKey(provider, proofID)
	bz := store.Get(proofKey)
	if bz == nil {
		return fmt.Errorf("proof not found")
	}
	
	var proof types.ServiceProof
	k.cdc.MustUnmarshal(bz, &proof)
	
	// Check if validator has already verified this proof
	for _, v := range proof.VerifiedBy {
		if v == validator {
			return fmt.Errorf("validator has already verified this proof")
		}
	}
	
	// Add validator to verified by list
	proof.VerifiedBy = append(proof.VerifiedBy, validator)
	
	// Check if proof is now verified (minimum verifications reached)
	params := k.GetServiceParams(ctx)
	if uint32(len(proof.VerifiedBy)) >= params.MinVerifications && isVerified {
		proof.Verified = true
		proof.Score = sdk.NewIntFromUint64(score)
		
		// Update service score
		k.updateServiceScore(ctx, provider, sdk.NewIntFromUint64(score))
	}
	
	// Update proof
	newBz := k.cdc.MustMarshal(&proof)
	store.Set(proofKey, newBz)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProofVerified,
			sdk.NewAttribute(types.AttributeKeyValidator, validator),
			sdk.NewAttribute(types.AttributeKeyProvider, provider),
			sdk.NewAttribute(types.AttributeKeyProofID, proofID),
			sdk.NewAttribute(types.AttributeKeyVerified, fmt.Sprintf("%t", isVerified)),
			sdk.NewAttribute(types.AttributeKeyScore, fmt.Sprintf("%d", score)),
		),
	)
	
	// Call hooks if set
	if k.hooks != nil && proof.Verified {
		k.hooks.AfterProofVerified(ctx, provider, proofID, proof.Score)
	}
	
	return nil
}

// updateServiceScore updates the service score for a provider
func (k Keeper) updateServiceScore(ctx sdk.Context, provider string, additionalScore sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	scoreKey := types.GetServiceScoreKey(provider)
	
	var serviceScore types.ServiceScore
	bz := store.Get(scoreKey)
	if bz == nil {
		serviceScore = types.ServiceScore{
			Provider:    provider,
			Score:       sdk.ZeroInt(),
			LastUpdated: 0,
		}
	} else {
		k.cdc.MustUnmarshal(bz, &serviceScore)
	}
	
	// Update score
	serviceScore.Score = serviceScore.Score.Add(additionalScore)
	serviceScore.LastUpdated = uint64(ctx.BlockHeight())
	
	// Store updated score
	newBz := k.cdc.MustMarshal(&serviceScore)
	store.Set(scoreKey, newBz)
	
	// Update total service score
	k.updateTotalServiceScore(ctx, additionalScore)
}

// updateTotalServiceScore updates the total service score
func (k Keeper) updateTotalServiceScore(ctx sdk.Context, additionalScore sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	
	var totalScore sdk.Int
	bz := store.Get(types.TotalServiceScoreKey)
	if bz == nil {
		totalScore = sdk.ZeroInt()
	} else {
		k.cdc.MustUnmarshal(bz, &totalScore)
	}
	
	// Update total score
	totalScore = totalScore.Add(additionalScore)
	
	// Store updated total score
	newBz := k.cdc.MustMarshal(&totalScore)
	store.Set(types.TotalServiceScoreKey, newBz)
}

// GetServiceScore returns the service score for a provider
func (k Keeper) GetServiceScore(ctx sdk.Context, provider string) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	scoreKey := types.GetServiceScoreKey(provider)
	
	bz := store.Get(scoreKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	
	var serviceScore types.ServiceScore
	k.cdc.MustUnmarshal(bz, &serviceScore)
	return serviceScore.Score
}

// GetTotalServiceScore returns the total service score
func (k Keeper) GetTotalServiceScore(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	
	bz := store.Get(types.TotalServiceScoreKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	
	var totalScore sdk.Int
	k.cdc.MustUnmarshal(bz, &totalScore)
	return totalScore
}

// DecayServiceScores decays all service scores based on the decay rate
func (k Keeper) DecayServiceScores(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	params := k.GetServiceParams(ctx)
	
	// Iterate through all service scores
	iterator := sdk.KVStorePrefixIterator(store, types.ServiceScorePrefix)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var serviceScore types.ServiceScore
		k.cdc.MustUnmarshal(iterator.Value(), &serviceScore)
		
		// Calculate decayed score
		decayedScore := sdk.NewDecFromInt(serviceScore.Score).
			Mul(sdk.OneDec().Sub(params.ScoreDecayRate)).
			TruncateInt()
		
		// Update score
		serviceScore.Score = decayedScore
		serviceScore.LastUpdated = uint64(ctx.BlockHeight())
		
		// Store updated score
		newBz := k.cdc.MustMarshal(&serviceScore)
		store.Set(iterator.Key(), newBz)
	}
	
	// Recalculate total service score
	k.recalculateTotalServiceScore(ctx)
}

// recalculateTotalServiceScore recalculates the total service score
func (k Keeper) recalculateTotalServiceScore(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	
	totalScore := sdk.ZeroInt()
	
	// Iterate through all service scores
	iterator := sdk.KVStorePrefixIterator(store, types.ServiceScorePrefix)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var serviceScore types.ServiceScore
		k.cdc.MustUnmarshal(iterator.Value(), &serviceScore)
		
		// Add to total score
		totalScore = totalScore.Add(serviceScore.Score)
	}
	
	// Store updated total score
	newBz := k.cdc.MustMarshal(&totalScore)
	store.Set(types.TotalServiceScoreKey, newBz)
}
