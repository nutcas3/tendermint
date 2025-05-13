package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state for the noderewards module.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		RewardModifier:   DefaultRewardModifier(),
		NodePerformances: []NodePerformance{},
	}
}

// GenesisState defines the noderewards module's genesis state.
type GenesisState struct {
	RewardModifier   RewardModifier    `json:"reward_modifier"`
	NodePerformances []NodePerformance `json:"node_performances"`
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	// Validate reward modifier parameters
	if gs.RewardModifier.ServiceScoreWeight.IsNegative() {
		return fmt.Errorf("service score weight cannot be negative: %s", gs.RewardModifier.ServiceScoreWeight)
	}
	
	if gs.RewardModifier.UptimeWeight.IsNegative() {
		return fmt.Errorf("uptime weight cannot be negative: %s", gs.RewardModifier.UptimeWeight)
	}
	
	if gs.RewardModifier.ResponseTimeWeight.IsNegative() {
		return fmt.Errorf("response time weight cannot be negative: %s", gs.RewardModifier.ResponseTimeWeight)
	}
	
	// Ensure weights sum to 1
	sumWeights := gs.RewardModifier.ServiceScoreWeight.Add(gs.RewardModifier.UptimeWeight).Add(gs.RewardModifier.ResponseTimeWeight)
	if !sumWeights.Equal(sdk.OneDec()) {
		return fmt.Errorf("weights must sum to 1, got: %s", sumWeights)
	}
	
	if gs.RewardModifier.MinModifier.IsNegative() {
		return fmt.Errorf("minimum modifier cannot be negative: %s", gs.RewardModifier.MinModifier)
	}
	
	if gs.RewardModifier.MaxModifier.IsNegative() {
		return fmt.Errorf("maximum modifier cannot be negative: %s", gs.RewardModifier.MaxModifier)
	}
	
	if gs.RewardModifier.MinModifier.GT(gs.RewardModifier.MaxModifier) {
		return fmt.Errorf("minimum modifier cannot be greater than maximum modifier: %s > %s", 
			gs.RewardModifier.MinModifier, gs.RewardModifier.MaxModifier)
	}
	
	// Validate node performances
	validatorAddresses := make(map[string]bool)
	for _, performance := range gs.NodePerformances {
		if _, exists := validatorAddresses[performance.ValidatorAddr]; exists {
			return fmt.Errorf("duplicate validator address: %s", performance.ValidatorAddr)
		}
		validatorAddresses[performance.ValidatorAddr] = true
		
		if performance.ServiceScore.IsNegative() {
			return fmt.Errorf("service score cannot be negative: %s", performance.ServiceScore)
		}
		
		if performance.UptimePercent.IsNegative() || performance.UptimePercent.GT(sdk.OneDec()) {
			return fmt.Errorf("uptime percent must be between 0 and 1: %s", performance.UptimePercent)
		}
		
		if performance.ResponseTime.IsNegative() {
			return fmt.Errorf("response time cannot be negative: %s", performance.ResponseTime)
		}
	}
	
	return nil
}
