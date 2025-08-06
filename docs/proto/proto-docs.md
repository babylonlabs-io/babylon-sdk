<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [babylonlabs/babylon/v1beta1/babylon.proto](#babylonlabs/babylon/v1beta1/babylon.proto)
    - [BSNContracts](#babylonlabs.babylon.v1beta1.BSNContracts)
    - [Params](#babylonlabs.babylon.v1beta1.Params)
  
- [babylonlabs/babylon/v1beta1/genesis.proto](#babylonlabs/babylon/v1beta1/genesis.proto)
    - [GenesisState](#babylonlabs.babylon.v1beta1.GenesisState)
  
- [babylonlabs/babylon/v1beta1/query.proto](#babylonlabs/babylon/v1beta1/query.proto)
    - [QueryBSNContractsRequest](#babylonlabs.babylon.v1beta1.QueryBSNContractsRequest)
    - [QueryBSNContractsResponse](#babylonlabs.babylon.v1beta1.QueryBSNContractsResponse)
    - [QueryParamsRequest](#babylonlabs.babylon.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#babylonlabs.babylon.v1beta1.QueryParamsResponse)
  
    - [Query](#babylonlabs.babylon.v1beta1.Query)
  
- [babylonlabs/babylon/v1beta1/tx.proto](#babylonlabs/babylon/v1beta1/tx.proto)
    - [MsgSetBSNContracts](#babylonlabs.babylon.v1beta1.MsgSetBSNContracts)
    - [MsgSetBSNContractsResponse](#babylonlabs.babylon.v1beta1.MsgSetBSNContractsResponse)
    - [MsgUpdateParams](#babylonlabs.babylon.v1beta1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#babylonlabs.babylon.v1beta1.MsgUpdateParamsResponse)
  
    - [Msg](#babylonlabs.babylon.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="babylonlabs/babylon/v1beta1/babylon.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonlabs/babylon/v1beta1/babylon.proto



<a name="babylonlabs.babylon.v1beta1.BSNContracts"></a>

### BSNContracts
BSNContracts holds all four contract addresses for the Babylon module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `babylon_contract` | [string](#string) |  |  |
| `btc_light_client_contract` | [string](#string) |  |  |
| `btc_staking_contract` | [string](#string) |  |  |
| `btc_finality_contract` | [string](#string) |  |  |






<a name="babylonlabs.babylon.v1beta1.Params"></a>

### Params
Params defines the parameters for the x/babylon module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_gas_begin_blocker` | [uint32](#uint32) |  | max_gas_begin_blocker defines the maximum gas that can be spent in a contract sudo callback for begin blocker |
| `max_gas_end_blocker` | [uint32](#uint32) |  | max_gas_end_blocker defines the maximum gas that can be spent in a contract sudo callback for end blocker |
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
| `params` | [Params](#babylonlabs.babylon.v1beta1.Params) |  |  |
| `bsn_contracts` | [BSNContracts](#babylonlabs.babylon.v1beta1.BSNContracts) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="babylonlabs/babylon/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonlabs/babylon/v1beta1/query.proto



<a name="babylonlabs.babylon.v1beta1.QueryBSNContractsRequest"></a>

### QueryBSNContractsRequest
QueryBSNContractsRequest is the request type for the
Query/BSNContracts RPC method






<a name="babylonlabs.babylon.v1beta1.QueryBSNContractsResponse"></a>

### QueryBSNContractsResponse
QueryBSNContractsResponse is the response type for the
Query/BSNContracts RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bsn_contracts` | [BSNContracts](#babylonlabs.babylon.v1beta1.BSNContracts) |  |  |






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
| `BSNContracts` | [QueryBSNContractsRequest](#babylonlabs.babylon.v1beta1.QueryBSNContractsRequest) | [QueryBSNContractsResponse](#babylonlabs.babylon.v1beta1.QueryBSNContractsResponse) | BSNContracts queries the contract addresses of x/babylon module. | GET|/babylonlabs/babylon/v1beta1/bsn-contracts|

 <!-- end services -->



<a name="babylonlabs/babylon/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonlabs/babylon/v1beta1/tx.proto



<a name="babylonlabs.babylon.v1beta1.MsgSetBSNContracts"></a>

### MsgSetBSNContracts
MsgSetBSNContracts is the Msg/SetBSNContracts request
type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address that controls the module (defaults to x/gov unless overwritten). |
| `contracts` | [BSNContracts](#babylonlabs.babylon.v1beta1.BSNContracts) |  | contracts holds all four contract addresses. |






<a name="babylonlabs.babylon.v1beta1.MsgSetBSNContractsResponse"></a>

### MsgSetBSNContractsResponse
MsgSetBSNContractsResponse is the Msg/SetBSNContracts
response type.






<a name="babylonlabs.babylon.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams is the Msg/UpdateParams request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address that controls the module (defaults to x/gov unless overwritten). |
| `params` | [Params](#babylonlabs.babylon.v1beta1.Params) |  | params defines the x/auth parameters to update.

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
| `SetBSNContracts` | [MsgSetBSNContracts](#babylonlabs.babylon.v1beta1.MsgSetBSNContracts) | [MsgSetBSNContractsResponse](#babylonlabs.babylon.v1beta1.MsgSetBSNContractsResponse) | SetBSNContracts defines an operation for instantiating the Cosmos BSN contracts. | |
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

