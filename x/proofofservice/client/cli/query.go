package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/serv-chain/serv/x/proofofservice/types"
)

// GetQueryCmd returns the query commands for the proofofservice module
func GetQueryCmd(queryRoute string) *cobra.Command {
	proofOfServiceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	proofOfServiceQueryCmd.AddCommand(
		GetCmdQueryServiceParams(),
		GetCmdQueryServiceProvider(),
		GetCmdQueryServiceProviders(),
		GetCmdQueryProof(),
		GetCmdQueryServiceScore(),
		GetCmdQueryTotalServiceScore(),
	)

	return proofOfServiceQueryCmd
}

// GetCmdQueryServiceParams implements the query service parameters command handler
func GetCmdQueryServiceParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current service parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ServiceParams(cmd.Context(), &types.QueryServiceParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryServiceProvider implements the query service provider command handler
func GetCmdQueryServiceProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider [address]",
		Short: "Query a service provider by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ServiceProvider(cmd.Context(), &types.QueryServiceProviderRequest{
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

// GetCmdQueryServiceProviders implements the query service providers command handler
func GetCmdQueryServiceProviders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers [service-type]",
		Short: "Query all service providers, optionally filtered by service type",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			
			var serviceType string
			if len(args) > 0 {
				serviceType = args[0]
			}
			
			res, err := queryClient.ServiceProviders(cmd.Context(), &types.QueryServiceProvidersRequest{
				ServiceType: serviceType,
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

// GetCmdQueryProof implements the query proof command handler
func GetCmdQueryProof() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proof [provider-address] [proof-id]",
		Short: "Query a proof of service",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Proof(cmd.Context(), &types.QueryProofRequest{
				Provider: args[0],
				ProofId:  args[1],
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

// GetCmdQueryServiceScore implements the query service score command handler
func GetCmdQueryServiceScore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "score [address]",
		Short: "Query the service score for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ServiceScore(cmd.Context(), &types.QueryServiceScoreRequest{
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

// GetCmdQueryTotalServiceScore implements the query total service score command handler
func GetCmdQueryTotalServiceScore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-score",
		Short: "Query the total service score",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TotalServiceScore(cmd.Context(), &types.QueryTotalServiceScoreRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
