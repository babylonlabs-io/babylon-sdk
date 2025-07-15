// Package clientcontroller/cosmwasm wraps the CosmWasm RPC/gRPC client for easy interaction with a Wasm enabled node.
// It simplifies querying and submitting transactions to Babylon SDK node (bcd)

// Core CosmWasm RPC/gRPC client lives under https://github.com/babylonlabs-io/cosmwasm-client

// Clientcontroller is adapted from:
// https://github.com/babylonlabs-io/finality-provider/blob/base/consumer-chain-support/clientcontroller/cosmwasm/consumer.go

package cosmwasm

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sdkErr "cosmossdk.io/errors"
	wasmdparams "github.com/CosmWasm/wasmd/app/params"
	wasmdtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	cwconfig "github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmos-integration-e2e/clientcontroller/config"
	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmos-integration-e2e/clientcontroller/types"
	cwcclient "github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmwasm-client/client"
	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmwasm-client/wasmclient"

	"github.com/babylonlabs-io/babylon/v3/crypto/eots"
	"github.com/babylonlabs-io/babylon/v3/testutil/datagen"
	bbntypes "github.com/babylonlabs-io/babylon/v3/types"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	cmtcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkquerytypes "github.com/cosmos/cosmos-sdk/types/query"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	"go.uber.org/zap"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

type CosmwasmConsumerController struct {
	cwClient *cwcclient.Client
	cfg      *cwconfig.CosmwasmConfig
	logger   *zap.Logger
}

func NewCosmwasmConsumerController(
	cfg *cwconfig.CosmwasmConfig,
	encodingCfg wasmdparams.EncodingConfig,
	logger *zap.Logger,
) (*CosmwasmConsumerController, error) {
	wasmdConfig := cfg.ToQueryClientConfig()

	if err := wasmdConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config for Wasmd client: %w", err)
	}

	wc, err := cwcclient.New(
		wasmdConfig,
		"wasmd",
		encodingCfg,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Wasmd client: %w", err)
	}

	return &CosmwasmConsumerController{
		wc,
		cfg,
		logger,
	}, nil
}

func (cc *CosmwasmConsumerController) sendMsg(msg sdk.Msg, expectedErrs []*sdkErr.Error, unrecoverableErrs []*sdkErr.Error) (*wasmclient.RelayerTxResponse, error) {
	return cc.sendMsgs([]sdk.Msg{msg}, expectedErrs, unrecoverableErrs)
}

func (cc *CosmwasmConsumerController) sendMsgs(msgs []sdk.Msg, expectedErrs []*sdkErr.Error, unrecoverableErrs []*sdkErr.Error) (*wasmclient.RelayerTxResponse, error) {
	return cc.cwClient.SendMsgs(
		context.Background(),
		msgs,
		expectedErrs,
		unrecoverableErrs,
	)
}

func (cc *CosmwasmConsumerController) reliablySendMsg(msg sdk.Msg, expectedErrs []*sdkErr.Error, unrecoverableErrs []*sdkErr.Error) (*wasmclient.RelayerTxResponse, error) {
	return cc.reliablySendMsgs([]sdk.Msg{msg}, expectedErrs, unrecoverableErrs)
}

func (cc *CosmwasmConsumerController) reliablySendMsgs(msgs []sdk.Msg, expectedErrs []*sdkErr.Error, unrecoverableErrs []*sdkErr.Error) (*wasmclient.RelayerTxResponse, error) {
	return cc.cwClient.ReliablySendMsgs(
		context.Background(),
		msgs,
		expectedErrs,
		unrecoverableErrs,
	)
}

// CommitPubRandList commits a list of Schnorr public randomness to contract deployed on Consumer Chain
// it returns tx hash and error
func (cc *CosmwasmConsumerController) CommitPubRandList(
	fpPk *btcec.PublicKey,
	startHeight uint64,
	numPubRand uint64,
	commitment []byte,
	sig *schnorr.Signature,
) (*types.TxResponse, error) {
	bip340Sig := bbntypes.NewBIP340SignatureFromBTCSig(sig).MustMarshal()

	// Construct the ExecMsg struct
	msg := ExecMsg{
		CommitPublicRandomness: &CommitPublicRandomness{
			FPPubKeyHex: bbntypes.NewBIP340PubKeyFromBTCPK(fpPk).MarshalHex(),
			StartHeight: startHeight,
			NumPubRand:  numPubRand,
			Commitment:  commitment,
			Signature:   bip340Sig,
		},
	}

	// Marshal the ExecMsg struct to JSON
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	res, err := cc.ExecuteFinalityContract(msgBytes)
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{TxHash: res.TxHash}, nil
}

