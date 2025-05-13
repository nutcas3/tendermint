package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

// RewardMetrics represents the metrics used to calculate rewards
type RewardMetrics struct {
	TotalServiceScore sdk.Int `json:"total_service_score"`
	TotalStaked       sdk.Int `json:"total_staked"`
	EpochNumber       uint64  `json:"epoch_number"`
}

// RewardParams represents the parameters for reward calculation
type RewardParams struct {
	ServiceScoreWeight sdk.Dec `json:"service_score_weight"`
	StakingWeight      sdk.Dec `json:"staking_weight"`
	RewardPerEpoch     sdk.Int `json:"reward_per_epoch"`
	EpochDuration      uint64  `json:"epoch_duration"` // In blocks
}

// AccumulatedRewards represents the rewards accumulated for an address
type AccumulatedRewards struct {
	Address  string  `json:"address"`
	Rewards  sdk.Int `json:"rewards"`
	LastClaim uint64  `json:"last_claim"` // Last epoch when rewards were claimed
}

// DefaultRewardParams returns default parameters for reward calculation
func DefaultRewardParams() RewardParams {
	return RewardParams{
		ServiceScoreWeight: sdk.NewDecWithPrec(6, 1), // 0.6
		StakingWeight:      sdk.NewDecWithPrec(4, 1), // 0.4
		RewardPerEpoch:     sdk.NewInt(1000000),      // 1 SERV (assuming 6 decimals)
		EpochDuration:      100,                      // 100 blocks per epoch
	}
}

// DefaultRewardMetrics returns default metrics for reward calculation
func DefaultRewardMetrics() RewardMetrics {
	return RewardMetrics{
		TotalServiceScore: sdk.ZeroInt(),
		TotalStaked:       sdk.ZeroInt(),
		EpochNumber:       0,
	}
}
