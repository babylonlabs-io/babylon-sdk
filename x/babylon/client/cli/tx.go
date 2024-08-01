package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

const flagAuthority = "authority"

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Babylon transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		SilenceUsage:               true,
	}
	txCmd.AddCommand()
	return txCmd
}
