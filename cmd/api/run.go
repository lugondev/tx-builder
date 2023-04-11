package api

import (
	authjwt "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/jwt"
	authkey "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/key"
	"os"

	"github.com/lugondev/tx-builder/cmd/flags"
	"github.com/lugondev/tx-builder/src/api"

	"github.com/lugondev/tx-builder/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdErr error

func newRunCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		RunE:  run,
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := errors.CombineErrors(cmdErr, cmd.Context().Err()); err != nil {
				os.Exit(1)
			}
		},
	}

	flags.NewAPIFlags(runCmd.Flags())

	return runCmd
}

func run(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	apiApp, err := api.New(ctx, flags.NewAPIConfig(viper.GetViper()))
	if err != nil {
		return err
	}

	//qkmClient, err := api.QKMClient(apiApp.GetConfig())
	//if err != nil {
	//	return  err
	//}
	//
	//postgresClient, err := gopg.New("orchestrate.api", apiApp.GetConfig().Postgres)
	//if err != nil {
	//	return  err
	//}

	authjwt.Init(ctx)
	authkey.Init(ctx)

	return apiApp.Run(ctx)
}
