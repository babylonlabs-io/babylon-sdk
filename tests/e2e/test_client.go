package e2e

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm/ibctesting"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/babylonlabs-io/babylon-sdk/demo/app"
	bbntypes "github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/stretchr/testify/require"
)

// Query is a query type used in tests only
type Query map[string]map[string]any

// QueryResponse is a response type used in tests only
type QueryResponse map[string]any

// To can be used to navigate through the map structure
func (q QueryResponse) To(path ...string) QueryResponse {
	r, ok := q[path[0]]
	if !ok {
		panic(fmt.Sprintf("key %q does not exist", path[0]))
	}
	var x QueryResponse = r.(map[string]any)
	if len(path) == 1 {
		return x
	}
	return x.To(path[1:]...)
}

func (q QueryResponse) Array(key string) []QueryResponse {
	val, ok := q[key]
	if !ok {
		panic(fmt.Sprintf("key %q does not exist", key))
	}
	sl := val.([]any)
	result := make([]QueryResponse, len(sl))
	for i, v := range sl {
		result[i] = v.(map[string]any)
	}
	return result
}

func Querier(t *testing.T, chain *ibctesting.TestChain) func(contract string, query Query) (QueryResponse, error) {
	return func(contract string, query Query) (QueryResponse, error) {
		qRsp := make(map[string]any)
		err := chain.SmartQuery(contract, query, &qRsp)
		if err != nil {
			return nil, err
		}
		return qRsp, nil
	}
}

type TestProviderClient struct {
	t     *testing.T
	Chain *ibctesting.TestChain
}

func NewProviderClient(t *testing.T, chain *ibctesting.TestChain) *TestProviderClient {
	return &TestProviderClient{t: t, Chain: chain}
}

func (p *TestProviderClient) Exec(contract sdk.AccAddress, payload []byte, funds ...sdk.Coin) (*abci.ExecTxResult, error) {
	rsp, err := p.Chain.SendMsgs(&wasmtypes.MsgExecuteContract{
		Sender:   p.Chain.SenderAccount.GetAddress().String(),
		Contract: contract.String(),
		Msg:      payload,
		Funds:    funds,
	})
	return rsp, err
}

type TestConsumerClient struct {
	t         *testing.T
	Chain     *ibctesting.TestChain
	Contracts ConsumerContract
	App       *app.ConsumerApp
}

func NewConsumerClient(t *testing.T, chain *ibctesting.TestChain) *TestConsumerClient {
	return &TestConsumerClient{t: t, Chain: chain, App: chain.App.(*app.ConsumerApp)}
}

type ConsumerContract struct {
	Babylon     sdk.AccAddress
	BTCStaking  sdk.AccAddress
	BTCFinality sdk.AccAddress
}

func (p *TestConsumerClient) GetSender() sdk.AccAddress {
	return p.Chain.SenderAccount.GetAddress()
}

