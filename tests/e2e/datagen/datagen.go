package datagen

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"math/rand"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/babylonlabs-io/babylon/v3/testutil/datagen"
	bbn "github.com/babylonlabs-io/babylon/v3/types"
	btclctypes "github.com/babylonlabs-io/babylon/v3/x/btclightclient/types"
	bstypes "github.com/babylonlabs-io/babylon/v3/x/btcstaking/types"
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
)

// BtcHeaderInfo represents the Bitcoin header information needed for contract instantiation
type BtcHeaderInfo struct {
	Header BtcHeader `json:"header"`
	Work   []byte    `json:"total_work"`
	Height int64     `json:"height"`
}

func (bhi *BtcHeaderInfo) ToHeaderInfo() (*btclctypes.BTCHeaderInfo, error) {
	prevHash, err := chainhash.NewHashFromStr(bhi.Header.PrevBlockhash)
	if err != nil {
		return nil, err
	}
	merkleRoot, err := chainhash.NewHashFromStr(bhi.Header.MerkleRoot)
	if err != nil {
		return nil, err
	}

	header := bbn.NewBTCHeaderBytesFromBlockHeader(&wire.BlockHeader{
		Version:    bhi.Header.Version,
		PrevBlock:  *prevHash,
		MerkleRoot: *merkleRoot,
		Timestamp:  time.Unix(int64(bhi.Header.Time), 0),
		Bits:       bhi.Header.Bits,
		Nonce:      bhi.Header.Nonce,
	})

	var workInt big.Int
	workInt.SetBytes(bhi.Work)

	work := sdkmath.NewUintFromBigInt(&workInt)

	return &btclctypes.BTCHeaderInfo{
		Header: &header,
		Hash:   header.Hash(),
		Height: uint32(bhi.Height),
		Work:   &work,
	}, nil
}

func GenBTCHeadersMsg(parent *btclctypes.BTCHeaderInfo) ([]*btclctypes.BTCHeaderInfo, BabylonExecuteMsg) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	var chain *datagen.BTCHeaderPartialChain
	if parent == nil {
		chain = datagen.NewBTCHeaderChainWithLength(r, 0, 0, 10)
	} else {
		chain = datagen.NewBTCHeaderChainFromParentInfo(r, parent, 10)
	}

	headers := []BtcHeader{}
	for _, header := range chain.Headers {
		headers = append(headers, BtcHeader{
			Version:       header.Version,
			PrevBlockhash: header.PrevBlock.String(),
			MerkleRoot:    header.MerkleRoot.String(),
			Time:          uint32(header.Timestamp.Unix()),
			Bits:          header.Bits,
			Nonce:         header.Nonce,
		})
	}

	msg := BabylonExecuteMsg{
		BtcHeaders: BTCHeadersMsg{
			Headers: headers,
		},
	}

	if parent == nil {
		firstHeader := chain.GetChainInfo()[0]
		msg.BtcHeaders.FirstHeight = &firstHeader.Height
		firstWork, _ := firstHeader.Work.Marshal()
		firstWorkHex := hex.EncodeToString(firstWork)
		msg.BtcHeaders.FirstWork = &firstWorkHex
	} else {
		msg.BtcHeaders.FirstHeight = nil
		msg.BtcHeaders.FirstWork = nil
	}

	return chain.GetChainInfo(), msg
}

