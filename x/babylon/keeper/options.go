package keeper

// option that is applied after keeper is setup with the VM. Used for decorators mainly.
type postOptsFn func(*Keeper)

func (f postOptsFn) apply(keeper *Keeper) {
	f(keeper)
}
