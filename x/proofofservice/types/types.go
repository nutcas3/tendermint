package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ServiceProvider represents a registered service provider
type ServiceProvider struct {
	Address     string    `json:"address"`
	ServiceType string    `json:"service_type"`
	Metadata    string    `json:"metadata"`
	RegisteredAt time.Time `json:"registered_at"`
	Active      bool      `json:"active"`
}

// ServiceProof represents a proof of service submission
type ServiceProof struct {
	ProofID     string    `json:"proof_id"`
	Provider    string    `json:"provider"`
	ServiceType string    `json:"service_type"`
	Evidence    string    `json:"evidence"` // Could be a hash of evidence data
	Timestamp   time.Time `json:"timestamp"`
	Verified    bool      `json:"verified"`
	VerifiedBy  []string  `json:"verified_by"` // List of validators who verified this proof
	Score       sdk.Int   `json:"score"`       // Score assigned to this proof
}

// ServiceScore represents the accumulated service score for a provider
type ServiceScore struct {
	Provider string  `json:"provider"`
	Score    sdk.Int `json:"score"`
	LastUpdated uint64 `json:"last_updated"` // Last epoch when score was updated
}

// ServiceParams represents the parameters for service validation
type ServiceParams struct {
	MinVerifications uint32  `json:"min_verifications"` // Minimum number of verifications required
	ScoreDecayRate   sdk.Dec `json:"score_decay_rate"`  // Rate at which scores decay over time
	ProofValidityPeriod uint64 `json:"proof_validity_period"` // Number of blocks a proof is valid for
	MaxProofsPerEpoch uint32 `json:"max_proofs_per_epoch"` // Maximum number of proofs a provider can submit per epoch
}

// DefaultServiceParams returns default parameters for service validation
func DefaultServiceParams() ServiceParams {
	return ServiceParams{
		MinVerifications:    3,
		ScoreDecayRate:      sdk.NewDecWithPrec(1, 1), // 0.1 (10% decay)
		ProofValidityPeriod: 100,                      // 100 blocks
		MaxProofsPerEpoch:   5,
	}
}
