package wasm

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/gogoproto/proto"
)

func RegisterGrpcQueries(bApp *baseapp.BaseApp, appCodec codec.Codec) []wasmkeeper.Option {
	queryRouter := bApp.GRPCQueryRouter()
	queryPluginOpt := wasmkeeper.WithQueryPlugins(
		&wasmkeeper.QueryPlugins{
			Stargate: wasmkeeper.AcceptListStargateQuerier(WhitelistedGrpcQuery(), queryRouter, appCodec),
			Grpc:     wasmkeeper.AcceptListGrpcQuerier(WhitelistedGrpcQuery(), queryRouter, appCodec),
		})

	return []wasmkeeper.Option{
		queryPluginOpt,
	}
}

// WhitelistedGrpcQuery returns the whitelisted Grpc queries
func WhitelistedGrpcQuery() wasmkeeper.AcceptedQueries {
	return wasmkeeper.AcceptedQueries{
		// mint
		"/cosmos.mint.v1beta1.Query/Params": func() proto.Message { return &minttypes.QueryParamsResponse{} },
	}
}
