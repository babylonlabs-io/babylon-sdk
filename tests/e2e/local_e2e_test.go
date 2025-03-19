package e2e

import (
	"encoding/json"
	"testing"

	"github.com/CosmWasm/wasmd/tests/ibctesting"
	"github.com/babylonlabs-io/babylon-sdk/demo/app"
	appparams "github.com/babylonlabs-io/babylon-sdk/demo/app/params"
	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/types"
	bbntypes "github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting2 "github.com/cosmos/ibc-go/v8/testing"
	"github.com/stretchr/testify/suite"
)

var testMsg types.ExecuteMessage

// In the Test function, we create and run the suite
func TestBabylonSDKTestSuite(t *testing.T) {
	suite.Run(t, new(BabylonSDKTestSuite))
}

// Define the test suite and include the s.Suite struct
type BabylonSDKTestSuite struct {
	suite.Suite

	// provider/consumer and their metadata
	Coordinator      *ibctesting.Coordinator
	ConsumerChain    *ibctesting.TestChain
	ProviderChain    *ibctesting.TestChain
	ConsumerApp      *app.ConsumerApp
	IbcPath          *ibctesting.Path
	ProviderDenom    string
	ConsumerDenom    string
	MyProvChainActor string

	// clients side information
	ProviderCli      *types.TestProviderClient
	ConsumerCli      *types.TestConsumerClient
	ConsumerContract *types.ConsumerContract
}

// SetupSuite runs once before the suite's tests are run
func (s *BabylonSDKTestSuite) SetupSuite() {
	// overwrite init messages in Babylon
	appparams.SetAddressPrefixes()

	// set up coordinator and chains
	t := s.T()
	coord := types.NewIBCCoordinator(t)
	provChain := coord.GetChain(ibctesting2.GetChainID(1))
	consChain := coord.GetChain(ibctesting2.GetChainID(2))

	s.Coordinator = coord
	s.ConsumerChain = consChain
	s.ProviderChain = provChain
	s.ConsumerApp = consChain.App.(*app.ConsumerApp)
	s.IbcPath = ibctesting.NewPath(consChain, provChain)
	s.ProviderDenom = sdk.DefaultBondDenom
	s.ConsumerDenom = sdk.DefaultBondDenom
	s.MyProvChainActor = provChain.SenderAccount.GetAddress().String()
}

func (s *BabylonSDKTestSuite) Test1ContractDeployment() {
	s.Coordinator.SetupConnections(s.IbcPath)

	// consumer client
	consumerCli := types.NewConsumerClient(s.T(), s.ConsumerChain)
	// setup contracts on consumer
	consumerContracts, err := consumerCli.BootstrapContracts()
	s.NoError(err)
	// provider client
	providerCli := types.NewProviderClient(s.T(), s.ProviderChain)

	s.NotEmpty(consumerCli.Chain.ChainID)
	s.NotEmpty(providerCli.Chain.ChainID)
	s.NotEmpty(consumerContracts.Babylon)
	s.NotEmpty(consumerContracts.BTCStaking)
	s.NotEmpty(consumerContracts.BTCFinality)

	s.ProviderCli = providerCli
	s.ConsumerCli = consumerCli
	s.ConsumerContract = consumerContracts

	// query admins
	adminRespStaking, err := s.ConsumerCli.Query(s.ConsumerContract.BTCStaking, types.Query{"admin": {}})
	s.NoError(err)
	s.Equal(adminRespStaking["admin"], s.ConsumerCli.GetSender().String())
	adminRespFinality, err := s.ConsumerCli.Query(s.ConsumerContract.BTCFinality, types.Query{"admin": {}})
	s.NoError(err)
	s.Equal(adminRespFinality["admin"], s.ConsumerCli.GetSender().String())

	// get contract addresses
	babylonContractAddress := s.ConsumerContract.Babylon.String()
	btcStakingContractAddress := s.ConsumerContract.BTCStaking.String()
	btcFinalityContractAddress := s.ConsumerContract.BTCFinality.String()

	// update the contract address in parameters
	msgUpdateParams := &bbntypes.MsgUpdateParams{
		Authority: s.ConsumerApp.BabylonKeeper.GetAuthority(),
		Params: bbntypes.Params{
			MaxGasBeginBlocker:         500_000,
			BabylonContractAddress:     babylonContractAddress,
			BtcStakingContractAddress:  btcStakingContractAddress,
			BtcFinalityContractAddress: btcFinalityContractAddress,
		},
	}
	s.ConsumerCli.MustExecGovProposal(msgUpdateParams)

	// assert the contract addresses are updated
	params := s.ConsumerApp.BabylonKeeper.GetParams(s.ConsumerChain.GetContext())
	s.Equal(babylonContractAddress, params.BabylonContractAddress)
	s.Equal(btcStakingContractAddress, params.BtcStakingContractAddress)
	s.Equal(btcFinalityContractAddress, params.BtcFinalityContractAddress)
}

