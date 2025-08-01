package cosmwasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ConsumerFpsResponse struct {
	Fps []SingleConsumerFpResponse `json:"fps"`
}

// SingleConsumerFpResponse represents the finality provider data returned by the contract query.
// For more details, refer to the following links:
// https://github.com/babylonchain/babylon-contract/blob/v0.5.3/packages/apis/src/btc_staking_api.rs
// https://github.com/babylonchain/babylon-contract/blob/v0.5.3/contracts/btc-staking/src/msg.rs
// https://github.com/babylonchain/babylon-contract/blob/v0.5.3/contracts/btc-staking/schema/btc-staking.json
type SingleConsumerFpResponse struct {
	BtcPkHex         string `json:"btc_pk_hex"`
	SlashedHeight    uint64 `json:"slashed_height"`
	SlashedBtcHeight uint32 `json:"slashed_btc_height"`
	ConsumerId       string `json:"consumer_id"`
}

type ConsumerDelegationsResponse struct {
	Delegations []SingleConsumerDelegationResponse `json:"delegations"`
}

type SingleConsumerDelegationResponse struct {
	StakerAddr           string                      `json:"staker_addr"`
	BtcPkHex             string                      `json:"btc_pk_hex"`
	FpBtcPkList          []string                    `json:"fp_btc_pk_list"`
	StartHeight          uint32                      `json:"start_height"`
	EndHeight            uint32                      `json:"end_height"`
	TotalSat             uint64                      `json:"total_sat"`
	StakingTx            []byte                      `json:"staking_tx"`
	SlashingTx           []byte                      `json:"slashing_tx"`
	DelegatorSlashingSig []byte                      `json:"delegator_slashing_sig"`
	CovenantSigs         []CovenantAdaptorSignatures `json:"covenant_sigs"`
	StakingOutputIdx     uint32                      `json:"staking_output_idx"`
	UnbondingTime        uint32                      `json:"unbonding_time"`
	UndelegationInfo     *BtcUndelegationInfo        `json:"undelegation_info"`
	ParamsVersion        uint32                      `json:"params_version"`
	Slashed              bool                        `json:"slashed"`
}
type ConsumerFpInfoResponse struct {
	BtcPkHex        string `json:"btc_pk_hex"`
	TotalActiveSats uint64 `json:"total_active_sats"`
	Slashed         bool   `json:"slashed"`
}

type ConsumerFpsByPowerResponse struct {
	Fps []ConsumerFpInfoResponse `json:"fps"`
}

type FinalitySignatureResponse struct {
	Signature []byte `json:"signature"`
}

type BlocksResponse struct {
	Blocks []IndexedBlock `json:"blocks"`
}

type IndexedBlock struct {
	Height    uint64 `json:"height"`
	AppHash   []byte `json:"app_hash"`
	Finalized bool   `json:"finalized"`
}

type NewFinalityProvider struct {
	Description *FinalityProviderDescription `json:"description,omitempty"`
	Commission  string                       `json:"commission"`
	Addr        string                       `json:"addr"`
	BtcPkHex    string                       `json:"btc_pk_hex"`
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
	BTCSig     []byte `json:"btc_sig"`
}

type CovenantAdaptorSignatures struct {
	CovPK       []byte   `json:"cov_pk"`
	AdaptorSigs [][]byte `json:"adaptor_sigs"`
}

type SignatureInfo struct {
	PK  []byte `json:"pk"`
	Sig []byte `json:"sig"`
}

type BtcUndelegationInfo struct {
	UnbondingTx              []byte                       `json:"unbonding_tx"`
	SlashingTx               []byte                       `json:"slashing_tx"`
	DelegatorSlashingSig     []byte                       `json:"delegator_slashing_sig"`
	CovenantSlashingSigs     []*CovenantAdaptorSignatures `json:"covenant_slashing_sigs"`
	CovenantUnbondingSigList []*SignatureInfo             `json:"covenant_unbonding_sig_list"`
	DelegatorUnbondingInfo   *DelegatorUnbondingInfo      `json:"delegator_unbonding_info"`
}

