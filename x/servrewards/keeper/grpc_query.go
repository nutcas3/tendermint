package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/servrewards/types"
)

// Querier is used for implementing the Query gRPC service
type Querier struct {
	Keeper
}

// NewQueryServer creates a new gRPC query server for the servrewards module
func NewQueryServer(k Keeper) types.QueryServer {
	return &Querier{Keeper: k}
}

// RewardMetrics implements the Query/RewardMetrics gRPC method
func (q Querier) RewardMetrics(c context.Context, req *types.QueryRewardMetricsRequest) (*types.QueryRewardMetricsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	metrics := q.Keeper.GetRewardMetrics(ctx)

	return &types.QueryRewardMetricsResponse{
		Metrics: &metrics,
	}, nil
}

// RewardParams implements the Query/RewardParams gRPC method
func (q Querier) RewardParams(c context.Context, req *types.QueryRewardParamsRequest) (*types.QueryRewardParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := q.Keeper.GetRewardParams(ctx)

	return &types.QueryRewardParamsResponse{
		Params: &params,
	}, nil
}

// AccumulatedRewards implements the Query/AccumulatedRewards gRPC method
func (q Querier) AccumulatedRewards(c context.Context, req *types.QueryAccumulatedRewardsRequest) (*types.QueryAccumulatedRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	rewards := q.Keeper.GetAccumulatedRewards(ctx, req.Address)

	return &types.QueryAccumulatedRewardsResponse{
		Rewards: &rewards,
	}, nil
}
