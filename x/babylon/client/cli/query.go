package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the Babylon module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		SilenceUsage:               true,
	}
	queryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryBSNContracts(),
	)
	return queryCmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current babylon parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as babylon parameters.

Example:
$ %s query babylon params
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryBSNContracts implements the contracts query command.
func GetCmdQueryBSNContracts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bsn-contracts",
		Args:  cobra.NoArgs,
		Short: "Query the contract addresses for the Babylon module",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query contract addresses set in the Babylon module.

Example:
$ %s query babylon contracts
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BSNContracts(cmd.Context(), &types.QueryBSNContractsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.BsnContracts)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
