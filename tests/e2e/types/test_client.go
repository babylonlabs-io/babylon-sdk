package types

import (
	"fmt"
	"testing"

	ibctesting "github.com/CosmWasm/wasmd/tests/wasmibctesting"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/stretchr/testify/require"

	"github.com/babylonlabs-io/babylon-sdk/demo/app"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/client/cli"
	bbntypes "github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
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

func Querier(t *testing.T, chain *ibctesting.WasmTestChain) func(contract string, query Query) (QueryResponse, error) {
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
	Chain *ibctesting.WasmTestChain
}

func NewProviderClient(t *testing.T, chain *ibctesting.WasmTestChain) *TestProviderClient {
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
	Chain     *ibctesting.WasmTestChain
	Contracts ConsumerContract
	App       *app.ConsumerApp
}

func NewConsumerClient(t *testing.T, chain *ibctesting.WasmTestChain) *TestConsumerClient {
	return &TestConsumerClient{t: t, Chain: chain, App: chain.App.(*app.ConsumerApp)}
}

type ConsumerContract struct {
	Babylon        sdk.AccAddress
	BTCLightClient sdk.AccAddress
	BTCStaking     sdk.AccAddress
	BTCFinality    sdk.AccAddress
}

func (p *TestConsumerClient) GetSender() sdk.AccAddress {
	return p.Chain.SenderAccount.GetAddress()
}

// TODO(babylon): deploy Babylon contracts
func (p *TestConsumerClient) BootstrapContracts() (*ConsumerContract, error) {
	babylonContractWasmId := p.Chain.StoreCodeFile("../testdata/babylon_contract.wasm").CodeID
	btcLightClientContractWasmId := p.Chain.StoreCodeFile("../testdata/btc_light_client.wasm").CodeID
	btcStakingContractWasmId := p.Chain.StoreCodeFile("../testdata/btc_staking.wasm").CodeID
	btcFinalityContractWasmId := p.Chain.StoreCodeFile("../testdata/btc_finality.wasm").CodeID

	// Instantiate the contract
	// instantiate Babylon contract using cli.Parse approach
	msgInit, err := cli.ParseInstantiateArgs(
		[]string{
			fmt.Sprintf("%d", babylonContractWasmId),
			fmt.Sprintf("%d", btcLightClientContractWasmId),
			fmt.Sprintf("%d", btcStakingContractWasmId),
			fmt.Sprintf("%d", btcFinalityContractWasmId),
			"regtest",
			"01020304",
			"1",
			"2",
			"false",
			fmt.Sprintf(`{"admin":"%s"}`, p.GetSender().String()),
			fmt.Sprintf(`{"admin":"%s"}`, p.GetSender().String()),
			`{"header": {"version": 536870912, "prev_blockhash": "000000c0a3841a6ae64c45864ae25314b40fd522bfb299a4b6bd5ef288cae74d", "merkle_root": "e666a9797b7a650597098ca6bf500bd0873a86ada05189f87073b6dfdbcaf4ee", "time": 1599332844, "bits": 503394215, "nonce": 9108535}, "height": 2016, "total_work": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAkY98OU="}`,
			"test-consumer",
			"test-consumer-description",
		},
		"",
		p.GetSender().String(),
		p.GetSender().String(),
	)
	if err != nil {
		return nil, err
	}
	_, err = p.Chain.SendMsgs(msgInit)
	if err != nil {
		return nil, err
	}

	params := p.App.BabylonKeeper.GetParams(p.Chain.GetContext())
	babylonAddr, btcLightClientAddr, btcStakingAddr, btcFinalityAddr, err := params.GetContractAddresses()
	if err != nil {
		return nil, err
	}

	r := ConsumerContract{
		Babylon:        babylonAddr,
		BTCLightClient: btcLightClientAddr,
		BTCStaking:     btcStakingAddr,
		BTCFinality:    btcFinalityAddr,
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

func submitGovProposal(t *testing.T, chain *ibctesting.WasmTestChain, msgs ...sdk.Msg) uint64 {
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

func voteAndPassGovProposal(t *testing.T, chain *ibctesting.WasmTestChain, proposalID uint64) {
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
	coord.CommitBlock(chain.TestChain)

	// ensure proposal is passed
	proposal, err := chainApp.GovKeeper.Proposals.Get(chain.GetContext(), proposalID)
	require.NoError(t, err)
	require.Equal(t, proposal.Status, govv1.ProposalStatus_PROPOSAL_STATUS_PASSED)
}
