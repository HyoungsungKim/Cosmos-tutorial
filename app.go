package app

import (
	"os"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/x/auth"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const appName = "nameservice"

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.nscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.nsd")
	ModuleBasics    sdk.ModuleBasicManager
)

type nameServiceApp struct {
	*bam.BaseApp
}

func NewNameServiceApp(longer log.Logger, db dbm.DB) *nameServiceApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &nameServiceApp{
		BaseApp: bApp,
		cdc:     cdc,
	}

	return app
}
