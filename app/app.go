package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"encoding/json"
	"github.com/dengwenqi123/magic-eye/types"
)

const (
	appName = "MagicEyeApp"
)

type MagicEyeApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	//Keys to accesss the substores
	keyMain    *sdk.KVStoreKey
	KeyAccount *sdk.KVStoreKey

	//Manage getting and setting accounts
	accountMapper sdk.AccountMapper
	coinKeeper    bank.Keeper
}

func NewMagicEyeApp(logger log.Logger, db dbm.DB) *MagicEyeApp {

	//Create app-level codec for txs and accounts.
	var cdc = MakeCodec()

	var app = &MagicEyeApp{
		BaseApp:    bam.NewBaseApp(appName, cdc, logger, db),
		cdc:        cdc,
		keyMain:    sdk.NewKVStoreKey("main"),
		KeyAccount: sdk.NewKVStoreKey("acc"),
	}

	//Define the accountMapper.
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.KeyAccount,      // target store
		&types.MgyAccount{}, // prototype
	)

	//
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.Router().AddRoute("bank", bank.NewHandler(app.coinKeeper))

	//
	app.SetInitChainer(app.initChainer)
	app.MountStoresIAVL(app.keyMain, app.KeyAccount)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, auth.BurnFeeHandler))
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}

//Custom tx codec
func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()
	wire.RegisterCrypto(cdc) //Register crypto.
	sdk.RegisterWire(cdc)    // Register Msgs
	bank.RegisterWire(cdc)

	//register custom MgyAccount
	cdc.RegisterInterface((*sdk.Account)(nil), nil)
	cdc.RegisterConcrete(&types.MgyAccount{}, "basecoin/Account", nil)
	return cdc

}

func (app *MagicEyeApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToMGyAccount()
		if err != nil {
			panic(err)
		}
		app.accountMapper.SetAccount(ctx, acc)
	}
	return abci.ResponseInitChain{}
}

//Custom logic for state export
func (app *MagicEyeApp) ExportAppStateJSON() (appState json.RawMessage, err error) {
	//ctx := app.NewContext(true,abci.Header{})
	//app.NewContext(true,abci.Header{})
	ctx := app.NewContext(true, abci.Header{})

	accounts := []*types.GenesisAccount{}
	appendAccount := func(acc sdk.Account) (stop bool) {
		account := &types.GenesisAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)

	genState := types.GenesisState{
		Accounts: accounts,
	}
	return wire.MarshalJSONIndent(app.cdc, genState)
}
