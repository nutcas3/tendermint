package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/serv-chain/serv/x/servrewards/types"
)

// GetQueryCmd returns the query commands for the servrewards module
func GetQueryCmd(queryRoute string) *cobra.Command {
	servRewardsQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	servRewardsQueryCmd.AddCommand(
		GetCmdQueryRewardMetrics(),
		GetCmdQueryRewardParams(),
		GetCmdQueryAccumulatedRewards(),
	)

	return servRewardsQueryCmd
}

// GetCmdQueryRewardMetrics implements the query reward metrics command handler
func GetCmdQueryRewardMetrics() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Query the current SERV reward metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RewardMetrics(cmd.Context(), &types.QueryRewardMetricsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryRewardParams implements the query reward parameters command handler
func GetCmdQueryRewardParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current SERV reward parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RewardParams(cmd.Context(), &types.QueryRewardParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryAccumulatedRewards implements the query accumulated rewards command handler
func GetCmdQueryAccumulatedRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards [address]",
		Short: "Query accumulated SERV rewards for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AccumulatedRewards(cmd.Context(), &types.QueryAccumulatedRewardsRequest{
				Address: args[0],
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
