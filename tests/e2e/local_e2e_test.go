package e2e

import (
	"encoding/json"
	"testing"

	"github.com/CosmWasm/wasmd/tests/ibctesting"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting2 "github.com/cosmos/ibc-go/v8/testing"
	"github.com/stretchr/testify/suite"

	"github.com/babylonlabs-io/babylon-sdk/demo/app"
	appparams "github.com/babylonlabs-io/babylon-sdk/demo/app/params"
	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/types"
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
	BSNChain         *ibctesting.TestChain
	ProviderChain    *ibctesting.TestChain
	BSNApp           *app.BSNApp
	IbcPath          *ibctesting.Path
	ProviderDenom    string
	BSNDenom         string
	MyProvChainActor string

	// clients side information
	ProviderCli *types.TestProviderClient
	BSNCli      *types.TestBSNClient
	BSNContract *types.BSNContract
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
	s.BSNChain = consChain
	s.ProviderChain = provChain
	s.BSNApp = consChain.App.(*app.BSNApp)
	s.IbcPath = ibctesting.NewPath(consChain, provChain)
	s.ProviderDenom = sdk.DefaultBondDenom
	s.BSNDenom = sdk.DefaultBondDenom
	s.MyProvChainActor = provChain.SenderAccount.GetAddress().String()
}

func (s *BabylonSDKTestSuite) Test1ContractDeployment() {
	s.Coordinator.SetupConnections(s.IbcPath)

	// BSN client
	bsnCli := types.NewBSNClient(s.T(), s.BSNChain)
	// setup contracts on BSN
	bsnContracts, err := bsnCli.BootstrapContracts()
	s.NoError(err)
	// provider client
	providerCli := types.NewProviderClient(s.T(), s.ProviderChain)

	s.NotEmpty(bsnCli.Chain.ChainID)
	s.NotEmpty(providerCli.Chain.ChainID)
	s.NotEmpty(bsnContracts.Babylon)
	s.NotEmpty(bsnContracts.BTCLightClient)
	s.NotEmpty(bsnContracts.BTCStaking)
	s.NotEmpty(bsnContracts.BTCFinality)

	s.ProviderCli = providerCli
	s.BSNCli = bsnCli
	s.BSNContract = bsnContracts

	// assert the contract addresses are updated
	params := s.BSNApp.BabylonKeeper.GetParams(s.BSNChain.GetContext())
	s.Equal(s.BSNContract.Babylon.String(), params.BabylonContractAddress)
	s.Equal(s.BSNContract.BTCLightClient.String(), params.BtcLightClientContractAddress)
	s.Equal(s.BSNContract.BTCStaking.String(), params.BtcStakingContractAddress)
	s.Equal(s.BSNContract.BTCFinality.String(), params.BtcFinalityContractAddress)

	// query admins
	adminRespStaking, err := s.BSNCli.Query(s.BSNContract.BTCStaking, types.Query{"admin": {}})
	s.NoError(err)
	s.Equal(adminRespStaking["admin"], s.BSNCli.GetSender().String())
	adminRespFinality, err := s.BSNCli.Query(s.BSNContract.BTCFinality, types.Query{"admin": {}})
	s.NoError(err)
	s.Equal(adminRespFinality["admin"], s.BSNCli.GetSender().String())
}

func (s *BabylonSDKTestSuite) Test2InsertBTCHeaders() {
	// generate headers
	headers, headersMsg := types.GenBTCHeadersMsg(nil)
	headersMsgBytes, err := json.Marshal(headersMsg)
	s.NoError(err)
	// send headers to the BTCLightClient contract. This is to ensure that the contract is
	// indexing BTC headers correctly.
	res, err := s.BSNCli.Exec(s.BSNContract.BTCLightClient, headersMsgBytes)
	s.NoError(err, res)

	// query the base header
	baseHeader, err := s.BSNCli.Query(s.BSNContract.BTCLightClient, types.Query{"btc_base_header": {}})
	s.NoError(err)
	s.NotEmpty(baseHeader)
	s.T().Logf("baseHeader: %v", baseHeader)

	// query the tip header
	tipHeader, err := s.BSNCli.Query(s.BSNContract.BTCLightClient, types.Query{"btc_tip_header": {}})
	s.NoError(err)
	s.NotEmpty(tipHeader)
	s.T().Logf("tipHeader: %v", tipHeader)

	// insert more headers
	_, headersMsg2 := types.GenBTCHeadersMsg(headers[len(headers)-1])
	headersMsgBytes2, err := json.Marshal(headersMsg2)
	s.NoError(err)
	res, err = s.BSNCli.Exec(s.BSNContract.BTCLightClient, headersMsgBytes2)
	s.NoError(err, res)

	// query the tip header again
	tipHeader2, err := s.BSNCli.Query(s.BSNContract.BTCLightClient, types.Query{"btc_tip_header": {}})
	s.NoError(err)
	s.NotEmpty(tipHeader2)
	s.T().Logf("tipHeader2: %v", tipHeader2)
}

func (s *BabylonSDKTestSuite) Test3MockBSNFpDelegation() {
	testMsg = types.GenExecMessage()
	msgBytes, err := json.Marshal(testMsg)
	s.NoError(err)

	// send msg to BTC staking contract via admin account
	_, err = s.BSNCli.Exec(s.BSNContract.BTCStaking, msgBytes)
	s.NoError(err)

	// ensure the finality provider is on consumer chain
	consumerFps, err := s.BSNCli.Query(s.BSNContract.BTCStaking, types.Query{"finality_providers": {}})
	s.NoError(err)
	s.NotEmpty(consumerFps)

	// ensure delegations are on consumer chain
	consumerDels, err := s.BSNCli.Query(s.BSNContract.BTCStaking, types.Query{"delegations": {}})
	s.NoError(err)
	s.NotEmpty(consumerDels)

	// ensure the BTC staking is activated
	resp, err := s.BSNCli.Query(s.BSNContract.BTCStaking, types.Query{"activated_height": {}})
	s.NoError(err)
	parsedActivatedHeight := resp["height"].(float64)
	currentHeight := s.BSNChain.GetContext().BlockHeight()
	s.Equal(uint64(parsedActivatedHeight), uint64(currentHeight))
}

func (s *BabylonSDKTestSuite) Test4BeginBlock() {
	err := s.BSNApp.BabylonKeeper.BeginBlocker(s.BSNChain.GetContext())
	s.NoError(err)
}

func (s *BabylonSDKTestSuite) Test4EndBlock() {
	_, err := s.BSNApp.BabylonKeeper.EndBlocker(s.BSNChain.GetContext())
	s.NoError(err)
}

func (s *BabylonSDKTestSuite) Test5NextBlock() {
	// get current height
	height := s.BSNChain.GetContext().BlockHeight()
	// ensure the current block is not indexed yet
	_, err := s.BSNCli.Query(s.BSNContract.BTCFinality, types.Query{
		"block": {
			"height": uint64(height),
		},
	})
	s.Error(err)

	// this triggers BeginBlock and EndBlock
	s.BSNChain.NextBlock()

	// ensure the current block is indexed
	_, err = s.BSNCli.Query(s.BSNContract.BTCFinality, types.Query{
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
