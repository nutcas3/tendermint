package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state for the proofofservice module.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ServiceParams:     DefaultServiceParams(),
		ServiceProviders:  []ServiceProvider{},
		ServiceProofs:     []ServiceProof{},
		TotalServiceScore: sdk.ZeroInt(),
	}
}

// GenesisState defines the proofofservice module's genesis state.
type GenesisState struct {
	ServiceParams     ServiceParams     `json:"service_params"`
	ServiceProviders  []ServiceProvider `json:"service_providers"`
	ServiceProofs     []ServiceProof    `json:"service_proofs"`
	TotalServiceScore sdk.Int           `json:"total_service_score"`
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	// Validate service parameters
	if gs.ServiceParams.MinVerifications == 0 {
		return fmt.Errorf("minimum verifications must be positive")
	}
	
	if gs.ServiceParams.ScoreDecayRate.IsNegative() || gs.ServiceParams.ScoreDecayRate.GT(sdk.OneDec()) {
		return fmt.Errorf("score decay rate must be between 0 and 1: %s", gs.ServiceParams.ScoreDecayRate)
	}
	
	if gs.ServiceParams.ProofValidityPeriod == 0 {
		return fmt.Errorf("proof validity period must be positive")
	}
	
	if gs.ServiceParams.MaxProofsPerEpoch == 0 {
		return fmt.Errorf("max proofs per epoch must be positive")
	}
	
	// Validate service providers
	providerAddresses := make(map[string]bool)
	for _, provider := range gs.ServiceProviders {
		if _, exists := providerAddresses[provider.Address]; exists {
			return fmt.Errorf("duplicate service provider address: %s", provider.Address)
		}
		providerAddresses[provider.Address] = true
		
		if provider.ServiceType == "" {
			return fmt.Errorf("service type cannot be empty for provider: %s", provider.Address)
		}
	}
	
	// Validate service proofs
	proofIDs := make(map[string]bool)
	for _, proof := range gs.ServiceProofs {
		proofKey := proof.Provider + proof.ProofID
		if _, exists := proofIDs[proofKey]; exists {
			return fmt.Errorf("duplicate proof ID for provider: %s, proofID: %s", proof.Provider, proof.ProofID)
		}
		proofIDs[proofKey] = true
		
		if proof.ServiceType == "" {
			return fmt.Errorf("service type cannot be empty for proof: %s", proof.ProofID)
		}
		
		if proof.Evidence == "" {
			return fmt.Errorf("evidence cannot be empty for proof: %s", proof.ProofID)
		}
		
		if proof.Score.IsNegative() {
			return fmt.Errorf("proof score cannot be negative: %s", proof.Score)
		}
	}
	
	// Validate total service score
	if gs.TotalServiceScore.IsNegative() {
		return fmt.Errorf("total service score cannot be negative: %s", gs.TotalServiceScore)
	}
	
	return nil
}