type DelegatorUnbondingInfo struct {
	SpendStakeTx []byte `json:"spend_stake_tx"`
}

type ActiveBtcDelegation struct {
	StakerAddr           string                      `json:"staker_addr"`
	BTCPkHex             string                      `json:"btc_pk_hex"`
	FpBtcPkList          []string                    `json:"fp_btc_pk_list"`
	StartHeight          uint32                      `json:"start_height"`
	EndHeight            uint32                      `json:"end_height"`
	TotalSat             uint64                      `json:"total_sat"`
	StakingTx            []byte                      `json:"staking_tx"`
	SlashingTx           []byte                      `json:"slashing_tx"`
	DelegatorSlashingSig []byte                      `json:"delegator_slashing_sig"`
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

type ExecMsg struct {
	BtcStaking              *BtcStaking              `json:"btc_staking,omitempty"`
	SubmitFinalitySignature *SubmitFinalitySignature `json:"submit_finality_signature,omitempty"`
	CommitPublicRandomness  *CommitPublicRandomness  `json:"commit_public_randomness,omitempty"`
	WithdrawRewards         *WithdrawRewards         `json:"withdraw_rewards,omitempty"`
}

type BtcStaking struct {
	NewFP       []NewFinalityProvider   `json:"new_fp"`
	ActiveDel   []ActiveBtcDelegation   `json:"active_del"`
	SlashedDel  []SlashedBtcDelegation  `json:"slashed_del"`
	UnbondedDel []UnbondedBtcDelegation `json:"unbonded_del"`
}

type FinalityProviderInfo struct {
	BtcPkHex string `json:"btc_pk_hex"`
	Height   uint64 `json:"height,omitempty"`
}

type SubmitFinalitySignature struct {
	FpPubkeyHex string `json:"fp_pubkey_hex"`
	Height      uint64 `json:"height"`
	PubRand     []byte `json:"pub_rand"`
	Proof       Proof  `json:"proof"` // nested struct
	BlockHash   []byte `json:"block_hash"`
	Signature   []byte `json:"signature"`
}

type Proof struct {
	Total    int64    `json:"total"`
	Index    int64    `json:"index"`
	LeafHash []byte   `json:"leaf_hash"`
	Aunts    [][]byte `json:"aunts"`
}

type CommitPublicRandomness struct {
	FPPubKeyHex string `json:"fp_pubkey_hex"`
	StartHeight uint64 `json:"start_height"`
	NumPubRand  uint64 `json:"num_pub_rand"`
	Commitment  []byte `json:"commitment"`
	Signature   []byte `json:"signature"`
}

type WithdrawRewards struct {
	StakerAddr  string `json:"staker_addr"`
	FpPubkeyHex string `json:"fp_pubkey_hex"`
}

type QueryMsgFinalityProviderInfo struct {
	FinalityProviderInfo FinalityProviderInfo `json:"finality_provider_info"`
}

type BlockQuery struct {
	Height uint64 `json:"height"`
}

type QueryMsgBlock struct {
	Block BlockQuery `json:"block"`
}

type QueryMsgBlocks struct {
	Blocks BlocksQuery `json:"blocks"`
}

type BtcHeadersQuery struct {
	Limit *uint32 `json:"limit,omitempty"`
}

type QueryMsgBtcHeaders struct {
	BtcHeaders BtcHeadersQuery `json:"btc_headers"`
}

type BtcHeadersResponse struct {
	Headers []BtcHeaderResponse `json:"headers"`
}

type BtcHeaderResponse struct {
	Header  BtcHeader `json:"header"`
	Hash    string    `json:"hash"`
	Height  uint32    `json:"height"`
	CumWork string    `json:"cum_work"` // Using string for Uint256
}

type BtcHeader struct {
	Version       int32  `json:"version"`
	PrevBlockhash string `json:"prev_blockhash"`
	MerkleRoot    string `json:"merkle_root"`
	Time          uint32 `json:"time"`
	Bits          uint32 `json:"bits"`
	Nonce         uint32 `json:"nonce"`
}

type BlocksQuery struct {
	StartAfter *uint64 `json:"start_after,omitempty"`
	Limit      *uint32 `json:"limit,omitempty"`
	Finalized  *bool   `json:"finalised,omitempty"` //TODO: finalised or finalized, typo in smart contract
	Reverse    *bool   `json:"reverse,omitempty"`
}

type QueryMsgFinalityProviderPower struct {
	FinalityProviderPower FinalityProviderPowerQuery `json:"finality_provider_power"`
}

type FinalityProviderPowerQuery struct {
	BtcPkHex string `json:"btc_pk_hex"`
	Height   uint64 `json:"height"`
}

type ConsumerFpPowerResponse struct {
	Power uint64 `json:"power"`
}

type QueryMsgActivatedHeight struct {
	ActivatedHeight struct{} `json:"activated_height"`
}

type QueryMsgFinalitySignature struct {
	FinalitySignature FinalitySignatureQuery `json:"finality_signature"`
}

type FinalitySignatureQuery struct {
	BtcPkHex string `json:"btc_pk_hex"`
	Height   uint64 `json:"height"`
}

type QueryMsgFinalityProviders struct {
	FinalityProviders struct{} `json:"finality_providers"`
}

type QueryMsgFinalityProvider struct {
	FinalityProvider FinalityProviderQuery `json:"finality_provider"`
}

type FinalityProviderQuery struct {
	BtcPkHex string `json:"btc_pk_hex"`
}

type QueryMsgDelegations struct {
	Delegations struct{} `json:"delegations"`
}

type QueryMsgFinalityProvidersByPower struct {
	FinalityProvidersByPower struct{} `json:"finality_providers_by_power"`
}

type QueryMsgLastPubRandCommit struct {
	LastPubRandCommit LastPubRandCommitQuery `json:"last_pub_rand_commit"`
}

type LastPubRandCommitQuery struct {
	BtcPkHex string `json:"btc_pk_hex"`
}

type QueryMsgPendingRewards struct {
	PendingRewards PendingRewardsQuery `json:"pending_rewards"`
}

type PendingRewardsQuery struct {
	StakerAddr  string `json:"staker_addr"`
	FpPubkeyHex string `json:"fp_pubkey_hex"`
}

type ConsumerPendingRewardsResponse struct {
	Rewards []SingleConsumerPendingRewardsResponse `json:"rewards"`
}

type SingleConsumerPendingRewardsResponse struct {
	Rewards sdk.Coin `json:"rewards"`
}

type QueryMsgAllPendingRewards struct {
	PendingRewards AllPendingRewardsQuery `json:"all_pending_rewards"`
}

type AllPendingRewardsQuery struct {
	StakerAddr string                        `json:"staker_addr"`
	StartAfter *SinglePendingRewardsResponse `json:"start_after,omitempty"`
	Limit      *uint32                       `json:"limit,omitempty"`
}

type ConsumerAllPendingRewardsResponse struct {
	Rewards []SinglePendingRewardsResponse `json:"rewards"`
}

type SinglePendingRewardsResponse struct {
	StakingTxHash []byte   `json:"staking_tx_hash"`
	FpPubkeyHex   string   `json:"fp_pubkey_hex"`
	Rewards       sdk.Coin `json:"rewards"`
}

type QueryMsgLastConsumerHeader struct {
	LastConsumerHeader struct{} `json:"last_consumer_header"`
}

type ConsumerHeaderResponse struct {
	ConsumerID          string `json:"consumer_id"`
	Hash                string `json:"hash"`
	Height              uint64 `json:"height"`
	Time                string `json:"time,omitempty"`
	BabylonHeaderHash   string `json:"babylon_header_hash"`
	BabylonHeaderHeight uint64 `json:"babylon_header_height"`
	BabylonEpoch        uint64 `json:"babylon_epoch"`
	BabylonTxHash       string `json:"babylon_tx_hash"`
}