func (cc *CosmwasmConsumerController) SubmitFinalitySig(
	fpSK *btcec.PrivateKey,
	fpBtcPk *btcec.PublicKey,
	privateRand *eots.PrivateRand,
	pubRand *bbntypes.SchnorrPubRand,
	proof *cmtcrypto.Proof,
	heightToVote uint64,
) (*types.TxResponse, error) {
	block, err := cc.GetCometBlock(int64(heightToVote))
	if err != nil {
		return nil, err
	}

	msgToSign := append(sdk.Uint64ToBigEndian(heightToVote), block.Block.AppHash...)
	sig, err := eots.Sign(fpSK, privateRand, msgToSign)
	if err != nil {
		return nil, err
	}
	eotsSig := bbntypes.NewSchnorrEOTSSigFromModNScalar(sig)

	submitFinalitySig := &SubmitFinalitySignature{
		FpPubkeyHex: bbntypes.NewBIP340PubKeyFromBTCPK(fpBtcPk).MarshalHex(),
		Height:      heightToVote,
		PubRand:     pubRand.MustMarshal(),
		Proof: Proof{
			Total:    proof.Total,
			Index:    proof.Index,
			LeafHash: proof.LeafHash,
			Aunts:    proof.Aunts,
		},
		BlockHash: block.Block.AppHash,
		Signature: eotsSig.MustMarshal(),
	}

	msg := ExecMsg{
		SubmitFinalitySignature: submitFinalitySig,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	res, err := cc.ExecuteFinalityContract(msgBytes)
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{TxHash: res.TxHash, Events: fromCosmosEventsToBytes(res.Events)}, nil
}

func (cc *CosmwasmConsumerController) SubmitInvalidFinalitySig(
	r *rand.Rand,
	fpSK *btcec.PrivateKey,
	fpBtcPk *btcec.PublicKey,
	privateRand *eots.PrivateRand,
	pubRand *bbntypes.SchnorrPubRand,
	proof *cmtcrypto.Proof,
	heightToVote uint64,
) (*types.TxResponse, error) {
	invalidAppHash := datagen.GenRandomByteArray(r, 32)
	invalidMsgToSign := append(sdk.Uint64ToBigEndian(uint64(heightToVote)), invalidAppHash...)
	invalidSig, err := eots.Sign(fpSK, privateRand, invalidMsgToSign)
	if err != nil {
		return nil, err
	}
	invalidEotsSig := bbntypes.NewSchnorrEOTSSigFromModNScalar(invalidSig)

	submitFinalitySig := &SubmitFinalitySignature{
		FpPubkeyHex: bbntypes.NewBIP340PubKeyFromBTCPK(fpBtcPk).MarshalHex(),
		Height:      heightToVote,
		PubRand:     pubRand.MustMarshal(),
		Proof: Proof{
			Total:    proof.Total,
			Index:    proof.Index,
			LeafHash: proof.LeafHash,
			Aunts:    proof.Aunts,
		},
		BlockHash: invalidAppHash,
		Signature: invalidEotsSig.MustMarshal(),
	}

	msg := ExecMsg{
		SubmitFinalitySignature: submitFinalitySig,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	res, err := cc.ExecuteFinalityContract(msgBytes)
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{TxHash: res.TxHash, Events: fromCosmosEventsToBytes(res.Events)}, nil
}

// SubmitBatchFinalitySigs submits a batch of finality signatures to Babylon
func (cc *CosmwasmConsumerController) SubmitBatchFinalitySigs(
	fpPk *btcec.PublicKey,
	blocks []*types.BlockInfo,
	pubRandList []*btcec.FieldVal,
	proofList [][]byte,
	sigs []*btcec.ModNScalar,
) (*types.TxResponse, error) {
	emptyErrs := []*sdkErr.Error{}

	msgs := make([]sdk.Msg, 0, len(blocks))
	for i, b := range blocks {
		cmtProof := cmtcrypto.Proof{}
		if err := cmtProof.Unmarshal(proofList[i]); err != nil {
			return nil, err
		}

		msg := ExecMsg{
			SubmitFinalitySignature: &SubmitFinalitySignature{
				FpPubkeyHex: bbntypes.NewBIP340PubKeyFromBTCPK(fpPk).MarshalHex(),
				Height:      b.Height,
				PubRand:     bbntypes.NewSchnorrPubRandFromFieldVal(pubRandList[i]).MustMarshal(),
				Proof: Proof{
					Total:    cmtProof.Total,
					Index:    cmtProof.Index,
					LeafHash: cmtProof.LeafHash,
					Aunts:    cmtProof.Aunts,
				},
				BlockHash: b.Hash,
				Signature: bbntypes.NewSchnorrEOTSSigFromModNScalar(sigs[i]).MustMarshal(),
			},
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}

		execMsg := &wasmdtypes.MsgExecuteContract{
			Sender:   cc.cwClient.MustGetAddr(),
			Contract: cc.MustQueryBabylonContracts().BtcFinalityContract,
			Msg:      msgBytes,
		}
		msgs = append(msgs, execMsg)
	}

	res, err := cc.sendMsgs(msgs, emptyErrs, emptyErrs)
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{TxHash: res.TxHash}, nil
}

// QueryFinalityProviderHasPower queries whether the finality provider has voting power at a given height
func (cc *CosmwasmConsumerController) QueryFinalityProviderHasPower(
	fpPk *btcec.PublicKey,
	blockHeight uint64,
) (bool, error) {
	fpBtcPkHex := bbntypes.NewBIP340PubKeyFromBTCPK(fpPk).MarshalHex()

	queryMsgStruct := QueryMsgFinalityProviderInfo{
		FinalityProviderInfo: FinalityProviderInfo{
			BtcPkHex: fpBtcPkHex,
			Height:   blockHeight,
		},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return false, fmt.Errorf("failed to marshal query message: %v", err)
	}
	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcStakingContract, string(queryMsgBytes))
	if err != nil {
		return false, err
	}

	var resp ConsumerFpInfoResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return false, err
	}

	return resp.Power > 0, nil
}

