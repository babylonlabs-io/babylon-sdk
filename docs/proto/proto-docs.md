<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [babylonlabs/babylon/v1beta1/babylon.proto](#babylonlabs/babylon/v1beta1/babylon.proto)
    - [Gauge](#babylonlabs.babylon.v1beta1.Gauge)
    - [Params](#babylonlabs.babylon.v1beta1.Params)
  
- [babylonlabs/babylon/v1beta1/genesis.proto](#babylonlabs/babylon/v1beta1/genesis.proto)
    - [GenesisState](#babylonlabs.babylon.v1beta1.GenesisState)
  
- [babylonlabs/babylon/v1beta1/query.proto](#babylonlabs/babylon/v1beta1/query.proto)
    - [QueryParamsRequest](#babylonlabs.babylon.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#babylonlabs.babylon.v1beta1.QueryParamsResponse)
  
    - [Query](#babylonlabs.babylon.v1beta1.Query)
  
- [babylonlabs/babylon/v1beta1/tx.proto](#babylonlabs/babylon/v1beta1/tx.proto)
    - [MsgInstantiateBabylonContracts](#babylonlabs.babylon.v1beta1.MsgInstantiateBabylonContracts)
    - [MsgInstantiateBabylonContractsResponse](#babylonlabs.babylon.v1beta1.MsgInstantiateBabylonContractsResponse)
    - [MsgUpdateParams](#babylonlabs.babylon.v1beta1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#babylonlabs.babylon.v1beta1.MsgUpdateParamsResponse)
  
    - [Msg](#babylonlabs.babylon.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="babylonlabs/babylon/v1beta1/babylon.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonlabs/babylon/v1beta1/babylon.proto



<a name="babylonlabs.babylon.v1beta1.Gauge"></a>

### Gauge
Gauge is an object that stores rewards to be distributed
code adapted from
https://github.com/osmosis-labs/osmosis/blob/v18.0.0/proto/osmosis/incentives/gauge.proto


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | coins are coins that have been in the gauge Can have multiple coin denoms |






<a name="babylonlabs.babylon.v1beta1.Params"></a>

### Params
Params defines the parameters for the x/babylon module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `babylon_contract_code_id` | [uint64](#uint64) |  | babylon_contract_code_id is the code ID of the Babylon contract |
| `btc_light_client_contract_code_id` | [uint64](#uint64) |  | btc_light_client_contract_code_id is the code ID of the BTC light client contract |
| `btc_staking_contract_code_id` | [uint64](#uint64) |  | btc_staking_contract_code_id is the code ID of the BTC staking contract |
| `btc_finality_contract_code_id` | [uint64](#uint64) |  | btc_finality_contract_code_id is the code ID of the BTC finality contract |
| `babylon_contract_address` | [string](#string) |  | babylon_contract_address is the address of the Babylon contract |
| `btc_light_client_contract_address` | [string](#string) |  | btc_light_client_contract_address is the address of the BTC light client contract |
| `btc_staking_contract_address` | [string](#string) |  | btc_staking_contract_address is the address of the BTC staking contract |
| `btc_finality_contract_address` | [string](#string) |  | btc_finality_contract_address is the address of the BTC finality contract |
| `max_gas_begin_blocker` | [uint32](#uint32) |  | max_gas_begin_blocker defines the maximum gas that can be spent in a contract sudo callback |
| `btc_staking_portion` | [string](#string) |  | btc_staking_portion is the portion of rewards that goes to Finality Providers/delegations NOTE: the portion of each Finality Provider/delegation is calculated by using its voting power and finality provider's commission |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="babylonlabs/babylon/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonlabs/babylon/v1beta1/genesis.proto



<a name="babylonlabs.babylon.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines babylon module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#babylonlabs.babylon.v1beta1.Params) |  | params defines all the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="babylonlabs/babylon/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonlabs/babylon/v1beta1/query.proto



<a name="babylonlabs.babylon.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the
Query/Params RPC method






<a name="babylonlabs.babylon.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the
Query/Params RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#babylonlabs.babylon.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="babylonlabs.babylon.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#babylonlabs.babylon.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#babylonlabs.babylon.v1beta1.QueryParamsResponse) | Params queries the parameters of x/babylon module. | GET|/babylonlabs/babylon/v1beta1/params|

 <!-- end services -->



<a name="babylonlabs/babylon/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonlabs/babylon/v1beta1/tx.proto



<a name="babylonlabs.babylon.v1beta1.MsgInstantiateBabylonContracts"></a>

### MsgInstantiateBabylonContracts
MsgInstantiateBabylonContracts is the Msg/InstantiateBabylonContracts request
type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer` | [string](#string) |  | signer is the address who submits the message. |
| `babylon_contract_code_id` | [uint64](#uint64) |  | babylon_contract_code_id is the code ID for the Babylon contract. |
| `btc_light_client_contract_code_id` | [uint64](#uint64) |  | btc_light_client_contract_code_id is the code ID for the BTC light client contract. |
| `btc_staking_contract_code_id` | [uint64](#uint64) |  | btc_staking_contract_code_id is the code ID for the BTC staking contract. |
| `btc_finality_contract_code_id` | [uint64](#uint64) |  | btc_finality_contract_code_id is the code ID for the BTC finality contract. |
| `network` | [string](#string) |  | network is the Bitcoin network to connect to (e.g. "regtest", "testnet", "mainnet") |
| `babylon_tag` | [string](#string) |  | babylon_tag is a unique identifier for this Babylon instance |
| `btc_confirmation_depth` | [uint32](#uint32) |  | btc_confirmation_depth is the number of confirmations required for Bitcoin transactions |
| `checkpoint_finalization_timeout` | [uint32](#uint32) |  | checkpoint_finalization_timeout is the timeout in blocks for checkpoint finalization |
| `notify_cosmos_zone` | [bool](#bool) |  | notify_cosmos_zone indicates whether to notify the Cosmos zone of events |
| `ibc_transfer_channel_id` | [string](#string) |  | ibc_transfer_channel_id is the IBC channel ID for the IBC transfer contract. If empty then the reward distribution will be done at the consumer side. |
| `btc_light_client_msg` | [bytes](#bytes) |  | btc_light_client_msg is the initialization message for the BTC light client contract |
| `btc_staking_msg` | [bytes](#bytes) |  | btc_staking_msg is the initialization message for the BTC staking contract |
| `btc_finality_msg` | [bytes](#bytes) |  | btc_finality_msg is the initialization message for the BTC finality contract |
| `consumer_name` | [string](#string) |  | consumer_name is the name of this consumer chain |
| `consumer_description` | [string](#string) |  | consumer_description is a description of this consumer chain |
| `admin` | [string](#string) |  | admin is the address that controls the Babylon module |






<a name="babylonlabs.babylon.v1beta1.MsgInstantiateBabylonContractsResponse"></a>

### MsgInstantiateBabylonContractsResponse
MsgInstantiateBabylonContractsResponse is the Msg/InstantiateBabylonContracts
response type.






<a name="babylonlabs.babylon.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams is the Msg/UpdateParams request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address that controls the module (defaults to x/gov unless overwritten). |
| `params` | [Params](#babylonlabs.babylon.v1beta1.Params) |  | params defines the x/babylon parameters to update.

NOTE: All parameters must be supplied. |






<a name="babylonlabs.babylon.v1beta1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="babylonlabs.babylon.v1beta1.Msg"></a>

### Msg
Msg defines the wasm Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `InstantiateBabylonContracts` | [MsgInstantiateBabylonContracts](#babylonlabs.babylon.v1beta1.MsgInstantiateBabylonContracts) | [MsgInstantiateBabylonContractsResponse](#babylonlabs.babylon.v1beta1.MsgInstantiateBabylonContractsResponse) | InstantiateBabylonContracts defines an operation for instantiating the Babylon contracts. | |
| `UpdateParams` | [MsgUpdateParams](#babylonlabs.babylon.v1beta1.MsgUpdateParams) | [MsgUpdateParamsResponse](#babylonlabs.babylon.v1beta1.MsgUpdateParamsResponse) | UpdateParams defines a (governance) operation for updating the x/auth module parameters. The authority defaults to the x/gov module account. | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

