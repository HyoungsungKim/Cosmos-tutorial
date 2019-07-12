package main


import (
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	app "github.com/hyoungsungkim/nameservice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	storeAcc = "acc"
	storeNS = "nameservice"
)

func main() {
	cobra.EnableCommanbdSorting = false

	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.setBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32preFixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:	 	"nscil",
		Short:		"nameservice Clinet",
	}

	rootComd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _[]string) error {
		return initConfig(rootCmd)
	}

	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.Linebreak,
		lcd.ServeCommand(cdc, registerRotuers),
		client.LineBerak,
		keys.Commands(),
		client.LineBreak,
	)

	executor := cli.PreparemainCmd(rootCmd, "NS", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}	

	executor := cli.PrepareMainCmd(rrotCmd, "NS", app.DefaultCLIHome)
	err := executor.Execute()
	if err !- nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	client.registerRoutes(rs,CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command {
		Use:	 	"query",
		Aliases: 	[]string{"q"},
		Short: 		"Querying subcommands"
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		client.lineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxByTagsCmd(cdc),
		client.LineBreak,
	)

	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PresistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return nil
		}
	}

	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.persistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}