func (cc *CosmwasmConsumerController) QueryFinalityProviderInfo(
	fpPk *btcec.PublicKey,
	opts ...uint64,
) (*ConsumerFpInfoResponse, error) {
	fpBtcPkHex := bbntypes.NewBIP340PubKeyFromBTCPK(fpPk).MarshalHex()

	queryMsgStruct := QueryMsgFinalityProviderInfo{
		FinalityProviderInfo: FinalityProviderInfo{
			BtcPkHex: fpBtcPkHex,
		},
	}

	if len(opts) > 0 {
		queryMsgStruct.FinalityProviderInfo.Height = opts[0]
	}

	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcStakingContract, string(queryMsgBytes))
	if err != nil {
		return nil, err
	}

	var resp ConsumerFpInfoResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &resp, nil
}

func (cc *CosmwasmConsumerController) QueryFinalityProvidersByPower() (*ConsumerFpsByPowerResponse, error) {
	queryMsgStruct := QueryMsgFinalityProvidersByPower{
		FinalityProvidersByPower: struct{}{},
	}

	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcStakingContract, string(queryMsgBytes))
	if err != nil {
		return nil, err
	}

	var resp ConsumerFpsByPowerResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cc *CosmwasmConsumerController) QueryLatestFinalizedBlock() (*types.BlockInfo, error) {
	isFinalized := true
	limit := uint32(1)
	blocks, err := cc.queryLatestBlocks(nil, &limit, &isFinalized, nil)
	if err != nil || len(blocks) == 0 {
		// do not return error here as FP handles this situation by
		// not running fast sync
		return nil, nil
	}

	return blocks[0], nil
}

func (cc *CosmwasmConsumerController) QueryBlocks(startHeight, endHeight, limit uint64) ([]*types.BlockInfo, error) {
	return cc.queryCometBlocksInRange(startHeight, endHeight)
}

