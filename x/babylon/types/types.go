package types

import (
	fmt "fmt"
	"os"

	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
)

func GetGZippedContractCode(path string) ([]byte, error) {
	wasm, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// gzip the wasm file
	if ioutils.IsWasm(wasm) {
		wasm, err = ioutils.GzipIt(wasm)

		if err != nil {
			return nil, err
		}
	} else if !ioutils.IsGzip(wasm) {
		return nil, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}

	return wasm, nil
}