// MustGetInitialHeader returns the initial BTC header for the babylon contract instantiation
func MustGetInitialHeader() BtcHeaderInfo {
	// Initial base header generated using GenBTCHeadersMsg
	// https://github.com/babylonlabs-io/babylon-sdk/blob/a8ad676746c1ab75ea849572fc82a3d0b8900c7a/tests/e2e/types/datagen.go#L20
	headerHex := "04000000f67ad7695d9b662a72ff3d8edbbb2de0bfa67b13974bb9910d116d5cbd863e683b48539be32ca95bd6e51615157ca70b6eff9993f11205780810022bfc0d0163cb088653ffff7f206fdbb6ed"

	// the initial height must be at difficulty adjustment boundary
	height := int64(2016)

	// Decode the header
	headerBytes, err := hex.DecodeString(headerHex)
	if err != nil {
		panic("static header hex is invalid") // shouldn't happen for hardcoded value
	}

	// Deserialize the header
	var header wire.BlockHeader
	err = header.Deserialize(bytes.NewReader(headerBytes))
	if err != nil {
		panic("static header bytes are invalid") // shouldn't happen for hardcoded value
	}

	// Calculate the work (difficulty in big-endian bytes)
	work := blockchain.CalcWork(header.Bits)
	// Convert work to 32-byte big-endian
	workBytes := make([]byte, 32)
	work.FillBytes(workBytes) // Pads with zeros if needed

	return BtcHeaderInfo{
		Header: BtcHeader{
			Version:       header.Version,
			PrevBlockhash: header.PrevBlock.String(),
			MerkleRoot:    header.MerkleRoot.String(),
			Time:          uint32(header.Timestamp.Unix()),
			Bits:          header.Bits,
			Nonce:         header.Nonce,
		},
		Height: height,
		Work:   workBytes,
	}
}

func MustGetInitialHeaderInStr() string {
	header := MustGetInitialHeader()
	jsonBytes := marshalHeaderToJSON(header)
	return string(jsonBytes)
}

// MustGetInitialHeaderInHex returns the initial BTC header encoded in json and hex-encoded form
func MustGetInitialHeaderInHex() string {
	header := MustGetInitialHeader()
	jsonBytes := marshalHeaderToJSON(header)
	return hex.EncodeToString(jsonBytes)
}

// marshalHeaderToJSON marshals a BtcHeaderInfo object into JSON bytes.
func marshalHeaderToJSON(header BtcHeaderInfo) []byte {
	jsonBytes, err := json.Marshal(header)
	if err != nil {
		panic("failed to marshal header to json") // shouldn't happen
	}
	return jsonBytes
}

func GenExecMessage() ExecuteMessage {
	_, newDel := genBTCDelegation()

	addr := datagen.GenRandomAccount().Address

	newFp := NewFinalityProvider{
		Description: &FinalityProviderDescription{
			Moniker:         "fp1",
			Identity:        "Finality Provider 1",
			Website:         "https://fp1.com",
			SecurityContact: "security_contact",
			Details:         "details",
		},
		Commission: "0.05",
		Addr:       addr,
		BTCPKHex:   newDel.FpBtcPkList[0],
		Pop: &ProofOfPossessionBtc{
			BTCSigType: 0,
			BTCSig:     base64.StdEncoding.EncodeToString([]byte("mock_pub_rand")),
		},
		ConsumerID: "osmosis-1",
	}

	// Create the ExecuteMessage instance
	executeMessage := ExecuteMessage{
		BtcStaking: BtcStaking{
			NewFP:       []NewFinalityProvider{newFp},
			ActiveDel:   []ActiveBtcDelegation{newDel},
			SlashedDel:  []SlashedBtcDelegation{},
			UnbondedDel: []UnbondedBtcDelegation{},
		},
	}

	return executeMessage
}

func genBTCDelegation() (*bstypes.Params, ActiveBtcDelegation) {
	net := &chaincfg.RegressionNetParams
	r := rand.New(rand.NewSource(time.Now().Unix()))
	t := &testing.T{}

	delSK, _, err := datagen.GenRandomBTCKeyPair(r)
	require.NoError(t, err)

	// restaked to a random number of finality providers
	numRestakedFPs := int(datagen.RandomInt(r, 10) + 1)
	_, fpPKs, err := datagen.GenRandomBTCKeyPairs(r, numRestakedFPs)
	require.NoError(t, err)
	fpBTCPKs := bbn.NewBIP340PKsFromBTCPKs(fpPKs)

	// (3, 5) covenant committee
	covenantSKs, covenantPKs, err := datagen.GenRandomBTCKeyPairs(r, 5)
	require.NoError(t, err)
	covenantQuorum := uint32(3)

	stakingTimeBlocks := uint16(5)
	stakingValue := int64(2 * 10e8)
	slashingAddress, err := datagen.GenRandomBTCAddress(r, net)
	require.NoError(t, err)
	slashingPkScript, err := txscript.PayToAddrScript(slashingAddress)
	require.NoError(t, err)

	slashingRate := sdkmath.LegacyNewDecWithPrec(int64(datagen.RandomInt(r, 41)+10), 2)
	unbondingTime := uint16(100) + 1
	slashingChangeLockTime := unbondingTime

	bsParams := &bstypes.Params{
		CovenantPks:      bbn.NewBIP340PKsFromBTCPKs(covenantPKs),
		CovenantQuorum:   covenantQuorum,
		SlashingPkScript: slashingPkScript,
	}

	// only the quorum of signers provided the signatures
	covenantSigners := covenantSKs[:covenantQuorum]

	// construct the BTC delegation with everything
	btcDel, err := datagen.GenRandomBTCDelegation(
		r,
		t,
		net,
		fpBTCPKs,
		delSK,
		"",
		covenantSigners,
		covenantPKs,
		covenantQuorum,
		slashingPkScript,
		uint32(stakingTimeBlocks),
		uint32(1000),
		uint32(1000+stakingTimeBlocks),
		uint64(stakingValue),
		slashingRate,
		slashingChangeLockTime,
	)
	require.NoError(t, err)

	activeDel := convertBTCDelegationToActiveBtcDelegation(btcDel)
	return bsParams, activeDel
}