func (cc *CosmwasmConsumerController) QueryBlock(height uint64) (*types.BlockInfo, error) {
	block, err := cc.cwClient.GetBlock(int64(height))
	if err != nil {
		return nil, err
	}
	return &types.BlockInfo{
		Height: uint64(block.Block.Header.Height),
		Hash:   block.Block.Header.AppHash,
	}, nil
}

// QueryLastPublicRandCommit returns the last public randomness commitments
func (cc *CosmwasmConsumerController) QueryLastPublicRandCommit(fpPk *btcec.PublicKey) (*types.PubRandCommit, error) {
	fpBtcPk := bbntypes.NewBIP340PubKeyFromBTCPK(fpPk)

	// Construct the query message
	queryMsgStruct := QueryMsgLastPubRandCommit{
		LastPubRandCommit: LastPubRandCommitQuery{
			BtcPkHex: fpBtcPk.MarshalHex(),
		},
	}

	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	// Query the smart contract state
	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcFinalityContract, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	if dataFromContract == nil || dataFromContract.Data == nil || len(dataFromContract.Data.Bytes()) == 0 || strings.Contains(string(dataFromContract.Data), "null") {
		// expected when there is no PR commit at all
		return nil, nil
	}

	// Define a response struct
	var commit types.PubRandCommit
	err = json.Unmarshal(dataFromContract.Data.Bytes(), &commit)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if err := commit.Validate(); err != nil {
		return nil, err
	}

	return &commit, nil
}

func (cc *CosmwasmConsumerController) QueryBtcHeaders(limit *uint32) (*BtcHeadersResponse, error) {
	queryMsgStruct := QueryMsgBtcHeaders{
		BtcHeaders: BtcHeadersQuery{
			Limit: limit,
		},
	}

	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %w", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcLightClientContract, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	var resp BtcHeadersResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

func (cc *CosmwasmConsumerController) QueryIsBlockFinalized(height uint64) (bool, error) {
	resp, err := cc.QueryIndexedBlock(height)
	if err != nil || resp == nil {
		return false, nil
	}

	return resp.Finalized, nil
}

func (cc *CosmwasmConsumerController) QueryActivatedHeight() (uint64, error) {
	// Construct the query message
	queryMsg := QueryMsgActivatedHeight{
		ActivatedHeight: struct{}{},
	}

	// Marshal the query message to JSON
	queryMsgBytes, err := json.Marshal(queryMsg)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal query message: %w", err)
	}

	// Query the smart contract state
	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcStakingContract, string(queryMsgBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	// Unmarshal the response
	var resp struct {
		Height uint64 `json:"height"`
	}
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if resp.Height == 0 {
		return 0, fmt.Errorf("BTC staking is not activated yet")
	}

	// Return the activated height
	return resp.Height, nil
}

func (cc *CosmwasmConsumerController) QueryLatestBlockHeight() (uint64, error) {
	block, err := cc.queryCometBestBlock()
	if err != nil {
		return 0, err
	}
	return block.Height, err
}

func (cc *CosmwasmConsumerController) QueryFinalitySignature(fpBtcPkHex string, height uint64) (*FinalitySignatureResponse, error) {
	queryMsgStruct := QueryMsgFinalitySignature{
		FinalitySignature: FinalitySignatureQuery{
			BtcPkHex: fpBtcPkHex,
			Height:   height,
		},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcFinalityContract, string(queryMsgBytes))
	if err != nil {
		return nil, err
	}

	var resp FinalitySignatureResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cc *CosmwasmConsumerController) QueryFinalityProviders() (*ConsumerFpsResponse, error) {
	queryMsgStruct := QueryMsgFinalityProviders{
		FinalityProviders: struct{}{},
	}

	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcStakingContract, string(queryMsgBytes))
	if err != nil {
		return nil, err
	}

	var resp ConsumerFpsResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cc *CosmwasmConsumerController) QueryFinalityProvider(btcPkHex string) (*SingleConsumerFpResponse, error) {
	queryMsgStruct := QueryMsgFinalityProvider{
		FinalityProvider: FinalityProviderQuery{
			BtcPkHex: btcPkHex,
		},
	}

	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcStakingContract, string(queryMsgBytes))
	if err != nil {
		return nil, err
	}

	var resp SingleConsumerFpResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cc *CosmwasmConsumerController) QueryDelegations() (*ConsumerDelegationsResponse, error) {
	queryMsgStruct := QueryMsgDelegations{
		Delegations: struct{}{},
	}

	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcStakingContract, string(queryMsgBytes))
	if err != nil {
		return nil, err
	}

	var resp ConsumerDelegationsResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cc *CosmwasmConsumerController) QueryPendingRewards(stakerAddress, fpPubkeyHex string) (*ConsumerPendingRewardsResponse, error) {
	queryMsgStruct := QueryMsgPendingRewards{
		PendingRewards: PendingRewardsQuery{
			StakerAddr:  stakerAddress,
			FpPubkeyHex: fpPubkeyHex,
		},
	}

	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcStakingContract, string(queryMsgBytes))
	if err != nil {
		return nil, err
	}

	var resp ConsumerPendingRewardsResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cc *CosmwasmConsumerController) QueryAllPendingRewards(stakerAddress string, startAfter *SinglePendingRewardsResponse, limit *uint32) (*ConsumerAllPendingRewardsResponse, error) {
	queryMsgStruct := QueryMsgAllPendingRewards{
		PendingRewards: AllPendingRewardsQuery{
			StakerAddr: stakerAddress,
			StartAfter: startAfter,
			Limit:      limit,
		},
	}

	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcStakingContract, string(queryMsgBytes))
	if err != nil {
		return nil, err
	}

	var resp ConsumerAllPendingRewardsResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// WithdrawRewards
