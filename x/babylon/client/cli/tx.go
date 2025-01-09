package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

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
	txCmd.AddCommand(
		NewInstantiateBabylonContractsCmd(),
	)
	return txCmd
}

// [babylon-contract-code-id] [btc-staking-contract-code-id] [btc-finality-contract-code-id] [btc-network] [babylon-tag] [btc-confirmation-depth] [checkpoint-finalization-timeout] [notify-cosmos-zone] [consumer-name] [consumer-description] [admin]
func ParseInstantiateArgs(args []string, sender string) (*types.MsgInstantiateBabylonContracts, error) {
	// get the id of the code to instantiate
	babylonContractCodeID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, err
	}
	btcStakingContractCodeID, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return nil, err
	}
	btcFinalityContractCodeID, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return nil, err
	}

	btcNetwork := args[3]
	babylonTag := args[4]
	btcConfirmationDepth, err := strconv.ParseUint(args[5], 10, 32)
	if err != nil {
		return nil, err
	}
	checkpointFinalizationTimeout, err := strconv.ParseUint(args[6], 10, 32)
	if err != nil {
		return nil, err
	}
	notifyCosmosZone, err := strconv.ParseBool(args[7])
	if err != nil {
		return nil, err
	}
	consumerName := args[8]
	consumerDescription := args[9]
	adminStr := args[10]

	// build and sign the transaction, then broadcast to Tendermint
	msg := types.MsgInstantiateBabylonContracts{
		Signer:                        sender,
		BabylonContractCodeId:         babylonContractCodeID,
		BtcStakingContractCodeId:      btcStakingContractCodeID,
		BtcFinalityContractCodeId:     btcFinalityContractCodeID,
		Network:                       btcNetwork,
		BabylonTag:                    babylonTag,
		BtcConfirmationDepth:          uint32(btcConfirmationDepth),
		CheckpointFinalizationTimeout: uint32(checkpointFinalizationTimeout),
		NotifyCosmosZone:              notifyCosmosZone,
		ConsumerName:                  consumerName,
		ConsumerDescription:           consumerDescription,
		Admin:                         adminStr,
	}
	return &msg, nil
}

func NewInstantiateBabylonContractsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "instantiate-babylon-contracts [babylon-contract-code-id] [btc-staking-contract-code-id] [btc-finality-contract-code-id] [btc-network] [babylon-tag] [btc-confirmation-depth] [checkpoint-finalization-timeout] [notify-cosmos-zone] [consumer-name] [consumer-description] [admin]",
		Short:   "Instantiate Babylon contracts",
		Long:    "Instantiate Babylon contracts",
		Aliases: []string{"i"},
		Args:    cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg, err := ParseInstantiateArgs(args, clientCtx.GetFromAddress().String())
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
		SilenceUsage: true,
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