func convertBTCDelegationToActiveBtcDelegation(mockDel *bstypes.BTCDelegation) ActiveBtcDelegation {
	var fpBtcPkList []string
	for _, pk := range mockDel.FpBtcPkList {
		fpBtcPkList = append(fpBtcPkList, pk.MarshalHex())
	}

	var covenantSigs []CovenantAdaptorSignatures
	for _, cs := range mockDel.CovenantSigs {
		var adaptorSigs []string
		for _, sig := range cs.AdaptorSigs {
			adaptorSigs = append(adaptorSigs, base64.StdEncoding.EncodeToString(sig))
		}
		covenantSigs = append(covenantSigs, CovenantAdaptorSignatures{
			CovPK:       cs.CovPk.MarshalHex(),
			AdaptorSigs: adaptorSigs,
		})
	}

	var covenantUnbondingSigs []SignatureInfo
	for _, sigInfo := range mockDel.BtcUndelegation.CovenantUnbondingSigList {
		covenantUnbondingSigs = append(covenantUnbondingSigs, SignatureInfo{
			PK:  sigInfo.Pk.MarshalHex(),
			Sig: base64.StdEncoding.EncodeToString(sigInfo.Sig.MustMarshal()),
		})
	}

	var covenantSlashingSigs []CovenantAdaptorSignatures
	for _, cs := range mockDel.BtcUndelegation.CovenantSlashingSigs {
		var adaptorSigs []string
		for _, sig := range cs.AdaptorSigs {
			adaptorSigs = append(adaptorSigs, base64.StdEncoding.EncodeToString(sig))
		}
		covenantSlashingSigs = append(covenantSlashingSigs, CovenantAdaptorSignatures{
			CovPK:       cs.CovPk.MarshalHex(),
			AdaptorSigs: adaptorSigs,
		})
	}

	undelegationInfo := BtcUndelegationInfo{
		UnbondingTx:           base64.StdEncoding.EncodeToString(mockDel.BtcUndelegation.UnbondingTx),
		SlashingTx:            base64.StdEncoding.EncodeToString(mockDel.BtcUndelegation.SlashingTx.MustMarshal()),
		DelegatorSlashingSig:  base64.StdEncoding.EncodeToString(mockDel.BtcUndelegation.DelegatorSlashingSig.MustMarshal()),
		CovenantUnbondingSigs: covenantUnbondingSigs,
		CovenantSlashingSigs:  covenantSlashingSigs,
	}

	return ActiveBtcDelegation{
		StakerAddr:           mockDel.StakerAddr,
		BTCPkHex:             mockDel.BtcPk.MarshalHex(),
		FpBtcPkList:          fpBtcPkList,
		StartHeight:          mockDel.StartHeight,
		EndHeight:            mockDel.EndHeight,
		TotalSat:             mockDel.TotalSat,
		StakingTx:            base64.StdEncoding.EncodeToString(mockDel.StakingTx),
		SlashingTx:           base64.StdEncoding.EncodeToString(mockDel.SlashingTx.MustMarshal()),
		DelegatorSlashingSig: base64.StdEncoding.EncodeToString(mockDel.DelegatorSig.MustMarshal()),
		CovenantSigs:         covenantSigs,
		StakingOutputIdx:     mockDel.StakingOutputIdx,
		UnbondingTime:        mockDel.UnbondingTime,
		UndelegationInfo:     undelegationInfo,
		ParamsVersion:        mockDel.ParamsVersion,
	}
}