func (cc *CosmwasmConsumerController) WithdrawRewards(stakerAddress, fpPubkeyHex string) (*types.TxResponse, error) {
	// Construct the ExecMsg struct
	msg := ExecMsg{
		WithdrawRewards: &WithdrawRewards{
			StakerAddr:  stakerAddress,
			FpPubkeyHex: fpPubkeyHex,
		},
	}

	// Marshal the ExecMsg struct to JSON
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	res, err := cc.ExecuteStakingContract(msgBytes)
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{TxHash: res.TxHash}, nil
}

func (cc *CosmwasmConsumerController) QueryBabylonContractBalances() (sdk.Coins, error) {
	return cc.QueryBalances(cc.MustQueryBabylonContracts().BtcStakingContract)
}

func (cc *CosmwasmConsumerController) QueryFinalityContractBalances() (sdk.Coins, error) {
	return cc.QueryBalances(cc.MustQueryBabylonContracts().BtcFinalityContract)
}

func (cc *CosmwasmConsumerController) QueryStakingContractBalances() (sdk.Coins, error) {
	return cc.QueryBalances(cc.MustQueryBabylonContracts().BtcStakingContract)
}

func (cc *CosmwasmConsumerController) QueryBalance(address string, denom string) (*sdk.Coin, error) {
	grpcConn, err := cc.createGrpcConnection()
	if err != nil {
		return nil, err
	}
	defer grpcConn.Close()

	// create a gRPC client to query the x/bank service
	bankClient := banktypes.NewQueryClient(grpcConn)
	bankRes, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: address, Denom: denom},
	)
	if err != nil {
		return nil, err
	}

	return bankRes.GetBalance(), nil
}

// QueryBalances returns balances at the address
func (cc *CosmwasmConsumerController) QueryBalances(address string) (sdk.Coins, error) {
	grpcConn, err := cc.createGrpcConnection()
	if err != nil {
		return nil, err
	}
	defer grpcConn.Close()

	// create a gRPC client to query the x/bank service.
	bankClient := banktypes.NewQueryClient(grpcConn)
	bankRes, err := bankClient.AllBalances(
		context.Background(),
		&banktypes.QueryAllBalancesRequest{Address: address},
	)
	if err != nil {
		return nil, err
	}
	return bankRes.GetBalances(), nil
}

