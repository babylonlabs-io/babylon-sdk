package babylon

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	wasmdtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/babylonlabs-io/babylon/v3/client/babylonclient"
	"github.com/babylonlabs-io/babylon/v3/testutil/datagen"
	bbntypes "github.com/babylonlabs-io/babylon/v3/types"
	cmtcrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkquerytypes "github.com/cosmos/cosmos-sdk/types/query"

	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/cosmos-integration-e2e/clientcontroller/types"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"

	babylonsdk "github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	"github.com/babylonlabs-io/babylon/v3/crypto/eots"
)

// QueryContractBtcHeaders queries BTC headers from contract
func (bc *BabylonController) QueryContractBtcHeaders(limit *uint32) (*BtcHeadersResponse, error) {
	queryMsgStruct := QueryMsgBtcHeaders{
		BtcHeaders: BtcHeadersQuery{
			Limit: limit,
		},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := bc.QuerySmartContractState(bc.MustQueryBabylonSDKParams().BtcLightClientContractAddress, string(queryMsgBytes))
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

// QueryContractFinalitySignature queries finality signature from contract
func (bc *BabylonController) QueryContractFinalitySignature(fpBtcPkHex string, height uint64) (*FinalitySignatureResponse, error) {
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

	dataFromContract, err := bc.QuerySmartContractState(bc.MustQueryBabylonSDKParams().BtcFinalityContractAddress, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	var resp FinalitySignatureResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// QueryContractFinalityProviders queries finality providers from contract
func (bc *BabylonController) QueryContractFinalityProviders() (*ConsumerFpsResponse, error) {
	queryMsgStruct := QueryMsgFinalityProviders{
		FinalityProviders: struct{}{},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := bc.QuerySmartContractState(bc.MustQueryBabylonSDKParams().BtcFinalityContractAddress, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	var resp ConsumerFpsResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// QueryContractFinalityProvider queries a single finality provider from contract
func (bc *BabylonController) QueryContractFinalityProvider(btcPkHex string) (*SingleConsumerFpResponse, error) {
	queryMsgStruct := QueryMsgFinalityProvider{
		FinalityProvider: FinalityProviderQuery{
			BtcPkHex: btcPkHex,
		},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := bc.QuerySmartContractState(bc.MustQueryBabylonSDKParams().BtcFinalityContractAddress, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	var resp SingleConsumerFpResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// QueryContractDelegations queries delegations from contract
func (bc *BabylonController) QueryContractDelegations() (*ConsumerDelegationsResponse, error) {
	queryMsgStruct := QueryMsgDelegations{
		Delegations: struct{}{},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := bc.QuerySmartContractState(bc.MustQueryBabylonSDKParams().BtcStakingContractAddress, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	var resp ConsumerDelegationsResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// QueryContractAllPendingRewards queries all pending rewards from contract
func (bc *BabylonController) QueryContractAllPendingRewards(stakerAddress string, startAfter *SinglePendingRewardsResponse, limit *uint32) (*ConsumerAllPendingRewardsResponse, error) {
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

	dataFromContract, err := bc.QuerySmartContractState(bc.MustQueryBabylonSDKParams().BtcStakingContractAddress, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	var resp ConsumerAllPendingRewardsResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// QueryContractLastBTCTimestampedHeader queries the last BTC timestamped header from contract
func (bc *BabylonController) QueryContractLastBTCTimestampedHeader() (*ConsumerHeaderResponse, error) {
	queryMsgStruct := QueryMsgLastConsumerHeader{
		LastConsumerHeader: struct{}{},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := bc.QuerySmartContractState(bc.MustQueryBabylonSDKParams().BabylonContractAddress, string(queryMsgBytes))
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

// QueryWasm executes a wasm query function using the babylon client's grpc infrastructure
func (bc *BabylonController) QueryWasm(f func(ctx context.Context, queryClient wasmdtypes.QueryClient) error) error {
	conn, err := bc.createGrpcConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	queryClient := wasmdtypes.NewQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), bc.cfg.Timeout)
	defer cancel()

	return f(ctx, queryClient)
}

// ListCodes lists all wasm codes with the given pagination
func (bc *BabylonController) ListCodes(pagination *sdkquerytypes.PageRequest) (*wasmdtypes.QueryCodesResponse, error) {
	var resp *wasmdtypes.QueryCodesResponse
	err := bc.QueryWasm(func(ctx context.Context, queryClient wasmdtypes.QueryClient) error {
		var err error
		req := &wasmdtypes.QueryCodesRequest{
			Pagination: pagination,
		}
		resp, err = queryClient.Codes(ctx, req)
		return err
	})

	return resp, err
}

// MustQueryBabylonSDKParams queries the Babylon SDK module parameters
func (bc *BabylonController) MustQueryBabylonSDKParams() *babylonsdk.Params {
	ctx := context.Background()

	grpcConn, err := bc.createGrpcConnection()
	if err != nil {
		panic(err)
	}
	defer grpcConn.Close()

	queryClient := babylonsdk.NewQueryClient(grpcConn)

	resp, err := queryClient.Params(ctx, &babylonsdk.QueryParamsRequest{})
	if err != nil {
		panic(err)
	}

	return &resp.Params
}

// ExecuteFinalityContract executes a message on the BTC finality contract
func (bc *BabylonController) ExecuteFinalityContract(msgBytes []byte) (*babylonclient.RelayerTxResponse, error) {
	execMsg := &wasmdtypes.MsgExecuteContract{
		Sender:   bc.MustGetTxSigner(),
		Contract: bc.MustQueryBabylonSDKParams().BtcFinalityContractAddress,
		Msg:      msgBytes,
	}

	return bc.sendMsg(execMsg, emptyErrs, emptyErrs)
}

// QuerySmartContractState queries the smart contract state
func (bc *BabylonController) QuerySmartContractState(contractAddress string, queryData string) (*wasmdtypes.QuerySmartContractStateResponse, error) {
	grpcConn, err := bc.createGrpcConnection()
	if err != nil {
		return nil, err
	}
	defer grpcConn.Close()

	queryClient := wasmdtypes.NewQueryClient(grpcConn)

	resp, err := queryClient.SmartContractState(
		context.Background(),
		&wasmdtypes.QuerySmartContractStateRequest{
			Address:   contractAddress,
			QueryData: []byte(queryData),
		},
	)
	return resp, err
}

// ContractCommitPubRandList commits a list of Schnorr public randomness to contract deployed on Consumer Chain
func (bc *BabylonController) ContractCommitPubRandList(
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

	res, err := bc.ExecuteFinalityContract(msgBytes)
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{TxHash: res.TxHash}, nil
}

// ContractSubmitFinalitySig submits a finality signature to the contract
func (bc *BabylonController) ContractSubmitFinalitySig(
	fpSK *btcec.PrivateKey,
	fpBtcPk *btcec.PublicKey,
	privateRand *eots.PrivateRand,
	pubRand *bbntypes.SchnorrPubRand,
	proof *cmtcrypto.Proof,
	heightToVote uint64,
) (*types.TxResponse, error) {
	block, err := bc.bbnClient.QueryClient.GetBlock(int64(heightToVote))
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

	res, err := bc.ExecuteFinalityContract(msgBytes)
	if err != nil {
		return nil, err
	}

	// Convert events to bytes if needed
	var eventsBytes []byte
	if res.Events != nil {
		eventsBytes, _ = json.Marshal(res.Events)
	}

	return &types.TxResponse{TxHash: res.TxHash, Events: eventsBytes}, nil
}

// ContractSubmitInvalidFinalitySig submits an invalid finality signature to the contract (for testing)
func (bc *BabylonController) ContractSubmitInvalidFinalitySig(
	r *rand.Rand,
	fpSK *btcec.PrivateKey,
	fpBtcPk *btcec.PublicKey,
	privateRand *eots.PrivateRand,
	pubRand *bbntypes.SchnorrPubRand,
	proof *cmtcrypto.Proof,
	heightToVote uint64,
) (*types.TxResponse, error) {
	// Use invalid message to create invalid signature
	invalidAppHash := datagen.GenRandomByteArray(r, 32)
	invalidMsgToSign := append(sdk.Uint64ToBigEndian(heightToVote), invalidAppHash...)
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

	res, err := bc.ExecuteFinalityContract(msgBytes)
	if err != nil {
		return nil, err
	}

	// Convert events to bytes if needed
	var eventsBytes []byte
	if res.Events != nil {
		eventsBytes, _ = json.Marshal(res.Events)
	}

	return &types.TxResponse{TxHash: res.TxHash, Events: eventsBytes}, nil
}

// QueryContractFinalityProviderInfo queries finality provider info from contract
func (bc *BabylonController) QueryContractFinalityProviderInfo(
	fpPk *btcec.PublicKey,
	opts ...uint64,
) (*ConsumerFpInfoResponse, error) {
	var height uint64
	if len(opts) > 0 {
		height = opts[0]
	}

	queryMsgStruct := QueryMsgFinalityProviderInfo{
		FinalityProviderInfo: FinalityProviderInfo{
			BtcPkHex: bbntypes.NewBIP340PubKeyFromBTCPK(fpPk).MarshalHex(),
			Height:   height,
		},
	}
	queryMsgBytes, err := json.Marshal(queryMsgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query message: %v", err)
	}

	dataFromContract, err := bc.QuerySmartContractState(bc.MustQueryBabylonSDKParams().BtcFinalityContractAddress, string(queryMsgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to query smart contract state: %w", err)
	}

	var resp ConsumerFpInfoResponse
	err = json.Unmarshal(dataFromContract.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

// ExecuteStakingContract executes a message on the BTC staking contract
func (bc *BabylonController) ExecuteStakingContract(msgBytes []byte) (*babylonclient.RelayerTxResponse, error) {
	execMsg := &wasmdtypes.MsgExecuteContract{
		Sender:   bc.MustGetTxSigner(),
		Contract: bc.MustQueryBabylonSDKParams().BtcStakingContractAddress,
		Msg:      msgBytes,
	}

	return bc.sendMsg(execMsg, emptyErrs, emptyErrs)
}

// WithdrawRewards withdraws rewards for a staker from a finality provider
func (bc *BabylonController) WithdrawRewards(stakerAddress, fpPubkeyHex string) (*types.TxResponse, error) {
	msg := ExecMsg{
		WithdrawRewards: &WithdrawRewards{
			StakerAddr:  stakerAddress,
			FpPubkeyHex: fpPubkeyHex,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	res, err := bc.ExecuteStakingContract(msgBytes)
	if err != nil {
		return nil, err
	}

	return &types.TxResponse{TxHash: res.TxHash}, nil
}

// QueryContractFinalityContractBalances queries the finality contract balances
func (bc *BabylonController) QueryContractFinalityContractBalances() (sdk.Coins, error) {
	return bc.QueryBalances(bc.MustQueryBabylonSDKParams().BtcFinalityContractAddress)
}

// QueryContractStakingContractBalances queries the staking contract balances
func (bc *BabylonController) QueryContractStakingContractBalances() (sdk.Coins, error) {
	return bc.QueryBalances(bc.MustQueryBabylonSDKParams().BtcStakingContractAddress)
}

// GetCometNodeStatus gets the tendermint node status
func (bc *BabylonController) GetCometNodeStatus() (*coretypes.ResultStatus, error) {
	return bc.bbnClient.QueryClient.GetStatus()
}
