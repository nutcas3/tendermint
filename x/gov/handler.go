package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/gov/types"
)

// NewHandler returns a handler for "gov" type messages.
func NewHandler() sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.ParameterChangeProposal:
			return handleParameterChangeProposal(ctx, msg)
		default:
			return nil, fmt.Errorf("unrecognized gov message type: %T", msg)
		}
	}
}

// handleParameterChangeProposal handles parameter change proposals.
func handleParameterChangeProposal(ctx sdk.Context, proposal types.ParameterChangeProposal) (*sdk.Result, error) {
	// Implement logic to change parameters based on the proposal
	for _, change := range proposal.Changes {
		// Apply each parameter change
		fmt.Printf("Changing parameter %s in subspace %s to %s\n", change.Key, change.Subspace, change.Value)
		// Here you would update the actual parameter in the store
	}

	return &sdk.Result{}, nil
}
