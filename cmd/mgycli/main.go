package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/dengwenqi123/magic-eye/app"
	"github.com/dengwenqi123/magic-eye/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tmlibs/cli"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mgycli",
		Short: "Mgycoin light-client",
	}
)

func main() {
	//disable sorting
	cobra.EnableCommandSorting = false

	//get the codec
	cdc := app.MakeCodec()

	rpc.AddCommands(rootCmd)
	rootCmd.AddCommand(client.LineBreak)
	tx.AddCommands(rootCmd, cdc)
	rootCmd.AddCommand(client.LineBreak)

	//add query/post commands (custom to binary)
	rootCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, types.GetAccountDecoder(cdc)),
		)...)

	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)

	executor := cli.PrepareMainCmd(rootCmd, "BC", os.ExpandEnv("$HOME/.mgycli"))
	executor.Execute()

}