func (s *BabylonSDKTestSuite) Test2MockConsumerFpDelegation() {
	// generate headers
	headersMsg := types.GenBTCHeadersMsg()
	headersMsgBytes, err := json.Marshal(headersMsg)
	s.NoError(err)
	// send headers to the Babylon contract. This is to ensure that the contract is
	// indexing BTC headers correctly.
	res, err := s.ConsumerCli.Exec(s.ConsumerContract.Babylon, headersMsgBytes)
	s.NoError(err, res)

	testMsg = types.GenExecMessage()
	msgBytes, err := json.Marshal(testMsg)
	s.NoError(err)

	// send msg to BTC staking contract via admin account
	_, err = s.ConsumerCli.Exec(s.ConsumerContract.BTCStaking, msgBytes)
	s.NoError(err)

	// ensure the finality provider is on consumer chain
	consumerFps, err := s.ConsumerCli.Query(s.ConsumerContract.BTCStaking, types.Query{"finality_providers": {}})
	s.NoError(err)
	s.NotEmpty(consumerFps)

	// ensure delegations are on consumer chain
	consumerDels, err := s.ConsumerCli.Query(s.ConsumerContract.BTCStaking, types.Query{"delegations": {}})
	s.NoError(err)
	s.NotEmpty(consumerDels)

	// ensure the BTC staking is activated
	resp, err := s.ConsumerCli.Query(s.ConsumerContract.BTCStaking, types.Query{"activated_height": {}})
	s.NoError(err)
	parsedActivatedHeight := resp["height"].(float64)
	currentHeight := s.ConsumerChain.GetContext().BlockHeight()
	s.Equal(uint64(parsedActivatedHeight), uint64(currentHeight))
}

func (s *BabylonSDKTestSuite) Test3BeginBlock() {
	err := s.ConsumerApp.BabylonKeeper.BeginBlocker(s.ConsumerChain.GetContext())
	s.NoError(err)
}

func (s *BabylonSDKTestSuite) Test4EndBlock() {
	_, err := s.ConsumerApp.BabylonKeeper.EndBlocker(s.ConsumerChain.GetContext())
	s.NoError(err)
}

func (s *BabylonSDKTestSuite) Test5NextBlock() {
	// get current height
	height := s.ConsumerChain.GetContext().BlockHeight()
	// ensure the current block is not indexed yet
	_, err := s.ConsumerCli.Query(s.ConsumerContract.BTCFinality, types.Query{
		"block": {
			"height": uint64(height),
		},
	})
	s.Error(err)

	// this triggers BeginBlock and EndBlock
	s.ConsumerChain.NextBlock()

	// ensure the current block is indexed
	_, err = s.ConsumerCli.Query(s.ConsumerContract.BTCFinality, types.Query{
		"block": {
			"height": uint64(height),
		},
	})
	s.NoError(err)
}

// TearDownSuite runs once after all the suite's tests have been run
func (s *BabylonSDKTestSuite) TearDownSuite() {
	// Cleanup code here
}
