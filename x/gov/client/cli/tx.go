package cli

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/serv-chain/serv/x/gov/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Governance transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govTxCmd.AddCommand(
		GetCmdSubmitParameterChangeProposal(),
	)

	return govTxCmd
}

// GetCmdSubmitParameterChangeProposal implements the command to submit a parameter change proposal
func GetCmdSubmitParameterChangeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-param-change [title] [description] [changes]",
		Short: "Submit a parameter change proposal",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title := args[0]
			description := args[1]
			changes := args[2] // This should be a JSON string representing the changes

			proposal := types.NewParameterChangeProposal(title, description, nil) // Parse changes into ParamChange slice

			if err := proposal.ValidateBasic(); err != nil {
				return err
			}

			msg := types.NewMsgSubmitProposal(proposal, clientCtx.GetFromAddress())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