func (cc *CosmwasmConsumerController) queryLatestBlocks(startAfter *uint64, limit *uint32, finalized, reverse *bool) ([]*types.BlockInfo, error) {
	// Construct the query message
	queryMsg := QueryMsgBlocks{
		Blocks: BlocksQuery{
			StartAfter: startAfter,
			Limit:      limit,
			Finalized:  finalized,
			Reverse:    reverse,
		},
	}

	// Marshal the query message to JSON
	queryMsgBytes, err := json.Marshal(queryMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %w", err)
	}

	// Query the smart contract state
	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcFinalityContract, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	// Unmarshal the response
	var resp BlocksResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Process the blocks and convert them to BlockInfo
	var blocks []*types.BlockInfo
	for _, b := range resp.Blocks {
		block := &types.BlockInfo{
			Height: b.Height,
			Hash:   b.AppHash,
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (cc *CosmwasmConsumerController) queryCometBestBlock() (*types.BlockInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cc.cfg.Timeout)
	defer cancel()

	// this will return 20 items at max in the descending order (highest first)
	chainInfo, err := cc.cwClient.RPCClient.BlockchainInfo(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	// Returning response directly, if header with specified number did not exist
	// at request will contain nil header
	return &types.BlockInfo{
		Height: uint64(chainInfo.BlockMetas[0].Header.Height),
		Hash:   chainInfo.BlockMetas[0].Header.AppHash,
	}, nil
}

func (cc *CosmwasmConsumerController) queryCometBlocksInRange(startHeight, endHeight uint64) ([]*types.BlockInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cc.cfg.Timeout)
	defer cancel()

	// this will return 20 items at max in the descending order (highest first)
	chainInfo, err := cc.cwClient.RPCClient.BlockchainInfo(ctx, int64(startHeight), int64(endHeight))
	if err != nil {
		return nil, err
	}

	// If no blocks found, return an empty slice
	if len(chainInfo.BlockMetas) == 0 {
		return nil, fmt.Errorf("no comet blocks found in the range")
	}

	// Process the blocks and convert them to BlockInfo
	var blocks []*types.BlockInfo
	for _, blockMeta := range chainInfo.BlockMetas {
		block := &types.BlockInfo{
			Height: uint64(blockMeta.Header.Height),
			Hash:   blockMeta.Header.AppHash,
		}
		blocks = append(blocks, block)
	}

	// Sort the blocks by height in ascending order
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Height < blocks[j].Height
	})

	return blocks, nil
}

func (cc *CosmwasmConsumerController) Close() error {
	if !cc.cwClient.IsRunning() {
		return nil
	}

	return cc.cwClient.Stop()
}