// TODO(babylon): deploy Babylon contracts
func (p *TestConsumerClient) BootstrapContracts() (*ConsumerContract, error) {
	babylonContractWasmId := p.Chain.StoreCodeFile("../testdata/babylon_contract.wasm").CodeID
	btcStakingContractWasmId := p.Chain.StoreCodeFile("../testdata/btc_staking.wasm").CodeID
	btcFinalityContractWasmId := p.Chain.StoreCodeFile("../testdata/btc_finality.wasm").CodeID

	// Instantiate the contract
	// TODO: parameterise
	btcStakingInitMsg := map[string]interface{}{
		"admin": p.GetSender().String(),
	}
	btcStakingInitMsgBytes, err := json.Marshal(btcStakingInitMsg)
	if err != nil {
		return nil, err
	}
	btcFinalityInitMsg := map[string]interface{}{
		"admin": p.GetSender().String(),
	}
	btcFinalityInitMsgBytes, err := json.Marshal(btcFinalityInitMsg)
	if err != nil {
		return nil, err
	}
	initMsg := map[string]interface{}{
		"network":                         "regtest",
		"babylon_tag":                     "01020304",
		"btc_confirmation_depth":          1,
		"checkpoint_finalization_timeout": 2,
		"notify_cosmos_zone":              false,
		"btc_staking_code_id":             btcStakingContractWasmId,
		"btc_staking_msg":                 btcStakingInitMsgBytes,
		"consumer_name":                   "test-consumer",
		"consumer_description":            "test-consumer-description",
		"btc_finality_code_id":            btcFinalityContractWasmId,
		"btc_finality_msg":                btcFinalityInitMsgBytes,
		"admin":                           p.GetSender().String(),
	}
	initMsgBytes, err := json.Marshal(initMsg)
	if err != nil {
		return nil, err
	}

	babylonContractAddr := InstantiateContract(p.t, p.Chain, babylonContractWasmId, initMsgBytes)
	res, err := p.Query(babylonContractAddr, Query{"config": {}})
	if err != nil {
		return nil, err
	}
	btcStakingContractAddr, ok := res["btc_staking"]
	if !ok {
		return nil, fmt.Errorf("failed to instantiate BTC staking contract")
	}
	btcFinalityContractAddr, ok := res["btc_finality"]
	if !ok {
		return nil, fmt.Errorf("failed to instantiate BTC finality contract")
	}

	r := ConsumerContract{
		Babylon:     babylonContractAddr,
		BTCStaking:  sdk.MustAccAddressFromBech32(btcStakingContractAddr.(string)),
		BTCFinality: sdk.MustAccAddressFromBech32(btcFinalityContractAddr.(string)),
	}
	p.Contracts = r
	return &r, nil
}

func (p *TestConsumerClient) Exec(contract sdk.AccAddress, payload []byte, funds ...sdk.Coin) (*abci.ExecTxResult, error) {
	rsp, err := p.Chain.SendMsgs(&wasmtypes.MsgExecuteContract{
		Sender:   p.GetSender().String(),
		Contract: contract.String(),
		Msg:      payload,
		Funds:    funds,
	})
	return rsp, err
}

func (p *TestConsumerClient) Query(contractAddr sdk.AccAddress, query Query) (QueryResponse, error) {
	return Querier(p.t, p.Chain)(contractAddr.String(), query)
}

// MustExecGovProposal submit and vote yes on proposal
func (p *TestConsumerClient) MustExecGovProposal(msg *bbntypes.MsgUpdateParams) {
	proposalID := submitGovProposal(p.t, p.Chain, msg)
	voteAndPassGovProposal(p.t, p.Chain, proposalID)
}

func submitGovProposal(t *testing.T, chain *ibctesting.TestChain, msgs ...sdk.Msg) uint64 {
	// get gov module parameters
	chainApp := chain.App.(*app.ConsumerApp)
	govParams, err := chainApp.GovKeeper.Params.Get(chain.GetContext())
	require.NoError(t, err)

	// construct proposal
	govMsg, err := govv1.NewMsgSubmitProposal(msgs, govParams.MinDeposit, chain.SenderAccount.GetAddress().String(), "", "my title", "my summary", false)
	require.NoError(t, err)

	// submit proposal
	_, err = chain.SendMsgs(govMsg)
	require.NoError(t, err)

	// get next proposal ID
	proposalID, err := chainApp.GovKeeper.ProposalID.Peek(chain.GetContext())
	require.NoError(t, err)

	return proposalID - 1
}

func voteAndPassGovProposal(t *testing.T, chain *ibctesting.TestChain, proposalID uint64) {
	// get gov module parameters
	chainApp := chain.App.(*app.ConsumerApp)
	govParams, err := chainApp.GovKeeper.Params.Get(chain.GetContext())
	require.NoError(t, err)

	// construct and submit vote
	vote := govv1.NewMsgVote(chain.SenderAccount.GetAddress(), proposalID, govv1.OptionYes, "testing")
	_, err = chain.SendMsgs(vote)
	require.NoError(t, err)

	// pass voting period
	coord := chain.Coordinator
	coord.IncrementTimeBy(*govParams.VotingPeriod)
	coord.CommitBlock(chain)

	// ensure proposal is passed
	proposal, err := chainApp.GovKeeper.Proposals.Get(chain.GetContext(), proposalID)
	require.NoError(t, err)
	require.Equal(t, proposal.Status, govv1.ProposalStatus_PROPOSAL_STATUS_PASSED)
}