type NewFinalityProvider struct {
	Description *FinalityProviderDescription `json:"description,omitempty"`
	Commission  string                       `json:"commission"`
	Addr        string                       `json:"addr,omitempty"`
	BTCPKHex    string                       `json:"btc_pk_hex"`
	Pop         *ProofOfPossessionBtc        `json:"pop,omitempty"`
	ConsumerID  string                       `json:"consumer_id"`
}

type FinalityProviderDescription struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type ProofOfPossessionBtc struct {
	BTCSigType int32  `json:"btc_sig_type"`
	BTCSig     string `json:"btc_sig"`
}

type CovenantAdaptorSignatures struct {
	CovPK       string   `json:"cov_pk"`
	AdaptorSigs []string `json:"adaptor_sigs"`
}

type SignatureInfo struct {
	PK  string `json:"pk"`
	Sig string `json:"sig"`
}

type BtcUndelegationInfo struct {
	UnbondingTx           string                      `json:"unbonding_tx"`
	DelegatorUnbondingSig string                      `json:"delegator_unbonding_sig"`
	CovenantUnbondingSigs []SignatureInfo             `json:"covenant_unbonding_sig_list"`
	SlashingTx            string                      `json:"slashing_tx"`
	DelegatorSlashingSig  string                      `json:"delegator_slashing_sig"`
	CovenantSlashingSigs  []CovenantAdaptorSignatures `json:"covenant_slashing_sigs"`
}

type ActiveBtcDelegation struct {
	StakerAddr           string                      `json:"staker_addr"`
	BTCPkHex             string                      `json:"btc_pk_hex"`
	FpBtcPkList          []string                    `json:"fp_btc_pk_list"`
	StartHeight          uint32                      `json:"start_height"`
	EndHeight            uint32                      `json:"end_height"`
	TotalSat             uint64                      `json:"total_sat"`
	StakingTx            string                      `json:"staking_tx"`
	SlashingTx           string                      `json:"slashing_tx"`
	DelegatorSlashingSig string                      `json:"delegator_slashing_sig"`
	CovenantSigs         []CovenantAdaptorSignatures `json:"covenant_sigs"`
	StakingOutputIdx     uint32                      `json:"staking_output_idx"`
	UnbondingTime        uint32                      `json:"unbonding_time"`
	UndelegationInfo     BtcUndelegationInfo         `json:"undelegation_info"`
	ParamsVersion        uint32                      `json:"params_version"`
}

type SlashedBtcDelegation struct {
	// Define fields as needed
}

type UnbondedBtcDelegation struct {
	// Define fields as needed
}

type BabylonExecuteMsg struct {
	BtcHeaders BTCHeadersMsg `json:"btc_headers"`
}

type BTCHeadersMsg struct {
	Headers     []BtcHeader `json:"headers"`
	FirstWork   *string     `json:"first_work,omitempty"`
	FirstHeight *uint32     `json:"first_height,omitempty"`
}

type BtcHeader struct {
	Version       int32  `json:"version"`
	PrevBlockhash string `json:"prev_blockhash"`
	MerkleRoot    string `json:"merkle_root"`
	Time          uint32 `json:"time"`
	Bits          uint32 `json:"bits"`
	Nonce         uint32 `json:"nonce"`
}

type ExecuteMessage struct {
	BtcStaking BtcStaking `json:"btc_staking"`
}

type BtcStaking struct {
	NewFP       []NewFinalityProvider   `json:"new_fp"`
	ActiveDel   []ActiveBtcDelegation   `json:"active_del"`
	SlashedDel  []SlashedBtcDelegation  `json:"slashed_del"`
	UnbondedDel []UnbondedBtcDelegation `json:"unbonded_del"`
}
