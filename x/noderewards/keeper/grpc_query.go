package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/noderewards/types"
)

// Querier is used for implementing the Query gRPC service
type Querier struct {
	Keeper
}

// NewQueryServer creates a new gRPC query server for the noderewards module
func NewQueryServer(k Keeper) types.QueryServer {
	return &Querier{Keeper: k}
}

// RewardModifier implements the Query/RewardModifier gRPC method
func (q Querier) RewardModifier(c context.Context, req *types.QueryRewardModifierRequest) (*types.QueryRewardModifierResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	modifier := q.Keeper.GetRewardModifier(ctx)

	return &types.QueryRewardModifierResponse{
		Modifier: &modifier,
	}, nil
}

// NodePerformance implements the Query/NodePerformance gRPC method
func (q Querier) NodePerformance(c context.Context, req *types.QueryNodePerformanceRequest) (*types.QueryNodePerformanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	performance := q.Keeper.GetNodePerformance(ctx, req.ValidatorAddr)

	return &types.QueryNodePerformanceResponse{
		Performance: &performance,
	}, nil
}

// RewardModifierForValidator implements the Query/RewardModifierForValidator gRPC method
func (q Querier) RewardModifierForValidator(c context.Context, req *types.QueryRewardModifierForValidatorRequest) (*types.QueryRewardModifierForValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	modifier := q.Keeper.CalculateRewardModifier(ctx, req.ValidatorAddr)

	return &types.QueryRewardModifierForValidatorResponse{
		Modifier: modifier,
	}, nil
}
