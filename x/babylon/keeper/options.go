package keeper

import wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

// option that is applied after keeper is setup with the VM. Used for decorators mainly.
type postOptsFn func(*Keeper)

func (f postOptsFn) apply(keeper *Keeper) {
	f(keeper)
}

// WithWasmKeeperDecorated can set a decorator to the wasm keeper
func WithWasmKeeperDecorated(cb func(*wasmkeeper.Keeper) *wasmkeeper.Keeper) Option {
	return postOptsFn(func(keeper *Keeper) {
		keeper.wasm = cb(keeper.wasm)
	})
}
