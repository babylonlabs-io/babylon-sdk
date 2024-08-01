package keeper

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/contract"
)

type (
	// abstract query keeper
	ViewKeeper interface {
		GetTest(ctx sdk.Context, actor sdk.AccAddress) string
	}
)

// NewQueryDecorator constructor to build a chained custom querier.
// The babylon custom query handler is placed at the first position
// and delegates to the next in chain for any queries that do not match
// the babylon custom query namespace.
//
// To be used with `wasmkeeper.WithQueryHandlerDecorator(BabylonKeeper.NewQueryDecorator(app.BabylonKeeper)))`
func NewQueryDecorator(k ViewKeeper) func(wasmkeeper.WasmVMQueryHandler) wasmkeeper.WasmVMQueryHandler {
	return func(next wasmkeeper.WasmVMQueryHandler) wasmkeeper.WasmVMQueryHandler {
		return ChainedCustomQuerier(k, next)
	}
}

// ChainedCustomQuerier implements the babylon custom query handler.
// The given WasmVMQueryHandler is receiving all unhandled queries and must therefore
// not be nil.
//
// This CustomQuerier is designed as an extension point. See the NewQueryDecorator impl how to
// set this up for wasmd.
func ChainedCustomQuerier(k ViewKeeper, next wasmkeeper.WasmVMQueryHandler) wasmkeeper.WasmVMQueryHandler {
	if k == nil {
		panic("ms keeper must not be nil")
	}
	if next == nil {
		panic("next handler must not be nil")
	}
	return QueryHandlerFn(func(ctx sdk.Context, caller sdk.AccAddress, request wasmvmtypes.QueryRequest) ([]byte, error) {
		if request.Custom == nil {
			return next.HandleQuery(ctx, caller, request)
		}
		var contractQuery contract.CustomQuery
		if err := json.Unmarshal(request.Custom, &contractQuery); err != nil {
			return nil, errorsmod.Wrap(err, "babylon query")
		}
		query := contractQuery.Test
		if query == nil {
			return next.HandleQuery(ctx, caller, request)
		}

		res := contract.TestResponse{
			Placeholder2: "hello world",
		}
		return json.Marshal(res)
	})
}

var _ wasmkeeper.WasmVMQueryHandler = QueryHandlerFn(nil)

// QueryHandlerFn helper type that implements wasmkeeper.WasmVMQueryHandler
type QueryHandlerFn func(ctx sdk.Context, caller sdk.AccAddress, request wasmvmtypes.QueryRequest) ([]byte, error)

// HandleQuery handles contract query
func (q QueryHandlerFn) HandleQuery(ctx sdk.Context, caller sdk.AccAddress, request wasmvmtypes.QueryRequest) ([]byte, error) {
	return q(ctx, caller, request)
}
