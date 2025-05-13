package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRegisterService = "register_service"
	TypeMsgSubmitProof     = "submit_proof"
	TypeMsgVerifyProof     = "verify_proof"
)

var _ sdk.Msg = &MsgRegisterService{}
var _ sdk.Msg = &MsgSubmitProof{}
var _ sdk.Msg = &MsgVerifyProof{}

// MsgRegisterService defines a message for registering a service provider
type MsgRegisterService struct {
	Provider    string `json:"provider"`
	ServiceType string `json:"service_type"`
	Metadata    string `json:"metadata"`
}

// NewMsgRegisterService creates a new MsgRegisterService instance
func NewMsgRegisterService(provider, serviceType, metadata string) *MsgRegisterService {
	return &MsgRegisterService{
		Provider:    provider,
		ServiceType: serviceType,
		Metadata:    metadata,
	}
}

// Route implements sdk.Msg
func (msg MsgRegisterService) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgRegisterService) Type() string {
	return TypeMsgRegisterService
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterService) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address: %s", err)
	}

	if msg.ServiceType == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service type cannot be empty")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgRegisterService) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Provider)
	return []sdk.AccAddress{addr}
}

// MsgSubmitProof defines a message for submitting proof of service
type MsgSubmitProof struct {
	Provider    string `json:"provider"`
	ServiceType string `json:"service_type"`
	ProofID     string `json:"proof_id"`
	Evidence    string `json:"evidence"`
}

// NewMsgSubmitProof creates a new MsgSubmitProof instance
func NewMsgSubmitProof(provider, serviceType, proofID, evidence string) *MsgSubmitProof {
	return &MsgSubmitProof{
		Provider:    provider,
		ServiceType: serviceType,
		ProofID:     proofID,
		Evidence:    evidence,
	}
}

// Route implements sdk.Msg
func (msg MsgSubmitProof) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgSubmitProof) Type() string {
	return TypeMsgSubmitProof
}

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitProof) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address: %s", err)
	}

	if msg.ServiceType == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service type cannot be empty")
	}

	if msg.ProofID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "proof ID cannot be empty")
	}

	if msg.Evidence == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "evidence cannot be empty")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgSubmitProof) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements sdk.Msg
func (msg MsgSubmitProof) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Provider)
	return []sdk.AccAddress{addr}
}

// MsgVerifyProof defines a message for verifying a proof of service
type MsgVerifyProof struct {
	Validator  string `json:"validator"`
	Provider   string `json:"provider"`
	ProofID    string `json:"proof_id"`
	IsVerified bool   `json:"is_verified"`
	Score      uint64 `json:"score"` // Score assigned to this proof (0-100)
}

// NewMsgVerifyProof creates a new MsgVerifyProof instance
func NewMsgVerifyProof(validator, provider, proofID string, isVerified bool, score uint64) *MsgVerifyProof {
	return &MsgVerifyProof{
		Validator:  validator,
		Provider:   provider,
		ProofID:    proofID,
		IsVerified: isVerified,
		Score:      score,
	}
}

// Route implements sdk.Msg
func (msg MsgVerifyProof) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg MsgVerifyProof) Type() string {
	return TypeMsgVerifyProof
}

// ValidateBasic implements sdk.Msg
func (msg MsgVerifyProof) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Validator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address: %s", err)
	}

	if msg.ProofID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "proof ID cannot be empty")
	}

	if msg.Score > 100 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "score must be between 0 and 100")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgVerifyProof) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements sdk.Msg
func (msg MsgVerifyProof) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Validator)
	return []sdk.AccAddress{addr}
}
