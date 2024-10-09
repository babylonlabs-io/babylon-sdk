<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [babylonchain/babylon/v1beta1/babylon.proto](#babylonchain/babylon/v1beta1/babylon.proto)
    - [Params](#babylonchain.babylon.v1beta1.Params)
  
- [babylonchain/babylon/v1beta1/genesis.proto](#babylonchain/babylon/v1beta1/genesis.proto)
    - [GenesisState](#babylonchain.babylon.v1beta1.GenesisState)
  
- [babylonchain/babylon/v1beta1/query.proto](#babylonchain/babylon/v1beta1/query.proto)
    - [QueryParamsRequest](#babylonchain.babylon.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#babylonchain.babylon.v1beta1.QueryParamsResponse)
  
    - [Query](#babylonchain.babylon.v1beta1.Query)
  
- [babylonchain/babylon/v1beta1/tx.proto](#babylonchain/babylon/v1beta1/tx.proto)
    - [MsgUpdateParams](#babylonchain.babylon.v1beta1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#babylonchain.babylon.v1beta1.MsgUpdateParamsResponse)
  
    - [Msg](#babylonchain.babylon.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="babylonchain/babylon/v1beta1/babylon.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonchain/babylon/v1beta1/babylon.proto



<a name="babylonchain.babylon.v1beta1.Params"></a>

### Params
Params defines the parameters for the x/babylon module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `babylon_contract_address` | [string](#string) |  | babylon_contract_address is the address of the Babylon contract |
| `btc_staking_contract_address` | [string](#string) |  | btc_staking_contract_address is the address of the BTC staking contract |
| `btc_finality_contract_address` | [string](#string) |  | btc_finality_contract_address is the address of the BTC finality contract |
| `max_gas_begin_blocker` | [uint32](#uint32) |  | max_gas_begin_blocker defines the maximum gas that can be spent in a contract sudo callback |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="babylonchain/babylon/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonchain/babylon/v1beta1/genesis.proto



<a name="babylonchain.babylon.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines babylon module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#babylonchain.babylon.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="babylonchain/babylon/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonchain/babylon/v1beta1/query.proto



<a name="babylonchain.babylon.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the
Query/Params RPC method






<a name="babylonchain.babylon.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the
Query/Params RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#babylonchain.babylon.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="babylonchain.babylon.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#babylonchain.babylon.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#babylonchain.babylon.v1beta1.QueryParamsResponse) | Params queries the parameters of x/babylon module. | GET|/babylonchain/babylon/v1beta1/params|

 <!-- end services -->



<a name="babylonchain/babylon/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## babylonchain/babylon/v1beta1/tx.proto



<a name="babylonchain.babylon.v1beta1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams is the Msg/UpdateParams request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address that controls the module (defaults to x/gov unless overwritten). |
| `params` | [Params](#babylonchain.babylon.v1beta1.Params) |  | params defines the x/auth parameters to update.

NOTE: All parameters must be supplied. |






<a name="babylonchain.babylon.v1beta1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="babylonchain.babylon.v1beta1.Msg"></a>

### Msg
Msg defines the wasm Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#babylonchain.babylon.v1beta1.MsgUpdateParams) | [MsgUpdateParamsResponse](#babylonchain.babylon.v1beta1.MsgUpdateParamsResponse) | UpdateParams defines a (governance) operation for updating the x/auth module parameters. The authority defaults to the x/gov module account. | |

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

