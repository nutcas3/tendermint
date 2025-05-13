package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/serv-chain/serv/x/proofofservice/types"
)

// GetTxCmd returns the transaction commands for the proofofservice module
func GetTxCmd() *cobra.Command {
	proofOfServiceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	proofOfServiceTxCmd.AddCommand(
		NewRegisterServiceCmd(),
		NewSubmitProofCmd(),
		NewVerifyProofCmd(),
	)

	return proofOfServiceTxCmd
}

// NewRegisterServiceCmd implements the register service command handler
func NewRegisterServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [service-type] [metadata]",
		Short: "Register as a service provider",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serviceType := args[0]
			metadata := args[1]

			msg := types.NewMsgRegisterService(
				clientCtx.GetFromAddress().String(),
				serviceType,
				metadata,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewSubmitProofCmd implements the submit proof command handler
func NewSubmitProofCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proof [service-type] [proof-id] [evidence]",
		Short: "Submit proof of service",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serviceType := args[0]
			proofID := args[1]
			evidence := args[2]

			msg := types.NewMsgSubmitProof(
				clientCtx.GetFromAddress().String(),
				serviceType,
				proofID,
				evidence,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewVerifyProofCmd implements the verify proof command handler
func NewVerifyProofCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-proof [provider-address] [proof-id] [is-verified] [score]",
		Short: "Verify a proof of service (validators only)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			providerAddr := args[0]
			proofID := args[1]
			
			isVerified, err := strconv.ParseBool(args[2])
			if err != nil {
				return fmt.Errorf("invalid is-verified flag: %w", err)
			}
			
			score, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid score: %w", err)
			}
			
			if score > 100 {
				return fmt.Errorf("score must be between 0 and 100")
			}

			msg := types.NewMsgVerifyProof(
				clientCtx.GetFromAddress().String(),
				providerAddr,
				proofID,
				isVerified,
				score,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
