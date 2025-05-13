package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/proofofservice/types"
)

// Querier is used for implementing the Query gRPC service
type Querier struct {
	Keeper
}

// NewQueryServer creates a new gRPC query server for the proofofservice module
func NewQueryServer(k Keeper) types.QueryServer {
	return &Querier{Keeper: k}
}

// ServiceParams implements the Query/ServiceParams gRPC method
func (q Querier) ServiceParams(c context.Context, req *types.QueryServiceParamsRequest) (*types.QueryServiceParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := q.Keeper.GetServiceParams(ctx)

	return &types.QueryServiceParamsResponse{
		Params: &params,
	}, nil
}

// ServiceProvider implements the Query/ServiceProvider gRPC method
func (q Querier) ServiceProvider(c context.Context, req *types.QueryServiceProviderRequest) (*types.QueryServiceProviderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	provider, found := q.Keeper.GetServiceProvider(ctx, req.Address)
	if !found {
		return nil, status.Error(codes.NotFound, "service provider not found")
	}

	return &types.QueryServiceProviderResponse{
		Provider: &provider,
	}, nil
}

// ServiceProviders implements the Query/ServiceProviders gRPC method
func (q Querier) ServiceProviders(c context.Context, req *types.QueryServiceProvidersRequest) (*types.QueryServiceProvidersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	// In a real implementation, we would iterate through all providers in the store
	// and filter by service type if provided
	// For now, we'll return an empty list
	providers := []*types.ServiceProvider{}

	return &types.QueryServiceProvidersResponse{
		Providers: providers,
	}, nil
}

// Proof implements the Query/Proof gRPC method
func (q Querier) Proof(c context.Context, req *types.QueryProofRequest) (*types.QueryProofResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Provider == "" {
		return nil, status.Error(codes.InvalidArgument, "provider address cannot be empty")
	}

	if req.ProofId == "" {
		return nil, status.Error(codes.InvalidArgument, "proof ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	proof, found := q.Keeper.GetProof(ctx, req.Provider, req.ProofId)
	if !found {
		return nil, status.Error(codes.NotFound, "proof not found")
	}

	return &types.QueryProofResponse{
		Proof: &proof,
	}, nil
}

// ServiceScore implements the Query/ServiceScore gRPC method
func (q Querier) ServiceScore(c context.Context, req *types.QueryServiceScoreRequest) (*types.QueryServiceScoreResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	score := q.Keeper.GetServiceScore(ctx, req.Address)

	return &types.QueryServiceScoreResponse{
		Score: score,
	}, nil
}

// TotalServiceScore implements the Query/TotalServiceScore gRPC method
func (q Querier) TotalServiceScore(c context.Context, req *types.QueryTotalServiceScoreRequest) (*types.QueryTotalServiceScoreResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	totalScore := q.Keeper.GetTotalServiceScore(ctx)

	return &types.QueryTotalServiceScoreResponse{
		TotalScore: totalScore,
	}, nil
}