func (cc *CosmwasmConsumerController) ExecuteStakingContract(msgBytes []byte) (*wasmclient.RelayerTxResponse, error) {
	emptyErrs := []*sdkErr.Error{}

	execMsg := &wasmdtypes.MsgExecuteContract{
		Sender:   cc.cwClient.MustGetAddr(),
		Contract: cc.MustQueryBabylonContracts().BtcStakingContract,
		Msg:      msgBytes,
	}

	res, err := cc.sendMsg(execMsg, emptyErrs, emptyErrs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (cc *CosmwasmConsumerController) ExecuteFinalityContract(msgBytes []byte) (*wasmclient.RelayerTxResponse, error) {
	emptyErrs := []*sdkErr.Error{}

	execMsg := &wasmdtypes.MsgExecuteContract{
		Sender:   cc.cwClient.MustGetAddr(),
		Contract: cc.MustQueryBabylonContracts().BtcFinalityContract,
		Msg:      msgBytes,
	}

	res, err := cc.sendMsg(execMsg, emptyErrs, emptyErrs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// QuerySmartContractState queries the smart contract state
// NOTE: this function is only meant to be used in tests.
func (cc *CosmwasmConsumerController) QuerySmartContractState(contractAddress string, queryData string) (*wasmdtypes.QuerySmartContractStateResponse, error) {
	return cc.cwClient.QuerySmartContractState(contractAddress, queryData)
}

// StoreWasmCode stores the wasm code on the consumer chain
// NOTE: this function is only meant to be used in tests.
func (cc *CosmwasmConsumerController) StoreWasmCode(wasmFile string) error {
	return cc.cwClient.StoreWasmCode(wasmFile)
}

// InstantiateContract instantiates a contract with the given code id and init msg
// NOTE: this function is only meant to be used in tests.
func (cc *CosmwasmConsumerController) InstantiateContract(codeID uint64, initMsg []byte) error {
	return cc.cwClient.InstantiateContract(codeID, initMsg)
}

// GetLatestCodeId returns the latest wasm code id.
// NOTE: this function is only meant to be used in tests.
func (cc *CosmwasmConsumerController) GetLatestCodeId() (uint64, error) {
	return cc.cwClient.GetLatestCodeId()
}

// ListContractsByCode lists all contracts by wasm code id
// NOTE: this function is only meant to be used in tests.
func (cc *CosmwasmConsumerController) ListContractsByCode(codeID uint64, pagination *sdkquerytypes.PageRequest) (*wasmdtypes.QueryContractsByCodeResponse, error) {
	return cc.cwClient.ListContractsByCode(codeID, pagination)
}

// MustGetValidatorAddress gets the validator address of the consumer chain
// NOTE: this function is only meant to be used in tests.
func (cc *CosmwasmConsumerController) MustGetValidatorAddress() string {
	return cc.cwClient.MustGetAddr()
}

// GetCometNodeStatus gets the tendermint node status of the consumer chain
// NOTE: this function is only meant to be used in tests.
func (cc *CosmwasmConsumerController) GetCometNodeStatus() (*coretypes.ResultStatus, error) {
	return cc.cwClient.GetStatus()
}

// GetCometBlock gets the tendermint block at a given height
// NOTE: this function is only meant to be used in tests.
func (cc *CosmwasmConsumerController) GetCometBlock(height int64) (*coretypes.ResultBlock, error) {
	return cc.cwClient.GetBlock(height)
}

// QueryIndexedBlock queries the indexed block at a given height
// NOTE: this function is only meant to be used in tests.
func (cc *CosmwasmConsumerController) QueryIndexedBlock(height uint64) (*IndexedBlock, error) {
	// Construct the query message
	queryMsgStruct := QueryMsgBlock{
		Block: BlockQuery{
			Height: height,
		},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	// Query the smart contract state
	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BtcFinalityContract, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	// Unmarshal the response
	var resp IndexedBlock
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

func fromCosmosEventsToBytes(events []wasmclient.RelayerEvent) []byte {
	bytes, err := json.Marshal(events)
	if err != nil {
		return nil
	}
	return bytes
}

func (cc *CosmwasmConsumerController) QueryNodeStatus() (*coretypes.ResultStatus, error) {
	return cc.cwClient.QueryClient.RPCClient.Status(context.Background())
}

func (cc *CosmwasmConsumerController) QueryChannelClientState(channelID, portID string) (*channeltypes.QueryChannelClientStateResponse, error) {
	var resp *channeltypes.QueryChannelClientStateResponse
	err := cc.cwClient.QueryClient.QueryIBCChannel(func(ctx context.Context, queryClient channeltypes.QueryClient) error {
		var err error
		req := &channeltypes.QueryChannelClientStateRequest{
			ChannelId: channelID,
			PortId:    portID,
		}
		resp, err = queryClient.ChannelClientState(ctx, req)
		return err
	})

	return resp, err
}

func (cc *CosmwasmConsumerController) QueryNextSequenceReceive(channelID, portID string) (*channeltypes.QueryNextSequenceReceiveResponse, error) {
	var resp *channeltypes.QueryNextSequenceReceiveResponse
	err := cc.cwClient.QueryClient.QueryIBCChannel(func(ctx context.Context, queryClient channeltypes.QueryClient) error {
		var err error
		req := &channeltypes.QueryNextSequenceReceiveRequest{
			ChannelId: channelID,
			PortId:    portID,
		}
		resp, err = queryClient.NextSequenceReceive(ctx, req)
		return err
	})
	return resp, err
}

// IBCChannels queries the IBC channels
func (cc *CosmwasmConsumerController) IBCChannels() (*channeltypes.QueryChannelsResponse, error) {
	return cc.cwClient.IBCChannels()
}

func (cc *CosmwasmConsumerController) createGrpcConnection() (*grpc.ClientConn, error) {
	// Create a connection to the gRPC server.
	parsedUrl, err := url.Parse(cc.cfg.GRPCAddr)
	if err != nil {
		return nil, fmt.Errorf("grpc-address is not correctly formatted: %w", err)
	}
	endpoint := fmt.Sprintf("%s:%s", parsedUrl.Hostname(), parsedUrl.Port())
	grpcConn, err := grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // The Cosmos SDK doesn't support any transport security mechanism.
		// This instantiates a general gRPC codec which handles proto bytes. We pass in a nil interface registry
		// if the request/response types contain interface instead of 'nil' you should pass the application specific codec.
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		return nil, err
	}
	return grpcConn, nil
}

// SendMsg is a public wrapper around sendMsg for external use
func (cc *CosmwasmConsumerController) SendMsg(msg sdk.Msg) (*wasmclient.RelayerTxResponse, error) {
	emptyErrs := []*sdkErr.Error{}
	return cc.sendMsg(msg, emptyErrs, emptyErrs)
}

// SubmitGovernanceProposal submits a real governance proposal and returns the proposal ID
func (cc *CosmwasmConsumerController) SubmitGovernanceProposal(msgs []sdk.Msg, title, summary string) (uint64, error) {
	// Query gov params for min deposit
	grpcConn, err := cc.createGrpcConnection()
	if err != nil {
		return 0, err
	}
	defer grpcConn.Close()
	govClient := govtypes.NewQueryClient(grpcConn)
	paramsResp, err := govClient.Params(context.Background(), &govtypes.QueryParamsRequest{ParamsType: "deposit"})
	if err != nil {
		return 0, err
	}
	minDeposit := paramsResp.Params.MinDeposit

	// Construct MsgSubmitProposal
	proposer := cc.MustGetValidatorAddress()
	govMsg, err := govtypes.NewMsgSubmitProposal(msgs, minDeposit, proposer, "", title, summary, false)
	if err != nil {
		return 0, err
	}

	// Submit proposal
	emptyErrs := []*sdkErr.Error{}
	_, err = cc.sendMsg(govMsg, emptyErrs, emptyErrs)
	if err != nil {
		return 0, err
	}

	// Query for the latest proposal ID
	proposalsResp, err := govClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{
		Depositor: proposer,
	})
	if err != nil {
		return 0, err
	}
	if len(proposalsResp.Proposals) == 0 {
		return 0, fmt.Errorf("no proposals found after submission")
	}
	// Return the highest proposal ID
	maxID := proposalsResp.Proposals[0].Id
	for _, p := range proposalsResp.Proposals {
		if p.Id > maxID {
			maxID = p.Id
		}
	}
	return maxID, nil
}

// VoteOnProposal votes on a governance proposal
func (cc *CosmwasmConsumerController) VoteOnProposal(proposalID uint64, option govtypes.VoteOption) error {
	voter := cc.MustGetValidatorAddress()
	voterAddr, err := sdk.AccAddressFromBech32(voter)
	if err != nil {
		return err
	}
	voteMsg := govtypes.NewMsgVote(voterAddr, proposalID, option, "")
	emptyErrs := []*sdkErr.Error{}
	_, err = cc.sendMsg(voteMsg, emptyErrs, emptyErrs)
	return err
}

// QueryProposalStatus queries the status of a proposal by ID
func (cc *CosmwasmConsumerController) QueryProposalStatus(proposalID uint64) (govtypes.ProposalStatus, error) {
	grpcConn, err := cc.createGrpcConnection()
	if err != nil {
		return govtypes.ProposalStatus_PROPOSAL_STATUS_UNSPECIFIED, err
	}
	defer grpcConn.Close()
	govClient := govtypes.NewQueryClient(grpcConn)
	resp, err := govClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: proposalID})
	if err != nil {
		return govtypes.ProposalStatus_PROPOSAL_STATUS_UNSPECIFIED, err
	}
	return resp.Proposal.Status, nil
}

func (cc *CosmwasmConsumerController) QueryLastBTCTimestampedHeader() (*ConsumerHeaderResponse, error) {
	queryMsgStruct := QueryMsgLastConsumerHeader{
		LastConsumerHeader: struct{}{},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := cc.QuerySmartContractState(cc.MustQueryBabylonContracts().BabylonContract, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	var resp ConsumerHeaderResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}
