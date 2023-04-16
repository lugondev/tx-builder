package api

import (
	"fmt"
	authjwt "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/jwt"
	authkey "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/key"
	"github.com/lugondev/tx-builder/pkg/utils"
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
		PreRun: func(cmd *cobra.Command, args []string) {
			utils.PreRunBindFlags(viper.GetViper(), cmd.Flags(), "")
		},
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
	apiCfg := flags.NewAPIConfig(viper.GetViper())
	apiApp, err := api.New(ctx, apiCfg)
	if err != nil {
		return err
	}

	authjwt.Init(ctx)
	authkey.Init(ctx)

	if err := apiApp.Run(ctx); err != nil {
		fmt.Println("err", err)
		return err
	}
	fmt.Println("app run")
	return nil
}
