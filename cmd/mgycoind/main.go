package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/dengwenqi123/magic-eye/app"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	"os"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "mgycoind",
		Short:             "Mgycoin Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, server.DefaultAppInit,
		server.ConstructAppCreator(newApp, "mgycoin"),
		server.ConstructAppExporter(exportAppState, "mgycoin"))

	//prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.mgycoind")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB) abci.Application {
	return app.NewMagicEyeApp(logger, db)
}

func exportAppState(logger log.Logger, db dbm.DB) (json.RawMessage, error) {
	mapp := app.NewMagicEyeApp(logger, db)
	return mapp.ExportAppStateJSON()
}
