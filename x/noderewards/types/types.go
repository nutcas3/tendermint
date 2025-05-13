package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NodePerformance represents performance metrics for a validator node
type NodePerformance struct {
	ValidatorAddr  string  `json:"validator_addr"`
	ServiceScore   sdk.Int `json:"service_score"`
	UptimePercent  sdk.Dec `json:"uptime_percent"`
	ResponseTime   sdk.Int `json:"response_time"` // In milliseconds
	LastUpdateHeight int64   `json:"last_update_height"`
}

// RewardModifier represents parameters for modifying staking rewards
type RewardModifier struct {
	ServiceScoreWeight sdk.Dec `json:"service_score_weight"`
	UptimeWeight       sdk.Dec `json:"uptime_weight"`
	ResponseTimeWeight sdk.Dec `json:"response_time_weight"`
	MinModifier        sdk.Dec `json:"min_modifier"` // Minimum reward modifier (e.g., 0.5 = 50% of base rewards)
	MaxModifier        sdk.Dec `json:"max_modifier"` // Maximum reward modifier (e.g., 2.0 = 200% of base rewards)
}

// DefaultRewardModifier returns default parameters for reward modification
func DefaultRewardModifier() RewardModifier {
	return RewardModifier{
		ServiceScoreWeight: sdk.NewDecWithPrec(5, 1), // 0.5
		UptimeWeight:       sdk.NewDecWithPrec(3, 1), // 0.3
		ResponseTimeWeight: sdk.NewDecWithPrec(2, 1), // 0.2
		MinModifier:        sdk.NewDecWithPrec(5, 1), // 0.5 (50% of base rewards)
		MaxModifier:        sdk.NewDec(2),            // 2.0 (200% of base rewards)
	}
}

// DefaultNodePerformance returns default performance metrics for a validator node
func DefaultNodePerformance(validatorAddr string) NodePerformance {
	return NodePerformance{
		ValidatorAddr:    validatorAddr,
		ServiceScore:     sdk.ZeroInt(),
		UptimePercent:    sdk.OneDec(), // 100% uptime initially
		ResponseTime:     sdk.NewInt(100), // 100ms default response time
		LastUpdateHeight: 0,
	}
}
