package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/serv-chain/serv/x/noderewards/types"
)

// GetQueryCmd returns the query commands for the noderewards module
func GetQueryCmd(queryRoute string) *cobra.Command {
	nodeRewardsQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nodeRewardsQueryCmd.AddCommand(
		GetCmdQueryRewardModifier(),
		GetCmdQueryNodePerformance(),
		GetCmdQueryRewardModifierForValidator(),
	)

	return nodeRewardsQueryCmd
}

// GetCmdQueryRewardModifier implements the query reward modifier parameters command handler
func GetCmdQueryRewardModifier() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current reward modifier parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RewardModifier(cmd.Context(), &types.QueryRewardModifierRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryNodePerformance implements the query node performance command handler
func GetCmdQueryNodePerformance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "performance [validator-address]",
		Short: "Query the performance metrics for a validator node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.NodePerformance(cmd.Context(), &types.QueryNodePerformanceRequest{
				ValidatorAddr: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryRewardModifierForValidator implements the query reward modifier for a validator command handler
func GetCmdQueryRewardModifierForValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modifier [validator-address]",
		Short: "Query the calculated reward modifier for a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RewardModifierForValidator(cmd.Context(), &types.QueryRewardModifierForValidatorRequest{
				ValidatorAddr: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
