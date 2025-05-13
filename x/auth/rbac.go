package auth

import (
	"errors"
)

// Role represents a user role in the system
// with associated permissions.
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleUser     Role = "user"
	RoleValidator Role = "validator"
)

// Permission represents an action that can be performed
// by a role.
type Permission string

const (
	PermSubmitProposal Permission = "submit_proposal"
	PermVoteProposal   Permission = "vote_proposal"
	PermManageParams   Permission = "manage_params"
)

// RolePermissions maps roles to their permissions.
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		PermSubmitProposal,
		PermVoteProposal,
		PermManageParams,
	},
	RoleUser: {
		PermSubmitProposal,
		PermVoteProposal,
	},
	RoleValidator: {
		PermVoteProposal,
	},
}

// CheckPermission checks if a role has the specified permission.
func CheckPermission(role Role, permission Permission) error {
	permissions, exists := RolePermissions[role]
	if !exists {
		return errors.New("role does not exist")
	}

	for _, perm := range permissions {
		if perm == permission {
			return nil
		}
	}

	return errors.New("permission denied")
}
