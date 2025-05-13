package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// ParameterChangeProposal defines a proposal to change parameters
// in the system.
type ParameterChangeProposal struct {
	Title       string
	Description string
	Changes     []ParamChange
}

// ParamChange defines a single parameter change
// with a subspace, key, and value.
type ParamChange struct {
	Subspace string
	Key      string
	Value    string
}

// NewParameterChangeProposal creates a new parameter change proposal
func NewParameterChangeProposal(title, description string, changes []ParamChange) ParameterChangeProposal {
	return ParameterChangeProposal{
		Title:       title,
		Description: description,
		Changes:     changes,
	}
}

// GetTitle returns the title of the proposal
func (p ParameterChangeProposal) GetTitle() string {
	return p.Title
}

// GetDescription returns the description of the proposal
func (p ParameterChangeProposal) GetDescription() string {
	return p.Description
}

// ProposalRoute returns the routing key of the proposal
func (p ParameterChangeProposal) ProposalRoute() string {
	return govtypes.RouterKey
}

// ProposalType returns the type of the proposal
func (p ParameterChangeProposal) ProposalType() string {
	return "ParameterChange"
}

// ValidateBasic performs basic validation on the proposal
func (p ParameterChangeProposal) ValidateBasic() error {
	if len(p.Title) == 0 {
		return fmt.Errorf("proposal title cannot be empty")
	}
	if len(p.Description) == 0 {
		return fmt.Errorf("proposal description cannot be empty")
	}
	if len(p.Changes) == 0 {
		return fmt.Errorf("proposal must contain at least one change")
	}
	return nil
}
