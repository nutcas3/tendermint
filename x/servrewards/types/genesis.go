package types

import (
	"fmt"
)

// DefaultGenesis returns the default genesis state for the servrewards module.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		RewardMetrics:      DefaultRewardMetrics(),
		RewardParams:       DefaultRewardParams(),
		AccumulatedRewards: []AccumulatedRewards{},
	}
}

// GenesisState defines the servrewards module's genesis state.
type GenesisState struct {
	RewardMetrics      RewardMetrics       `json:"reward_metrics"`
	RewardParams       RewardParams        `json:"reward_params"`
	AccumulatedRewards []AccumulatedRewards `json:"accumulated_rewards"`
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	// Validate reward parameters
	if gs.RewardParams.ServiceScoreWeight.IsNegative() {
		return fmt.Errorf("service score weight cannot be negative: %s", gs.RewardParams.ServiceScoreWeight)
	}
	
	if gs.RewardParams.StakingWeight.IsNegative() {
		return fmt.Errorf("staking weight cannot be negative: %s", gs.RewardParams.StakingWeight)
	}
	
	if gs.RewardParams.RewardPerEpoch.IsNegative() {
		return fmt.Errorf("reward per epoch cannot be negative: %s", gs.RewardParams.RewardPerEpoch)
	}
	
	if gs.RewardParams.EpochDuration == 0 {
		return fmt.Errorf("epoch duration must be positive")
	}
	
	// Ensure weights sum to 1
	sumWeights := gs.RewardParams.ServiceScoreWeight.Add(gs.RewardParams.StakingWeight)
	if !sumWeights.Equal(OneDec()) {
		return fmt.Errorf("service score weight and staking weight must sum to 1, got: %s", sumWeights)
	}
	
	// Validate accumulated rewards
	for _, reward := range gs.AccumulatedRewards {
		if reward.Rewards.IsNegative() {
			return fmt.Errorf("accumulated rewards cannot be negative: %s", reward.Rewards)
		}
	}
	
	return nil
}
