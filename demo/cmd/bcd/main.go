package main

import (
	"os"

	"cosmossdk.io/log"
	"github.com/babylonlabs-io/babylon-sdk/demo/app"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		log.NewLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
