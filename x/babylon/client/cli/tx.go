package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

const (
	flagAdmin                = "admin"
	flagIbcTransferChannelId = "ibc-transfer-channel-id"
)

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

// [babylon-contract-code-id] [btc-staking-contract-code-id] [btc-finality-contract-code-id] [btc-network] [babylon-tag] [btc-confirmation-depth] [checkpoint-finalization-timeout] [notify-cosmos-zone] [btc-staking-init-msg-json] [btc-finality-init-msg-json] [consumer-name] [consumer-description]
func ParseInstantiateArgs(args []string, ibcTransferChannelId string, sender string, admin string) (*types.MsgInstantiateBabylonContracts, error) {
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
	btcStakingInitMsg := []byte(args[8])
	if !json.Valid(btcStakingInitMsg) {
		return nil, fmt.Errorf("invalid JSON in btc staking init message")
	}
	btcFinalityInitMsg := []byte(args[9])
	if !json.Valid(btcFinalityInitMsg) {
		return nil, fmt.Errorf("invalid JSON in btc finality init message")
	}
	consumerName := args[10]
	consumerDescription := args[11]

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
		BtcStakingMsg:                 btcStakingInitMsg,
		BtcFinalityMsg:                btcFinalityInitMsg,
		ConsumerName:                  consumerName,
		ConsumerDescription:           consumerDescription,
	}
	if len(ibcTransferChannelId) > 0 {
		msg.IbcTransferChannelId = ibcTransferChannelId
	}
	if len(admin) > 0 {
		msg.Admin = admin
	}
	return &msg, nil
}

func NewInstantiateBabylonContractsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "instantiate-babylon-contracts [babylon-contract-code-id] [btc-staking-contract-code-id] [btc-finality-contract-code-id] [btc-network] [babylon-tag] [btc-confirmation-depth] [checkpoint-finalization-timeout] [notify-cosmos-zone] [btc-staking-init-msg-json] [btc-finality-init-msg-json] [consumer-name] [consumer-description]",
		Short:   "Instantiate Babylon contracts",
		Long:    "Instantiate Babylon contracts",
		Aliases: []string{"i"},
		Args:    cobra.ExactArgs(12),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ibcTransferChannelId, err := cmd.Flags().GetString(flagIbcTransferChannelId)
			if err != nil {
				return err
			}
			admin, err := cmd.Flags().GetString(flagAdmin)
			if err != nil {
				return err
			}

			msg, err := ParseInstantiateArgs(args, ibcTransferChannelId, clientCtx.GetFromAddress().String(), admin)
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
		SilenceUsage: true,
	}

	cmd.Flags().String(flagAdmin, "", "Admin address for the contracts")
	cmd.Flags().String(flagIbcTransferChannelId, "", "IBC transfer channel ID")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
