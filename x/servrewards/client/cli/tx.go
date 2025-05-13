package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/serv-chain/serv/x/servrewards/types"
)

// GetTxCmd returns the transaction commands for the servrewards module
func GetTxCmd() *cobra.Command {
	servRewardsTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	servRewardsTxCmd.AddCommand(
		NewClaimRewardCmd(),
		NewUpdateRewardParamsCmd(),
	)

	return servRewardsTxCmd
}

// NewClaimRewardCmd implements the claim reward command handler
func NewClaimRewardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-reward",
		Short: "Claim accumulated SERV rewards",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimReward(clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUpdateRewardParamsCmd implements the update reward parameters command handler
// This would typically be a governance proposal, but we provide a direct command for testing
func NewUpdateRewardParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [service-score-weight] [staking-weight] [reward-per-epoch] [epoch-duration]",
		Short: "Update SERV reward parameters (governance)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serviceScoreWeight, err := sdk.NewDecFromStr(args[0])
			if err != nil {
				return fmt.Errorf("invalid service score weight: %w", err)
			}

			stakingWeight, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return fmt.Errorf("invalid staking weight: %w", err)
			}

			rewardPerEpoch, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid reward per epoch: %s", args[2])
			}

			epochDuration, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid epoch duration: %w", err)
			}

			msg := types.NewMsgUpdateRewardParams(
				clientCtx.GetFromAddress().String(),
				serviceScoreWeight,
				stakingWeight,
				rewardPerEpoch,
				epochDuration,
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
